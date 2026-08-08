package main

import (
	"fmt"
	"strings"
)

// CompareCategory holds one stat comparison
type CompareCategory struct {
	Label  string  `json:"label"`
	Val1   float64 `json:"val1"`
	Val2   float64 `json:"val2"`
	Unit   string  `json:"unit"`
	Winner int     `json:"winner"` // 1 = player1, 2 = player2, 0 = tie
}

// CompareResult holds the full comparison
type CompareResult struct {
	Player1    string            `json:"player1"`
	Player2    string            `json:"player2"`
	Season1    string            `json:"season1"`
	Season2    string            `json:"season2"`
	Sport      string            `json:"sport"`
	Categories []CompareCategory `json:"categories"`
	Score1     int               `json:"score1"`
	Score2     int               `json:"score2"`
	Winner     string            `json:"winner"`
}

// compareVal returns who wins a stat (higher is better unless lowerBetter)
func compareVal(label string, v1 float64, v2 float64, unit string, lowerBetter bool) CompareCategory {
	winner := 0
	if v1 > v2 {
		if lowerBetter {
			winner = 2
		} else {
			winner = 1
		}
	} else if v2 > v1 {
		if lowerBetter {
			winner = 1
		} else {
			winner = 2
		}
	}
	return CompareCategory{Label: label, Val1: v1, Val2: v2, Unit: unit, Winner: winner}
}

// CompareNBAPlayers compares two NBA players by name and season
func CompareNBAPlayers(name1 string, season1 string, name2 string, season2 string) (*CompareResult, error) {
	p1, err := getNBAPlayerSeason(name1, season1)
	if err != nil {
		return nil, fmt.Errorf("player 1 not found: %w", err)
	}
	p2, err := getNBAPlayerSeason(name2, season2)
	if err != nil {
		return nil, fmt.Errorf("player 2 not found: %w", err)
	}

	categories := []CompareCategory{
		compareVal("Points Per Game", p1.Points, p2.Points, "PPG", false),
		compareVal("Assists Per Game", p1.Assists, p2.Assists, "APG", false),
		compareVal("Rebounds Per Game", p1.Rebounds, p2.Rebounds, "RPG", false),
		compareVal("Steals Per Game", p1.Steals, p2.Steals, "SPG", false),
		compareVal("Blocks Per Game", p1.Blocks, p2.Blocks, "BPG", false),
		compareVal("Field Goal %", p1.FGPct*100, p2.FGPct*100, "%", false),
		compareVal("3-Point %", p1.ThreePct*100, p2.ThreePct*100, "%", false),
		compareVal("Games Played", float64(p1.Games), float64(p2.Games), "G", false),
		compareVal("Turnovers Per Game", p1.Turnovers, p2.Turnovers, "TPG", true),
	}

	return buildResult(p1.Name, p2.Name, p1.Season, p2.Season, "NBA", categories), nil
}

// CompareNFLPassers compares two QBs
func CompareNFLPassers(name1 string, season1 string, name2 string, season2 string) (*CompareResult, error) {
	p1, err := getNFLPasserSeason(name1, season1)
	if err != nil {
		return nil, fmt.Errorf("player 1 not found: %w", err)
	}
	p2, err := getNFLPasserSeason(name2, season2)
	if err != nil {
		return nil, fmt.Errorf("player 2 not found: %w", err)
	}

	categories := []CompareCategory{
		compareVal("Passing Yards", float64(p1.Yards), float64(p2.Yards), "YDS", false),
		compareVal("Touchdowns", float64(p1.Touchdowns), float64(p2.Touchdowns), "TD", false),
		compareVal("Interceptions", float64(p1.Interceptions), float64(p2.Interceptions), "INT", true),
		compareVal("Completion %", p1.CompPct, p2.CompPct, "%", false),
		compareVal("Yards Per Game", p1.YardsPerGame, p2.YardsPerGame, "YPG", false),
		compareVal("Passer Rating", p1.Rating, p2.Rating, "RTG", false),
		compareVal("Games Played", float64(p1.Games), float64(p2.Games), "G", false),
	}

	return buildResult(p1.Name, p2.Name, p1.Season, p2.Season, "NFL", categories), nil
}

// CompareNFLRushers compares two RBs
func CompareNFLRushers(name1 string, season1 string, name2 string, season2 string) (*CompareResult, error) {
	p1, err := getNFLRusherSeason(name1, season1)
	if err != nil {
		return nil, fmt.Errorf("player 1 not found: %w", err)
	}
	p2, err := getNFLRusherSeason(name2, season2)
	if err != nil {
		return nil, fmt.Errorf("player 2 not found: %w", err)
	}

	categories := []CompareCategory{
		compareVal("Rushing Yards", float64(p1.Yards), float64(p2.Yards), "YDS", false),
		compareVal("Touchdowns", float64(p1.Touchdowns), float64(p2.Touchdowns), "TD", false),
		compareVal("Attempts", float64(p1.Attempts), float64(p2.Attempts), "ATT", false),
		compareVal("Yards Per Carry", p1.YardsPerCarry, p2.YardsPerCarry, "YPC", false),
		compareVal("Yards Per Game", p1.YardsPerGame, p2.YardsPerGame, "YPG", false),
		compareVal("Fumbles", float64(p1.Fumbles), float64(p2.Fumbles), "FUM", true),
		compareVal("Games Played", float64(p1.Games), float64(p2.Games), "G", false),
	}

	return buildResult(p1.Name, p2.Name, p1.Season, p2.Season, "NFL", categories), nil
}

// CompareNFLReceivers compares two WRs/TEs
func CompareNFLReceivers(name1 string, season1 string, name2 string, season2 string) (*CompareResult, error) {
	p1, err := getNFLReceiverSeason(name1, season1)
	if err != nil {
		return nil, fmt.Errorf("player 1 not found: %w", err)
	}
	p2, err := getNFLReceiverSeason(name2, season2)
	if err != nil {
		return nil, fmt.Errorf("player 2 not found: %w", err)
	}

	categories := []CompareCategory{
		compareVal("Receiving Yards", float64(p1.Yards), float64(p2.Yards), "YDS", false),
		compareVal("Touchdowns", float64(p1.Touchdowns), float64(p2.Touchdowns), "TD", false),
		compareVal("Receptions", float64(p1.Receptions), float64(p2.Receptions), "REC", false),
		compareVal("Targets", float64(p1.Targets), float64(p2.Targets), "TGT", false),
		compareVal("Yards Per Reception", p1.YardsPerRec, p2.YardsPerRec, "YPR", false),
		compareVal("Yards Per Game", p1.YardsPerGame, p2.YardsPerGame, "YPG", false),
		compareVal("Games Played", float64(p1.Games), float64(p2.Games), "G", false),
	}

	return buildResult(p1.Name, p2.Name, p1.Season, p2.Season, "NFL", categories), nil
}

// CompareNFLDefenders compares two defensive players
func CompareNFLDefenders(name1 string, season1 string, name2 string, season2 string) (*CompareResult, error) {
	p1, err := getNFLDefenderSeason(name1, season1)
	if err != nil {
		return nil, fmt.Errorf("player 1 not found: %w", err)
	}
	p2, err := getNFLDefenderSeason(name2, season2)
	if err != nil {
		return nil, fmt.Errorf("player 2 not found: %w", err)
	}

	categories := []CompareCategory{
		compareVal("Sacks", p1.Sacks, p2.Sacks, "SACKS", false),
		compareVal("Tackles", float64(p1.Tackles), float64(p2.Tackles), "TKL", false),
		compareVal("Interceptions", float64(p1.Interceptions), float64(p2.Interceptions), "INT", false),
		compareVal("Pass Deflections", float64(p1.PassDef), float64(p2.PassDef), "PD", false),
		compareVal("Forced Fumbles", float64(p1.ForcedFum), float64(p2.ForcedFum), "FF", false),
		compareVal("QB Hits", float64(p1.QBHits), float64(p2.QBHits), "QBH", false),
		compareVal("Games Played", float64(p1.Games), float64(p2.Games), "G", false),
	}

	return buildResult(p1.Name, p2.Name, p1.Season, p2.Season, "NFL", categories), nil
}

// buildResult calculates scores and overall winner
func buildResult(name1 string, name2 string, season1 string, season2 string, sport string, categories []CompareCategory) *CompareResult {
	score1, score2 := 0, 0
	for _, c := range categories {
		if c.Winner == 1 {
			score1++
		} else if c.Winner == 2 {
			score2++
		}
	}

	winner := "Tie"
	if score1 > score2 {
		winner = name1
	} else if score2 > score1 {
		winner = name2
	}

	return &CompareResult{
		Player1:    name1,
		Player2:    name2,
		Season1:    season1,
		Season2:    season2,
		Sport:      sport,
		Categories: categories,
		Score1:     score1,
		Score2:     score2,
		Winner:     winner,
	}
}

// --- Database lookup helpers ---

func getNBAPlayerSeason(name string, season string) (*DBPlayer, error) {
	q := "%" + strings.ToLower(name) + "%"
	normalized := "%" + strings.ToLower(normalizeText(name)) + "%"

	if season == "career" {
		rows, err := db.Query(`
			SELECT games, points, assists, rebounds, steals, blocks,
			fg_pct, three_pct, ft_pct, turnovers, name, team, position
			FROM players
			WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
			AND sport = 'NBA'
			ORDER BY season ASC
		`, q, normalized)
		if err != nil {
			return nil, fmt.Errorf("player not found: %s", name)
		}
		defer rows.Close()

		var tG, tPts, tAst, tReb, tStl, tBlk, tFg, tFg3, tFt, tTov float64
		totalGames := 0
		pName, team, pos := "", "", ""

		for rows.Next() {
			var games int
			var pts, ast, reb, stl, blk, fg, fg3, ft, tov float64
			var n, t, p string
			err := rows.Scan(&games, &pts, &ast, &reb, &stl, &blk, &fg, &fg3, &ft, &tov, &n, &t, &p)
			if err != nil {
				continue
			}
			g := float64(games)
			tPts += pts * g
			tAst += ast * g
			tReb += reb * g
			tStl += stl * g
			tBlk += blk * g
			tFg += fg * g
			tFg3 += fg3 * g
			tFt += ft * g
			tTov += tov * g
			tG += g
			totalGames += games
			pName = n
			team = t
			pos = p
		}

		if totalGames == 0 {
			return nil, fmt.Errorf("player not found: %s", name)
		}

		return &DBPlayer{
			Name: pName, Team: team, Position: pos,
			Sport: "NBA", Season: "Career",
			Games:     totalGames,
			Points:    tPts / tG,
			Assists:   tAst / tG,
			Rebounds:  tReb / tG,
			Steals:    tStl / tG,
			Blocks:    tBlk / tG,
			FGPct:     tFg / tG,
			ThreePct:  tFg3 / tG,
			FTPct:     tFt / tG,
			Turnovers: tTov / tG,
		}, nil
	}

	if season == "" {
		row := db.QueryRow(`
			SELECT id, name, age, team, position, sport, season,
			games, points, assists, rebounds, steals, blocks,
			fg_pct, three_pct, ft_pct, turnovers
			FROM players
			WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
			AND sport = 'NBA'
			ORDER BY season DESC LIMIT 1
		`, q, normalized)
		r := &DBPlayer{}
		err := row.Scan(&r.ID, &r.Name, &r.Age, &r.Team, &r.Position,
			&r.Sport, &r.Season, &r.Games, &r.Points, &r.Assists,
			&r.Rebounds, &r.Steals, &r.Blocks, &r.FGPct,
			&r.ThreePct, &r.FTPct, &r.Turnovers)
		if err != nil {
			return nil, fmt.Errorf("player not found: %s", name)
		}
		return r, nil
	}

	row := db.QueryRow(`
		SELECT id, name, age, team, position, sport, season,
		games, points, assists, rebounds, steals, blocks,
		fg_pct, three_pct, ft_pct, turnovers
		FROM players
		WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
		AND sport = 'NBA' AND season = ?
		LIMIT 1
	`, q, normalized, season)
	r := &DBPlayer{}
	err := row.Scan(&r.ID, &r.Name, &r.Age, &r.Team, &r.Position,
		&r.Sport, &r.Season, &r.Games, &r.Points, &r.Assists,
		&r.Rebounds, &r.Steals, &r.Blocks, &r.FGPct,
		&r.ThreePct, &r.FTPct, &r.Turnovers)
	if err != nil {
		return nil, fmt.Errorf("player not found: %s %s", name, season)
	}
	return r, nil
}

func getNFLPasserSeason(name string, season string) (*NFLPassing, error) {
	q := "%" + strings.ToLower(name) + "%"
	normalized := "%" + strings.ToLower(normalizeText(name)) + "%"

	if season == "career" {
		rows, err := db.Query(`
			SELECT games, completions, attempts, yards, touchdowns,
			interceptions, comp_pct, yards_per_game, rating, name, team, position
			FROM nfl_passing
			WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
			ORDER BY season ASC
		`, q, normalized)
		if err != nil {
			return nil, fmt.Errorf("passer not found: %s", name)
		}
		defer rows.Close()

		var tG, tCmp, tAtt, tYds, tTd, tInt int
		var tRating float64
		seasons := 0
		pName, team, pos := "", "", ""

		for rows.Next() {
			var games, cmp, att, yds, td, interc int
			var compPct, ypg, rating float64
			var n, t, p string
			err := rows.Scan(&games, &cmp, &att, &yds, &td, &interc, &compPct, &ypg, &rating, &n, &t, &p)
			if err != nil {
				continue
			}
			tG += games
			tCmp += cmp
			tAtt += att
			tYds += yds
			tTd += td
			tInt += interc
			tRating += rating
			seasons++
			pName = n
			team = t
			pos = p
		}

		if seasons == 0 {
			return nil, fmt.Errorf("passer not found: %s", name)
		}

		compPct := 0.0
		if tAtt > 0 {
			compPct = float64(tCmp) / float64(tAtt) * 100
		}

		return &NFLPassing{
			Name: pName, Team: team, Position: pos,
			Season: "Career", Games: tG,
			Completions: tCmp, Attempts: tAtt,
			Yards: tYds, Touchdowns: tTd,
			Interceptions: tInt, CompPct: compPct,
			YardsPerGame: float64(tYds) / float64(tG),
			Rating:       tRating / float64(seasons),
		}, nil
	}

	var query string
	var args []interface{}
	if season == "" {
		query = `SELECT id, name, age, team, position, season, games,
			completions, attempts, yards, touchdowns, interceptions,
			comp_pct, yards_per_game, rating
			FROM nfl_passing
			WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
			ORDER BY season DESC LIMIT 1`
		args = []interface{}{q, normalized}
	} else {
		query = `SELECT id, name, age, team, position, season, games,
			completions, attempts, yards, touchdowns, interceptions,
			comp_pct, yards_per_game, rating
			FROM nfl_passing
			WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
			AND season = ? LIMIT 1`
		args = []interface{}{q, normalized, season}
	}

	r := db.QueryRow(query, args...)
	p := &NFLPassing{}
	err := r.Scan(&p.ID, &p.Name, &p.Age, &p.Team, &p.Position, &p.Season,
		&p.Games, &p.Completions, &p.Attempts, &p.Yards,
		&p.Touchdowns, &p.Interceptions, &p.CompPct, &p.YardsPerGame, &p.Rating)
	if err != nil {
		return nil, fmt.Errorf("passer not found: %s", name)
	}
	return p, nil
}

func getNFLRusherSeason(name string, season string) (*NFLRushing, error) {
	q := "%" + strings.ToLower(name) + "%"
	normalized := "%" + strings.ToLower(normalizeText(name)) + "%"

	if season == "career" {
		rows, err := db.Query(`
			SELECT games, attempts, yards, touchdowns, fumbles, name, team, position
			FROM nfl_rushing
			WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
			ORDER BY season ASC
		`, q, normalized)
		if err != nil {
			return nil, fmt.Errorf("rusher not found: %s", name)
		}
		defer rows.Close()

		var tG, tAtt, tYds, tTd, tFum int
		pName, team, pos := "", "", ""

		for rows.Next() {
			var games, att, yds, td, fum int
			var n, t, p string
			err := rows.Scan(&games, &att, &yds, &td, &fum, &n, &t, &p)
			if err != nil {
				continue
			}
			tG += games
			tAtt += att
			tYds += yds
			tTd += td
			tFum += fum
			pName = n
			team = t
			pos = p
		}

		if tG == 0 {
			return nil, fmt.Errorf("rusher not found: %s", name)
		}

		return &NFLRushing{
			Name: pName, Team: team, Position: pos,
			Season: "Career", Games: tG,
			Attempts: tAtt, Yards: tYds,
			Touchdowns: tTd, Fumbles: tFum,
			YardsPerCarry: float64(tYds) / float64(tAtt),
			YardsPerGame:  float64(tYds) / float64(tG),
		}, nil
	}

	var query string
	var args []interface{}
	if season == "" {
		query = `SELECT id, name, age, team, position, season, games,
			attempts, yards, touchdowns, yards_per_carry, yards_per_game, fumbles
			FROM nfl_rushing
			WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
			ORDER BY season DESC LIMIT 1`
		args = []interface{}{q, normalized}
	} else {
		query = `SELECT id, name, age, team, position, season, games,
			attempts, yards, touchdowns, yards_per_carry, yards_per_game, fumbles
			FROM nfl_rushing
			WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
			AND season = ? LIMIT 1`
		args = []interface{}{q, normalized, season}
	}

	r := db.QueryRow(query, args...)
	p := &NFLRushing{}
	err := r.Scan(&p.ID, &p.Name, &p.Age, &p.Team, &p.Position, &p.Season,
		&p.Games, &p.Attempts, &p.Yards, &p.Touchdowns,
		&p.YardsPerCarry, &p.YardsPerGame, &p.Fumbles)
	if err != nil {
		return nil, fmt.Errorf("rusher not found: %s", name)
	}
	return p, nil
}

func getNFLReceiverSeason(name string, season string) (*NFLReceiving, error) {
	q := "%" + strings.ToLower(name) + "%"
	normalized := "%" + strings.ToLower(normalizeText(name)) + "%"

	if season == "career" {
		rows, err := db.Query(`
			SELECT games, receptions, targets, yards, touchdowns, fumbles, name, team, position
			FROM nfl_receiving
			WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
			ORDER BY season ASC
		`, q, normalized)
		if err != nil {
			return nil, fmt.Errorf("receiver not found: %s", name)
		}
		defer rows.Close()

		var tG, tRec, tTgt, tYds, tTd, tFum int
		pName, team, pos := "", "", ""

		for rows.Next() {
			var games, rec, tgt, yds, td, fum int
			var n, t, p string
			err := rows.Scan(&games, &rec, &tgt, &yds, &td, &fum, &n, &t, &p)
			if err != nil {
				continue
			}
			tG += games
			tRec += rec
			tTgt += tgt
			tYds += yds
			tTd += td
			tFum += fum
			pName = n
			team = t
			pos = p
		}

		if tG == 0 {
			return nil, fmt.Errorf("receiver not found: %s", name)
		}

		return &NFLReceiving{
			Name: pName, Team: team, Position: pos,
			Season: "Career", Games: tG,
			Receptions: tRec, Targets: tTgt,
			Yards: tYds, Touchdowns: tTd, Fumbles: tFum,
			YardsPerRec:  float64(tYds) / float64(tRec),
			YardsPerGame: float64(tYds) / float64(tG),
		}, nil
	}

	var query string
	var args []interface{}
	if season == "" {
		query = `SELECT id, name, age, team, position, season, games,
			receptions, targets, yards, touchdowns, yards_per_rec, yards_per_game, fumbles
			FROM nfl_receiving
			WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
			ORDER BY season DESC LIMIT 1`
		args = []interface{}{q, normalized}
	} else {
		query = `SELECT id, name, age, team, position, season, games,
			receptions, targets, yards, touchdowns, yards_per_rec, yards_per_game, fumbles
			FROM nfl_receiving
			WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
			AND season = ? LIMIT 1`
		args = []interface{}{q, normalized, season}
	}

	r := db.QueryRow(query, args...)
	p := &NFLReceiving{}
	err := r.Scan(&p.ID, &p.Name, &p.Age, &p.Team, &p.Position, &p.Season,
		&p.Games, &p.Receptions, &p.Targets, &p.Yards,
		&p.Touchdowns, &p.YardsPerRec, &p.YardsPerGame, &p.Fumbles)
	if err != nil {
		return nil, fmt.Errorf("receiver not found: %s", name)
	}
	return p, nil
}

func getNFLDefenderSeason(name string, season string) (*NFLDefense, error) {
	q := "%" + strings.ToLower(name) + "%"
	normalized := "%" + strings.ToLower(normalizeText(name)) + "%"

	if season == "career" {
		rows, err := db.Query(`
			SELECT games, interceptions, sacks, tackles, pass_def, forced_fum, qb_hits, name, team, position
			FROM nfl_defense
			WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
			ORDER BY season ASC
		`, q, normalized)
		if err != nil {
			return nil, fmt.Errorf("defender not found: %s", name)
		}
		defer rows.Close()

		var tG, tInt, tTackles, tPD, tFF, tQBH int
		var tSacks float64
		pName, team, pos := "", "", ""

		for rows.Next() {
			var games, interc, tackles, pd, ff, qbh int
			var sacks float64
			var n, t, p string
			err := rows.Scan(&games, &interc, &sacks, &tackles, &pd, &ff, &qbh, &n, &t, &p)
			if err != nil {
				continue
			}
			tG += games
			tInt += interc
			tSacks += sacks
			tTackles += tackles
			tPD += pd
			tFF += ff
			tQBH += qbh
			pName = n
			team = t
			pos = p
		}

		if tG == 0 {
			return nil, fmt.Errorf("defender not found: %s", name)
		}

		return &NFLDefense{
			Name: pName, Team: team, Position: pos,
			Season: "Career", Games: tG,
			Interceptions: tInt, Sacks: tSacks,
			Tackles: tTackles, PassDef: tPD,
			ForcedFum: tFF, QBHits: tQBH,
		}, nil
	}

	var query string
	var args []interface{}
	if season == "" {
		query = `SELECT id, name, age, team, position, season, games,
			interceptions, sacks, tackles, pass_def, forced_fum, qb_hits
			FROM nfl_defense
			WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
			ORDER BY season DESC LIMIT 1`
		args = []interface{}{q, normalized}
	} else {
		query = `SELECT id, name, age, team, position, season, games,
			interceptions, sacks, tackles, pass_def, forced_fum, qb_hits
			FROM nfl_defense
			WHERE (LOWER(name) LIKE ? OR LOWER(search_name) LIKE ?)
			AND season = ? LIMIT 1`
		args = []interface{}{q, normalized, season}
	}

	r := db.QueryRow(query, args...)
	p := &NFLDefense{}
	err := r.Scan(&p.ID, &p.Name, &p.Age, &p.Team, &p.Position, &p.Season,
		&p.Games, &p.Interceptions, &p.Sacks, &p.Tackles,
		&p.PassDef, &p.ForcedFum, &p.QBHits)
	if err != nil {
		return nil, fmt.Errorf("defender not found: %s", name)
	}
	return p, nil
}
