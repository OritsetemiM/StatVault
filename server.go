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

// NBA handlers
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

// NFL handlers
func handleNFLSearch(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		respondError(w, 400, "search query required")
		return
	}
	results, err := searchNFLAll(search)
	if err != nil {
		respondError(w, 500, "Search failed")
		return
	}
	respondJSON(w, results)
}

func handleNFLProfile(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		respondError(w, 400, "name required")
		return
	}
	profile, err := getNFLProfile(name)
	if err != nil {
		respondError(w, 404, "Player not found")
		return
	}
	respondJSON(w, profile)
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

func handleNFLDefense(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		respondError(w, 400, "search query required")
		return
	}
	results, err := searchNFLDefense(search)
	if err != nil {
		respondError(w, 500, "Search failed")
		return
	}
	respondJSON(w, results)
}

func handleCompare(w http.ResponseWriter, r *http.Request) {
	sport := r.URL.Query().Get("sport")
	statType := r.URL.Query().Get("stat_type")
	name1 := r.URL.Query().Get("player1")
	name2 := r.URL.Query().Get("player2")
	season1 := r.URL.Query().Get("season1")
	season2 := r.URL.Query().Get("season2")

	if name1 == "" || name2 == "" {
		respondError(w, 400, "Both player names required")
		return
	}

	var result *CompareResult
	var err error

	if sport == "NBA" {
		result, err = CompareNBAPlayers(name1, season1, name2, season2)
	} else if sport == "NFL" {
		switch statType {
		case "passing":
			result, err = CompareNFLPassers(name1, season1, name2, season2)
		case "rushing":
			result, err = CompareNFLRushers(name1, season1, name2, season2)
		case "receiving":
			result, err = CompareNFLReceivers(name1, season1, name2, season2)
		case "defense":
			result, err = CompareNFLDefenders(name1, season1, name2, season2)
		default:
			respondError(w, 400, "stat_type required for NFL (passing/rushing/receiving/defense)")
			return
		}
	} else {
		respondError(w, 400, "sport must be NBA or NFL")
		return
	}

	if err != nil {
		respondError(w, 404, err.Error())
		return
	}

	respondJSON(w, result)
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

	// API routes
	http.HandleFunc("/api/players", handlePlayers)
	http.HandleFunc("/api/health", handleHealth)
	http.HandleFunc("/api/career", handleCareer)
	http.HandleFunc("/api/nfl/search", handleNFLSearch)
	http.HandleFunc("/api/nfl/profile", handleNFLProfile)
	http.HandleFunc("/api/nfl/passing", handleNFLPassing)
	http.HandleFunc("/api/nfl/rushing", handleNFLRushing)
	http.HandleFunc("/api/nfl/receiving", handleNFLReceiving)
	http.HandleFunc("/api/nfl/defense", handleNFLDefense)
	http.HandleFunc("/api/compare", handleCompare)

	// static frontend
	http.Handle("/", http.FileServer(http.Dir("./static")))

	port := "8080"
	fmt.Printf("StatVault server running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}