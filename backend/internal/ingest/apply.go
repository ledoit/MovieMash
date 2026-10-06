package ingest

import (
	"database/sql"
	"fmt"
	"html"

	"github.com/lib/pq"
)

func Import(db *sql.DB, doc Document) (int, error) {
	imported := 0
	for _, profile := range doc.Profiles {
		if len(profile.Top4) != 4 {
			continue
		}
		if err := importProfile(db, profile); err != nil {
			return imported, fmt.Errorf("%s: %w", profile.Username, err)
		}
		imported++
	}
	return imported, nil
}

func importProfile(db *sql.DB, profile Profile) error {
	profile.Username = html.UnescapeString(profile.Username)
	profile.DisplayName = html.UnescapeString(profile.DisplayName)

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	userID, err := upsertProfile(tx, profile)
	if err != nil {
		return err
	}

	ids := make([]int, 0, 4)
	for _, film := range profile.Top4 {
		id, err := upsertMovie(tx, film)
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}
	if err := replaceTop4(tx, userID, ids); err != nil {
		return err
	}
	for _, film := range profile.Watches {
		if film.Year == 0 || film.LetterboxdID == "" {
			continue
		}
		movieID, err := upsertMovie(tx, film)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(
			`INSERT INTO watches (user_id, movie_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			userID, movieID,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func upsertProfile(tx *sql.Tx, profile Profile) (int, error) {
	var id int
	var control string
	err := tx.QueryRow(
		`SELECT id, control FROM users WHERE letterboxd_username = $1`,
		profile.Username,
	).Scan(&id, &control)
	if err == sql.ErrNoRows {
		err = tx.QueryRow(
			`INSERT INTO users (username, letterboxd_username, display_name, control)
			 VALUES ($1, $1, $2, 'simulator')
			 RETURNING id`,
			profile.Username, profile.DisplayName,
		).Scan(&id)
		return id, err
	}
	if err != nil {
		return 0, err
	}
	if control == "human" {
		_, err = tx.Exec(`UPDATE users SET display_name = $1 WHERE id = $2`, profile.DisplayName, id)
		return id, err
	}
	_, err = tx.Exec(
		`UPDATE users SET display_name = $1, control = 'simulator' WHERE id = $2 AND control <> 'human'`,
		profile.DisplayName, id,
	)
	return id, err
}

func upsertMovie(tx *sql.Tx, film Film) (int, error) {
	film.Title = html.UnescapeString(film.Title)
	var id int
	err := tx.QueryRow(`SELECT id FROM movies WHERE letterboxd_id = $1`, film.LetterboxdID).Scan(&id)
	if err == nil {
		_, err = tx.Exec(
			`UPDATE movies
			 SET title = $1, year = $2, poster_url = COALESCE(NULLIF($3, ''), poster_url)
			 WHERE id = $4`,
			film.Title, film.Year, film.Poster, id,
		)
		return id, err
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	err = tx.QueryRow(
		`SELECT id FROM movies WHERE lower(title) = lower($1) AND year = $2 ORDER BY id LIMIT 1`,
		film.Title, film.Year,
	).Scan(&id)
	if err == nil {
		_, err = tx.Exec(
			`UPDATE movies
			 SET letterboxd_id = COALESCE(letterboxd_id, $1),
			     poster_url = COALESCE(NULLIF($2, ''), poster_url)
			 WHERE id = $3`,
			film.LetterboxdID, film.Poster, id,
		)
		return id, err
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	err = tx.QueryRow(
		`INSERT INTO movies (title, year, poster_url, letterboxd_id)
		 VALUES ($1, $2, NULLIF($3, ''), $4)
		 RETURNING id`,
		film.Title, film.Year, film.Poster, film.LetterboxdID,
	).Scan(&id)
	return id, err
}

func replaceTop4(tx *sql.Tx, userID int, movieIDs []int) error {
	var setID int
	var existing pq.Int64Array
	err := tx.QueryRow(
		`SELECT id, movie_ids FROM top4_sets WHERE user_id = $1 AND is_current ORDER BY id DESC LIMIT 1`,
		userID,
	).Scan(&setID, &existing)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil && sameIDs(existing, movieIDs) {
		return nil
	}
	if err == nil {
		if _, err := tx.Exec(`UPDATE top4_sets SET is_current = FALSE WHERE id = $1`, setID); err != nil {
			return err
		}
	}
	_, err = tx.Exec(
		`INSERT INTO top4_sets (user_id, movie_ids, is_current) VALUES ($1, $2, TRUE)`,
		userID, pq.Array(movieIDs),
	)
	return err
}

func sameIDs(existing pq.Int64Array, next []int) bool {
	if len(existing) != len(next) {
		return false
	}
	for i := range next {
		if int(existing[i]) != next[i] {
			return false
		}
	}
	return true
}
