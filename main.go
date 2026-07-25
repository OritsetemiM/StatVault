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
	ID        int    `json:"id"`
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

// fetchPlayers pulls players, optionally filtering by name
func fetchPlayers(apiKey string, search string) ([]Player, error) {
	url := "https://api.balldontlie.io/v1/players?per_page=25"
	if search != "" {
		url += "&search=" + search
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey = os.Getenv("BALLDONTLIE_API_KEY")

	// test stats lookup
	results := searchPlayerStats("lebron")
	if len(results) == 0 {
		fmt.Println("No players found")
		return
	}

	p := results[0]
	fmt.Printf("%s %s (%s)\n", p.FirstName, p.LastName, p.Team)
	fmt.Printf("Games: %d\n", p.Games)
	fmt.Printf("PPG: %.1f | APG: %.1f | RPG: %.1f\n",
		p.Points, p.Assists, p.Rebounds)
	fmt.Printf("FG%%: %.1f%% | 3P%%: %.1f%%\n",
		p.FGPct*100, p.ThreePct*100)
}