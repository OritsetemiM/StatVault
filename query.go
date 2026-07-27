package main

import (
	"fmt"
	"strings"
)

// DBPlayer represents a player row from the database
type DBPlayer struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Age       int     `json:"age"`
	Team      string  `json:"team"`
	Position  string  `json:"position"`
	Sport     string  `json:"sport"`
	Season    string  `json:"season"`
	Games     int     `json:"games"`
	Points    float64 `json:"points"`
	Assists   float64 `json:"assists"`
	Rebounds  float64 `json:"rebounds"`
	Steals    float64 `json:"steals"`
	Blocks    float64 `json:"blocks"`
	FGPct     float64 `json:"fg_pct"`
	ThreePct  float64 `json:"three_pct"`
	FTPct     float64 `json:"ft_pct"`
	Turnovers float64 `json:"turnovers"`
}

// CareerStats holds calculated career averages
type CareerStats struct {
	Name     string  `json:"name"`
	Sport    string  `json:"sport"`
	Seasons  int     `json:"seasons"`
	Games    int     `json:"games"`
	Points   float64 `json:"points"`
	Assists  float64 `json:"assists"`
	Rebounds float64 `json:"rebounds"`
	Steals   float64 `json:"steals"`
	Blocks   float64 `json:"blocks"`
	FGPct    float64 `json:"fg_pct"`
	ThreePct float64 `json:"three_pct"`
}

// searchDB searches for players by name in both name and search_name columns
func searchDB(query string, sport string) ([]DBPlayer, error) {
	q := "%" + strings.ToLower(query) + "%"
	normalized := "%" + strings.ToLower(normalizeText(query)) + "%"

	rows, err := db.Query(`
		SELECT id, name, age, team, position, sport, season,
		games, points, assists, rebounds, steals, blocks,
		fg_pct, three_pct, ft_pct, turnovers
		FROM players
		WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
		AND sport = ?
		ORDER BY season DESC, points DESC
		LIMIT 20
	`, q, normalized, sport)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	defer rows.Close()

	var players []DBPlayer
	for rows.Next() {
		var p DBPlayer
		err := rows.Scan(
			&p.ID, &p.Name, &p.Age, &p.Team, &p.Position,
			&p.Sport, &p.Season, &p.Games, &p.Points,
			&p.Assists, &p.Rebounds, &p.Steals, &p.Blocks,
			&p.FGPct, &p.ThreePct, &p.FTPct, &p.Turnovers,
		)
		if err != nil {
			continue
		}
		players = append(players, p)
	}
	return players, nil
}

// getPlayerByName gets a specific player's most recent season
func getPlayerByName(name string, sport string) (*DBPlayer, error) {
	q := "%" + strings.ToLower(name) + "%"
	normalized := "%" + strings.ToLower(normalizeText(name)) + "%"

	row := db.QueryRow(`
		SELECT id, name, age, team, position, sport, season,
		games, points, assists, rebounds, steals, blocks,
		fg_pct, three_pct, ft_pct, turnovers
		FROM players
		WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
		AND sport = ?
		ORDER BY season DESC
		LIMIT 1
	`, q, normalized, sport)

	var p DBPlayer
	err := row.Scan(
		&p.ID, &p.Name, &p.Age, &p.Team, &p.Position,
		&p.Sport, &p.Season, &p.Games, &p.Points,
		&p.Assists, &p.Rebounds, &p.Steals, &p.Blocks,
		&p.FGPct, &p.ThreePct, &p.FTPct, &p.Turnovers,
	)
	if err != nil {
		return nil, fmt.Errorf("player not found: %w", err)
	}
	return &p, nil
}

// getCareerAverages calculates career averages across all seasons
func getCareerAverages(name string, sport string) (*CareerStats, error) {
	q := "%" + strings.ToLower(name) + "%"
	normalized := "%" + strings.ToLower(normalizeText(name)) + "%"

	rows, err := db.Query(`
		SELECT games, points, assists, rebounds, steals, blocks,
		fg_pct, three_pct
		FROM players
		WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
		AND sport = ?
		ORDER BY season ASC
	`, q, normalized, sport)
	if err != nil {
		return nil, fmt.Errorf("career query failed: %w", err)
	}
	defer rows.Close()

	var totalGames int
	var totalPts, totalAst, totalReb, totalStl, totalBlk, totalFg, totalFg3 float64
	seasons := 0

	for rows.Next() {
		var games int
		var pts, ast, reb, stl, blk, fg, fg3 float64
		err := rows.Scan(&games, &pts, &ast, &reb, &stl, &blk, &fg, &fg3)
		if err != nil {
			continue
		}
		totalPts += pts * float64(games)
		totalAst += ast * float64(games)
		totalReb += reb * float64(games)
		totalStl += stl * float64(games)
		totalBlk += blk * float64(games)
		totalFg += fg * float64(games)
		totalFg3 += fg3 * float64(games)
		totalGames += games
		seasons++
	}

	if seasons == 0 {
		return nil, fmt.Errorf("no career data found for %s", name)
	}

	g := float64(totalGames)
	return &CareerStats{
		Name:     name,
		Sport:    sport,
		Seasons:  seasons,
		Games:    totalGames,
		Points:   totalPts / g,
		Assists:  totalAst / g,
		Rebounds: totalReb / g,
		Steals:   totalStl / g,
		Blocks:   totalBlk / g,
		FGPct:    totalFg / g,
		ThreePct: totalFg3 / g,
	}, nil
}