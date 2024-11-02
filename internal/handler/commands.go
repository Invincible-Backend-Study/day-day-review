package handler

import "github.com/bwmarrin/discordgo"

// 명령어 상수
const (
	commandRegisterUser = "회원-등록"
)

const (
	cIdRegisterUserModal         = "nickname_modal"
	cIdRegisterUserNicknameInput = "nickname_input"
)

// 명령어 목록
var commands = []*discordgo.ApplicationCommand{
	{
		Name:        commandRegisterUser,
		Description: "닉네임을 등록합니다",
	},
}
