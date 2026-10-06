package elo

import (
	"database/sql"

	"github.com/lib/pq"
)

type setRow struct {
	ID      int
	UserID  sql.NullInt64
	Card    Card
	MovieID []int
}

// Apply updates Elo for the two sets, each of their movies, and the two owning profiles.
// A profile is skipped when both sets belong to the same user.
func Apply(tx *sql.Tx, winnerSetID, loserSetID int) (winnerElo, loserElo float64, err error) {
	sets, err := lockSets(tx, winnerSetID, loserSetID)
	if err != nil {
		return 0, 0, err
	}
	winner := sets[winnerSetID]
	loser := sets[loserSetID]

	movieIDs := uniqueInts(append(append([]int{}, winner.MovieID...), loser.MovieID...))
	movies, err := lockMovies(tx, movieIDs)
	if err != nil {
		return 0, 0, err
	}

	wMovies := cardsFor(movies, winner.MovieID)
	lMovies := cardsFor(movies, loser.MovieID)
	wMovies, lMovies = Movies(wMovies, lMovies)
	for _, card := range append(wMovies, lMovies...) {
		if err := writeMovie(tx, card); err != nil {
			return 0, 0, err
		}
	}

	wSet, lSet := Bout(winner.Card, loser.Card, 1)
	if err := writeSet(tx, wSet); err != nil {
		return 0, 0, err
	}
	if err := writeSet(tx, lSet); err != nil {
		return 0, 0, err
	}

	if winner.UserID.Valid && loser.UserID.Valid && winner.UserID.Int64 != loser.UserID.Int64 {
		users, err := lockUsers(tx, []int{int(winner.UserID.Int64), int(loser.UserID.Int64)})
		if err != nil {
			return 0, 0, err
		}
		wu := users[int(winner.UserID.Int64)]
		lu := users[int(loser.UserID.Int64)]
		wu, lu = Bout(wu, lu, 1)
		if err := writeUser(tx, wu); err != nil {
			return 0, 0, err
		}
		if err := writeUser(tx, lu); err != nil {
			return 0, 0, err
		}
	}

	return wSet.Elo, lSet.Elo, nil
}

func lockSets(tx *sql.Tx, a, b int) (map[int]setRow, error) {
	rows, err := tx.Query(`
		SELECT id, user_id, elo, match_count, wins, movie_ids
		FROM top4_sets
		WHERE id IN ($1, $2)
		ORDER BY id
		FOR UPDATE
	`, a, b)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[int]setRow{}
	for rows.Next() {
		var row setRow
		var ids pq.Int64Array
		if err := rows.Scan(&row.ID, &row.UserID, &row.Card.Elo, &row.Card.Matches, &row.Card.Wins, &ids); err != nil {
			return nil, err
		}
		row.Card.ID = row.ID
		for _, id := range ids {
			row.MovieID = append(row.MovieID, int(id))
		}
		out[row.ID] = row
	}
	if len(out) != 2 {
		return nil, sql.ErrNoRows
	}
	return out, rows.Err()
}

func lockMovies(tx *sql.Tx, ids []int) (map[int]Card, error) {
	rows, err := tx.Query(`
		SELECT id, elo, match_count, wins
		FROM movies
		WHERE id = ANY($1)
		ORDER BY id
		FOR UPDATE
	`, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[int]Card{}
	for rows.Next() {
		var card Card
		if err := rows.Scan(&card.ID, &card.Elo, &card.Matches, &card.Wins); err != nil {
			return nil, err
		}
		out[card.ID] = card
	}
	return out, rows.Err()
}

func lockUsers(tx *sql.Tx, ids []int) (map[int]Card, error) {
	rows, err := tx.Query(`
		SELECT id, elo, match_count, wins
		FROM users
		WHERE id = ANY($1)
		ORDER BY id
		FOR UPDATE
	`, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[int]Card{}
	for rows.Next() {
		var card Card
		if err := rows.Scan(&card.ID, &card.Elo, &card.Matches, &card.Wins); err != nil {
			return nil, err
		}
		out[card.ID] = card
	}
	return out, rows.Err()
}

func cardsFor(movies map[int]Card, ids []int) []Card {
	out := make([]Card, 0, len(ids))
	seen := map[int]bool{}
	for _, id := range ids {
		if seen[id] {
			continue
		}
		seen[id] = true
		if card, ok := movies[id]; ok {
			out = append(out, card)
		}
	}
	return out
}

func uniqueInts(ids []int) []int {
	seen := map[int]bool{}
	out := make([]int, 0, len(ids))
	for _, id := range ids {
		if seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	return out
}

func writeMovie(tx *sql.Tx, card Card) error {
	_, err := tx.Exec(
		`UPDATE movies SET elo = $1, match_count = $2, wins = $3 WHERE id = $4`,
		card.Elo, card.Matches, card.Wins, card.ID,
	)
	return err
}

func writeSet(tx *sql.Tx, card Card) error {
	_, err := tx.Exec(
		`UPDATE top4_sets SET elo = $1, match_count = $2, wins = $3 WHERE id = $4`,
		card.Elo, card.Matches, card.Wins, card.ID,
	)
	return err
}

func writeUser(tx *sql.Tx, card Card) error {
	_, err := tx.Exec(
		`UPDATE users SET elo = $1, match_count = $2, wins = $3 WHERE id = $4`,
		card.Elo, card.Matches, card.Wins, card.ID,
	)
	return err
}
