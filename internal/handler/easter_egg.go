package handler

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

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
