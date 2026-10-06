package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"moviemash/backend/internal/database"
	"moviemash/backend/internal/ingest"
)

func main() {
	jsonPath := flag.String("json", "", "path to a Letterboxd export")
	fetchUser := flag.String("fetch", "", "try one public profile page and stop if Letterboxd blocks it")
	flag.Parse()

	if *fetchUser != "" {
		if err := tryPublicPage(*fetchUser); err != nil {
			log.Fatal(err)
		}
		return
	}
	if *jsonPath == "" {
		log.Fatal("pass -json or -fetch")
	}

	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()

	body, err := os.ReadFile(*jsonPath)
	if err != nil {
		log.Fatal(err)
	}
	var doc ingest.Document
	if err := json.Unmarshal(body, &doc); err != nil {
		log.Fatal(err)
	}
	n, err := ingest.Import(database.DB, doc)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("imported %d profiles", n)
}

func tryPublicPage(username string) error {
	req, err := http.NewRequest(http.MethodGet, "https://letterboxd.com/"+username+"/", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "MovieMash/0.1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil {
		return err
	}
	text := string(body)
	if resp.StatusCode == http.StatusForbidden || strings.Contains(text, "cf-mitigated") || strings.Contains(text, "challenge-platform") {
		return fmt.Errorf("letterboxd returned a challenge for %s; save the public profile and import it with -json. this command does not bypass the challenge", username)
	}
	films := ingest.ParseFavorites(text)
	fmt.Printf("%s favorites: %d\n", username, len(films))
	for _, film := range films {
		fmt.Printf("- %s (%d) %s\n", film.Title, film.Year, film.LetterboxdID)
	}
	return nil
}
