package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	databaseURL := getEnv("DATABASE_URL", "postgres://moviemash:password@localhost:5432/proji?sslmode=disable")

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Sample movies
	movies := []struct {
		Title    string
		Year     int
		Director string
		Genres   []string
	}{
		{"The Shawshank Redemption", 1994, "Frank Darabont", []string{"Drama"}},
		{"The Godfather", 1972, "Francis Ford Coppola", []string{"Crime", "Drama"}},
		{"Pulp Fiction", 1994, "Quentin Tarantino", []string{"Crime", "Drama"}},
		{"The Dark Knight", 2008, "Christopher Nolan", []string{"Action", "Crime", "Drama"}},
		{"Fight Club", 1999, "David Fincher", []string{"Drama"}},
		{"Inception", 2010, "Christopher Nolan", []string{"Action", "Sci-Fi", "Thriller"}},
		{"Goodfellas", 1990, "Martin Scorsese", []string{"Biography", "Crime", "Drama"}},
		{"The Matrix", 1999, "Lana Wachowski, Lilly Wachowski", []string{"Action", "Sci-Fi"}},
		{"Interstellar", 2014, "Christopher Nolan", []string{"Adventure", "Drama", "Sci-Fi"}},
		{"Parasite", 2019, "Bong Joon-ho", []string{"Comedy", "Drama", "Thriller"}},
		{"Whiplash", 2014, "Damien Chazelle", []string{"Drama", "Music"}},
		{"Mad Max: Fury Road", 2015, "George Miller", []string{"Action", "Adventure", "Sci-Fi"}},
		{"The Grand Budapest Hotel", 2014, "Wes Anderson", []string{"Adventure", "Comedy", "Drama"}},
		{"Her", 2013, "Spike Jonze", []string{"Drama", "Romance", "Sci-Fi"}},
		{"La La Land", 2016, "Damien Chazelle", []string{"Comedy", "Drama", "Music"}},
		{"Get Out", 2017, "Jordan Peele", []string{"Horror", "Mystery", "Thriller"}},
		{"Moonlight", 2016, "Barry Jenkins", []string{"Drama"}},
		{"The Social Network", 2010, "David Fincher", []string{"Biography", "Drama"}},
		{"No Country for Old Men", 2007, "Joel Coen, Ethan Coen", []string{"Crime", "Drama", "Thriller"}},
		{"There Will Be Blood", 2007, "Paul Thomas Anderson", []string{"Drama"}},
	}

	fmt.Println("Inserting movies...")
	movieIDs := make([]int64, 0, len(movies))
	for _, movie := range movies {
		var movieID int64
		err := db.QueryRow(
			`INSERT INTO movies (title, year, director, genres, letterboxd_id, created_at)
			 VALUES ($1, $2, $3, $4, $5, NOW())
			 ON CONFLICT (letterboxd_id) DO UPDATE SET title = $1
			 RETURNING id`,
			movie.Title, movie.Year, movie.Director, fmt.Sprintf("{%s}", movie.Genres[0]), 
			fmt.Sprintf("film/%s-%d", movie.Title, movie.Year),
		).Scan(&movieID)
		if err != nil {
			log.Printf("Error inserting movie %s: %v", movie.Title, err)
			continue
		}
		movieIDs = append(movieIDs, movieID)
		fmt.Printf("  ✓ %s (%d)\n", movie.Title, movie.Year)
	}

	fmt.Printf("\nInserted %d movies\n\n", len(movieIDs))

	// Create some top 4 sets
	fmt.Println("Creating top 4 sets...")
	rand.Seed(time.Now().UnixNano())
	top4Sets := make([]int64, 0)

	for i := 0; i < 10; i++ {
		// Randomly select 4 movies
		selected := make([]int64, 4)
		used := make(map[int]bool)
		for j := 0; j < 4; j++ {
			idx := rand.Intn(len(movieIDs))
			for used[idx] {
				idx = rand.Intn(len(movieIDs))
			}
			used[idx] = true
			selected[j] = movieIDs[idx]
		}

		var setID int64
		err := db.QueryRow(
			`INSERT INTO top4_sets (user_letterboxd_id, movie_ids, scraped_at, created_at)
			 VALUES ($1, $2, NOW(), NOW())
			 RETURNING id`,
			fmt.Sprintf("user%d", i+1),
			fmt.Sprintf("{%d,%d,%d,%d}", selected[0], selected[1], selected[2], selected[3]),
		).Scan(&setID)
		if err != nil {
			log.Printf("Error creating top 4 set: %v", err)
			continue
		}
		top4Sets = append(top4Sets, setID)
		fmt.Printf("  ✓ Created top 4 set %d\n", setID)
	}

	fmt.Printf("\nCreated %d top 4 sets\n\n", len(top4Sets))

	// Create some comparisons
	fmt.Println("Creating comparisons...")
	comparisons := make([]int64, 0)
	for i := 0; i < 5; i++ {
		if len(top4Sets) < 2 {
			break
		}
		setA := top4Sets[rand.Intn(len(top4Sets))]
		setB := top4Sets[rand.Intn(len(top4Sets))]
		for setB == setA {
			setB = top4Sets[rand.Intn(len(top4Sets))]
		}

		var compID int64
		err := db.QueryRow(
			`INSERT INTO comparisons (set_a_id, set_b_id, votes_a, votes_b, created_at, expires_at)
			 VALUES ($1, $2, 0, 0, NOW(), NOW() + INTERVAL '1 hour')
			 RETURNING id`,
			setA, setB,
		).Scan(&compID)
		if err != nil {
			log.Printf("Error creating comparison: %v", err)
			continue
		}
		comparisons = append(comparisons, compID)
		fmt.Printf("  ✓ Created comparison %d (set %d vs set %d)\n", compID, setA, setB)
	}

	fmt.Printf("\nCreated %d comparisons\n", len(comparisons))
	fmt.Println("\n✅ Seed data created successfully!")
	fmt.Println("\nYou can now:")
	fmt.Println("  1. Visit http://localhost:4200 to see the frontend")
	fmt.Println("  2. Use the Gladiators tab to see comparisons")
	fmt.Println("  3. Vote on comparisons to generate data")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

