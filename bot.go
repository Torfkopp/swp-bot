package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

func SessionCreate(tok string) *discordgo.Session {
	// Create a new Discord session
	bot, err := discordgo.New("Bot " + tok)
	if err != nil {
		panic(err)
	}

	return bot
}

// This function needs improvement, like a lot
func StartBot(a *API, bot *discordgo.Session) {
	// Register events
	bot.AddHandler(Ready)
	bot.AddHandler(MessageCreate)

	err := bot.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}

func Ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "Catan")
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Content == "!reviewers" {
		s.ChannelMessageSend(m.ChannelID, GetReviewers()) // TODO: Get the used API into this handler
	}

	// Just for testing purpose
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "pong")
	}

}
