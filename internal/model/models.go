package model

import (
	"github.com/google/uuid"
	"time"
)

type FeelScore struct {
	Score  int
	Reason string
}

type Scrum struct {
	UserId     uuid.UUID
	Goal       string
	Commitment string
	Feels      FeelScore
	CreatedAt  time.Time
}

type Retrospection struct {
	UserId       uuid.UUID
	GoalAchieved string
	Learned      string
	Feels        FeelScore
	CreatedAt    time.Time
}

type User struct {
	UserId     uuid.UUID
	Name       string
	DiscordKey string
}
