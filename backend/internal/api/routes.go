package api

import (
	"net/http"
)

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

