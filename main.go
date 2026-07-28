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

		// import NFL data
	fmt.Println("Importing NFL 2025 passing...")
	err = importNFLPassing("data/nfl_2025_passing.csv", "2025")
	if err != nil {
		log.Printf("Warning: could not import NFL passing: %v\n", err)
	}

	fmt.Println("Importing NFL 2025 rushing...")
	err = importNFLRushing("data/nfl_2025_rushing.csv", "2025")
	if err != nil {
		log.Printf("Warning: could not import NFL rushing: %v\n", err)
	}

	fmt.Println("Importing NFL 2025 receiving...")
	err = importNFLReceiving("data/nfl_2025_receiving.csv", "2025")
	if err != nil {
		log.Printf("Warning: could not import NFL receiving: %v\n", err)
	}

	fmt.Println("All seasons imported. Starting server...")

	// start the web server
	startServer()


}