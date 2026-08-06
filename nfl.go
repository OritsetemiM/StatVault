package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
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
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Age           int     `json:"age"`
	Team          string  `json:"team"`
	Position      string  `json:"position"`
	Season        string  `json:"season"`
	Games         int     `json:"games"`
	Attempts      int     `json:"attempts"`
	Yards         int     `json:"yards"`
	Touchdowns    int     `json:"touchdowns"`
	YardsPerCarry float64 `json:"yards_per_carry"`
	YardsPerGame  float64 `json:"yards_per_game"`
	Fumbles       int     `json:"fumbles"`
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

// NFLDefense holds defensive player stats for one season
type NFLDefense struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Age           int     `json:"age"`
	Team          string  `json:"team"`
	Position      string  `json:"position"`
	Season        string  `json:"season"`
	Games         int     `json:"games"`
	Interceptions int     `json:"interceptions"`
	Sacks         float64 `json:"sacks"`
	Tackles       int     `json:"tackles"`
	PassDef       int     `json:"pass_def"`
	ForcedFum     int     `json:"forced_fum"`
	QBHits        int     `json:"qb_hits"`
}

// NFLPlayerProfile holds all stats across all tables for one player
type NFLPlayerProfile struct {
	Name      string         `json:"name"`
	Position  string         `json:"position"`
	Team      string         `json:"team"`
	Passing   []NFLPassing   `json:"passing,omitempty"`
	Rushing   []NFLRushing   `json:"rushing,omitempty"`
	Receiving []NFLReceiving `json:"receiving,omitempty"`
	Defense   []NFLDefense   `json:"defense,omitempty"`
}

// initNFLDB creates all four NFL tables
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

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS nfl_defense (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			search_name TEXT,
			age INTEGER,
			team TEXT,
			position TEXT,
			season TEXT,
			games INTEGER,
			interceptions INTEGER,
			sacks REAL,
			tackles INTEGER,
			pass_def INTEGER,
			forced_fum INTEGER,
			qb_hits INTEGER,
			UNIQUE(name, season)
		)
	`)
	if err != nil {
		fmt.Println("Error creating nfl_defense table:", err)
	}

	fmt.Println("NFL tables initialized")
}

// firstColIndex finds the first occurrence of a column name
func firstColIndex(headers []string, name string) int {
	for i, h := range headers {
		if strings.TrimSpace(h) == name {
			return i
		}
	}
	return -1
}

// cleanName removes PFR suffixes and trims whitespace
func cleanName(name string) string {
	name = strings.TrimRight(name, "*+")
	return strings.TrimSpace(name)
}

// importNFLPassing reads a PFR passing CSV
func importNFLPassing(filename string, season string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("could not read headers: %w", err)
	}

	colIdx := map[string]int{}
	for i, h := range headers {
		colIdx[strings.TrimSpace(h)] = i
	}

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

		if _, err := strconv.Atoi(row[0]); err != nil {
			continue
		}

		name := cleanName(row[colIdx["Player"]])
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true

		age := parseInt(row[colIdx["Age"]])
		team := strings.TrimSpace(row[colIdx["Team"]])
		pos := strings.TrimSpace(row[colIdx["Pos"]])
		games := parseInt(row[colIdx["G"]])
		cmp := parseInt(row[colIdx["Cmp"]])
		att := parseInt(row[colIdx["Att"]])
		yds := 0
		if ydsIdx >= 0 && ydsIdx < len(row) {
			yds = parseInt(row[ydsIdx])
		}
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
func importNFLRushing(filename string, season string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	// read first row to detect format
	firstRow, err := reader.Read()
	if err != nil {
		return fmt.Errorf("could not read first row: %w", err)
	}

	// detect if first row is headers, category labels, or data
	hasHeaders := false
	var headers []string

	if firstRow[0] == "Rk" {
		// normal header row
		headers = firstRow
		hasHeaders = true
	} else if _, numErr := strconv.Atoi(firstRow[0]); numErr == nil {
		// first row is already data — no headers at all
		// use standard PFR rushing column order
		headers = []string{
    		"Rk","Player","Age","Team","Pos","G","GS",
    	"Tgt","Rec","Yds","Y/R","TD","1D","Succ%","Lng","Y/Tgt","R/G","Y/G","Ctch%","YScm","RRTD","Fmb","Awards","Player-additional",
		}
		hasHeaders = false
	} else {
		// category labels row — read next row for real headers
		headers, err = reader.Read()
		if err != nil {
			return fmt.Errorf("could not read headers: %w", err)
		}
		hasHeaders = true
	}

	// build column index map
	colIdx := map[string]int{}
	for i, h := range headers {
		colIdx[strings.TrimSpace(h)] = i
	}

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

	// process first row if it was data
	processRow := func(row []string) {
		if len(row) < 8 {
			return
		}
		if row[0] == "Rk" || row[0] == "" {
			return
		}
		if _, err := strconv.Atoi(row[0]); err != nil {
			return
		}

		name := cleanName(row[colIdx["Player"]])
		if name == "" || seen[name] {
			return
		}
		seen[name] = true

		age := parseInt(row[colIdx["Age"]])
		team := strings.TrimSpace(row[colIdx["Team"]])
		pos := strings.TrimSpace(row[colIdx["Pos"]])
		games := parseInt(row[colIdx["G"]])
		att := parseInt(row[colIdx["Att"]])
		yds := 0
		if ydsIdx >= 0 && ydsIdx < len(row) {
			yds = parseInt(row[ydsIdx])
		}
		td := parseInt(row[colIdx["TD"]])
		ypc := parseFloat(row[colIdx["Y/A"]])
		ypg := parseFloat(row[colIdx["Y/G"]])
		fum := parseInt(row[colIdx["Fmb"]])

		stmt.Exec(
			name, normalizeText(name), age, team, pos, season,
			games, att, yds, td, ypc, ypg, fum,
		)
		count++
	}

	// if first row was data, process it
	if !hasHeaders {
		processRow(firstRow)
	}

	// process remaining rows
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		processRow(row)
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
	reader.FieldsPerRecord = -1

	firstRow, err := reader.Read()
	if err != nil {
		return fmt.Errorf("could not read first row: %w", err)
	}

	hasHeaders := false
	var headers []string

	if firstRow[0] == "Rk" {
		headers = firstRow
		hasHeaders = true
	} else if _, numErr := strconv.Atoi(firstRow[0]); numErr == nil {
		headers = []string{
			"Rk","Player","Age","Team","Pos","G","GS",
			"Tgt","Rec","Yds","Y/R","TD","1D","Succ%","Lng","Y/Tgt","R/G","Y/G","Ctch%","YScm","RRTD","Fmb","Awards","Player-additional",
		}
		hasHeaders = false
	} else {
		headers, err = reader.Read()
		if err != nil {
			return fmt.Errorf("could not read headers: %w", err)
		}
		hasHeaders = true
	}

	colIdx := map[string]int{}
	for i, h := range headers {
		colIdx[strings.TrimSpace(h)] = i
	}

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

	processRow := func(row []string) {
		if len(row) < 10 {
			return
		}
		if row[0] == "Rk" || row[0] == "" {
			return
		}
		if _, err := strconv.Atoi(row[0]); err != nil {
			return
		}

		name := cleanName(row[colIdx["Player"]])
		if name == "" || seen[name] {
			return
		}
		seen[name] = true

		age := parseInt(row[colIdx["Age"]])
		team := strings.TrimSpace(row[colIdx["Team"]])
		pos := strings.TrimSpace(row[colIdx["Pos"]])
		games := parseInt(row[colIdx["G"]])
		rec := parseInt(row[colIdx["Rec"]])
		tgt := parseInt(row[colIdx["Tgt"]])
		yds := 0
		if ydsIdx >= 0 && ydsIdx < len(row) {
			yds = parseInt(row[ydsIdx])
		}
		td := parseInt(row[colIdx["TD"]])
		ypr := parseFloat(row[colIdx["Y/R"]])
		ypg := parseFloat(row[colIdx["Y/G"]])
		fum := parseInt(row[colIdx["Fmb"]])

		stmt.Exec(
			name, normalizeText(name), age, team, pos, season,
			games, rec, tgt, yds, td, ypr, ypg, fum,
		)
		count++
	}

	if !hasHeaders {
		processRow(firstRow)
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		processRow(row)
	}

	fmt.Printf("Imported %d receivers for %s season\n", count, season)
	return nil
}
// importNFLDefense reads a PFR defense CSV (no header row)
func importNFLDefense(filename string, season string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	_, err = db.Exec("DELETE FROM nfl_defense WHERE season = ?", season)
	if err != nil {
		return fmt.Errorf("could not clear old data: %w", err)
	}

	stmt, err := db.Prepare(`
		INSERT OR IGNORE INTO nfl_defense
		(name, search_name, age, team, position, season, games,
		interceptions, sacks, tackles, pass_def, forced_fum, qb_hits)
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

		// defense CSV has no headers — data starts immediately
		// columns: Rk, Player, Age, Team, Pos, G, GS, Int, Yds, TD, PD, FF, FR, Yds, TD, Sk, Comb, Solo, Ast, TFL, QBHits, Sfty, Awards, Player-additional
		if len(row) < 16 {
			continue
		}

		// skip non-numeric rank rows
		if _, err := strconv.Atoi(row[0]); err != nil {
			continue
		}

		name := cleanName(row[1])
		if name == "" {
			continue
		}

		// keep first occurrence (combined stats for traded players)
		if seen[name] {
			continue
		}
		seen[name] = true

		age := parseInt(row[2])
		team := strings.TrimSpace(row[3])
		pos := strings.TrimSpace(row[4])
		// skip pure offensive positions unless they have real defensive stats
	offensivePositions := map[string]bool{
		"QB": true, "RB": true, "FB": true, "TE": true,
		"OL": true, "OT": true, "OG": true, "C": true, "G": true, "T": true,
		"K": true, "P": true, "LS": true,
	}

	if offensivePositions[pos] {
		continue
	}

	// for WRs — only include if they have meaningful defensive stats
	// (Travis Hunter type players)
	if pos == "WR" {
		tackles := parseInt(row[16])
		interc := parseInt(row[7])
		sacks := parseFloat(row[15])
		passDef := parseInt(row[10])
		if tackles == 0 && interc == 0 && sacks == 0 && passDef == 0 {
			continue
		}
}
		games := parseInt(row[5])
		interc := parseInt(row[7])
		passDef := parseInt(row[10])
		forcedFum := parseInt(row[11])
		sacks := parseFloat(row[15])
		tackles := parseInt(row[16])

		qbHits := 0
		if len(row) > 20 {
			qbHits = parseInt(row[20])
		}

		_, err = stmt.Exec(
			name, normalizeText(name), age, team, pos, season,
			games, interc, sacks, tackles, passDef, forcedFum, qbHits,
		)
		if err != nil {
			continue
		}
		count++
	}

	fmt.Printf("Imported %d defensive players for %s season\n", count, season)
	return nil
}

// searchNFLAll searches all 4 NFL tables and returns combined results
// for players with stats in multiple tables, shows their best stat type
func searchNFLAll(query string) ([]map[string]interface{}, error) {
	passing, _ := searchNFLPassing(query)
	rushing, _ := searchNFLRushing(query)
	receiving, _ := searchNFLReceiving(query)
	defense, _ := searchNFLDefense(query)

	// track which stat types each player has
	playerStats := map[string][]string{}
	playerInfo := map[string]map[string]interface{}{}

	for _, p := range passing {
		playerStats[p.Name] = append(playerStats[p.Name], "passing")
		if _, ok := playerInfo[p.Name]; !ok {
			playerInfo[p.Name] = map[string]interface{}{
				"name": p.Name, "position": p.Position,
				"team": p.Team, "season": p.Season,
				"yards": p.Yards, "stat_type": "passing",
			}
		}
	}
	for _, p := range rushing {
		playerStats[p.Name] = append(playerStats[p.Name], "rushing")
		if _, ok := playerInfo[p.Name]; !ok {
			playerInfo[p.Name] = map[string]interface{}{
				"name": p.Name, "position": p.Position,
				"team": p.Team, "season": p.Season,
				"yards": p.Yards, "stat_type": "rushing",
			}
		}
	}
	for _, p := range receiving {
		playerStats[p.Name] = append(playerStats[p.Name], "receiving")
		if _, ok := playerInfo[p.Name]; !ok {
			playerInfo[p.Name] = map[string]interface{}{
				"name": p.Name, "position": p.Position,
				"team": p.Team, "season": p.Season,
				"yards": p.Yards, "stat_type": "receiving",
			}
		}
	}
	for _, p := range defense {
		playerStats[p.Name] = append(playerStats[p.Name], "defense")
		if _, ok := playerInfo[p.Name]; !ok {
			playerInfo[p.Name] = map[string]interface{}{
				"name": p.Name, "position": p.Position,
				"team": p.Team, "season": p.Season,
				"yards": 0, "stat_type": "defense",
				"sacks": p.Sacks,
			}
		}
	}

	// build results with stat_types list
	var results []map[string]interface{}
	for name, info := range playerInfo {
		info["stat_types"] = playerStats[name]
		results = append(results, info)
	}

	return results, nil
}

// getNFLProfile gets ALL stats for a player across all tables
func getNFLProfile(name string) (*NFLPlayerProfile, error) {
	passing, _ := searchNFLPassing(name)
	rushing, _ := searchNFLRushing(name)
	receiving, _ := searchNFLReceiving(name)
	defense, _ := searchNFLDefense(name)

	// filter to exact name matches only
	exactPassing := filterExactPassing(passing, name)
	exactRushing := filterExactRushing(rushing, name)
	exactReceiving := filterExactReceiving(receiving, name)
	exactDefense := filterExactDefense(defense, name)

	if len(exactPassing) == 0 && len(exactRushing) == 0 &&
		len(exactReceiving) == 0 && len(exactDefense) == 0 {
		return nil, fmt.Errorf("player not found: %s", name)
	}

	// get position and team from most recent data
	pos, team := "", ""
	if len(exactPassing) > 0 {
		pos = exactPassing[0].Position
		team = exactPassing[0].Team
	} else if len(exactRushing) > 0 {
		pos = exactRushing[0].Position
		team = exactRushing[0].Team
	} else if len(exactReceiving) > 0 {
		pos = exactReceiving[0].Position
		team = exactReceiving[0].Team
	} else if len(exactDefense) > 0 {
		pos = exactDefense[0].Position
		team = exactDefense[0].Team
	}

	return &NFLPlayerProfile{
		Name:      name,
		Position:  pos,
		Team:      team,
		Passing:   exactPassing,
		Rushing:   exactRushing,
		Receiving: exactReceiving,
		Defense:   exactDefense,
	}, nil
}

// exact match filters
func filterExactPassing(players []NFLPassing, name string) []NFLPassing {
	name = strings.ToLower(name)
	var result []NFLPassing
	for _, p := range players {
		if strings.ToLower(p.Name) == name {
			result = append(result, p)
		}
	}
	return result
}

func filterExactRushing(players []NFLRushing, name string) []NFLRushing {
	name = strings.ToLower(name)
	var result []NFLRushing
	for _, p := range players {
		if strings.ToLower(p.Name) == name {
			result = append(result, p)
		}
	}
	return result
}

func filterExactReceiving(players []NFLReceiving, name string) []NFLReceiving {
	name = strings.ToLower(name)
	var result []NFLReceiving
	for _, p := range players {
		if strings.ToLower(p.Name) == name {
			result = append(result, p)
		}
	}
	return result
}

func filterExactDefense(players []NFLDefense, name string) []NFLDefense {
	name = strings.ToLower(name)
	var result []NFLDefense
	for _, p := range players {
		if strings.ToLower(p.Name) == name {
			result = append(result, p)
		}
	}
	return result
}

// search functions
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

func searchNFLDefense(query string) ([]NFLDefense, error) {
	q := "%" + strings.ToLower(query) + "%"
	normalized := "%" + strings.ToLower(normalizeText(query)) + "%"

	rows, err := db.Query(`
		SELECT id, name, age, team, position, season, games,
		interceptions, sacks, tackles, pass_def, forced_fum, qb_hits
		FROM nfl_defense
		WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
		ORDER BY season DESC, sacks DESC
		LIMIT 20
	`, q, normalized)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []NFLDefense
	for rows.Next() {
		var p NFLDefense
		err := rows.Scan(
			&p.ID, &p.Name, &p.Age, &p.Team, &p.Position, &p.Season,
			&p.Games, &p.Interceptions, &p.Sacks, &p.Tackles,
			&p.PassDef, &p.ForcedFum, &p.QBHits,
		)
		if err != nil {
			continue
		}
		results = append(results, p)
	}
	return results, nil
}