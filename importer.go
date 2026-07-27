package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// parseFloat safely parses a float, returning 0 on error
func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" || s == "." {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

// parseInt safely parses an int, returning 0 on error
func parseInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

// importNBACSV reads a basketball-reference CSV and loads it into SQLite
func importNBACSV(filename string, season string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	// read header row
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("could not read headers: %w", err)
	}

	// map column names to indices
	colIdx := map[string]int{}
	for i, h := range headers {
		colIdx[strings.TrimSpace(h)] = i
	}

	// clear existing data for this season
	_, err = db.Exec("DELETE FROM players WHERE season = ? AND sport = 'NBA'", season)
	if err != nil {
		return fmt.Errorf("could not clear old data: %w", err)
	}

	// prepare insert statement
	stmt, err := db.Prepare(`
		INSERT OR IGNORE INTO players 
		(name, search_name, age, team, position, sport, season, games, games_started, 
		minutes_per_game, points, assists, rebounds, steals, blocks, 
		fg_pct, three_pct, ft_pct, turnovers)
		VALUES (?, ?, ?, ?, ?, 'NBA', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}
	defer stmt.Close()

	count := 0
	seenPlayers := map[string]bool{}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		// skip header rows that repeat in the middle of the table
		if len(row) == 0 || row[0] == "Rk" || row[0] == "" {
			continue
		}

		// get player name
		name := strings.TrimSpace(row[colIdx["Player"]])
		if name == "" {
			continue
		}

		// skip duplicates — keep first (total stats for traded players)
		if seenPlayers[name] {
			continue
		}
		seenPlayers[name] = true

		// extract stats
		age := parseInt(row[colIdx["Age"]])
		team := strings.TrimSpace(row[colIdx["Team"]])
		pos := strings.TrimSpace(row[colIdx["Pos"]])
		games := parseInt(row[colIdx["G"]])
		gamesStarted := parseInt(row[colIdx["GS"]])
		mp := parseFloat(row[colIdx["MP"]])
		pts := parseFloat(row[colIdx["PTS"]])
		ast := parseFloat(row[colIdx["AST"]])
		reb := parseFloat(row[colIdx["TRB"]])
		stl := parseFloat(row[colIdx["STL"]])
		blk := parseFloat(row[colIdx["BLK"]])
		fgPct := parseFloat(row[colIdx["FG%"]])
		threePct := parseFloat(row[colIdx["3P%"]])
		ftPct := parseFloat(row[colIdx["FT%"]])
		tov := parseFloat(row[colIdx["TOV"]])

		_, err = stmt.Exec(
			name, normalizeText(name), age, team, pos, season,
			games, gamesStarted, mp,
			pts, ast, reb, stl, blk,
			fgPct, threePct, ftPct, tov,
		)
		if err != nil {
			fmt.Printf("Warning: could not insert %s: %v\n", name, err)
			continue
		}
		count++
	}

	fmt.Printf("Imported %d players for %s season\n", count, season)
	return nil
}