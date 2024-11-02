package handler

import (
	"database/sql"
	"day-day-review/internal/model"
	"day-day-review/internal/repository"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/mattn/go-sqlite3"
)

// Manager 봇의 핸들러를 관리합니다.
type Manager struct {
	db      *sql.DB
	guildId string
}

// NewHandlerManager 새로운 핸들러 매니저를 생성합니다.
func NewHandlerManager(db *sql.DB, guildId string) *Manager {
	return &Manager{
		db: db, guildId: guildId,
	}
}

// RegisterCommands 봇에 명령어를 등록합니다. 명령어는 commands.go에 정의되어 있습니다.
func (manager *Manager) RegisterCommands(s *discordgo.Session, _ *discordgo.Ready) {
	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, manager.guildId, cmd)
		if err != nil {
			log.Printf("Cannot create command: %v\n", err)
		}
	}
}

// RegisterInteractions 봇의 상호작용을 처리합니다.
func (manager *Manager) RegisterInteractions(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.Type {
	case discordgo.InteractionApplicationCommand:
		manager.handleApplicationCommand(session, interaction)
	case discordgo.InteractionModalSubmit:
		manager.handleModalSubmit(session, interaction)
	}
}

// handleApplicationCommand 봇의 명령어를 처리합니다.
func (manager *Manager) handleApplicationCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.ApplicationCommandData().Name {
	case commandRegisterUser:
		err := session.InteractionRespond(interaction.Interaction, createRegisterUserModal())
		if err != nil {
			log.Printf("Error responding with modal: %v", err)
		}
	}
}

// handleModalSubmit 모달의 제출을 처리합니다.
func (manager *Manager) handleModalSubmit(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.ModalSubmitData().CustomID {
	case cIdRegisterUserModal:
		manager.interactionRegisterModal(session, interaction)
	}
}

// interactionRegisterModal 사용자 등록 모달의 상호작용을 처리합니다.
func (manager *Manager) interactionRegisterModal(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	nickname, err := extractValueFromComponent(interaction.ModalSubmitData().Components, cIdRegisterUserNicknameInput)
	if err != nil {
		log.Println("Error extracting value from component: ", err)
		sendEphemeralMessage(session, interaction, "닉네임을 입력해주세요.")
		return
	}
	userID := interaction.Member.User.ID

	err = repository.InsertUser(manager.db, model.User{Name: nickname, DiscordUserId: userID})
	if err != nil {
		log.Println("Failed to insert user: ", err)
		sqliteErr, ok := err.(sqlite3.Error)
		log.Println(sqliteErr)
		if ok && sqliteErr.Code == sqlite3.ErrConstraint {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				sendEphemeralMessage(session, interaction, "이미 등록한 사용자입니다.")
				log.Println("discord_user_id already exists: %w", err)
			}
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				sendEphemeralMessage(session, interaction, "이미 등록된 이름입니다. 다른 이름을 입력해주세요.")
				log.Println("name already exists: %w", err)
			}
		}
		return
	}
	response := fmt.Sprintf("닉네임 '%s' 등록 완료!", nickname)

	sendEphemeralMessage(session, interaction, response)
	log.Println("Received nickname:", nickname)
}

// extractValueFromComponent 컴포넌트에서 값을 추출합니다.
func extractValueFromComponent(components []discordgo.MessageComponent, customID string) (string, error) {
	for _, component := range components {
		if actionRow, ok := component.(*discordgo.ActionsRow); ok {
			for _, item := range actionRow.Components {
				if input, ok := item.(*discordgo.TextInput); ok && input.CustomID == customID {
					return input.Value, nil
				}
			}
		}
	}
	return "", fmt.Errorf("component with customID %v not found", customID)
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
