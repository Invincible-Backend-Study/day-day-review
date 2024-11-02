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
		discord_user_id CHAR(30) PRIMARY KEY,
	    name CHAR(20) NOT NULL
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalln("Table Creation fail: ", err)
		return nil, err
	}

	log.Println("Database Initialize Completed")
	return db, nil
}
