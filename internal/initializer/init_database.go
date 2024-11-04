package initializer

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase(filePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatalf("Database Connection fail: %v", err)
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS user (
		discord_user_id CHAR(30) PRIMARY KEY,
	    name CHAR(20) NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS scrum (
		user_id INTEGER NOT NULL,
		goal TEXT NOT NULL,
		commitment TEXT,
		feel_score INTEGER,
		feel_reason TEXT,
		created_at TIMESTAMP DEFAULT (datetime('now', '+09:00')),
		
		PRIMARY KEY (user_id, created_at)
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Println("Table Creation fail: ", err)
		return nil, err
	}

	log.Println("Database Initialize Completed")
	return db, nil
}
