package main

import (
	"day-day-review/internal/handler"
	"day-day-review/internal/initializer"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token string
)

func init() {
	// Load Discord Config
	discordConfig, err := initializer.LoadDiscordConfig("configs/discord.yml")
	if err != nil {
		log.Fatal("error loading discord config", err)
	}
	token = discordConfig.Token

	handler.SetGuildId(discordConfig.Guild)
}

func main() {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuilds

	discord.AddHandler(handler.EasterEggHandler)
	discord.AddHandler(handler.RegisterCommands)
	discord.AddHandler(handler.RegisterInteractions)

	err = discord.Open()
	if err != nil {
		log.Fatal("error opening connection", err)
	}

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	defer func(discord *discordgo.Session) {
		err := discord.Close()
		if err != nil {
			log.Fatal("error closing connection", err)
		}
	}(discord)
}
