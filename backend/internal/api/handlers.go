package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ComparisonResponse struct {
	ID      int64   `json:"id"`
	SetA    Top4Set `json:"set_a"`
	SetB    Top4Set `json:"set_b"`
	VotesA  int     `json:"votes_a"`
	VotesB  int     `json:"votes_b"`
}

type Top4Set struct {
	ID     int64   `json:"id"`
	Movies []Movie `json:"movies"`
}

type Movie struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
	Poster   string `json:"poster,omitempty"`
}

type MovieRanking struct {
	Rank        int    `json:"rank"`
	Movie       Movie  `json:"movie"`
	Wins        int    `json:"wins"`
	Appearances int    `json:"appearances"`
	WinRate     float64 `json:"win_rate"`
}

type VoteRequest struct {
	ComparisonID int64 `json:"comparison_id" binding:"required"`
	WinnerSetID  int64 `json:"winner_set_id" binding:"required"`
}

func (h *Handler) GetComparison(c *gin.Context) {
	// Get an active comparison or create a new one
	var comparisonID, setAID, setBID, votesA, votesB int64
	var expiresAt time.Time

	// Try to find an active (non-expired) comparison
	err := h.db.QueryRow(
		`SELECT id, set_a_id, set_b_id, votes_a, votes_b, expires_at 
		 FROM comparisons 
		 WHERE expires_at > NOW() 
		 ORDER BY RANDOM() 
		 LIMIT 1`,
	).Scan(&comparisonID, &setAID, &setBID, &votesA, &votesB, &expiresAt)

	// If no active comparison, create a new one
	if err == sql.ErrNoRows {
		// Get two random top 4 sets
		var setA, setB int64
		err = h.db.QueryRow(
			`SELECT id FROM top4_sets ORDER BY RANDOM() LIMIT 1`,
		).Scan(&setA)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No top 4 sets available"})
			return
		}

		err = h.db.QueryRow(
			`SELECT id FROM top4_sets WHERE id != $1 ORDER BY RANDOM() LIMIT 1`,
			setA,
		).Scan(&setB)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not enough top 4 sets available"})
			return
		}

		// Create new comparison
		err = h.db.QueryRow(
			`INSERT INTO comparisons (set_a_id, set_b_id, votes_a, votes_b, created_at, expires_at)
			 VALUES ($1, $2, 0, 0, NOW(), NOW() + INTERVAL '1 hour')
			 RETURNING id, set_a_id, set_b_id, votes_a, votes_b`,
			setA, setB,
		).Scan(&comparisonID, &setAID, &setBID, &votesA, &votesB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comparison"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query comparison"})
		return
	}

	// Fetch the actual top 4 sets with movies
	setA, err := h.fetchTop4Set(setAID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch set A"})
		return
	}

	setB, err := h.fetchTop4Set(setBID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch set B"})
		return
	}

	c.JSON(http.StatusOK, ComparisonResponse{
		ID:     comparisonID,
		SetA:   setA,
		SetB:   setB,
		VotesA: int(votesA),
		VotesB: int(votesB),
	})
}

func (h *Handler) fetchTop4Set(setID int64) (Top4Set, error) {
	var movieIDsStr string
	err := h.db.QueryRow(
		`SELECT movie_ids::text FROM top4_sets WHERE id = $1`,
		setID,
	).Scan(&movieIDsStr)
	if err != nil {
		return Top4Set{}, err
	}

	// Parse PostgreSQL array format: {1,2,3,4}
	movieIDsStr = movieIDsStr[1 : len(movieIDsStr)-1] // Remove { and }
	var movieIDs []int64
	if movieIDsStr != "" {
		parts := strings.Split(movieIDsStr, ",")
		for _, part := range parts {
			var id int64
			if _, err := fmt.Sscanf(strings.TrimSpace(part), "%d", &id); err == nil {
				movieIDs = append(movieIDs, id)
			}
		}
	}

	// Fetch movie details
	movies := make([]Movie, 0, len(movieIDs))
	for _, movieID := range movieIDs {
		var movie Movie
		var poster sql.NullString
		err := h.db.QueryRow(
			`SELECT id, title, year, director, poster_url FROM movies WHERE id = $1`,
			movieID,
		).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Director, &poster)
		if err != nil {
			continue // Skip missing movies
		}
		if poster.Valid {
			movie.Poster = poster.String
		}
		movies = append(movies, movie)
	}

	return Top4Set{
		ID:     setID,
		Movies: movies,
	}, nil
}

func (h *Handler) CreateVote(c *gin.Context) {
	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get or create anonymous user (TODO: implement auth middleware)
	var userID int64
	err := h.db.QueryRow(
		"SELECT id FROM users WHERE username = 'anonymous' LIMIT 1",
	).Scan(&userID)
	if err == sql.ErrNoRows {
		// Create anonymous user if doesn't exist
		err = h.db.QueryRow(
			"INSERT INTO users (username, email, password_hash) VALUES ('anonymous', 'anonymous@moviemash.com', 'dummy') RETURNING id",
		).Scan(&userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Query comparison to get set IDs
	var setAID, setBID int64
	err = h.db.QueryRow(
		"SELECT set_a_id, set_b_id FROM comparisons WHERE id = $1",
		req.ComparisonID,
	).Scan(&setAID, &setBID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comparison not found"})
		return
	}

	// Store vote directly in database (user_id can be NULL for anonymous votes)
	_, err = h.db.Exec(
		"INSERT INTO votes (user_id, comparison_id, winner_set_id, timestamp) VALUES ($1, $2, $3, $4)",
		userID, req.ComparisonID, req.WinnerSetID, time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record vote: " + err.Error()})
		return
	}

	// Update comparison vote counts
	var votesA, votesB int
	if req.WinnerSetID == setAID {
		_, err = h.db.Exec(
			"UPDATE comparisons SET votes_a = votes_a + 1 WHERE id = $1",
			req.ComparisonID,
		)
	} else {
		_, err = h.db.Exec(
			"UPDATE comparisons SET votes_b = votes_b + 1 WHERE id = $1",
			req.ComparisonID,
		)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vote counts"})
		return
	}

	// Update cache
	err = h.db.QueryRow(
		"SELECT votes_a, votes_b FROM comparisons WHERE id = $1",
		req.ComparisonID,
	).Scan(&votesA, &votesB)
	if err == nil {
		h.cache.UpdateVotes(req.ComparisonID, votesA, votesB)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote recorded"})
}

func (h *Handler) GetRecommendations(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// TODO: Query recommendations from database
	// If stale/missing, trigger ML service

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"recommendations": []Movie{},
	})
}

func (h *Handler) GetLeaderboard(c *gin.Context) {
	// Query to get movie rankings based on votes
	// This aggregates:
	// 1. Movies in winning sets (wins)
	// 2. Movies in all sets that were compared (appearances)
	// 3. Calculate win rate
	query := `
		WITH winning_movies AS (
			-- Count wins: distinct comparisons where the movie's set won
			SELECT 
				UNNEST(t.movie_ids) as movie_id,
				COUNT(DISTINCT v.comparison_id) as wins
			FROM votes v
			JOIN top4_sets t ON v.winner_set_id = t.id
			GROUP BY movie_id
		),
		all_appearances AS (
			-- Count appearances: distinct comparisons where the movie appeared
			SELECT 
				UNNEST(t.movie_ids) as movie_id,
				COUNT(DISTINCT c.id) as appearances
			FROM comparisons c
			JOIN top4_sets t ON (c.set_a_id = t.id OR c.set_b_id = t.id)
			GROUP BY movie_id
		)
		SELECT 
			m.id,
			m.title,
			m.year,
			m.director,
			COALESCE(m.poster_url, '') as poster,
			COALESCE(wm.wins, 0) as wins,
			COALESCE(aa.appearances, 0) as appearances,
			CASE 
				WHEN COALESCE(aa.appearances, 0) > 0 
				THEN COALESCE(wm.wins, 0)::float / aa.appearances::float
				ELSE 0.0
			END as win_rate
		FROM movies m
		LEFT JOIN winning_movies wm ON m.id = wm.movie_id
		LEFT JOIN all_appearances aa ON m.id = aa.movie_id
		WHERE COALESCE(aa.appearances, 0) > 0
		ORDER BY wins DESC, win_rate DESC, appearances DESC
		LIMIT 100
	`

	rows, err := h.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query leaderboard"})
		return
	}
	defer rows.Close()

	var rankings []MovieRanking
	rank := 1
	for rows.Next() {
		var ranking MovieRanking
		var poster sql.NullString
		
		err := rows.Scan(
			&ranking.Movie.ID,
			&ranking.Movie.Title,
			&ranking.Movie.Year,
			&ranking.Movie.Director,
			&poster,
			&ranking.Wins,
			&ranking.Appearances,
			&ranking.WinRate,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan leaderboard data"})
			return
		}

		if poster.Valid {
			ranking.Movie.Poster = poster.String
		}
		ranking.Rank = rank
		rankings = append(rankings, ranking)
		rank++
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating leaderboard results"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rankings": rankings,
	})
}

func (h *Handler) GetMovie(c *gin.Context) {
	movieIDStr := c.Param("id")
	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	// TODO: Query movie from database

	c.JSON(http.StatusOK, gin.H{
		"id": movieID,
		"message": "Get movie - TODO",
	})
}

func (h *Handler) Register(c *gin.Context) {
	// TODO: Implement user registration
	c.JSON(http.StatusOK, gin.H{"message": "Register - TODO"})
}

func (h *Handler) Login(c *gin.Context) {
	// TODO: Implement user login with JWT
	c.JSON(http.StatusOK, gin.H{"message": "Login - TODO"})
}

