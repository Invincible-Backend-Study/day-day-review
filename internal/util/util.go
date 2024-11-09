package util

import (
	"fmt"
	"io"
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
