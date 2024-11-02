package service

import (
	"database/sql"
	"day-day-review/internal/model"
	"day-day-review/internal/repository"
	"day-day-review/internal/util"
	"fmt"
	"log"

	"github.com/mattn/go-sqlite3"
)

// RegisterUser 함수는 주어진 nickname과 userId로 새로운 사용자 레코드를 데이터베이스에 추가합니다.
func RegisterUser(db *sql.DB, nickname string, userId string) string {
	err := repository.InsertUser(db, model.User{Name: nickname, DiscordUserId: userId})
	if err != nil {
		log.Println("Failed to insert user: ", err)

		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.Code == sqlite3.ErrConstraint {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintPrimaryKey:
				log.Printf("discord_user_id already exists: %v - %s", err, userId)
				return "이미 등록한 사용자입니다."
			case sqlite3.ErrConstraintUnique:
				log.Printf("name already exists: %v - %s", err, nickname)
				return "이미 등록된 이름입니다. 다른 이름을 입력해주세요."
			}
		}
		return "에러가 발생했습니다."
	}

	return fmt.Sprintf("닉네임 '%s' 등록 완료!", nickname)
}

// 주어진 내용들로 사용자의 오늘의 다짐 레코드를 데이터베이스에 추가합니다.
func CreateTodayScrumByUserId(db *sql.DB, userId, goal, commitment, feelReason string, feelScore int) string {
	scrum := model.Scrum{
		UserId:     userId,
		Goal:       goal,
		Commitment: commitment,
		Feels: model.FeelScore{
			Score:  feelScore,
			Reason: feelReason,
		},
		CreatedAt: util.GetTodayInKST(),
	}

	err := repository.InsertScrum(db, scrum)
	if err != nil {
		log.Println("Error inserting scrum data:", err)
		return "에러가 발생했습니다."
	}

	return "오늘의 다짐을 작성했습니다!"
}
