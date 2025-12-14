package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"moviemash/backend/internal/api"
	"moviemash/backend/internal/database"
)

func main() {
	// Load .env file (try current dir and parent dir)
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			if err := godotenv.Load("../../.env"); err != nil {
				log.Println("No .env file found, using environment variables")
			}
		}
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	log.Println("Database connected successfully")

	// Setup routes
	mux := api.SetupRoutes()

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("API available at http://0.0.0.0:%s/api/v1", port)

	// Start server - listen on all interfaces for Railway
	if err := http.ListenAndServe("0.0.0.0:"+port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

