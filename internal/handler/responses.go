package handler

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

// sendMessage 메시지를 전송합니다.
func sendMessage(session *discordgo.Session, interaction *discordgo.InteractionCreate, content string) {
	if err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	}); err != nil {
		log.Println("Error responding with message: ", err)
	}
}

// sendEphemeralMessage 개인 메시지를 전송합니다.
func sendEphemeralMessage(session *discordgo.Session, interaction *discordgo.InteractionCreate, content string) {
	if err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		log.Println("Error responding with ephemeral message: ", err)
	}
}

// logErrorAndSendMessage 에러를 로그에 기록하고 메시지를 전송합니다.
func logErrorAndSendMessage(session *discordgo.Session, interaction *discordgo.InteractionCreate, message string, err error) {
	log.Printf("%s: %v", message, err)
	sendEphemeralMessage(session, interaction, message)
}
