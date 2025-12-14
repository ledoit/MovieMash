package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"moviemash/backend/internal/database"
)

type TMDBResponse struct {
	Results []TMDBMovie `json:"results"`
}

type TMDBMovie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	PosterPath  string `json:"poster_path"`
}

type TMDBMovieDetails struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	PosterPath  string `json:"poster_path"`
	Director    string `json:"director"`
}

type CreditsResponse struct {
	Crew []CrewMember `json:"crew"`
}

type CrewMember struct {
	Job  string `json:"job"`
	Name string `json:"name"`
}

func main() {
	// Load .env file (try multiple locations)
	envPaths := []string{".env", "../.env", "../../.env", "../../../.env"}
	loaded := false
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			loaded = true
			break
		}
	}
	if !loaded {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		log.Fatal("TMDB_API_KEY not found in environment")
	}

	// Fetch movies that need poster updates
	// Check for missing, invalid, or duplicate poster URLs
	query := `
		SELECT id, title, year, director, poster_url
		FROM movies
		WHERE poster_url IS NULL 
		   OR poster_url = ''
		   OR poster_url LIKE '%placeholder%'
		   OR poster_url LIKE '%Z4Z4Z%'
		   OR poster_url NOT LIKE 'https://image.tmdb.org%'
		   OR poster_url IN (
			   SELECT poster_url 
			   FROM movies 
			   WHERE poster_url IS NOT NULL 
			   GROUP BY poster_url 
			   HAVING COUNT(*) > 1
		   )
		ORDER BY id
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Fatal("Failed to query movies:", err)
	}
	defer rows.Close()

	type MovieRow struct {
		ID        int
		Title     string
		Year      int
		Director  string
		PosterURL string
	}

	var movies []MovieRow
	for rows.Next() {
		var m MovieRow
		var posterURL sql.NullString
		err := rows.Scan(&m.ID, &m.Title, &m.Year, &m.Director, &posterURL)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		if posterURL.Valid {
			m.PosterURL = posterURL.String
		} else {
			m.PosterURL = ""
		}
		movies = append(movies, m)
	}

	if len(movies) == 0 {
		log.Println("No movies need poster updates")
		return
	}

	log.Printf("Found %d movies needing poster updates\n", len(movies))

	client := &http.Client{Timeout: 10 * time.Second}
	updated := 0

	for _, movie := range movies {
		log.Printf("Processing: %s (%d)", movie.Title, movie.Year)

		// Search for movie on TMDB
		searchURL := fmt.Sprintf(
			"https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s&year=%d",
			apiKey,
			url.QueryEscape(movie.Title),
			movie.Year,
		)

		resp, err := client.Get(searchURL)
		if err != nil {
			log.Printf("  Error searching: %v", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("  Error reading response: %v", err)
			continue
		}

		var tmdbResp TMDBResponse
		if err := json.Unmarshal(body, &tmdbResp); err != nil {
			log.Printf("  Error parsing JSON: %v", err)
			continue
		}

		if len(tmdbResp.Results) == 0 {
			log.Printf("  No results found")
			continue
		}

		// Use first result (best match)
		tmdbMovie := tmdbResp.Results[0]

		// Build poster URL
		var posterURL string
		if tmdbMovie.PosterPath != "" && tmdbMovie.PosterPath != "null" {
			// Ensure poster path starts with / and build full URL
			posterPath := tmdbMovie.PosterPath
			if !strings.HasPrefix(posterPath, "/") {
				posterPath = "/" + posterPath
			}
			posterURL = fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", posterPath)
		} else {
			log.Printf("  No poster path found")
			continue
		}

		// Update database
		updateQuery := `UPDATE movies SET poster_url = $1 WHERE id = $2`
		_, err = database.DB.Exec(updateQuery, posterURL, movie.ID)
		if err != nil {
			log.Printf("  Error updating database: %v", err)
			continue
		}

		log.Printf("  ✓ Updated poster: %s", posterURL)
		updated++

		// Rate limiting - be nice to TMDB API
		time.Sleep(250 * time.Millisecond)
	}

	log.Printf("\n✓ Updated %d/%d movies", updated, len(movies))
}

