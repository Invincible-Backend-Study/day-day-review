package util

import (
	"fmt"
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