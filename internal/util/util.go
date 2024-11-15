package util

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	kstLocation *time.Location
)

func init() {
	var err error
	kstLocation, err = time.LoadLocation("Asia/Seoul")
	if err != nil {
		fmt.Println("Failed to load timezone:", err)
		kstLocation = time.FixedZone("KST", 9*60*60)
	}
}

func GetTodayInKST() time.Time {
	now := time.Now().In(kstLocation)

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, kstLocation)
	return today
}

func ParseDate(dateStr string) (time.Time, error) {
	parsedTime, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(dateStr), kstLocation)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

// LoadFile 파일을 읽어와 바이트 배열로 반환
func LoadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Failed to close file:", err)
		}
	}(file)
	arr, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return arr, nil
}

// PickRandomNumber 0부터 upperBound-1 사이의 랜덤한 숫자를 반환 (upperBound는 포함하지 않음) upperBound가 0 이하일 경우 0을 반환
func PickRandomNumber(upperBound int) int {
	if upperBound <= 0 {
		return 0
	}
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	return r.Intn(upperBound - 1)
}
