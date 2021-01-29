package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/clinet/discordgo-embed"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func SessionCreate(tok string) *discordgo.Session {
	// Create a new Discord session
	bot, err := discordgo.New("Bot " + tok)
	if err != nil {
		log.Fatal(err)
	}

	return bot
}

func StartBot(a *API, bot *discordgo.Session) {
	// Register events
	bot.AddHandler(Ready)
	bot.AddHandler(MessageCreate)
	//bot.AddHandler(PeriodicMessage)

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

func Ready(s *discordgo.Session, _ *discordgo.Ready) {
	s.UpdateStatus(0, "Catan")
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	cfg := ReadConfig()

	if m.Author.ID == cfg["VIP"] {
		switch m.Content {
		case "!help":
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced(
				"__These awe the x3 suppowted intewactive commands:__",
				"**!allpullrequests:** Shows the UwU status of aww active puww wequests.\n"+
					"**!mypullrequests:** Shows the x3 status of youw own active puww wequests.\n"+
					"**!myreviews:**  Shows aww puww wequests which you'we a weviewew of, nya.\n"+
					"**!comments:** Shows the x3 comments >w< undew youw active puww *boops your nose* wequests. *(TODO)*",
				color))
		case "!allpullrequests":
			s.ChannelMessageSendEmbed(m.ChannelID, GetAllPullRequestsVIP(api))
		case "!mypullrequests":
			s.ChannelMessageSendEmbed(m.ChannelID, GetMyPullRequestsVIP(api, m.Author.ID))
		case "!myreviews":
			s.ChannelMessageSendEmbed(m.ChannelID, GetMyReviewsVIP(api, m.Author.ID))
		case "!comments":
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("*Nyot impwemented yet owo*", "", color))
		default:
			return
		}
	} else {
		switch m.Content {
		case "!help":
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced(
				"__These are the supported interactive commands:__",
				"**!allpullrequests:** Shows the status of all active pull requests.\n"+
					"**!mypullrequests:** Shows the status of your own active pull requests.\n"+
					"**!myreviews:** Shows all pull requests which you're a reviewer of.\n"+
					"**!comments:** Shows the comments under your active pull requests. *(TODO)*",
				color))
		case "!allpullrequests":
			s.ChannelMessageSendEmbed(m.ChannelID, GetAllPullRequests(api))
		case "!mypullrequests":
			s.ChannelMessageSendEmbed(m.ChannelID, GetMyPullRequests(api, m.Author.ID))
		case "!myreviews":
			s.ChannelMessageSendEmbed(m.ChannelID, GetMyReviews(api, m.Author.ID))
		case "!comments":
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("*Not implemented yet*", "", color))
		default:
			return
		}
	}
}
