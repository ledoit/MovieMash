package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://user:password@localhost:5432/proji?sslmode=disable"
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Read migration files
	// Try multiple possible paths
	wd, _ := os.Getwd()
	var migrationsDir string
	
	// Try from backend root (if running from backend/)
	if _, err := os.Stat(filepath.Join(wd, "migrations")); err == nil {
		migrationsDir = filepath.Join(wd, "migrations")
	} else if _, err := os.Stat(filepath.Join(wd, "..", "migrations")); err == nil {
		// Try from cmd/migrate directory
		migrationsDir = filepath.Join(wd, "..", "migrations")
	} else {
		// Try from backend/cmd/migrate
		migrationsDir = filepath.Join(wd, "..", "..", "migrations")
	}
	
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatal("Failed to read migrations directory:", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationPath := filepath.Join(migrationsDir, file.Name())
			sql, err := os.ReadFile(migrationPath)
			if err != nil {
				log.Printf("Failed to read %s: %v", file.Name(), err)
				continue
			}

			if _, err := db.Exec(string(sql)); err != nil {
				log.Printf("Failed to execute %s: %v", file.Name(), err)
				continue
			}

			fmt.Printf("Applied migration: %s\n", file.Name())
		}
	}

	fmt.Println("Migrations completed")
}

