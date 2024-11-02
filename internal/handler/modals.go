package handler

import "github.com/bwmarrin/discordgo"

func createRegisterUserModal() *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:    "닉네임 입력",
			CustomID: cIdRegisterUserModal,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    cIdRegisterUserNicknameInput,
							Label:       "닉네임을 입력하세요",
							Style:       discordgo.TextInputShort,
							Placeholder: "닉네임",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func createRegisterScrumModal() *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:    "Scrum 기록 입력",
			CustomID: cIdRegisterScrumModal,
			Components: []discordgo.MessageComponent{
				// Goal 입력 필드
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    cIdRegisterScrumGoalInput,
							Label:       "오늘의 목표를 입력하세요",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "오늘의 목표",
							Required:    true,
						},
					},
				},
				// Commitment 입력 필드
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    cIdRegisterScrumCommitmentInput,
							Label:       "오늘의 다짐을 입력하세요",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "오늘의 다짐",
							Required:    false,
						},
					},
				},
				// Score 입력 필드
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    cIdRegisterScrumScoreInput,
							Label:       "오늘의 기분 점수를 입력하세요 (숫자 0 ~ 10)",
							Style:       discordgo.TextInputShort,
							Placeholder: "기분 점수 (예: 7)",
							Required:    true,
						},
					},
				},
				// Reason 입력 필드
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    cIdRegisterScrumReasonInput,
							Label:       "기분 점수의 이유를 입력하세요",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "점수의 이유",
							Required:    false,
						},
					},
				},
			},
		},
	}
}
