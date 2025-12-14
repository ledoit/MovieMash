package api

import (
	"net/http"
	"os"
)

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins in production, or specific origin from env
		origin := r.Header.Get("Origin")
		allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
		
		if allowedOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		} else {
			// Allow any origin if not specified (for development)
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// SetupRoutes configures all API routes
func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// API v1 routes with CORS
	mux.HandleFunc("/api/v1/comparison", corsMiddleware(GetComparison))
	mux.HandleFunc("/api/v1/votes", corsMiddleware(CreateVote))
	mux.HandleFunc("/api/v1/leaderboard/top4", corsMiddleware(GetTop4Leaderboard))
	mux.HandleFunc("/api/v1/leaderboard/movies", corsMiddleware(GetMoviesLeaderboard))

	return mux
}

