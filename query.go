package main

import (
	"fmt"
	"strings"
)

// DBPlayer represents a player row from the database
type DBPlayer struct {
	ID       int
	Name     string
	Age      int
	Team     string
	Position string
	Sport    string
	Season   string
	Games    int
	Points   float64
	Assists  float64
	Rebounds float64
	Steals   float64
	Blocks   float64
	FGPct    float64
	ThreePct float64
	FTPct    float64
	Turnovers float64
}

// searchDB searches for players by name in the database
func searchDB(query string, sport string) ([]DBPlayer, error) {
	query = "%" + strings.ToLower(query) + "%"

	rows, err := db.Query(`
		SELECT id, name, age, team, position, sport, season,
		games, points, assists, rebounds, steals, blocks,
		fg_pct, three_pct, ft_pct, turnovers
		FROM players
		WHERE LOWER(name) LIKE ?
		AND sport = ?
		ORDER BY points DESC
		LIMIT 10
	`, query, sport)
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

// getPlayerByName gets a specific player's full stats
func getPlayerByName(name string, sport string) (*DBPlayer, error) {
	query := "%" + strings.ToLower(name) + "%"

	row := db.QueryRow(`
		SELECT id, name, age, team, position, sport, season,
		games, points, assists, rebounds, steals, blocks,
		fg_pct, three_pct, ft_pct, turnovers
		FROM players
		WHERE LOWER(name) LIKE ?
		AND sport = ?
		ORDER BY points DESC
		LIMIT 1
	`, query, sport)

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

// CareerStats holds calculated career averages from multiple seasons
type CareerStats struct {
	Name     string
	Sport    string
	Seasons  int
	Games    int
	Points   float64
	Assists  float64
	Rebounds float64
	Steals   float64
	Blocks   float64
	FGPct    float64
	ThreePct float64
}

// getCareerAverages calculates career averages across all seasons
func getCareerAverages(name string, sport string) (*CareerStats, error) {
	query := "%" + strings.ToLower(name) + "%"

	rows, err := db.Query(`
		SELECT games, points, assists, rebounds, steals, blocks,
		fg_pct, three_pct
		FROM players
		WHERE LOWER(name) LIKE ?
		AND sport = ?
		ORDER BY season ASC
	`, query, sport)
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
		// weight averages by games played
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