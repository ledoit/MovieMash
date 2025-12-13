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
	"time"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

type TMDBMovie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Year     string `json:"release_date"`
	Poster   string `json:"poster_path"`
}

type TMDBResponse struct {
	Results []TMDBMovie `json:"results"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	databaseURL := getEnv("DATABASE_URL", "postgres://moviemash:password@localhost:5432/proji?sslmode=disable")
	tmdbAPIKey := getEnv("TMDB_API_KEY", "")

	if tmdbAPIKey == "" {
		log.Fatal("TMDB_API_KEY environment variable is required. Get a free API key from https://www.themoviedb.org/settings/api")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Get all movies without poster URLs
	rows, err := db.Query(`
		SELECT id, title, year 
		FROM movies 
		WHERE poster_url IS NULL OR poster_url = ''
		ORDER BY id
	`)
	if err != nil {
		log.Fatal("Failed to query movies:", err)
	}
	defer rows.Close()

	var movies []struct {
		ID    int
		Title string
		Year  int
	}

	for rows.Next() {
		var m struct {
			ID    int
			Title string
			Year  int
		}
		if err := rows.Scan(&m.ID, &m.Title, &m.Year); err != nil {
			log.Printf("Error scanning movie: %v", err)
			continue
		}
		movies = append(movies, m)
	}

	if len(movies) == 0 {
		fmt.Println("No movies need poster URLs")
		return
	}

	fmt.Printf("Found %d movies without poster URLs\n\n", len(movies))

	client := &http.Client{Timeout: 10 * time.Second}
	baseURL := "https://api.themoviedb.org/3/search/movie"

	successCount := 0
	for i, movie := range movies {
		fmt.Printf("[%d/%d] Fetching poster for: %s (%d)... ", i+1, len(movies), movie.Title, movie.Year)

		// Build search URL
		searchURL := fmt.Sprintf("%s?api_key=%s&query=%s&year=%d", 
			baseURL, 
			tmdbAPIKey, 
			url.QueryEscape(movie.Title),
			movie.Year)

		resp, err := client.Get(searchURL)
		if err != nil {
			fmt.Printf("ERROR: Request failed - %v\n", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("ERROR: Failed to read response - %v\n", err)
			continue
		}

		var tmdbResp TMDBResponse
		if err := json.Unmarshal(body, &tmdbResp); err != nil {
			fmt.Printf("ERROR: Failed to parse JSON - %v\n", err)
			continue
		}

		if len(tmdbResp.Results) == 0 {
			fmt.Printf("NOT FOUND\n")
			continue
		}

		// Use the first result (usually the most relevant)
		result := tmdbResp.Results[0]
		if result.Poster == "" {
			fmt.Printf("NO POSTER\n")
			continue
		}

		// Build full poster URL (TMDB provides relative paths)
		posterURL := fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", result.Poster)

		// Update database
		_, err = db.Exec(
			`UPDATE movies SET poster_url = $1, updated_at = NOW() WHERE id = $2`,
			posterURL, movie.ID,
		)
		if err != nil {
			fmt.Printf("ERROR: Failed to update database - %v\n", err)
			continue
		}

		fmt.Printf("✓ Updated\n")
		successCount++

		// Be nice to TMDB API - rate limit
		time.Sleep(250 * time.Millisecond)
	}

	fmt.Printf("\n✅ Successfully updated %d/%d movies with poster URLs\n", successCount, len(movies))
	fmt.Println("\nPoster URLs are now available in the database!")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

