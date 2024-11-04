package handler

import (
	"day-day-review/internal/model"
	"day-day-review/internal/service"
	"day-day-review/internal/util"
	"fmt"
	"log"
	"strconv"
	"strings"

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
		userId := interaction.Member.User.ID
		todayScrumExists, err := service.ExistTodayScrum(userId)
		if err != nil {
			log.Println("Error select today scrum: ", err)
			sendEphemeralMessage(session, interaction, fmt.Sprint("%w", err))
			return
		}
		if todayScrumExists {
			sendEphemeralMessage(session, interaction, "오늘의 다짐을 이미 작성하셨습니다.")
			return
		}

		err = session.InteractionRespond(interaction.Interaction, createRegisterScrumModal())
		if err != nil {
			log.Printf("Error responding with modal: %v", err)
		}
	case commandGetTodayScrums:
		scrums, err := service.GetTodayScrums()
		if err != nil {
			log.Println("Error select today scrums: ", err)
			sendEphemeralMessage(session, interaction, fmt.Sprint("%w", err))
			return
		}
		sendMessage(session, interaction, scrumsToString(scrums))
	}
}

// scrumsToString scrum 목록을 문자열로 변환합니다.
func scrumsToString(scrums []model.ScrumDto) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("## 오늘(%s)의 다짐 목록: \n", util.GetTodayInKST().Format("2006-01-02")))
	for _, scrum := range scrums {
		result.WriteString(fmt.Sprintf("\n### %s\n```\n오늘의 목표: %s\n오늘의 다짐: %s\n기분 점수: %d\n이유: %s\n```",
			scrum.Name, scrum.Goal, scrum.Commitment, scrum.FeelScore, scrum.FeelReason))
	}
	return result.String()
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
	inputs := map[string]string{
		"goal":       cIdRegisterScrumGoalInput,
		"commitment": cIdRegisterScrumCommitmentInput,
		"feelScore":  cIdRegisterScrumScoreInput,
		"feelReason": cIdRegisterScrumReasonInput,
	}
	data := make(map[string]string)
	for key, componentId := range inputs {
		value, err := extractValueFromComponent(interaction.ModalSubmitData().Components, componentId)
		if err != nil {
			log.Println("Error extracting value from component: ", err)
			sendEphemeralMessage(session, interaction, fmt.Sprintf("Error in %s input: %v", key, err))
			return
		}
		data[key] = value
	}
	feelScore, err := strconv.Atoi(data["feelScore"])
	if err != nil || feelScore < 0 || feelScore > 10 {
		log.Println("Error converting string to int:", err)
		sendEphemeralMessage(session, interaction, fmt.Sprintf("Error in %d input: %v", feelScore, err))
		return
	}
	response := service.CreateTodayScrum(interaction.Member.User.ID, data["goal"], data["commitment"], data["feelReason"], feelScore)
	sendEphemeralMessage(session, interaction, response)
}
