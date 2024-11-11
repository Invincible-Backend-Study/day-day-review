package repository

import (
	"database/sql"
	"day-day-review/internal/model"
	"fmt"
	"log"
	"sync"
	"time"
)

type SQLiteFileRepository struct {
	db *sql.DB
}

var (
	instance *SQLiteFileRepository
	once     sync.Once
)

// NewSQLiteFileRepository 파일 경로를 받아 SQLiteFileRepository 싱글톤 인스턴스를 반환합니다.
func NewSQLiteFileRepository(filePath string) *SQLiteFileRepository {
	once.Do(func() {
		// 데이터베이스 연결 생성
		var db *sql.DB
		db, err := sql.Open("sqlite3", filePath)
		if err != nil {
			log.Fatalf("database connection failed: %v", err)
			return
		}

		// 테이블 생성 등 데이터베이스 초기화
		instance = &SQLiteFileRepository{db: db}
		if _, err := db.Exec(createTableQuery); err != nil {
			log.Fatalf("failed to create tables: %v", err)
			return
		}

		log.Println("Database initialization completed")
	})

	return instance
}

func (r *SQLiteFileRepository) InsertUser(user model.User) error {
	stmt, err := r.db.Prepare(insertUserQuery)
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

func (r *SQLiteFileRepository) ExistUserByUserId(userId string) (bool, error) {
	stmt, err := r.db.Prepare(existsUserQuery)
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("failed to close statement:", err)
		}
	}(stmt)

	var result int
	err = stmt.QueryRow(userId).Scan(&result)
	if err != nil {
		return false, fmt.Errorf("failed to execute statement: %v", err)
	}
	return result == 1, nil
}

func (r *SQLiteFileRepository) InsertScrum(scrum *model.Scrum) (*model.Scrum, error) {
	// Prepared statement 생성
	stmt, err := r.db.Prepare(insertScrumQuery)
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
	return scrum, nil
}

func (r *SQLiteFileRepository) InsertRetrospective(retrospective *model.Retrospective) (*model.Retrospective, error) {
	// Prepared statement 생성
	stmt, err := r.db.Prepare(insertRetrospectiveQuery)
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
	_, err = stmt.Exec(retrospective.UserId, retrospective.GoalAchieved, retrospective.Learned, retrospective.Feels.Score, retrospective.Feels.Reason, retrospective.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}
	return retrospective, nil
}

func (r *SQLiteFileRepository) ExistScrumByUserId(userId string, today time.Time) (bool, error) {
	stmt, err := r.db.Prepare(existScrumQuery)
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
	return count > 0, nil
}

func (r *SQLiteFileRepository) ExistRetrospectiveByUserId(userId string, today time.Time) (bool, error) {
	stmt, err := r.db.Prepare(existRetrospectiveQuery)
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
	return count > 0, nil
}

func (r *SQLiteFileRepository) SelectScrumListByDate(date time.Time) ([]*model.ScrumDto, error) {
	stmt, err := r.db.Prepare(selectTodayScrumQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("failed to close statement:", err)
		}
	}(stmt)

	rows, err := stmt.Query(date)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("failed to close rows:", err)
		}
	}(rows)

	var scrums []*model.ScrumDto
	for rows.Next() {
		var scrum model.ScrumDto
		err := rows.Scan(&scrum.Name, &scrum.Goal, &scrum.Commitment, &scrum.FeelScore, &scrum.FeelReason)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		scrums = append(scrums, &scrum)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return scrums, nil
}

func (r *SQLiteFileRepository) SelectRetrospectiveListByDate(date time.Time) ([]*model.RetrospectiveDto, error) {
	stmt, err := r.db.Prepare(selectTodayRetrospectiveQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("failed to close statement:", err)
		}
	}(stmt)

	rows, err := stmt.Query(date)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("failed to close rows:", err)
		}
	}(rows)

	var scrums []*model.RetrospectiveDto
	for rows.Next() {
		var scrum model.RetrospectiveDto
		err := rows.Scan(&scrum.Name, &scrum.GoalAchieved, &scrum.Learned, &scrum.FeelScore, &scrum.FeelReason)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		scrums = append(scrums, &scrum)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return scrums, nil
}
