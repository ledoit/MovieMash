package elo

import "math"

const (
	Start = 1500.0
	K     = 32.0
)

// Card is one rated side: a profile, a top 4, or a movie.
type Card struct {
	ID      int
	Elo     float64
	Matches int
	Wins    int
}

func Expected(a, b float64) float64 {
	return 1.0 / (1.0 + math.Pow(10, (b-a)/400.0))
}

// Bout applies one match. aScore is 1 when a wins and 0 when a loses.
func Bout(a, b Card, aScore float64) (Card, Card) {
	aElo, bElo := a.Elo, b.Elo
	a.Elo += K * (aScore - Expected(aElo, bElo))
	b.Elo += K * ((1 - aScore) - Expected(bElo, aElo))
	a.Matches++
	b.Matches++
	if aScore >= 1 {
		a.Wins++
	} else if aScore <= 0 {
		b.Wins++
	}
	return a, b
}

// Movies scores every winner against every loser using pre-match ratings.
// Each movie records one match. Its Elo change is the mean of those pairings.
func Movies(winners, losers []Card) ([]Card, []Card) {
	outW := append([]Card(nil), winners...)
	outL := append([]Card(nil), losers...)

	for i, w := range winners {
		sum := 0.0
		n := 0
		for _, l := range losers {
			if w.ID == l.ID {
				continue
			}
			sum += K * (1 - Expected(w.Elo, l.Elo))
			n++
		}
		if n == 0 {
			continue
		}
		outW[i].Elo = w.Elo + sum/float64(n)
		outW[i].Matches++
		outW[i].Wins++
	}

	for j, l := range losers {
		sum := 0.0
		n := 0
		for _, w := range winners {
			if w.ID == l.ID {
				continue
			}
			sum += K * (0 - Expected(l.Elo, w.Elo))
			n++
		}
		if n == 0 {
			continue
		}
		outL[j].Elo = l.Elo + sum/float64(n)
		outL[j].Matches++
	}

	return outW, outL
}

func WinRate(matches, wins int) *float64 {
	if matches == 0 {
		return nil
	}
	rate := float64(wins) / float64(matches)
	return &rate
}
