package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// NFLPassing holds QB stats for one season
type NFLPassing struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Age           int     `json:"age"`
	Team          string  `json:"team"`
	Position      string  `json:"position"`
	Season        string  `json:"season"`
	Games         int     `json:"games"`
	Completions   int     `json:"completions"`
	Attempts      int     `json:"attempts"`
	Yards         int     `json:"yards"`
	Touchdowns    int     `json:"touchdowns"`
	Interceptions int     `json:"interceptions"`
	CompPct       float64 `json:"comp_pct"`
	YardsPerGame  float64 `json:"yards_per_game"`
	Rating        float64 `json:"rating"`
}

// NFLRushing holds RB stats for one season
type NFLRushing struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Age          int     `json:"age"`
	Team         string  `json:"team"`
	Position     string  `json:"position"`
	Season       string  `json:"season"`
	Games        int     `json:"games"`
	Attempts     int     `json:"attempts"`
	Yards        int     `json:"yards"`
	Touchdowns   int     `json:"touchdowns"`
	YardsPerCarry float64 `json:"yards_per_carry"`
	YardsPerGame float64 `json:"yards_per_game"`
	Fumbles      int     `json:"fumbles"`
}

// NFLReceiving holds WR/TE stats for one season
type NFLReceiving struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Age          int     `json:"age"`
	Team         string  `json:"team"`
	Position     string  `json:"position"`
	Season       string  `json:"season"`
	Games        int     `json:"games"`
	Receptions   int     `json:"receptions"`
	Targets      int     `json:"targets"`
	Yards        int     `json:"yards"`
	Touchdowns   int     `json:"touchdowns"`
	YardsPerRec  float64 `json:"yards_per_rec"`
	YardsPerGame float64 `json:"yards_per_game"`
	Fumbles      int     `json:"fumbles"`
}

// initNFLDB creates the three NFL tables
func initNFLDB() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS nfl_passing (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			search_name TEXT,
			age INTEGER,
			team TEXT,
			position TEXT,
			season TEXT,
			games INTEGER,
			completions INTEGER,
			attempts INTEGER,
			yards INTEGER,
			touchdowns INTEGER,
			interceptions INTEGER,
			comp_pct REAL,
			yards_per_game REAL,
			rating REAL,
			UNIQUE(name, season)
		)
	`)
	if err != nil {
		fmt.Println("Error creating nfl_passing table:", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS nfl_rushing (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			search_name TEXT,
			age INTEGER,
			team TEXT,
			position TEXT,
			season TEXT,
			games INTEGER,
			attempts INTEGER,
			yards INTEGER,
			touchdowns INTEGER,
			yards_per_carry REAL,
			yards_per_game REAL,
			fumbles INTEGER,
			UNIQUE(name, season)
		)
	`)
	if err != nil {
		fmt.Println("Error creating nfl_rushing table:", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS nfl_receiving (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			search_name TEXT,
			age INTEGER,
			team TEXT,
			position TEXT,
			season TEXT,
			games INTEGER,
			receptions INTEGER,
			targets INTEGER,
			yards INTEGER,
			touchdowns INTEGER,
			yards_per_rec REAL,
			yards_per_game REAL,
			fumbles INTEGER,
			UNIQUE(name, season)
		)
	`)
	if err != nil {
		fmt.Println("Error creating nfl_receiving table:", err)
	}

	fmt.Println("NFL tables initialized")
}

// firstColIndex finds the index of the first occurrence of a column name
func firstColIndex(headers []string, name string) int {
	for i, h := range headers {
		if strings.TrimSpace(h) == name {
			return i
		}
	}
	return -1
}

// importNFLPassing reads a PFR passing CSV and loads it into SQLite
func importNFLPassing(filename string, season string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("could not read headers: %w", err)
	}

	// map column names to indices
	colIdx := map[string]int{}
	for i, h := range headers {
		colIdx[strings.TrimSpace(h)] = i
	}

	// Yds appears twice in passing CSV — first one is passing yards
	ydsIdx := firstColIndex(headers, "Yds")

	_, err = db.Exec("DELETE FROM nfl_passing WHERE season = ?", season)
	if err != nil {
		return fmt.Errorf("could not clear old data: %w", err)
	}

	stmt, err := db.Prepare(`
		INSERT OR IGNORE INTO nfl_passing
		(name, search_name, age, team, position, season, games,
		completions, attempts, yards, touchdowns, interceptions,
		comp_pct, yards_per_game, rating)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}
	defer stmt.Close()

	count := 0
	seen := map[string]bool{}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if len(row) == 0 || row[0] == "Rk" || row[0] == "" {
			continue
		}

		name := strings.TrimSpace(row[colIdx["Player"]])
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true

		// remove * and + suffixes PFR adds to Pro Bowl players
		name = strings.TrimRight(name, "*+")
		name = strings.TrimSpace(name)

		age := parseInt(row[colIdx["Age"]])
		team := strings.TrimSpace(row[colIdx["Team"]])
		pos := strings.TrimSpace(row[colIdx["Pos"]])
		games := parseInt(row[colIdx["G"]])
		cmp := parseInt(row[colIdx["Cmp"]])
		att := parseInt(row[colIdx["Att"]])
		yds := parseInt(row[ydsIdx]) // use first Yds column (passing yards)
		td := parseInt(row[colIdx["TD"]])
		interc := parseInt(row[colIdx["Int"]])
		cmpPct := parseFloat(row[colIdx["Cmp%"]])
		ypg := parseFloat(row[colIdx["Y/G"]])
		rating := parseFloat(row[colIdx["Rate"]])

		_, err = stmt.Exec(
			name, normalizeText(name), age, team, pos, season,
			games, cmp, att, yds, td, interc, cmpPct, ypg, rating,
		)
		if err != nil {
			continue
		}
		count++
	}

	fmt.Printf("Imported %d passers for %s season\n", count, season)
	return nil
}

// importNFLRushing reads a PFR rushing CSV
func importNFLRushing(filename string, season string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("could not read first row: %w", err)
	}

	// second row has the actual column names
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("could not read headers: %w", err)
	}

	colIdx := map[string]int{}
	for i, h := range headers {
		colIdx[strings.TrimSpace(h)] = i
	}

	// first Yds column is rushing yards
	ydsIdx := firstColIndex(headers, "Yds")

	_, err = db.Exec("DELETE FROM nfl_rushing WHERE season = ?", season)
	if err != nil {
		return fmt.Errorf("could not clear old data: %w", err)
	}

	stmt, err := db.Prepare(`
		INSERT OR IGNORE INTO nfl_rushing
		(name, search_name, age, team, position, season, games,
		attempts, yards, touchdowns, yards_per_carry, yards_per_game, fumbles)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}
	defer stmt.Close()

	count := 0
	seen := map[string]bool{}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if len(row) == 0 || row[0] == "Rk" || row[0] == "" {
			continue
		}

		name := strings.TrimSpace(row[colIdx["Player"]])
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true

		name = strings.TrimRight(name, "*+")
		name = strings.TrimSpace(name)

		age := parseInt(row[colIdx["Age"]])
		team := strings.TrimSpace(row[colIdx["Team"]])
		pos := strings.TrimSpace(row[colIdx["Pos"]])
		games := parseInt(row[colIdx["G"]])
		att := parseInt(row[colIdx["Att"]])
		yds := parseInt(row[ydsIdx])
		td := parseInt(row[colIdx["TD"]])
		ypc := parseFloat(row[colIdx["Y/A"]])
		ypg := parseFloat(row[colIdx["Y/G"]])
		fum := parseInt(row[colIdx["Fmb"]])

		_, err = stmt.Exec(
			name, normalizeText(name), age, team, pos, season,
			games, att, yds, td, ypc, ypg, fum,
		)
		if err != nil {
			continue
		}
		count++
	}

	fmt.Printf("Imported %d rushers for %s season\n", count, season)
	return nil
}

// importNFLReceiving reads a PFR receiving CSV
func importNFLReceiving(filename string, season string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	// skip the first row (category labels like "Rushing, Rushing...")
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("could not read first row: %w", err)
	}

	// second row has the actual column names
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("could not read headers: %w", err)
	}

	colIdx := map[string]int{}
	for i, h := range headers {
		colIdx[strings.TrimSpace(h)] = i
	}

	// first Yds column is receiving yards
	ydsIdx := firstColIndex(headers, "Yds")

	_, err = db.Exec("DELETE FROM nfl_receiving WHERE season = ?", season)
	if err != nil {
		return fmt.Errorf("could not clear old data: %w", err)
	}

	stmt, err := db.Prepare(`
		INSERT OR IGNORE INTO nfl_receiving
		(name, search_name, age, team, position, season, games,
		receptions, targets, yards, touchdowns, yards_per_rec, yards_per_game, fumbles)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}
	defer stmt.Close()

	count := 0
	seen := map[string]bool{}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if len(row) == 0 || row[0] == "Rk" || row[0] == "" {
			continue
		}

		name := strings.TrimSpace(row[colIdx["Player"]])
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true

		name = strings.TrimRight(name, "*+")
		name = strings.TrimSpace(name)

		age := parseInt(row[colIdx["Age"]])
		team := strings.TrimSpace(row[colIdx["Team"]])
		pos := strings.TrimSpace(row[colIdx["Pos"]])
		games := parseInt(row[colIdx["G"]])
		rec := parseInt(row[colIdx["Rec"]])
		tgt := parseInt(row[colIdx["Tgt"]])
		yds := parseInt(row[ydsIdx])
		td := parseInt(row[colIdx["TD"]])
		ypr := parseFloat(row[colIdx["Y/R"]])
		ypg := parseFloat(row[colIdx["Y/G"]])
		fum := parseInt(row[colIdx["Fmb"]])

		_, err = stmt.Exec(
			name, normalizeText(name), age, team, pos, season,
			games, rec, tgt, yds, td, ypr, ypg, fum,
		)
		if err != nil {
			continue
		}
		count++
	}

	fmt.Printf("Imported %d receivers for %s season\n", count, season)
	return nil
}

// searchNFLPassing searches for QBs by name
func searchNFLPassing(query string) ([]NFLPassing, error) {
	q := "%" + strings.ToLower(query) + "%"
	normalized := "%" + strings.ToLower(normalizeText(query)) + "%"

	rows, err := db.Query(`
		SELECT id, name, age, team, position, season, games,
		completions, attempts, yards, touchdowns, interceptions,
		comp_pct, yards_per_game, rating
		FROM nfl_passing
		WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
		ORDER BY season DESC, yards DESC
		LIMIT 20
	`, q, normalized)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []NFLPassing
	for rows.Next() {
		var p NFLPassing
		err := rows.Scan(
			&p.ID, &p.Name, &p.Age, &p.Team, &p.Position, &p.Season,
			&p.Games, &p.Completions, &p.Attempts, &p.Yards,
			&p.Touchdowns, &p.Interceptions, &p.CompPct,
			&p.YardsPerGame, &p.Rating,
		)
		if err != nil {
			continue
		}
		results = append(results, p)
	}
	return results, nil
}

// searchNFLRushing searches for RBs by name
func searchNFLRushing(query string) ([]NFLRushing, error) {
	q := "%" + strings.ToLower(query) + "%"
	normalized := "%" + strings.ToLower(normalizeText(query)) + "%"

	rows, err := db.Query(`
		SELECT id, name, age, team, position, season, games,
		attempts, yards, touchdowns, yards_per_carry, yards_per_game, fumbles
		FROM nfl_rushing
		WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
		ORDER BY season DESC, yards DESC
		LIMIT 20
	`, q, normalized)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []NFLRushing
	for rows.Next() {
		var p NFLRushing
		err := rows.Scan(
			&p.ID, &p.Name, &p.Age, &p.Team, &p.Position, &p.Season,
			&p.Games, &p.Attempts, &p.Yards, &p.Touchdowns,
			&p.YardsPerCarry, &p.YardsPerGame, &p.Fumbles,
		)
		if err != nil {
			continue
		}
		results = append(results, p)
	}
	return results, nil
}

// searchNFLReceiving searches for WRs/TEs by name
func searchNFLReceiving(query string) ([]NFLReceiving, error) {
	q := "%" + strings.ToLower(query) + "%"
	normalized := "%" + strings.ToLower(normalizeText(query)) + "%"

	rows, err := db.Query(`
		SELECT id, name, age, team, position, season, games,
		receptions, targets, yards, touchdowns, yards_per_rec, yards_per_game, fumbles
		FROM nfl_receiving
		WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
		ORDER BY season DESC, yards DESC
		LIMIT 20
	`, q, normalized)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []NFLReceiving
	for rows.Next() {
		var p NFLReceiving
		err := rows.Scan(
			&p.ID, &p.Name, &p.Age, &p.Team, &p.Position, &p.Season,
			&p.Games, &p.Receptions, &p.Targets, &p.Yards,
			&p.Touchdowns, &p.YardsPerRec, &p.YardsPerGame, &p.Fumbles,
		)
		if err != nil {
			continue
		}
		results = append(results, p)
	}
	return results, nil
}