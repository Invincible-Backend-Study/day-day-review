package model

type ScrumDto struct {
	Name       string
	Goal       string
	Commitment string
	FeelScore  int
	FeelReason string
}

type RetrospectiveDto struct {
	Name         string
	GoalAchieved string
	Learned      string
	FeelScore    int
	FeelReason   string
}
