package ingest

type Film struct {
	LetterboxdID string `json:"letterboxd_id"`
	Title        string `json:"title"`
	Year         int    `json:"year"`
	Poster       string `json:"poster,omitempty"`
}

type Profile struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Top4        []Film `json:"top4"`
	Watches     []Film `json:"watches"`
}

type Document struct {
	Profiles []Profile `json:"profiles"`
}
