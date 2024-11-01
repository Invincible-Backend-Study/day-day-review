package intializer

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase() (*sql.DB, error) {
	dbFile := "../configs/dayday.db"
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Database Connection fail: %v", err)
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS user (
		user_id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		discord_user_id TEXT NOT NULL UNIQUE
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Table Creation fail: %v", err)
		return nil, err
	}

	log.Println("Database Initialize Completed")
	return db, nil
}
