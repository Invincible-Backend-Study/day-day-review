package service

import (
	"day-day-review/internal/model"
	"day-day-review/internal/repository"
	"day-day-review/internal/util"
	"errors"
	"fmt"
	"log"
	"time"

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

func ExistUser(userId string) bool {
	exist, err := repository.ExistUserByUserId(userId)
	if err != nil {
		log.Println("Failed to select user: ", err)
		return false
	}
	return exist
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

	_, err := repository.InsertScrum(&scrum)
	if err != nil {
		log.Println("Error inserting scrum data:", err)
		return "에러가 발생했습니다."
	}

	return "오늘의 다짐을 작성했습니다!"
}

// CreateTodayRetrospectives 주어진 내용들로 사용자의 오늘의 다짐 레코드를 데이터베이스에 추가합니다.
func CreateTodayRetrospectives(userId, goalAchieved, learned, feelReason string, feelScore int) string {
	retrospective := model.Retrospective{
		UserId:       userId,
		GoalAchieved: goalAchieved,
		Learned:      learned,
		Feels: model.FeelScore{
			Score:  feelScore,
			Reason: feelReason,
		},
		CreatedAt: util.GetTodayInKST(),
	}

	_, err := repository.InsertRetrospective(&retrospective)
	if err != nil {
		log.Println("Error inserting Retrospectives data:", err)
		return "에러가 발생했습니다."
	}

	return "오늘의 회고를 작성했습니다!"
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

// ExistTodayRetrospective 주어진 사용자가 오늘의 회고를 작성했는지 여부를 반환합니다.
func ExistTodayRetrospective(userId string) (bool, error) {
	today := util.GetTodayInKST()
	result, err := repository.ExistRetrospectiveByUserId(userId, today)
	if err != nil {
		log.Printf("Error select retrospective data: %v - %s", err, userId)
		return false, fmt.Errorf("failed to select retrospective data: %w", err)
	}
	return result, nil
}

// GetTodayScrums 작성된 오늘의 다짐을 모두 반환합니다.
func GetTodayScrums() ([]*model.ScrumDto, error) {
	scrums, err := repository.SelectScrumListByDate(util.GetTodayInKST())
	if err != nil {
		return nil, err
	}
	return scrums, nil
}

// GetTodayRetrospectives 작성된 오늘의 회고를 모두 반환합니다.
func GetTodayRetrospectives() ([]*model.RetrospectiveDto, error) {
	scrums, err := repository.SelectRetrospectiveListByDate(util.GetTodayInKST())
	if err != nil {
		return nil, err
	}
	return scrums, nil
}

// GetScrumsByDate 작성된 특정 날짜의 다짐을 모두 반환합니다.
func GetScrumsByDate(date time.Time) ([]*model.ScrumDto, error) {
	scrums, err := repository.SelectScrumListByDate(date)
	if err != nil {
		return nil, err
	}
	return scrums, nil
}

// GetRetrospectivesByDate 작성된 특정 날짜의 회고를 모두 반환합니다.
func GetRetrospectivesByDate(date time.Time) ([]*model.RetrospectiveDto, error) {
	scrums, err := repository.SelectRetrospectiveListByDate(date)
	if err != nil {
		return nil, err
	}
	return scrums, nil
}
