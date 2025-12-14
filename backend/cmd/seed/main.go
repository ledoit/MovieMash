package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://moviemash:password@localhost:5432/moviemash?sslmode=disable"
		log.Printf("Using default DATABASE_URL: %s", databaseURL)
	}

	// Connect to database
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Connected to database successfully!")

	// Read and execute seed SQL file
	seedSQL, err := os.ReadFile("../../migrations/002_seed_data.sql")
	if err != nil {
		log.Fatal("Failed to read seed file:", err)
	}

	// Execute seed SQL
	if _, err := db.Exec(string(seedSQL)); err != nil {
		log.Fatal("Failed to execute seed SQL:", err)
	}

	fmt.Println("Seed data inserted successfully!")
	fmt.Println("Sample movies, top 4 sets, and comparisons are ready.")
}

