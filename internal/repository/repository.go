package repository

import (
	"day-day-review/internal/model"
	"time"
)

// Repository는 데이터베이스와 상호작용하는 기능을 추상화한 인터페이스입니다.
type Repository interface {
	// InsertUser는 새로운 사용자 정보를 데이터베이스에 삽입합니다.
	// user: 삽입할 사용자 정보
	// 반환값: 에러가 발생할 경우 해당 에러 반환
	InsertUser(user model.User) error

	// ExistUserByUserId는 주어진 사용자 ID로 사용자가 존재하는지 확인합니다.
	// userId: 확인할 사용자 ID
	// 반환값: 사용자가 존재할 경우 true와 에러 정보 반환
	ExistUserByUserId(userId string) (bool, error)

	// InsertScrum은 새로운 스크럼 정보를 데이터베이스에 삽입합니다.
	// scrum: 삽입할 스크럼 정보
	// 반환값: 삽입된 스크럼 정보와 에러가 발생할 경우 해당 에러 반환
	InsertScrum(scrum *model.Scrum) (*model.Scrum, error)

	// ExistScrumByUserId는 주어진 사용자 ID와 날짜로 스크럼이 존재하는지 확인합니다.
	// userId: 확인할 사용자 ID
	// today: 확인할 날짜
	// 반환값: 스크럼이 존재할 경우 true와 에러 정보 반환
	ExistScrumByUserId(userId string, date time.Time) (bool, error)

	// SelectScrumListByDate는 주어진 날짜에 해당하는 스크럼 목록을 조회합니다.
	// date: 조회할 날짜
	// 반환값: 해당 날짜의 스크럼 목록과 에러가 발생할 경우 해당 에러 반환
	SelectScrumListByDate(date time.Time) ([]*model.ScrumDto, error)

	// InsertRetrospective는 새로운 회고 정보를 데이터베이스에 삽입합니다.
	// retrospective: 삽입할 회고 정보
	// 반환값: 삽입된 회고 정보와 에러가 발생할 경우 해당 에러 반환
	InsertRetrospective(retrospective *model.Retrospective) (*model.Retrospective, error)

	// ExistRetrospectiveByUserId는 주어진 사용자 ID와 날짜로 회고가 존재하는지 확인합니다.
	// userId: 확인할 사용자 ID
	// today: 확인할 날짜
	// 반환값: 회고가 존재할 경우 true와 에러 정보 반환
	ExistRetrospectiveByUserId(userId string, date time.Time) (bool, error)

	// SelectRetrospectiveListByDate는 주어진 날짜에 해당하는 회고 목록을 조회합니다.
	// date: 조회할 날짜
	// 반환값: 해당 날짜의 회고 목록과 에러가 발생할 경우 해당 에러 반환
	SelectRetrospectiveListByDate(date time.Time) ([]*model.RetrospectiveDto, error)
}
