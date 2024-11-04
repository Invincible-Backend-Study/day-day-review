package handler

import (
	"day-day-review/internal/service"
	"fmt"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var (
	guildId string
)

func SetGuildId(id string) {
	guildId = id
}

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

// handleApplicationCommand 봇의 명령어를 처리합니다.
func handleApplicationCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.ApplicationCommandData().Name {
	case commandRegisterUser:
		err := session.InteractionRespond(interaction.Interaction, createRegisterUserModal())
		if err != nil {
			log.Printf("Error responding with modal: %v", err)
		}
	case commandRegisterTodayScrum:
		err := session.InteractionRespond(interaction.Interaction, createRegisterScrumModal())
		if err != nil {
			log.Printf("Error responding with modal: %v", err)
		}
	}
}

// handleModalSubmit 모달의 제출을 처리합니다.
func handleModalSubmit(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.ModalSubmitData().CustomID {
	case cIdRegisterUserModal:
		interactionRegisterUserModal(session, interaction)
	case cIdRegisterScrumModal:
		interactionRegisterScrumModal(session, interaction)
	}
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

// interactionRegisterUserModal 사용자 등록 모달의 상호작용을 처리합니다.
func interactionRegisterUserModal(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	nickname, err := extractValueFromComponent(interaction.
		ModalSubmitData().Components, cIdRegisterUserNicknameInput)
	if err != nil {
		log.Println("Error extracting value from component: ", err)
		sendEphemeralMessage(session, interaction, "닉네임을 입력해주세요.")
		return
	}
	log.Println("Received nickname:", nickname)

	userId := interaction.Member.User.ID

	response := service.AddUser(nickname, userId)

	sendEphemeralMessage(session, interaction, response)
}

// interactionRegisterScrumModal 오늘의 다짐 등록 모달의 상호작용을 처리합니다.
func interactionRegisterScrumModal(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	goal, err := extractValueFromComponent(interaction.
		ModalSubmitData().Components, cIdRegisterScrumGoalInput)
	if err != nil {
		log.Println("Error extracting value from component: ", err)
		sendEphemeralMessage(session, interaction, "닉네임을 입력해주세요.")
		return
	}

	commitment, err := extractValueFromComponent(interaction.
		ModalSubmitData().Components, cIdRegisterScrumCommitmentInput)
	if err != nil {
		log.Println("Error extracting value from component: ", err)
		sendEphemeralMessage(session, interaction, "오늘의 다짐을 입력해주세요.")
		return
	}

	feelScoreStr, err := extractValueFromComponent(interaction.
		ModalSubmitData().Components, cIdRegisterScrumScoreInput)
	if err != nil {
		log.Println("Error extracting value from component: ", err)
		sendEphemeralMessage(session, interaction, "기분 점수를 입력해주세요.")
		return
	}
	feelScore, err := strconv.Atoi(feelScoreStr)
	if err != nil {
		log.Println("Error converting string to int:", err)
		sendEphemeralMessage(session, interaction, "기분 점수는 0 이상 10 이하 숫자로 입력해주세요.")
		return
	}

	feelReason, err := extractValueFromComponent(interaction.
		ModalSubmitData().Components, cIdRegisterScrumReasonInput)
	if err != nil {
		log.Println("Error extracting value from component: ", err)
		sendEphemeralMessage(session, interaction, "기분 점수의 이유를 입력해주세요.")
		return
	}

	// log.Println(goal, commitment, feelScore, feelReason)
	userId := interaction.Member.User.ID

	response := service.CreateTodayScrum(userId, goal, commitment, feelReason, feelScore)

	sendEphemeralMessage(session, interaction, response)
}
