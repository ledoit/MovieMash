package ingest

import "testing"

const sample = `
<h2>Favorite films</h2>
<div data-item-name="Alien (1979)" data-item-slug="alien"></div>
<div data-item-name="Lee Cronin&#039;s The Mummy (2026)" data-item-slug="lee-cronins-the-mummy"></div>
<div data-item-name="Sunshine (2007)" data-item-slug="sunshine-2007"></div>
<div data-item-name="Interstellar (2014)" data-item-slug="interstellar"></div>
<div data-item-name="Extra (2001)" data-item-slug="extra"></div>
<h2>Recent activity</h2>
<div data-item-name="Not a favorite (1999)" data-item-slug="not-a-favorite"></div>
`

func TestParseFavorites(t *testing.T) {
	films := ParseFavorites(sample)
	if len(films) != 4 {
		t.Fatalf("got %d films", len(films))
	}
	if films[0].LetterboxdID != "alien" || films[0].Year != 1979 {
		t.Fatalf("first film %+v", films[0])
	}
	if films[1].Title != "Lee Cronin's The Mummy" {
		t.Fatalf("title %q", films[1].Title)
	}
	if films[3].LetterboxdID != "interstellar" {
		t.Fatalf("fourth %+v", films[3])
	}
}
