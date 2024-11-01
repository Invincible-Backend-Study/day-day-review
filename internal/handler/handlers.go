package handler

import (
	"database/sql"
	"day-day-review/internal/model"
	"day-day-review/internal/repository"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

type Manager struct {
	db      *sql.DB
	guildId string
}

func NewHandlerManager(db *sql.DB, guildId string) *Manager {
	return &Manager{
		db: db, guildId: guildId,
	}
}

func (manager *Manager) RegisterCommands(s *discordgo.Session, _ *discordgo.Ready) {
	_, err := s.ApplicationCommandCreate(s.State.User.ID, manager.guildId, &discordgo.ApplicationCommand{
		Name:        "회원-등록",
		Description: "닉네임을 등록합니다",
	})
	if err != nil {
		log.Fatalf("Cannot create command: %v", err)
	}
}

func (manager *Manager) RegisterHandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == "회원-등록" {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				Title:    "닉네임 입력",
				CustomID: "nickname_modal",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "nickname_input",
								Label:       "닉네임을 입력하세요",
								Style:       discordgo.TextInputShort,
								Placeholder: "닉네임",
								Required:    true,
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Printf("Error responding with modal: %v", err)
		}
	} else if i.Type == discordgo.InteractionModalSubmit && i.ModalSubmitData().CustomID == "nickname_modal" {
		// 닉네임 입력 후 실행할 작업을 여기에 추가
		nickname := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		userID := i.Member.User.ID

		err := repository.InsertUser(manager.db, model.User{Name: nickname, DiscordUserId: userID})
		if err != nil {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "이미 등록된 이름입니다. 다른 이름을 입력해주세요.",
					Flags:   discordgo.MessageFlagsEphemeral, // 사용자에게만 보이는 메시지
				},
			})
			if err != nil {
				log.Printf("응답 실패: %v", err)
			}
			return
		}
		response := fmt.Sprintf("닉네임 '%s' 등록 완료!", nickname)

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
				Flags:   discordgo.MessageFlagsEphemeral, // 사용자에게만 보이는 메시지
			},
		})
		if err != nil {
			log.Println("Error responding to modal submit:", err)
		}

		log.Println("Received nickname:", nickname)
	}
}
