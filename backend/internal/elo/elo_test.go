package elo

import "math"
import "testing"

func close(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 0.0001 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestEvenBout(t *testing.T) {
	a := Card{Elo: Start}
	b := Card{Elo: Start}
	a, b = Bout(a, b, 1)
	close(t, a.Elo, 1516)
	close(t, b.Elo, 1484)
	if a.Wins != 1 || a.Matches != 1 || b.Wins != 0 || b.Matches != 1 {
		t.Fatalf("counts a=%+v b=%+v", a, b)
	}
}

func TestFavoriteGainsLess(t *testing.T) {
	fav := Card{Elo: 1700}
	dog := Card{Elo: Start}
	fav, _ = Bout(fav, dog, 1)
	if fav.Elo-1700 >= 16 {
		t.Fatalf("favorite gained %v", fav.Elo-1700)
	}
}

func TestMoviesUsePrematchRatings(t *testing.T) {
	winners := []Card{{ID: 1, Elo: Start}, {ID: 2, Elo: Start}}
	losers := []Card{{ID: 3, Elo: Start}, {ID: 4, Elo: Start}}
	winners, losers = Movies(winners, losers)
	close(t, winners[0].Elo, 1516)
	close(t, winners[1].Elo, 1516)
	close(t, losers[0].Elo, 1484)
	close(t, losers[1].Elo, 1484)
	if winners[0].Matches != 1 || losers[0].Matches != 1 {
		t.Fatalf("each movie should record one match")
	}
}
