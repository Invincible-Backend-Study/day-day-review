package main

import (
	"database/sql"
	"day-day-review/internal/handler"
	"day-day-review/internal/initializer"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token   string
	db      *sql.DB
	guildID string
	manager *handler.Manager
)

func init() {
	flag.StringVar(&token, "t", "", "Bot token")
	flag.StringVar(&guildID, "g", "", "Guild ID")
	flag.Parse()

	var err error
	db, err = initializer.InitDatabase()
	if err != nil {
		log.Fatal("error initializing database", err)
	}

	manager = handler.NewHandlerManager(db, guildID)
}

func main() {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuilds

	discord.AddHandler(handler.EasterEggHandler)
	discord.AddHandler(manager.RegisterCommands)
	discord.AddHandler(manager.RegisterInteractions)

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
