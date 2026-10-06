package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	"moviemash/backend/internal/database"
)

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()

	if _, err := database.DB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)`); err != nil {
		log.Fatal(err)
	}

	dir := "migrations"
	if _, err := os.Stat(dir); err != nil {
		dir = filepath.Join("..", "..", "migrations")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("read migrations: %v", err)
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		var applied bool
		if err := database.DB.QueryRow(
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename = $1)`,
			name,
		).Scan(&applied); err != nil {
			log.Fatal(err)
		}
		if applied {
			log.Printf("skip %s", name)
			continue
		}

		body, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			log.Fatal(err)
		}

		tx, err := database.DB.Begin()
		if err != nil {
			log.Fatal(err)
		}
		if _, err := tx.Exec(string(body)); err != nil {
			_ = tx.Rollback()
			log.Fatalf("apply %s: %v", name, err)
		}
		if _, err := tx.Exec(
			`INSERT INTO schema_migrations (filename) VALUES ($1)`,
			name,
		); err != nil {
			_ = tx.Rollback()
			log.Fatal(err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}
		log.Printf("applied %s", name)
	}
}
