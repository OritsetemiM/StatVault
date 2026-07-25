package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// global API key so all handlers can use it
var apiKey string

// respondJSON sends a JSON response back to the browser
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error message as JSON
func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func handlePlayers(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		respondError(w, 400, "search query required")
		return
	}

	players, err := searchDB(search, "NBA")
	if err != nil {
		respondError(w, 500, "Search failed")
		return
	}

	respondJSON(w, players)
}

// handleHealth is a simple health check endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{"status": "ok"})
}

func startServer() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey = os.Getenv("BALLDONTLIE_API_KEY")
	if apiKey == "" {
		log.Fatal("API key not found in .env file")
	}

	// API routes MUST be registered before the static file server
	http.HandleFunc("/api/players", handlePlayers)
	http.HandleFunc("/api/health", handleHealth)

	// static frontend files — registered LAST
	http.Handle("/", http.FileServer(http.Dir("./static")))

	port := "8080"
	fmt.Printf("StatVault server running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}