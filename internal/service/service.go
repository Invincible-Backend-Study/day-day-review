package service

import (
	"day-day-review/internal/model"
	"day-day-review/internal/repository"
	"day-day-review/internal/util"
	"errors"
	"fmt"
	"log"

	"github.com/mattn/go-sqlite3"
)

func init() {
	repository.Initialize("configs/dayday.db")
}

// AddUser 함수는 주어진 nickname과 userId로 새로운 사용자 레코드를 데이터베이스에 추가합니다.
func AddUser(nickname string, userId string) string {
	err := repository.InsertUser(model.User{Name: nickname, DiscordUserId: userId})
	if err != nil {
		log.Println("Failed to insert user: ", err)

		var sqliteErr sqlite3.Error
		if ok := errors.As(err, &sqliteErr); ok && errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			switch {
			case errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintPrimaryKey):
				log.Printf("discord_user_id already exists: %v - %s", err, userId)
				return "이미 등록한 사용자입니다."
			case errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique):
				log.Printf("name already exists: %v - %s", err, nickname)
				return "이미 등록된 이름입니다. 다른 이름을 입력해주세요."
			}
		}
		return "에러가 발생했습니다."
	}

	return fmt.Sprintf("닉네임 '%s' 등록 완료!", nickname)
}

// CreateTodayScrum 주어진 내용들로 사용자의 오늘의 다짐 레코드를 데이터베이스에 추가합니다.
func CreateTodayScrum(userId, goal, commitment, feelReason string, feelScore int) string {
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

	_, err := repository.InsertScrum(scrum)
	if err != nil {
		log.Println("Error inserting scrum data:", err)
		return "에러가 발생했습니다."
	}

	return "오늘의 다짐을 작성했습니다!"
}

// ExistTodayScrum 주어진 사용자가 오늘의 다짐을 작성했는지 여부를 반환합니다.
func ExistTodayScrum(userId string) (bool, error) {
	today := util.GetTodayInKST()
	result, err := repository.ExistScrumByUserId(userId, today)
	if err != nil {
		log.Printf("Error select scrum data: %v - %s", err, userId)
		return false, fmt.Errorf("failed to select scrum data: %w", err)
	}
	return result, nil
}

// GetTodayScrums 작성된 오늘의 다짐을 모두 반환합니다.
func GetTodayScrums() ([]model.ScrumDto, error) {
	scrums, err := repository.SelectTodayScrumList(util.GetTodayInKST())
	if err != nil {
		return nil, err
	}
	return scrums, nil
}
