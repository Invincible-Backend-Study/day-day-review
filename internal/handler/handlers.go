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
	guildId       string
	modalHandlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		cIdRegisterUserModal:  interactionRegisterUserModal,
		cIdRegisterScrumModal: interactionRegisterScrumModal,
	}
	commandHandlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		commandRegisterUser:       registerUser,
		commandRegisterTodayScrum: registerTodayScrum,
		commandGetTodayScrums:     getTodayScrums,
	}
)

func SetGuildId(id string) {
	guildId = id
}

// handleModalSubmit 모달의 제출을 처리합니다.
func handleModalSubmit(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	modalHandlers[interaction.ModalSubmitData().CustomID](session, interaction)
}

// handleApplicationCommand 봇의 명령어를 처리합니다.
func handleApplicationCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	commandHandlers[interaction.ApplicationCommandData().Name](session, interaction)
}

func getTodayScrums(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	scrums, err := service.GetTodayScrums()
	if err != nil {
		logErrorAndSendMessage(session, interaction, "오늘의 다짐을 불러오는 중 오류가 발생했습니다.", err)
	}
	sendMessage(session, interaction, scrumsToString(scrums))
}

func registerUser(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interaction.Interaction, createRegisterUserModal())
	if err != nil {
		log.Printf("Error responding with modal: %v", err)
	}
}

func registerTodayScrum(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	userId := interaction.Member.User.ID
	todayScrumExists, err := service.ExistTodayScrum(userId)
	if err != nil {
		logErrorAndSendMessage(session, interaction, "오늘의 다짐을 등록하는 중 오류가 발생했습니다.", err)
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
}

// scrumsToString scrum 목록을 문자열로 변환합니다.
func scrumsToString(scrums []model.ScrumDto) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("## 오늘(%s)의 다짐 목록: \n", util.GetTodayInKST().Format("2006-01-02")))
	for _, scrum := range scrums {
		builder.WriteString(fmt.Sprintf("\n### %s\n```\n오늘의 목표: %s\n오늘의 다짐: %s\n기분 점수: %d\n이유: %s\n```",
			scrum.Name, scrum.Goal, scrum.Commitment, scrum.FeelScore, scrum.FeelReason))
	}
	return builder.String()
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

// interactionRegisterUserModal 사용자 등록 모달의 상호작용을 처리합니다.
func interactionRegisterUserModal(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	nickname, err := extractValueFromComponent(interaction.
		ModalSubmitData().Components, cIdRegisterUserNicknameInput)
	if err != nil {
		logErrorAndSendMessage(session, interaction, "닉네임을 추출하는 중 오류가 발생했습니다.", err)
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
			logErrorAndSendMessage(session, interaction, fmt.Sprintf("Error in %s input", key), err)
			return
		}
		data[key] = value
	}
	feelScore, err := strconv.Atoi(data["feelScore"])
	if err != nil || feelScore < 0 || feelScore > 10 {
		logErrorAndSendMessage(session, interaction, "기분 점수가 올바르지 않습니다.", err)
		return
	}
	response := service.CreateTodayScrum(interaction.Member.User.ID, data["goal"], data["commitment"], data["feelReason"], feelScore)
	sendEphemeralMessage(session, interaction, response)
}
