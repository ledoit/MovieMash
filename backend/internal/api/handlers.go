package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/lib/pq"
	"moviemash/backend/internal/database"
)

type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
	Poster   string `json:"poster"`
}

type Top4Set struct {
	ID     int     `json:"id"`
	Movies []Movie `json:"movies"`
}

type Comparison struct {
	ID    int     `json:"id"`
	SetA  Top4Set `json:"set_a"`
	SetB  Top4Set `json:"set_b"`
}

type VoteRequest struct {
	ComparisonID int `json:"comparison_id"`
	WinnerSetID  int `json:"winner_set_id"`
}

// GetComparison returns a random comparison of two top 4 sets
func GetComparison(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT 
			c.id as comparison_id,
			sa.id as set_a_id,
			sb.id as set_b_id
		FROM comparisons c
		JOIN top4_sets sa ON c.set_a_id = sa.id
		JOIN top4_sets sb ON c.set_b_id = sb.id
		ORDER BY RANDOM()
		LIMIT 1
	`

	var comparisonID, setAID, setBID int
	err := database.DB.QueryRow(query).Scan(&comparisonID, &setAID, &setBID)
	if err == sql.ErrNoRows {
		// No comparisons exist, create one from random sets
		err = createRandomComparison(&comparisonID, &setAID, &setBID)
		if err != nil {
			http.Error(w, "Failed to create comparison: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		log.Printf("Error fetching comparison: %v", err)
		http.Error(w, "Failed to fetch comparison", http.StatusInternalServerError)
		return
	}

	// Fetch movies for set A
	setA, err := getTop4Set(setAID)
	if err != nil {
		http.Error(w, "Failed to fetch set A: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch movies for set B
	setB, err := getTop4Set(setBID)
	if err != nil {
		http.Error(w, "Failed to fetch set B: "+err.Error(), http.StatusInternalServerError)
		return
	}

	comparison := Comparison{
		ID:   comparisonID,
		SetA: setA,
		SetB: setB,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comparison)
}

// CreateVote handles vote submission
func CreateVote(w http.ResponseWriter, r *http.Request) {
	var req VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get or create anonymous user
	anonUserID, err := getOrCreateAnonymousUser()
	if err != nil {
		http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert vote
	insertQuery := `
		INSERT INTO votes (comparison_id, user_id, winner_set_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var voteID int
	err = database.DB.QueryRow(insertQuery, req.ComparisonID, anonUserID, req.WinnerSetID).Scan(&voteID)
	if err != nil {
		log.Printf("Error inserting vote: %v", err)
		http.Error(w, "Failed to create vote", http.StatusInternalServerError)
		return
	}

	// Create a new comparison for next time
	// (This ensures fresh comparisons)
	go createRandomComparison(nil, nil, nil)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": voteID,
		"message": "Vote recorded",
	})
}

// GetTop4Leaderboard returns all top 4 sets (for now, without ranking)
func GetTop4Leaderboard(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id
		FROM top4_sets
		ORDER BY id
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error fetching top4 sets: %v", err)
		http.Error(w, "Failed to fetch top4 sets", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var sets []Top4Set
	for rows.Next() {
		var setID int
		if err := rows.Scan(&setID); err != nil {
			log.Printf("Error scanning set: %v", err)
			continue
		}

		set, err := getTop4Set(setID)
		if err != nil {
			log.Printf("Error fetching set %d: %v", setID, err)
			continue
		}

		sets = append(sets, set)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sets)
}

// GetMoviesLeaderboard returns all movies (for now, without ranking)
func GetMoviesLeaderboard(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT DISTINCT ON (title, year, director) id, title, year, director, poster_url
		FROM movies
		ORDER BY title, year, director, id
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error fetching movies: %v", err)
		http.Error(w, "Failed to fetch movies", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		var posterURL sql.NullString
		if err := rows.Scan(&m.ID, &m.Title, &m.Year, &m.Director, &posterURL); err != nil {
			log.Printf("Error scanning movie: %v", err)
			continue
		}

		if posterURL.Valid {
			m.Poster = posterURL.String
		} else {
			m.Poster = ""
		}

		movies = append(movies, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// Helper functions

func getTop4Set(setID int) (Top4Set, error) {
	query := `
		SELECT 
			t.id,
			array_agg(m.id ORDER BY array_position(t.movie_ids, m.id)) as movie_ids,
			array_agg(m.title ORDER BY array_position(t.movie_ids, m.id)) as titles,
			array_agg(m.year ORDER BY array_position(t.movie_ids, m.id)) as years,
			array_agg(COALESCE(m.director, '') ORDER BY array_position(t.movie_ids, m.id)) as directors,
			array_agg(COALESCE(m.poster_url, '') ORDER BY array_position(t.movie_ids, m.id)) as posters
		FROM top4_sets t
		JOIN movies m ON m.id = ANY(t.movie_ids)
		WHERE t.id = $1
		GROUP BY t.id
	`

	var set Top4Set
	var movieIDs pq.Int32Array
	var titles, directors, posters pq.StringArray
	var years pq.Int32Array

	err := database.DB.QueryRow(query, setID).Scan(
		&set.ID,
		&movieIDs,
		&titles,
		&years,
		&directors,
		&posters,
	)
	if err != nil {
		return set, err
	}

	set.Movies = make([]Movie, len(movieIDs))
	for i := range movieIDs {
		set.Movies[i] = Movie{
			ID:       int(movieIDs[i]),
			Title:    titles[i],
			Year:     int(years[i]),
			Director: directors[i],
			Poster:   posters[i],
		}
	}

	return set, nil
}

func createRandomComparison(comparisonID, setAID, setBID *int) error {
	// Get two random top 4 sets
	query := `
		SELECT id FROM top4_sets
		ORDER BY RANDOM()
		LIMIT 2
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
	}

	if len(ids) < 2 {
		return sql.ErrNoRows
	}

	// Create comparison
	insertQuery := `
		INSERT INTO comparisons (set_a_id, set_b_id)
		VALUES ($1, $2)
		RETURNING id
	`

	var newComparisonID int
	err = database.DB.QueryRow(insertQuery, ids[0], ids[1]).Scan(&newComparisonID)
	if err != nil {
		return err
	}

	if comparisonID != nil {
		*comparisonID = newComparisonID
	}
	if setAID != nil {
		*setAID = ids[0]
	}
	if setBID != nil {
		*setBID = ids[1]
	}

	return nil
}

func getOrCreateAnonymousUser() (int, error) {
	// Try to get existing anonymous user
	var userID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE username = 'anonymous'").Scan(&userID)
	if err == nil {
		return userID, nil
	}

	// Create anonymous user if doesn't exist
	err = database.DB.QueryRow(
		"INSERT INTO users (username, email) VALUES ('anonymous', 'anonymous@moviemash.local') RETURNING id",
	).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

