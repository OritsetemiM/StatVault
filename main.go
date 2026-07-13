package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Player represents an NBA player from the API
type Player struct {
	ID       int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Position  string `json:"position"`
	Team      struct {
		FullName string `json:"full_name"`
		City     string `json:"city"`
	} `json:"team"`
}

// APIResponse is the wrapper the balldontlie API sends back
type APIResponse struct {
	Data []Player `json:"data"`
}

// fetchPlayers pulls players from the NBA API
func fetchPlayers(apiKey string) ([]Player, error) {
	url := "https://api.balldontlie.io/v1/players?per_page=25"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// add the API key to the request header
	req.Header.Add("Authorization", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result APIResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func main() {
	// load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("BALLDONTLIE_API_KEY")
	if apiKey == "" {
		log.Fatal("API key not found in .env file")
	}

	fmt.Println("Connecting to NBA API...")
	players, err := fetchPlayers(apiKey)
	if err != nil {
		log.Fatal("Error fetching players:", err)
	}

	fmt.Printf("Got %d players:\n\n", len(players))
	for _, p := range players {
		fmt.Printf("%s %s — %s (%s)\n",
			p.FirstName,
			p.LastName,
			p.Position,
			p.Team.FullName,
		)
	}
}