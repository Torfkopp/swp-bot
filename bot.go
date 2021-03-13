package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
func StartBot(bot *discordgo.Session, bitbucketAPI1 *API, jiraAPI1 *API) {
	// Register events
	bot.AddHandler(Ready)
	bot.AddHandler(MessageCreate(bitbucketAPI1, jiraAPI1))

	// Create connection to bot
	err := bot.Open()
	if err != nil {
		log.Fatal("Error opening Discord session: ", err)
	}

	// Routine to listen to webhooks
	go ReceiveBitbucketWebhook(bot)

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
func Ready(session *discordgo.Session, _ *discordgo.Ready) {
	err := session.UpdateStatus(0, "") // Add whatever game you want here
	if err != nil {
		log.Fatal(err)
	}
}

// MessageCreate is called upon a received message and conditionally answers it
func MessageCreate(bitbucketAPI1 *API, jiraAPI1 *API) func(session *discordgo.Session, message *discordgo.MessageCreate) {
	// We need to return a function here so we can pass over the api object
	return func(session *discordgo.Session, message *discordgo.MessageCreate) {
		var err error
		if strings.HasPrefix(message.Content, "!post ") {
			_, err = session.ChannelMessageSend(cfg["RELAY_CHANNEL"], PostMessage(message.Content))
		}
		if strings.HasPrefix(message.Content, "!uwu ") {
			_, err = session.ChannelMessageSend(cfg["RELAY_CHANNEL"], PostUwUMessage(message.Content))
		}
		// This part is just for shits and giggles
		if message.Author.ID == cfg["VIP"] {
			switch message.Content {
			case "!help":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, HelpMessage(true))
			case "!allpullrequests":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, GetAllPullRequests(bitbucketAPI1, true))
			case "!mypullrequests":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, GetMyPullRequests(bitbucketAPI1, message.Author.ID, true))
			case "!myreviews":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, GetMyReviews(bitbucketAPI1, message.Author.ID, true))
			case "!about":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, AboutMessage(true))
			default:
				return
			}
		} else {
			// This part is the more serious portion of code here
			switch message.Content {
			case "!help":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, HelpMessage(false))
			case "!allpullrequests":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, GetAllPullRequests(bitbucketAPI1, false))
			case "!mypullrequests":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, GetMyPullRequests(bitbucketAPI1, message.Author.ID, false))
			case "!myreviews":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, GetMyReviews(bitbucketAPI1, message.Author.ID, false))
			//case "!sprint":
			//	_, err = session.ChannelMessageSendEmbed(message.ChannelID, GetActiveSprintMessage(jiraAPI1, false))
			case "!about":
				_, err = session.ChannelMessageSendEmbed(message.ChannelID, AboutMessage(false))
			default:
				return
			}
		}
		if err != nil {
			log.Println(err)
		}
	}
}
