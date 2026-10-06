package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"moviemash/backend/internal/database"
	"moviemash/backend/internal/elo"
)

type Movie struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Year        int      `json:"year"`
	Director    string   `json:"director"`
	Poster      string   `json:"poster"`
	Elo         float64  `json:"elo"`
	Appearances int      `json:"appearances"`
	Wins        int      `json:"wins"`
	WinRate     *float64 `json:"win_rate,omitempty"`
}

type Top4Set struct {
	ID          int      `json:"id"`
	Movies      []Movie  `json:"movies"`
	Elo         float64  `json:"elo"`
	Appearances int      `json:"appearances"`
	Wins        int      `json:"wins"`
	WinRate     *float64 `json:"win_rate,omitempty"`
}

type Comparison struct {
	ID   int     `json:"id"`
	SetA Top4Set `json:"set_a"`
	SetB Top4Set `json:"set_b"`
}

type VoteRequest struct {
	ComparisonID int `json:"comparison_id"`
	WinnerSetID  int `json:"winner_set_id"`
}

type Profile struct {
	ID          int      `json:"id"`
	Username    string   `json:"username"`
	Elo         float64  `json:"elo"`
	Appearances int      `json:"appearances"`
	Wins        int      `json:"wins"`
	WinRate     *float64 `json:"win_rate,omitempty"`
	Control     string   `json:"control"`
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func GetComparison(c *gin.Context) {
	voter := ensureVoter(c)
	if voter == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start voter session"})
		return
	}

	setAID, setBID, err := pickPair(voter)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not enough top 4s"})
		return
	}
	if err != nil {
		log.Printf("Error picking comparison: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comparison"})
		return
	}

	var comparisonID int
	err = database.DB.QueryRow(
		`INSERT INTO comparisons (set_a_id, set_b_id) VALUES ($1, $2) RETURNING id`,
		setAID, setBID,
	).Scan(&comparisonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comparison"})
		return
	}

	setA, err := getTop4Set(setAID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch set A"})
		return
	}
	setB, err := getTop4Set(setBID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch set B"})
		return
	}

	c.JSON(http.StatusOK, Comparison{ID: comparisonID, SetA: setA, SetB: setB})
}

func CreateVote(c *gin.Context) {
	voter := ensureVoter(c)
	if voter == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start voter session"})
		return
	}

	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vote"})
		return
	}
	defer tx.Rollback()

	var setA, setB int
	err = tx.QueryRow(
		`SELECT set_a_id, set_b_id FROM comparisons WHERE id = $1 FOR UPDATE`,
		req.ComparisonID,
	).Scan(&setA, &setB)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comparison not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vote"})
		return
	}
	if req.WinnerSetID != setA && req.WinnerSetID != setB {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Winner is not in this comparison"})
		return
	}
	loserID := setB
	if req.WinnerSetID == setB {
		loserID = setA
	}

	var seen bool
	if err := tx.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM votes v
			JOIN comparisons c ON c.id = v.comparison_id
			WHERE v.voter_key = $1
			  AND (
			    (c.set_a_id = $2 AND c.set_b_id = $3)
			    OR (c.set_a_id = $3 AND c.set_b_id = $2)
			  )
		)`, voter, setA, setB).Scan(&seen); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vote"})
		return
	}
	if seen {
		c.JSON(http.StatusConflict, gin.H{"error": "Pair already voted"})
		return
	}

	var voteID int
	err = tx.QueryRow(`
		INSERT INTO votes (comparison_id, voter_key, actor, winner_set_id, loser_set_id)
		VALUES ($1, $2, 'guest', $3, $4)
		RETURNING id
	`, req.ComparisonID, voter, req.WinnerSetID, loserID).Scan(&voteID)
	if err != nil {
		log.Printf("Error inserting vote: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vote"})
		return
	}

	winnerElo, loserElo, err := elo.Apply(tx, req.WinnerSetID, loserID)
	if err != nil {
		log.Printf("Error updating elo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ratings"})
		return
	}
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             voteID,
		"message":        "Vote recorded",
		"winner_set_elo": winnerElo,
		"loser_set_elo":  loserElo,
	})
}

func GetTop4Leaderboard(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT t.id
		FROM top4_sets t
		JOIN users u ON u.id = t.user_id
		WHERE t.is_current
		  AND (
		    NOT EXISTS (
		      SELECT 1 FROM users p
		      JOIN top4_sets ts ON ts.user_id = p.id AND ts.is_current
		      WHERE p.control IN ('simulator', 'human')
		    )
		    OR u.control IN ('simulator', 'human')
		  )
		ORDER BY t.elo DESC, t.match_count DESC, t.id
	`)
	if err != nil {
		log.Printf("Error fetching top4 sets: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top4 sets"})
		return
	}
	defer rows.Close()

	sets := []Top4Set{}
	for rows.Next() {
		var setID int
		if err := rows.Scan(&setID); err != nil {
			continue
		}
		set, err := getTop4Set(setID)
		if err != nil {
			log.Printf("Error fetching set %d: %v", setID, err)
			continue
		}
		sets = append(sets, set)
	}
	c.JSON(http.StatusOK, sets)
}

func GetMoviesLeaderboard(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, title, year, COALESCE(director, ''), COALESCE(poster_url, ''),
		       elo, match_count, wins
		FROM movies
		ORDER BY elo DESC, match_count DESC, title, id
	`)
	if err != nil {
		log.Printf("Error fetching movies: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}
	defer rows.Close()

	movies := []Movie{}
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Year, &m.Director, &m.Poster, &m.Elo, &m.Appearances, &m.Wins); err != nil {
			continue
		}
		m.WinRate = elo.WinRate(m.Appearances, m.Wins)
		movies = append(movies, m)
	}
	c.JSON(http.StatusOK, movies)
}

func GetProfilesLeaderboard(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id,
		       COALESCE(NULLIF(display_name, ''), letterboxd_username, username, ''),
		       elo, match_count, wins, control
		FROM users
		WHERE control <> 'system'
		ORDER BY elo DESC, match_count DESC, id
	`)
	if err != nil {
		log.Printf("Error fetching profiles: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profiles"})
		return
	}
	defer rows.Close()

	profiles := []Profile{}
	for rows.Next() {
		var p Profile
		if err := rows.Scan(&p.ID, &p.Username, &p.Elo, &p.Appearances, &p.Wins, &p.Control); err != nil {
			continue
		}
		p.WinRate = elo.WinRate(p.Appearances, p.Wins)
		profiles = append(profiles, p)
	}
	c.JSON(http.StatusOK, profiles)
}

func pickPair(voter string) (int, int, error) {
	query := `
		SELECT a.id, b.id
		FROM top4_sets a
		JOIN users ua ON ua.id = a.user_id
		JOIN top4_sets b ON a.id < b.id
		JOIN users ub ON ub.id = b.user_id
		WHERE a.is_current AND b.is_current
		  AND (
		    NOT EXISTS (
		      SELECT 1 FROM users p
		      JOIN top4_sets ts ON ts.user_id = p.id AND ts.is_current
		      WHERE p.control IN ('simulator', 'human')
		    )
		    OR (ua.control IN ('simulator', 'human') AND ub.control IN ('simulator', 'human'))
		  )
		  AND NOT EXISTS (
			SELECT 1
			FROM votes v
			JOIN comparisons c ON c.id = v.comparison_id
			WHERE v.voter_key = $1
			  AND (
			    (c.set_a_id = a.id AND c.set_b_id = b.id)
			    OR (c.set_a_id = b.id AND c.set_b_id = a.id)
			  )
		)
		ORDER BY RANDOM()
		LIMIT 1
	`
	var a, b int
	err := database.DB.QueryRow(query, voter).Scan(&a, &b)
	if err == sql.ErrNoRows {
		err = database.DB.QueryRow(`
			SELECT a.id, b.id
			FROM top4_sets a
			JOIN users ua ON ua.id = a.user_id
			JOIN top4_sets b ON a.id < b.id
			JOIN users ub ON ub.id = b.user_id
			WHERE a.is_current AND b.is_current
			  AND (
			    NOT EXISTS (
			      SELECT 1 FROM users p
			      JOIN top4_sets ts ON ts.user_id = p.id AND ts.is_current
			      WHERE p.control IN ('simulator', 'human')
			    )
			    OR (ua.control IN ('simulator', 'human') AND ub.control IN ('simulator', 'human'))
			  )
			ORDER BY RANDOM()
			LIMIT 1
		`).Scan(&a, &b)
	}
	return a, b, err
}

func getTop4Set(setID int) (Top4Set, error) {
	query := `
		SELECT
			t.id,
			t.elo,
			t.match_count,
			t.wins,
			array_agg(m.id ORDER BY array_position(t.movie_ids, m.id)),
			array_agg(m.title ORDER BY array_position(t.movie_ids, m.id)),
			array_agg(m.year ORDER BY array_position(t.movie_ids, m.id)),
			array_agg(COALESCE(m.director, '') ORDER BY array_position(t.movie_ids, m.id)),
			array_agg(COALESCE(m.poster_url, '') ORDER BY array_position(t.movie_ids, m.id)),
			array_agg(m.elo ORDER BY array_position(t.movie_ids, m.id)),
			array_agg(m.match_count ORDER BY array_position(t.movie_ids, m.id)),
			array_agg(m.wins ORDER BY array_position(t.movie_ids, m.id))
		FROM top4_sets t
		JOIN movies m ON m.id = ANY(t.movie_ids)
		WHERE t.id = $1
		GROUP BY t.id
	`

	var set Top4Set
	var movieIDs, years, appearances, wins pq.Int32Array
	var titles, directors, posters pq.StringArray
	var elos pq.Float64Array

	err := database.DB.QueryRow(query, setID).Scan(
		&set.ID,
		&set.Elo,
		&set.Appearances,
		&set.Wins,
		&movieIDs,
		&titles,
		&years,
		&directors,
		&posters,
		&elos,
		&appearances,
		&wins,
	)
	if err != nil {
		return set, err
	}
	set.WinRate = elo.WinRate(set.Appearances, set.Wins)

	set.Movies = make([]Movie, len(movieIDs))
	for i := range movieIDs {
		movie := Movie{
			ID:          int(movieIDs[i]),
			Title:       titles[i],
			Year:        int(years[i]),
			Director:    directors[i],
			Poster:      posters[i],
			Elo:         elos[i],
			Appearances: int(appearances[i]),
			Wins:        int(wins[i]),
		}
		movie.WinRate = elo.WinRate(movie.Appearances, movie.Wins)
		set.Movies[i] = movie
	}
	return set, nil
}
