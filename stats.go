package main

import (
	"fmt"
	"strings"
)

// SeasonStats holds one season of stats for a player
type SeasonStats struct {
	Season   string  `json:"season"`
	Age      int     `json:"age"`
	Team     string  `json:"team"`
	Games    int     `json:"games"`
	Points   float64 `json:"points"`
	Assists  float64 `json:"assists"`
	Rebounds float64 `json:"rebounds"`
	Steals   float64 `json:"steals"`
	Blocks   float64 `json:"blocks"`
	FGPct    float64 `json:"fg_pct"`
	ThreePct float64 `json:"three_pct"`
}

// PlayerProfile holds a player's full career
type PlayerProfile struct {
	ID        int           `json:"id"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	Position  string        `json:"position"`
	Sport     string        `json:"sport"`
	Seasons   []SeasonStats `json:"seasons"`
}

// PlayerStats holds a single season snapshot (used for comparisons)
type PlayerStats struct {
	PlayerID  int     `json:"player_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Team      string  `json:"team"`
	Position  string  `json:"position"`
	Sport     string  `json:"sport"`
	Season    string  `json:"season"`
	Points    float64 `json:"points"`
	Assists   float64 `json:"assists"`
	Rebounds  float64 `json:"rebounds"`
	Steals    float64 `json:"steals"`
	Blocks    float64 `json:"blocks"`
	FGPct     float64 `json:"fg_pct"`
	ThreePct  float64 `json:"three_pct"`
	Games     int     `json:"games"`
}

// CareerAverages calculates career averages across all seasons
func CareerAverages(p PlayerProfile) PlayerStats {
	var totalPts, totalAst, totalReb, totalStl, totalBlk, totalFg, totalFg3 float64
	totalGames := 0

	for _, s := range p.Seasons {
		totalPts += s.Points * float64(s.Games)
		totalAst += s.Assists * float64(s.Games)
		totalReb += s.Rebounds * float64(s.Games)
		totalStl += s.Steals * float64(s.Games)
		totalBlk += s.Blocks * float64(s.Games)
		totalFg += s.FGPct * float64(s.Games)
		totalFg3 += s.ThreePct * float64(s.Games)
		totalGames += s.Games
	}

	g := float64(totalGames)
	return PlayerStats{
		PlayerID:  p.ID,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Position:  p.Position,
		Sport:     p.Sport,
		Season:    "Career",
		Games:     totalGames,
		Points:    totalPts / g,
		Assists:   totalAst / g,
		Rebounds:  totalReb / g,
		Steals:    totalStl / g,
		Blocks:    totalBlk / g,
		FGPct:     totalFg / g,
		ThreePct:  totalFg3 / g,
	}
}

// LatestSeason returns the most recent season as a PlayerStats
func LatestSeason(p PlayerProfile) PlayerStats {
	s := p.Seasons[len(p.Seasons)-1]
	return PlayerStats{
		PlayerID:  p.ID,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Team:      s.Team,
		Position:  p.Position,
		Sport:     p.Sport,
		Season:    s.Season,
		Points:    s.Points,
		Assists:   s.Assists,
		Rebounds:  s.Rebounds,
		Steals:    s.Steals,
		Blocks:    s.Blocks,
		FGPct:     s.FGPct,
		ThreePct:  s.ThreePct,
		Games:     s.Games,
	}
}

// player database
var playerDB = map[string]PlayerProfile{
	"lebron james": {
		ID: 237, FirstName: "LeBron", LastName: "James",
		Position: "F", Sport: "NBA",
		Seasons: []SeasonStats{
			{"2003-04", 18, "Cleveland Cavaliers", 79, 20.9, 5.9, 5.5, 1.6, 0.7, 0.417, 0.290},
			{"2004-05", 19, "Cleveland Cavaliers", 80, 27.2, 7.2, 7.4, 2.2, 0.7, 0.472, 0.351},
			{"2005-06", 20, "Cleveland Cavaliers", 79, 31.4, 6.6, 7.0, 1.6, 0.8, 0.480, 0.335},
			{"2006-07", 21, "Cleveland Cavaliers", 78, 27.3, 6.0, 6.7, 1.6, 0.7, 0.476, 0.319},
			{"2007-08", 22, "Cleveland Cavaliers", 75, 30.0, 7.2, 7.9, 1.8, 1.1, 0.484, 0.315},
			{"2008-09", 23, "Cleveland Cavaliers", 81, 28.4, 7.2, 7.6, 1.7, 1.1, 0.489, 0.344},
			{"2009-10", 24, "Cleveland Cavaliers", 76, 29.7, 8.6, 7.3, 1.6, 1.0, 0.503, 0.333},
			{"2010-11", 25, "Miami Heat", 79, 26.7, 7.0, 7.5, 1.6, 0.6, 0.510, 0.330},
			{"2011-12", 26, "Miami Heat", 62, 27.1, 6.2, 7.9, 1.9, 0.8, 0.531, 0.362},
			{"2012-13", 27, "Miami Heat", 76, 26.8, 7.3, 8.0, 1.7, 0.9, 0.565, 0.406},
			{"2013-14", 28, "Miami Heat", 77, 27.1, 6.3, 6.9, 1.6, 0.3, 0.567, 0.379},
			{"2014-15", 29, "Cleveland Cavaliers", 69, 25.3, 7.4, 6.0, 1.6, 0.7, 0.488, 0.354},
			{"2015-16", 30, "Cleveland Cavaliers", 76, 25.3, 6.8, 7.4, 1.4, 0.6, 0.520, 0.309},
			{"2016-17", 31, "Cleveland Cavaliers", 74, 26.4, 8.7, 8.6, 1.2, 0.6, 0.548, 0.363},
			{"2017-18", 32, "Cleveland Cavaliers", 82, 27.5, 9.1, 8.6, 1.4, 0.9, 0.542, 0.367},
			{"2018-19", 33, "Los Angeles Lakers", 55, 27.4, 8.3, 8.5, 1.3, 0.6, 0.510, 0.339},
			{"2019-20", 34, "Los Angeles Lakers", 67, 25.3, 10.2, 7.8, 1.2, 0.5, 0.493, 0.348},
			{"2020-21", 35, "Los Angeles Lakers", 45, 25.0, 7.8, 7.7, 1.1, 0.6, 0.513, 0.365},
			{"2021-22", 36, "Los Angeles Lakers", 56, 30.3, 6.2, 8.2, 1.3, 1.1, 0.524, 0.354},
			{"2022-23", 37, "Los Angeles Lakers", 55, 28.9, 6.8, 8.3, 0.9, 0.6, 0.500, 0.326},
			{"2023-24", 38, "Los Angeles Lakers", 71, 25.7, 8.3, 7.3, 1.3, 0.5, 0.540, 0.410},
		},
	},
	"stephen curry": {
		ID: 115, FirstName: "Stephen", LastName: "Curry",
		Position: "G", Sport: "NBA",
		Seasons: []SeasonStats{
			{"2009-10", 21, "Golden State Warriors", 80, 17.5, 5.9, 4.5, 1.9, 0.2, 0.462, 0.437},
			{"2010-11", 22, "Golden State Warriors", 74, 18.6, 5.8, 3.9, 1.5, 0.3, 0.480, 0.442},
			{"2011-12", 23, "Golden State Warriors", 26, 14.7, 5.3, 3.4, 1.5, 0.3, 0.490, 0.455},
			{"2012-13", 24, "Golden State Warriors", 78, 22.9, 6.9, 4.0, 1.6, 0.2, 0.451, 0.453},
			{"2013-14", 25, "Golden State Warriors", 78, 24.0, 8.5, 4.3, 1.6, 0.2, 0.471, 0.421},
			{"2014-15", 26, "Golden State Warriors", 80, 23.8, 7.7, 4.3, 2.0, 0.2, 0.487, 0.443},
			{"2015-16", 27, "Golden State Warriors", 79, 30.1, 6.7, 5.4, 2.1, 0.2, 0.504, 0.454},
			{"2016-17", 28, "Golden State Warriors", 79, 25.3, 6.6, 4.5, 1.8, 0.2, 0.468, 0.411},
			{"2017-18", 29, "Golden State Warriors", 51, 26.4, 6.1, 5.1, 1.6, 0.2, 0.490, 0.423},
			{"2018-19", 30, "Golden State Warriors", 69, 27.3, 5.2, 5.3, 1.3, 0.4, 0.474, 0.437},
			{"2019-20", 31, "Golden State Warriors", 5, 20.8, 6.6, 5.0, 1.0, 0.4, 0.409, 0.245},
			{"2020-21", 32, "Golden State Warriors", 63, 32.0, 5.8, 5.5, 1.2, 0.4, 0.482, 0.421},
			{"2021-22", 33, "Golden State Warriors", 64, 25.5, 6.3, 5.2, 1.3, 0.4, 0.437, 0.380},
			{"2022-23", 34, "Golden State Warriors", 56, 29.4, 6.3, 6.1, 0.9, 0.4, 0.493, 0.427},
			{"2023-24", 35, "Golden State Warriors", 74, 26.4, 5.1, 4.5, 0.9, 0.4, 0.450, 0.408},
		},
	},
	"giannis antetokounmpo": {
		ID: 15, FirstName: "Giannis", LastName: "Antetokounmpo",
		Position: "F", Sport: "NBA",
		Seasons: []SeasonStats{
			{"2013-14", 18, "Milwaukee Bucks", 77, 6.8, 2.2, 4.4, 0.8, 0.5, 0.416, 0.347},
			{"2014-15", 19, "Milwaukee Bucks", 81, 12.7, 2.6, 6.7, 0.9, 1.0, 0.491, 0.157},
			{"2015-16", 20, "Milwaukee Bucks", 80, 16.9, 4.3, 7.7, 1.2, 1.4, 0.508, 0.171},
			{"2016-17", 21, "Milwaukee Bucks", 80, 22.9, 5.4, 8.8, 1.6, 1.9, 0.520, 0.272},
			{"2017-18", 22, "Milwaukee Bucks", 75, 26.9, 4.8, 10.0, 1.5, 1.4, 0.529, 0.307},
			{"2018-19", 23, "Milwaukee Bucks", 72, 27.7, 5.9, 12.5, 1.3, 1.5, 0.578, 0.258},
			{"2019-20", 24, "Milwaukee Bucks", 63, 29.6, 5.8, 13.7, 1.0, 1.0, 0.553, 0.304},
			{"2020-21", 25, "Milwaukee Bucks", 61, 28.1, 5.9, 11.0, 1.2, 1.2, 0.569, 0.303},
			{"2021-22", 26, "Milwaukee Bucks", 67, 29.9, 5.8, 11.6, 1.1, 1.4, 0.553, 0.290},
			{"2022-23", 27, "Milwaukee Bucks", 63, 31.1, 5.7, 11.8, 0.8, 0.8, 0.553, 0.275},
			{"2023-24", 28, "Milwaukee Bucks", 73, 30.4, 6.5, 11.5, 1.2, 1.1, 0.611, 0.274},
		},
	},
	"nikola jokic": {
		ID: 222, FirstName: "Nikola", LastName: "Jokic",
		Position: "C", Sport: "NBA",
		Seasons: []SeasonStats{
			{"2015-16", 20, "Denver Nuggets", 80, 10.0, 5.0, 7.3, 1.2, 0.6, 0.514, 0.333},
			{"2016-17", 21, "Denver Nuggets", 73, 16.7, 6.8, 9.8, 1.2, 0.8, 0.574, 0.324},
			{"2017-18", 22, "Denver Nuggets", 75, 18.5, 6.1, 10.7, 1.4, 0.8, 0.497, 0.396},
			{"2018-19", 23, "Denver Nuggets", 80, 20.2, 7.3, 10.8, 1.4, 0.7, 0.511, 0.307},
			{"2019-20", 24, "Denver Nuggets", 73, 20.2, 6.9, 10.2, 1.2, 0.6, 0.528, 0.314},
			{"2020-21", 25, "Denver Nuggets", 72, 26.4, 8.3, 10.8, 1.3, 0.7, 0.566, 0.388},
			{"2021-22", 26, "Denver Nuggets", 74, 27.1, 7.9, 13.8, 1.5, 0.9, 0.583, 0.336},
			{"2022-23", 27, "Denver Nuggets", 69, 24.5, 9.8, 11.8, 1.3, 0.7, 0.632, 0.383},
			{"2023-24", 28, "Denver Nuggets", 79, 26.4, 9.0, 12.4, 1.4, 0.9, 0.583, 0.359},
		},
	},
	"luka doncic": {
		ID: 434, FirstName: "Luka", LastName: "Doncic",
		Position: "G", Sport: "NBA",
		Seasons: []SeasonStats{
			{"2018-19", 19, "Dallas Mavericks", 72, 21.2, 6.0, 7.8, 1.1, 0.4, 0.424, 0.327},
			{"2019-20", 20, "Dallas Mavericks", 61, 28.8, 8.8, 9.4, 1.0, 0.2, 0.463, 0.319},
			{"2020-21", 21, "Dallas Mavericks", 66, 27.7, 8.6, 8.0, 1.0, 0.5, 0.478, 0.350},
			{"2021-22", 22, "Dallas Mavericks", 65, 28.4, 8.7, 9.1, 1.2, 0.5, 0.457, 0.350},
			{"2022-23", 23, "Dallas Mavericks", 66, 32.4, 8.0, 8.6, 1.4, 0.5, 0.496, 0.356},
			{"2023-24", 24, "Dallas Mavericks", 70, 33.9, 9.8, 9.2, 1.4, 0.5, 0.487, 0.382},
		},
	},
	"kevin durant": {
		ID: 140, FirstName: "Kevin", LastName: "Durant",
		Position: "F", Sport: "NBA",
		Seasons: []SeasonStats{
			{"2007-08", 19, "Seattle SuperSonics", 80, 20.0, 2.4, 4.4, 1.0, 0.9, 0.435, 0.287},
			{"2008-09", 20, "Oklahoma City Thunder", 74, 25.3, 2.8, 6.5, 1.3, 0.7, 0.476, 0.422},
			{"2009-10", 21, "Oklahoma City Thunder", 82, 30.1, 2.8, 7.6, 1.4, 1.0, 0.476, 0.365},
			{"2010-11", 22, "Oklahoma City Thunder", 78, 27.7, 2.7, 6.8, 1.1, 1.0, 0.496, 0.350},
			{"2011-12", 23, "Oklahoma City Thunder", 66, 28.0, 3.5, 8.0, 1.3, 1.2, 0.496, 0.387},
			{"2012-13", 24, "Oklahoma City Thunder", 81, 28.1, 4.6, 7.9, 1.4, 1.3, 0.510, 0.416},
			{"2013-14", 25, "Oklahoma City Thunder", 81, 32.0, 5.5, 7.4, 1.3, 0.7, 0.503, 0.391},
			{"2015-16", 27, "Oklahoma City Thunder", 72, 28.2, 5.0, 8.2, 1.0, 1.2, 0.505, 0.387},
			{"2016-17", 28, "Golden State Warriors", 62, 25.1, 4.8, 8.3, 1.1, 1.6, 0.537, 0.375},
			{"2017-18", 29, "Golden State Warriors", 68, 26.4, 5.4, 6.8, 0.7, 1.8, 0.516, 0.419},
			{"2018-19", 30, "Golden State Warriors", 78, 26.0, 5.9, 6.4, 0.7, 1.1, 0.521, 0.352},
			{"2020-21", 32, "Brooklyn Nets", 35, 26.9, 5.6, 7.1, 0.7, 1.3, 0.531, 0.450},
			{"2021-22", 33, "Brooklyn Nets", 55, 29.9, 6.4, 7.4, 0.9, 0.9, 0.513, 0.383},
			{"2022-23", 34, "Phoenix Suns", 39, 29.1, 5.0, 6.7, 0.8, 1.4, 0.560, 0.383},
			{"2023-24", 35, "Phoenix Suns", 75, 27.1, 5.0, 6.6, 0.9, 1.2, 0.524, 0.413},
		},
	},
	"patrick mahomes": {
		ID: 500, FirstName: "Patrick", LastName: "Mahomes",
		Position: "QB", Sport: "NFL",
		Seasons: []SeasonStats{
			{"2017", 21, "Kansas City Chiefs", 1, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			{"2018", 22, "Kansas City Chiefs", 16, 50.0, 0.0, 0.0, 0.0, 0.0, 0.660, 0.0},
			{"2019", 23, "Kansas City Chiefs", 14, 26.3, 0.0, 0.0, 0.0, 0.0, 0.655, 0.0},
			{"2020", 24, "Kansas City Chiefs", 15, 38.0, 0.0, 0.0, 0.0, 0.0, 0.664, 0.0},
			{"2021", 25, "Kansas City Chiefs", 17, 37.0, 0.0, 0.0, 0.0, 0.0, 0.663, 0.0},
			{"2022", 26, "Kansas City Chiefs", 17, 41.0, 0.0, 0.0, 0.0, 0.0, 0.674, 0.0},
			{"2023", 27, "Kansas City Chiefs", 16, 27.0, 0.0, 0.0, 0.0, 0.0, 0.671, 0.0},
		},
	},
	"tom brady": {
		ID: 501, FirstName: "Tom", LastName: "Brady",
		Position: "QB", Sport: "NFL",
		Seasons: []SeasonStats{
			{"2001", 23, "New England Patriots", 15, 18.0, 0.0, 0.0, 0.0, 0.0, 0.635, 0.0},
			{"2002", 24, "New England Patriots", 16, 28.0, 0.0, 0.0, 0.0, 0.0, 0.620, 0.0},
			{"2003", 25, "New England Patriots", 16, 23.0, 0.0, 0.0, 0.0, 0.0, 0.600, 0.0},
			{"2004", 26, "New England Patriots", 16, 28.0, 0.0, 0.0, 0.0, 0.0, 0.602, 0.0},
			{"2005", 27, "New England Patriots", 16, 26.0, 0.0, 0.0, 0.0, 0.0, 0.614, 0.0},
			{"2006", 28, "New England Patriots", 16, 24.0, 0.0, 0.0, 0.0, 0.0, 0.616, 0.0},
			{"2007", 29, "New England Patriots", 16, 50.0, 0.0, 0.0, 0.0, 0.0, 0.685, 0.0},
			{"2008", 30, "New England Patriots", 1, 0.0, 0.0, 0.0, 0.0, 0.0, 0.000, 0.0},
			{"2009", 31, "New England Patriots", 16, 28.0, 0.0, 0.0, 0.0, 0.0, 0.655, 0.0},
			{"2010", 32, "New England Patriots", 16, 36.0, 0.0, 0.0, 0.0, 0.0, 0.658, 0.0},
			{"2011", 33, "New England Patriots", 16, 39.0, 0.0, 0.0, 0.0, 0.0, 0.654, 0.0},
			{"2012", 34, "New England Patriots", 16, 34.0, 0.0, 0.0, 0.0, 0.0, 0.627, 0.0},
			{"2013", 35, "New England Patriots", 16, 25.0, 0.0, 0.0, 0.0, 0.0, 0.608, 0.0},
			{"2014", 36, "New England Patriots", 16, 33.0, 0.0, 0.0, 0.0, 0.0, 0.651, 0.0},
			{"2015", 37, "New England Patriots", 16, 36.0, 0.0, 0.0, 0.0, 0.0, 0.649, 0.0},
			{"2016", 38, "New England Patriots", 12, 28.0, 0.0, 0.0, 0.0, 0.0, 0.659, 0.0},
			{"2017", 40, "New England Patriots", 16, 32.0, 0.0, 0.0, 0.0, 0.0, 0.633, 0.0},
			{"2018", 41, "New England Patriots", 16, 29.0, 0.0, 0.0, 0.0, 0.0, 0.636, 0.0},
			{"2019", 42, "New England Patriots", 16, 24.0, 0.0, 0.0, 0.0, 0.0, 0.600, 0.0},
			{"2020", 43, "Tampa Bay Buccaneers", 16, 40.0, 0.0, 0.0, 0.0, 0.0, 0.659, 0.0},
			{"2021", 44, "Tampa Bay Buccaneers", 17, 43.0, 0.0, 0.0, 0.0, 0.0, 0.671, 0.0},
			{"2022", 45, "Tampa Bay Buccaneers", 17, 25.0, 0.0, 0.0, 0.0, 0.0, 0.660, 0.0},
		},
	},
}

// searchPlayerStats finds players by name
func searchPlayerStats(query string) []PlayerProfile {
	var results []PlayerProfile
	query = strings.ToLower(query)
	for key, profile := range playerDB {
		if strings.Contains(key, query) {
			results = append(results, profile)
		}
	}
	return results
}

// fetchPlayerStats gets the latest season stats for a player by ID
func fetchPlayerStats(apiKey string, playerID int) (*PlayerStats, error) {
	for _, profile := range playerDB {
		if profile.ID == playerID {
			s := LatestSeason(profile)
			return &s, nil
		}
	}
	return nil, fmt.Errorf("player not found")
}