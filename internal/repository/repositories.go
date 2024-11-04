package repository

import (
	"database/sql"
	"day-day-review/internal/model"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	database *sql.DB
	once     sync.Once
)

func Initialize(filePath string) {
	once.Do(func() {
		initRepository(filePath)
	})
}

func initRepository(filePath string) {
	var err error
	database, err = sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatalf("database Connection fail: %v", err)
	}
	_, err = database.Exec(createTableQuery)
	if err != nil {
		log.Println("Table Creation fail: ", err)
	}
	log.Println("Database Initialize Completed")
}

func InsertUser(user model.User) error {
	stmt, err := database.Prepare(insertUserQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("failed to close statement:", err)
		}
	}(stmt)

	_, err = stmt.Exec(user.Name, user.DiscordUserId)
	if err != nil {
		return err
	}
	return nil
}

func InsertScrum(scrum model.Scrum) (*model.Scrum, error) {
	// Prepared statement 생성
	stmt, err := database.Prepare(insertScrumQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("failed to close statement:", err)
		}
	}(stmt)
	// Prepared statement 실행
	_, err = stmt.Exec(scrum.UserId, scrum.Goal, scrum.Commitment, scrum.Feels.Score, scrum.Feels.Reason, scrum.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}
	return &scrum, nil
}

func ExistScrumByUserId(userId string, today time.Time) (bool, error) {
	stmt, err := database.Prepare(existScrumQuery)
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("failed to close statement:", err)
		}
	}(stmt)
	var count int
	err = stmt.QueryRow(userId, today).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to execute statement: %w", err)
	}
	log.Printf("count: %d", count)
	return count > 0, nil
}
