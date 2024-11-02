package repository

import (
	"database/sql"
	"day-day-review/internal/model"
	"fmt"
	"log"
)

func InsertUser(db *sql.DB, user model.User) error {
	query := `INSERT INTO User (name, discord_user_id) VALUES (?, ?)`
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("failed to close statement:", err)
		}
	}(stmt)

	result, err := stmt.Exec(user.Name, user.DiscordUserId)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert id: %v", err)
	}

	user.UserId = id
	fmt.Printf("Inserted user with ID %d\n", user.UserId)
	return nil
}
