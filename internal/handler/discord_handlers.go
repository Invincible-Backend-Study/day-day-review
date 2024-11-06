package handler

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

// RegisterCommands 봇에 명령어를 등록합니다. 명령어는 commands.go에 정의되어 있습니다.
func RegisterCommands(s *discordgo.Session, _ *discordgo.Ready) {
	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guildId, cmd)
		if err != nil {
			log.Printf("Cannot create command: %v\n", err)
		}
	}
}

// RegisterInteractions 봇의 상호작용을 처리합니다.
func RegisterInteractions(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.Type {
	case discordgo.InteractionApplicationCommand:
		handleApplicationCommand(session, interaction)
	case discordgo.InteractionModalSubmit:
		handleModalSubmit(session, interaction)
	}
}

func EasterEggHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			log.Println("Error sending message: ", err)
		}
	}
	if m.Content == "pong" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Ping!")
		if err != nil {
			log.Println("Error sending message: ", err)
		}
	}
}
