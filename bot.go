package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// SessionCreate creates a new Discord session
func SessionCreate(token string) *discordgo.Session {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	return bot
}

// StartBot adds event handlers and starts the main bot function
func StartBot(bot *discordgo.Session, api *API) {
	// Register events
	bot.AddHandler(Ready)
	bot.AddHandler(MessageCreate(api))

	// Create connection to bot
	err := bot.Open()
	if err != nil {
		log.Fatal("Error opening Discord session: ", err)
	}

	// Routine to periodically post messages
	PeriodicMessage(bot, api)

	// Wait here until CTRL-C or other term signal is received.
	// TODO Only call this if running in debug mode otherwise create background thread
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	err = bot.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// Ready is called upon a ready event and sets the bots status
func Ready(s *discordgo.Session, _ *discordgo.Ready) {
	err := s.UpdateStatus(0, "") // Add whatever game you want here
	if err != nil {
		log.Fatal(err)
	}
}

// MessageCreate is called upon a received message and conditionally answers it
func MessageCreate(api *API) func(session *discordgo.Session, message *discordgo.MessageCreate) {
	// We need to return a function here so we can pass over the api object
	return func(session *discordgo.Session, message *discordgo.MessageCreate) {
		var err error
		if strings.HasPrefix(message.Content, "!post ") {
			_, err = session.ChannelMessageSend(cfg["PING_CHANNEL"], PostMessage(message.Content))
		}
		// This part is just for shits and giggles
		if message.Author.ID == cfg["VIP"] {
			switch message.Content {
			case "!help":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, HelpMessageVIP())
			case "!allpullrequests":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, GetAllPullRequestsVIP(api))
			case "!mypullrequests":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, GetMyPullRequestsVIP(api, message.Author.ID))
			case "!myreviews":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, GetMyReviewsVIP(api, message.Author.ID))
			case "!about":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, AboutMessageVIP())
			default:
				return
			}
		} else {
			// This part is the more serious portion of code here
			switch message.Content {
			case "!help":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, HelpMessage())
			case "!allpullrequests":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, GetAllPullRequests(api))
			case "!mypullrequests":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, GetMyPullRequests(api, message.Author.ID))
			case "!myreviews":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, GetMyReviews(api, message.Author.ID))
			case "!about":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, AboutMessage())
			default:
				return
			}
		}
		if err != nil {
			log.Println(err)
		}
	}
}

// PeriodicMessage conditionally sends a message to a specified channel every 3 minutes
// TODO Make this a non-blocking routine
func PeriodicMessage(session *discordgo.Session, api *API) {
	var err error
	for {
		time.Sleep(3 * time.Minute)
		// Only send message if there are actual news
		if CheckNewPullRequest(api) {
			_, err = session.ChannelMessageSendEmbed(cfg["PING_CHANNEL"], NewPullRequestCreated(api))
			_, err = session.ChannelMessageSend(cfg["PING_CHANNEL"], NewPullRequestPing(api))
			if err != nil {
				log.Println(err)
			}
		}
	}
}
