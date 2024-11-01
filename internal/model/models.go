package model

import (
	"time"
)

type FeelScore struct {
	Score  int
	Reason string
}

type Scrum struct {
	UserId     int64
	Goal       string
	Commitment string
	Feels      FeelScore
	CreatedAt  time.Time
}

type Retrospection struct {
	UserId       int64
	GoalAchieved string
	Learned      string
	Feels        FeelScore
	CreatedAt    time.Time
}

type User struct {
	UserId        int64
	Name          string
	DiscordUserId string
}
