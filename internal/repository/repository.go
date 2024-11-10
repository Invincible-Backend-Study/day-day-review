package repository

import (
	"day-day-review/internal/model"
	"time"
)

// Repository 데이터베이스 상호작용을 추상화한 인터페이스입니다.
type Repository interface {
	// User
	InsertUser(user model.User) error
	ExistUserByUserId(userId string) (bool, error)

	// Scrum
	InsertScrum(scrum *model.Scrum) (*model.Scrum, error)
	ExistScrumByUserId(userId string, today time.Time) (bool, error)
	SelectScrumListByDate(date time.Time) ([]*model.ScrumDto, error)

	// Retrospective
	InsertRetrospective(retrospective *model.Retrospective) (*model.Retrospective, error)
	ExistRetrospectiveByUserId(userId string, today time.Time) (bool, error)
	SelectRetrospectiveListByDate(date time.Time) ([]*model.RetrospectiveDto, error)
}
