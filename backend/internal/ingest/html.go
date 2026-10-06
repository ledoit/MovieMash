package ingest

import (
	"html"
	"regexp"
	"strconv"
	"strings"
)

var itemPattern = regexp.MustCompile(`data-item-name="([^"]+)"[^>]*data-item-slug="([^"]+)"`)
var yearPattern = regexp.MustCompile(`^(.*) \((\d{4})\)$`)

func ParseItems(page string) []Film {
	matches := itemPattern.FindAllStringSubmatch(page, -1)
	out := make([]Film, 0, len(matches))
	seen := map[string]bool{}
	for _, match := range matches {
		slug := match[2]
		if seen[slug] {
			continue
		}
		seen[slug] = true
		name := html.UnescapeString(match[1])
		film := Film{LetterboxdID: slug, Title: name}
		if parts := yearPattern.FindStringSubmatch(name); parts != nil {
			film.Title = strings.TrimSpace(parts[1])
			film.Year, _ = strconv.Atoi(parts[2])
		}
		if film.Year == 0 {
			continue
		}
		out = append(out, film)
	}
	return out
}

func ParseFavorites(page string) []Film {
	start := strings.Index(page, "Favorite films")
	if start < 0 {
		return nil
	}
	section := page[start:]
	if next := strings.Index(section[len("Favorite films"):], "<h2"); next >= 0 {
		section = section[:len("Favorite films")+next]
	}
	films := ParseItems(section)
	if len(films) > 4 {
		films = films[:4]
	}
	return films
}
