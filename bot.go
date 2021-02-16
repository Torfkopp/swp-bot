package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/clinet/discordgo-embed"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func SessionCreate(tok string) *discordgo.Session {
	// Create a new Discord session
	bot, err := discordgo.New("Bot " + tok)
	if err != nil {
		log.Fatal(err)
	}

	return bot
}

func StartBot(bot *discordgo.Session) {
	// Register events
	bot.AddHandler(Ready)
	bot.AddHandler(MessageCreate)

	err := bot.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	PeriodicMessage(bot)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}

func Ready(s *discordgo.Session, _ *discordgo.Ready) {
	s.UpdateStatus(1, "")
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == cfg["VIP"] {
		switch m.Content {
		case "!help":
			s.ChannelMessageSendEmbed(m.ChannelID, HelpMessageVIP())
		case "!allpullrequests":
			s.ChannelMessageSendEmbed(m.ChannelID, GetAllPullRequestsVIP(bAPI1))
		case "!mypullrequests":
			s.ChannelMessageSendEmbed(m.ChannelID, GetMyPullRequestsVIP(bAPI1, m.Author.ID))
		case "!myreviews":
			s.ChannelMessageSendEmbed(m.ChannelID, GetMyReviewsVIP(bAPI1, m.Author.ID))
		case "!comments":
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("*Nyot impwemented yet owo*", "", color))
		case "!about":
			s.ChannelMessageSendEmbed(m.ChannelID, AboutMessageVIP())
		default:
			return
		}
	} else {
		switch m.Content {
		case "!help":
			s.ChannelMessageSendEmbed(m.ChannelID, HelpMessage())
		case "!allpullrequests":
			s.ChannelMessageSendEmbed(m.ChannelID, GetAllPullRequests(bAPI1))
		case "!mypullrequests":
			s.ChannelMessageSendEmbed(m.ChannelID, GetMyPullRequests(bAPI1, m.Author.ID))
		case "!myreviews":
			s.ChannelMessageSendEmbed(m.ChannelID, GetMyReviews(bAPI1, m.Author.ID))
		case "!comments":
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("*Not implemented yet*", "", color))
		case "!about":
			s.ChannelMessageSendEmbed(m.ChannelID, AboutMessage())
		default:
			return
		}
	}
}

func PeriodicMessage(s *discordgo.Session) {
	for {
		time.Sleep(3 * time.Minute)
		if CheckNewPullRequest(bAPI1) {
			s.ChannelMessageSendEmbed(cfg["PING_CHANNEL"], NewPullRequestCreated(bAPI1))
			s.ChannelMessageSend(cfg["PING_CHANNEL"], NewPullRequestPing(bAPI1))
		}
	}
}
