package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey = os.Getenv("BALLDONTLIE_API_KEY")
	_ = apiKey

	// initialize database
	initDB()
	defer db.Close()

	// import the CSV
	fmt.Println("Importing NBA 2025-26 season stats...")
	err = importNBACSV("nba2026.csv", "2025-26")
	if err != nil {
		log.Fatal("Import failed:", err)
	}

	fmt.Println("Done! Testing a search...")

	// test search
	players, err := searchDB("LeBron", "NBA")
	if err != nil {
		log.Fatal("Search failed:", err)
	}

	for _, p := range players {
		fmt.Printf("%s — %s | %.1f PPG | %.1f APG | %.1f RPG\n",
			p.Name, p.Team, p.Points, p.Assists, p.Rebounds)
	}
}