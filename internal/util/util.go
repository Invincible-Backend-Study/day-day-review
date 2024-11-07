package util

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	kstLocation *time.Location
	once        sync.Once
)

func initKSTLocation() {
	var err error
	kstLocation, err = time.LoadLocation("Asia/Seoul")
	if err != nil {
		fmt.Println("Failed to load timezone:", err)
		kstLocation = time.FixedZone("KST", 9*60*60)
	}
}

func GetTodayInKST() time.Time {
	once.Do(initKSTLocation)

	now := time.Now().In(kstLocation)

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, kstLocation)
	return today
}

func ParseDate(dateStr string) (time.Time, error) {
	// 날짜 문자열을 KST 타임존으로 파싱
	parsedTime, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(dateStr), kstLocation)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}
