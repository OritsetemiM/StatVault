package main

import (
	"fmt"
	"strings"
)

// PlayerStats holds a player's season averages
type PlayerStats struct {
	PlayerID  int     `json:"player_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Team      string  `json:"team"`
	Position  string  `json:"position"`
	Sport     string  `json:"sport"`
	Points    float64 `json:"points"`
	Assists   float64 `json:"assists"`
	Rebounds  float64 `json:"rebounds"`
	Steals    float64 `json:"steals"`
	Blocks    float64 `json:"blocks"`
	FGPct     float64 `json:"fg_pct"`
	ThreePct  float64 `json:"three_pct"`
	Games     int     `json:"games"`
}

// real 2023-24 NBA season averages
var nbaStats = map[string]PlayerStats{
	"lebron james": {
		PlayerID: 237, FirstName: "LeBron", LastName: "James",
		Team: "Los Angeles Lakers", Position: "F", Sport: "NBA",
		Points: 25.7, Assists: 8.3, Rebounds: 7.3,
		Steals: 1.3, Blocks: 0.5, FGPct: 0.540, ThreePct: 0.410, Games: 71,
	},
	"stephen curry": {
		PlayerID: 115, FirstName: "Stephen", LastName: "Curry",
		Team: "Golden State Warriors", Position: "G", Sport: "NBA",
		Points: 26.4, Assists: 5.1, Rebounds: 4.5,
		Steals: 0.9, Blocks: 0.4, FGPct: 0.450, ThreePct: 0.408, Games: 74,
	},
	"giannis antetokounmpo": {
		PlayerID: 15, FirstName: "Giannis", LastName: "Antetokounmpo",
		Team: "Milwaukee Bucks", Position: "F", Sport: "NBA",
		Points: 30.4, Assists: 6.5, Rebounds: 11.5,
		Steals: 1.2, Blocks: 1.1, FGPct: 0.611, ThreePct: 0.274, Games: 73,
	},
	"nikola jokic": {
		PlayerID: 222, FirstName: "Nikola", LastName: "Jokic",
		Team: "Denver Nuggets", Position: "C", Sport: "NBA",
		Points: 26.4, Assists: 9.0, Rebounds: 12.4,
		Steals: 1.4, Blocks: 0.9, FGPct: 0.583, ThreePct: 0.359, Games: 79,
	},
	"kevin durant": {
		PlayerID: 140, FirstName: "Kevin", LastName: "Durant",
		Team: "Phoenix Suns", Position: "F", Sport: "NBA",
		Points: 27.1, Assists: 5.0, Rebounds: 6.6,
		Steals: 0.9, Blocks: 1.2, FGPct: 0.524, ThreePct: 0.413, Games: 75,
	},
	"jayson tatum": {
		PlayerID: 434, FirstName: "Jayson", LastName: "Tatum",
		Team: "Boston Celtics", Position: "F", Sport: "NBA",
		Points: 26.9, Assists: 4.9, Rebounds: 8.1,
		Steals: 1.0, Blocks: 0.6, FGPct: 0.471, ThreePct: 0.371, Games: 74,
	},
	"luka doncic": {
		PlayerID: 434, FirstName: "Luka", LastName: "Doncic",
		Team: "Dallas Mavericks", Position: "G", Sport: "NBA",
		Points: 33.9, Assists: 9.8, Rebounds: 9.2,
		Steals: 1.4, Blocks: 0.5, FGPct: 0.487, ThreePct: 0.382, Games: 70,
	},
	"joel embiid": {
		PlayerID: 145, FirstName: "Joel", LastName: "Embiid",
		Team: "Philadelphia 76ers", Position: "C", Sport: "NBA",
		Points: 34.7, Assists: 5.6, Rebounds: 11.0,
		Steals: 1.2, Blocks: 1.7, FGPct: 0.528, ThreePct: 0.378, Games: 39,
	},
}

// fetchPlayerStats looks up a player by name from our stats database
func fetchPlayerStats(apiKey string, playerID int) (*PlayerStats, error) {
	for _, stats := range nbaStats {
		if stats.PlayerID == playerID {
			s := stats
			return &s, nil
		}
	}
	return nil, fmt.Errorf("player not found")
}

// searchPlayerStats finds a player by name search
func searchPlayerStats(query string) []PlayerStats {
	var results []PlayerStats
	query = strings.ToLower(query)
	for key, stats := range nbaStats {
		if strings.Contains(key, query) {
			results = append(results, stats)
		}
	}
	return results
}