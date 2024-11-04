package handler

import "github.com/bwmarrin/discordgo"

// 명령어 상수
const (
	commandRegisterUser       = "회원-등록"
	commandRegisterTodayScrum = "오늘-다짐"
	commandGetTodayScrums     = "오늘-다짐-보기"
)

const (
	cIdRegisterUserModal         = "nickname_modal"
	cIdRegisterUserNicknameInput = "nickname_input"

	cIdRegisterScrumModal           = "scrum_modal"
	cIdRegisterScrumGoalInput       = "scrum_goal_input"
	cIdRegisterScrumCommitmentInput = "scrum_commitment_input"
	cIdRegisterScrumScoreInput      = "scrum_score_input"
	cIdRegisterScrumReasonInput     = "scrum_reason_input"
)

// 명령어 목록
var commands = []*discordgo.ApplicationCommand{
	{
		Name:        commandRegisterUser,
		Description: "닉네임을 등록합니다",
	},
	{
		Name:        commandRegisterTodayScrum,
		Description: "오늘의 다짐을 등록합니다",
	},
	{
		Name:        commandGetTodayScrums,
		Description: "오늘의 다짐을 모두 보여줍니다",
	},
}
