package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(data)
}

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

func handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{"status": "ok"})
}

func handleCareer(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	sport := r.URL.Query().Get("sport")
	if name == "" {
		respondError(w, 400, "name query required")
		return
	}
	if sport == "" {
		sport = "NBA"
	}

	career, err := getCareerAverages(name, sport)
	if err != nil {
		respondError(w, 404, "Player not found")
		return
	}

	respondJSON(w, career)
}

func handleNFLPassing(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		respondError(w, 400, "search query required")
		return
	}
	results, err := searchNFLPassing(search)
	if err != nil {
		respondError(w, 500, "Search failed")
		return
	}
	respondJSON(w, results)
}

func handleNFLRushing(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		respondError(w, 400, "search query required")
		return
	}
	results, err := searchNFLRushing(search)
	if err != nil {
		respondError(w, 500, "Search failed")
		return
	}
	respondJSON(w, results)
}

func handleNFLReceiving(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		respondError(w, 400, "search query required")
		return
	}
	results, err := searchNFLReceiving(search)
	if err != nil {
		respondError(w, 500, "Search failed")
		return
	}
	respondJSON(w, results)
}

func startServer() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey = os.Getenv("BALLDONTLIE_API_KEY")
	if apiKey == "" {
		log.Print("Warning: no API key found")
	}

	// API routes first
	http.HandleFunc("/api/players", handlePlayers)
	http.HandleFunc("/api/health", handleHealth)
	http.HandleFunc("/api/career", handleCareer)

	// static frontend
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/api/nfl/passing", handleNFLPassing)
	http.HandleFunc("/api/nfl/rushing", handleNFLRushing)
	http.HandleFunc("/api/nfl/receiving", handleNFLReceiving)

	port := "8080"
	fmt.Printf("StatVault server running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}