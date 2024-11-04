package repository

import (
	"database/sql"
	"day-day-review/internal/model"
	"fmt"
	"log"
)

const (
	insertUserQuery  = `INSERT INTO User (name, discord_user_id) VALUES (?, ?)`
	insertScrumQuery = `INSERT INTO Scrum (user_id, goal, commitment, feel_score, feel_reason, created_at) VALUES (?, ?, ?, ?, ?, ?)`
)

func InsertUser(db *sql.DB, user model.User) error {
	stmt, err := db.Prepare(insertUserQuery)
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

func InsertScrum(db *sql.DB, scrum model.Scrum) (*model.Scrum, error) {
	// Prepared statement 생성
	stmt, err := db.Prepare(insertScrumQuery)
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
