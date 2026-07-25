package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// initDB opens or creates the SQLite database
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./statvault.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS players (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			age INTEGER,
			team TEXT,
			position TEXT,
			sport TEXT DEFAULT 'NBA',
			season TEXT,
			games INTEGER,
			games_started INTEGER,
			minutes_per_game REAL,
			points REAL,
			assists REAL,
			rebounds REAL,
			steals REAL,
			blocks REAL,
			fg_pct REAL,
			three_pct REAL,
			ft_pct REAL,
			turnovers REAL,
			UNIQUE(name, season, sport)
		)
	`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	fmt.Println("Database initialized successfully")
}