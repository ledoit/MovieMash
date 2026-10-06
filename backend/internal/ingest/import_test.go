package ingest

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"moviemash/backend/internal/database"
	"moviemash/backend/internal/elo"
)

func TestImportAndRate(t *testing.T) {
	dir := t.TempDir()
	port := uint32(55432)
	pg := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Username("moviemash").
		Password("password").
		Database("moviemash").
		Version(embeddedpostgres.V16).
		Locale("C").
		Encoding("UTF8").
		Port(port).
		RuntimePath(filepath.Join(dir, "runtime")).
		DataPath(filepath.Join(dir, "data")))
	if err := pg.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = pg.Stop() })

	t.Setenv("DATABASE_URL", "postgres://moviemash:password@127.0.0.1:55432/moviemash?sslmode=disable")
	if err := database.InitDB(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = database.CloseDB() })

	if err := applyMigrations(database.DB); err != nil {
		t.Fatal(err)
	}

	body, err := os.ReadFile(repoPath("ingest", "letterboxd.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc Document
	if err := json.Unmarshal(body, &doc); err != nil {
		t.Fatal(err)
	}
	n, err := Import(database.DB, doc)
	if err != nil {
		t.Fatal(err)
	}
	if n != 6 {
		t.Fatalf("imported %d", n)
	}

	var profiles int
	if err := database.DB.QueryRow(`SELECT count(*) FROM users WHERE control = 'simulator'`).Scan(&profiles); err != nil {
		t.Fatal(err)
	}
	if profiles != 6 {
		t.Fatalf("profiles %d", profiles)
	}

	var interstellarID int
	var letterboxdID string
	err = database.DB.QueryRow(
		`SELECT id, letterboxd_id FROM movies WHERE title = 'Interstellar' AND year = 2014`,
	).Scan(&interstellarID, &letterboxdID)
	if err != nil {
		t.Fatal(err)
	}
	if letterboxdID != "interstellar" {
		t.Fatalf("seed movie was not reused: %s", letterboxdID)
	}

	var setA, setB int
	err = database.DB.QueryRow(`
		SELECT a.id, b.id
		FROM top4_sets a
		JOIN users ua ON ua.id = a.user_id AND ua.control = 'simulator'
		JOIN top4_sets b ON a.id < b.id
		JOIN users ub ON ub.id = b.user_id AND ub.control = 'simulator'
		WHERE a.is_current AND b.is_current
		ORDER BY a.id, b.id
		LIMIT 1
	`).Scan(&setA, &setB)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := database.DB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	winner, loser, err := elo.Apply(tx, setA, setB)
	if err != nil {
		_ = tx.Rollback()
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	if winner <= 1500 || loser >= 1500 {
		t.Fatalf("elo did not move: winner %v loser %v", winner, loser)
	}
}

func applyMigrations(db *sql.DB) error {
	dir := repoPath("backend", "migrations")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	for _, entry := range entries {
		body, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(body)); err != nil {
			return err
		}
	}
	return nil
}

func repoPath(parts ...string) string {
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(file), "..", "..", "..")
	return filepath.Join(append([]string{root}, parts...)...)
}
