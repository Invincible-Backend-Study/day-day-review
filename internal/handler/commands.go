package handler

import "github.com/bwmarrin/discordgo"

// 명령어 상수
const (
	commandRegisterUser               = "회원-등록"
	commandRegisterTodayScrum         = "오늘-다짐"
	commandGetTodayScrums             = "오늘-다짐-보기"
	commandRegisterTodayRetrospection = "오늘-회고"
	commandGetTodayRetrospectives     = "오늘-회고-보기"
	commandGetScrumByDate             = "다짐-보기"
	commandGetRetrospectionByDate     = "회고-보기"
)

const (
	cIdRegisterUserModal         = "nickname_modal"
	cIdRegisterUserNicknameInput = "nickname_input"

	cIdRegisterScrumModal           = "scrum_modal"
	cIdRegisterScrumGoalInput       = "scrum_goal_input"
	cIdRegisterScrumCommitmentInput = "scrum_commitment_input"
	cIdRegisterScrumScoreInput      = "scrum_score_input"
	cIdRegisterScrumReasonInput     = "scrum_reason_input"

	cIdRegisterRetrospectionModal             = "retrospection_modal"
	cIdRegisterRetrospectionGoalAchievedInput = "retrospection_goal_achieved_input"
	cIdRegisterRetrospectionLearnedInput      = "retrospection_learned_input"
	cIdRegisterRetrospectionScoreInput        = "retrospection_score_input"
	cIdRegisterRetrospectionReasonInput       = "retrospection_reason_input"
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
		Name:        commandRegisterTodayRetrospection,
		Description: "오늘의 회고를 등록합니다.",
	},
	{
		Name:        commandGetTodayScrums,
		Description: "오늘의 다짐을 모두 보여줍니다",
	},
	{
		Name:        commandGetTodayRetrospectives,
		Description: "오늘의 회고를 모두 보여줍니다.",
	},
	{
		Name:        commandGetScrumByDate,
		Description: "특정 날짜의 다짐을 보여줍니다.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date",
				Description: "보고 싶은 날짜를 입력해주세요. (YYYY-MM-DD)",
				Required:    true,
			},
		},
	},
	{
		Name:        commandGetRetrospectionByDate,
		Description: "특정 날짜의 회고를 보여줍니다.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date",
				Description: "보고 싶은 날짜를 입력해주세요. (YYYY-MM-DD)",
				Required:    true,
			},
		},
	},
}
