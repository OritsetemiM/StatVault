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