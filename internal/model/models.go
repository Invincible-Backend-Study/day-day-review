package model

import (
	"time"
)

type FeelScore struct {
	Score  int
	Reason string
}

type Scrum struct {
	UserId     string
	Goal       string
	Commitment string
	Feels      FeelScore
	CreatedAt  time.Time
}

type Retrospective struct {
	UserId       string
	GoalAchieved string
	Learned      string
	Feels        FeelScore
	CreatedAt    time.Time
}

type User struct {
	DiscordUserId string
	Name          string
}
