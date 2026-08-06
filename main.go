package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var apiKey string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey = os.Getenv("BALLDONTLIE_API_KEY")
	_ = apiKey

	// initialize database
	initDB()

	// import all seasons
	seasons := []struct {
		file   string
		season string
	}{
		{"data/nba2020.csv", "2019-20"},
		{"data/nba2021.csv", "2020-21"},
		{"data/nba2022.csv", "2021-22"},
		{"data/nba2023.csv", "2022-23"},
		{"data/nba2024.csv", "2023-24"},
		{"data/nba2025.csv", "2024-25"},
		{"data/nba2026.csv", "2025-26"},
	}

	for _, s := range seasons {
		fmt.Printf("Importing %s...\n", s.season)
		err = importNBACSV(s.file, s.season)
		if err != nil {
			log.Printf("Warning: could not import %s: %v\n", s.file, err)
		}
	}

		// NFL seasons
nflSeasons := []string{"2020", "2021", "2022", "2023", "2024", "2025"}

	for _, s := range nflSeasons {
    	fmt.Printf("Importing NFL %s passing...\n", s)
    	err = importNFLPassing(fmt.Sprintf("data/nfl_%s_passing.csv", s), s)
    	if err != nil {
        	log.Printf("Warning: NFL passing %s: %v\n", s, err)
		}

    	fmt.Printf("Importing NFL %s rushing...\n", s)
    	err = importNFLRushing(fmt.Sprintf("data/nfl_%s_rushing.csv", s), s)
    	if err != nil {
        	log.Printf("Warning: NFL rushing %s: %v\n", s, err)
    	}

    	fmt.Printf("Importing NFL %s receiving...\n", s)
    	err = importNFLReceiving(fmt.Sprintf("data/nfl_%s_receiving.csv", s), s)
    	if err != nil {
        	log.Printf("Warning: NFL receiving %s: %v\n", s, err)
    	}

    	fmt.Printf("Importing NFL %s defense...\n", s)
    	err = importNFLDefense(fmt.Sprintf("data/nfl_%s_defense.csv", s), s)
    	if err != nil {
        	log.Printf("Warning: NFL defense %s: %v\n", s, err)
    	}
	}
	// start the web server
	startServer()
}