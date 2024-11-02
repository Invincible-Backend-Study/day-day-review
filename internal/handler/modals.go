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
