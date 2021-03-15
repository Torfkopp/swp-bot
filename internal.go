package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	bitbucketserver "github.com/go-playground/webhooks/bitbucket-server"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// NewAPI implements API constructor
func NewAPI(location string, token string, tokenType int) (*API, error) {
	// Check if url isn't empty
	if len(location) == 0 {
		return nil, errors.New("url empty")
	}

	// Parse URL
	endPoint, err := url.ParseRequestURI(location)
	if err != nil {
		return nil, err
	}

	// Create new API object
	api := new(API)
	api.endPoint = endPoint
	api.tokenType = tokenType
	api.token = token

	// Make sure we use a valid and secure connection
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}

	// Set up http connection with a reasonable timeout
	api.client = &http.Client{Transport: transport, Timeout: time.Minute}

	return api, nil
}

// Auth implements token auth
func (api *API) Auth(request *http.Request) {
	// Supports unauthenticated access as well:
	// If token is not set, no authorization header is added
	if api.tokenType == 1 && api.token != "" {
		request.Header.Set("Authorization", "Bearer "+api.token)
	}
	if api.tokenType == 2 && api.token != "" {
		request.Header.Set("Authorization", "Basic "+api.token)
		request.Header.Set("Content-Type", "application/json")
	}
}

// ReadConfig reads the provided config file and turns it into a map
func ReadConfig() map[string]string {
	// Open config file
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Defer file closure so it runs after the return
	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Decode config file from json
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func ReceiveBitbucketWebhook(session *discordgo.Session) {
	hook, _ := bitbucketserver.New()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		payload, err := hook.Parse(request,
			bitbucketserver.PullRequestOpenedEvent,
			bitbucketserver.PullRequestMergedEvent,
			bitbucketserver.PullRequestCommentAddedEvent,
			bitbucketserver.PullRequestReviewerApprovedEvent,
			bitbucketserver.PullRequestReviewerNeedsWorkEvent)
		if err != nil {
			log.Println(err)
		}
		switch payload.(type) {
		case bitbucketserver.PullRequestOpenedPayload:
			event := payload.(bitbucketserver.PullRequestOpenedPayload)
			_, err = session.ChannelMessageSendEmbed(cfg["PING_CHANNEL"], NewPullRequestCreated(event))
			_, err = session.ChannelMessageSend(cfg["PING_CHANNEL"], NewPullRequestPing(event))
			go ReviewTimer(session, event)
		case bitbucketserver.PullRequestMergedPayload:
			event := payload.(bitbucketserver.PullRequestMergedPayload)
			_, err = session.ChannelMessageSendEmbed(cfg["PING_CHANNEL"], PullRequestMerged(event))
		case bitbucketserver.PullRequestReviewerApprovedPayload:
			event := payload.(bitbucketserver.PullRequestReviewerApprovedPayload)
			_, err = session.ChannelMessageSendEmbed(cfg["PING_CHANNEL"], PullRequestApproved(event))
		case bitbucketserver.PullRequestReviewerNeedsWorkPayload:
			event := payload.(bitbucketserver.PullRequestReviewerNeedsWorkPayload)
			_, err = session.ChannelMessageSendEmbed(cfg["PING_CHANNEL"], PullRequestNeedsWork(event))
		}
		if err != nil || err != bitbucketserver.ErrEventNotSpecifiedToParse {
			log.Println(err)
		}
	})
	err := http.ListenAndServe(":"+listenPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// MakeEmbed takes a title, body and fields to turn them into a Discord embed
func MakeEmbed(title, body string, fieldTitles, fieldBodies []string) *discordgo.MessageEmbed {
	embedObject := embed.NewEmbed().SetColor(color).SetTitle(title).SetDescription(body)

	if (fieldTitles != nil || fieldBodies != nil) && len(fieldTitles) == len(fieldBodies) {
		for i, fieldTitle := range fieldTitles {
			embedObject.AddField(fieldTitle, fieldBodies[i])
		}
	}

	return embedObject.MessageEmbed
}

// Debug outputs debug messages
func Debug(msg interface{}) {
	if debugFlag {
		log.Printf("%+v\n", msg)
	}
}

// Brace yourself for the next lines - Marvin

// UwU? What is dis?
// An UwUifiew by Mawio
// UwUify makes strings majestic
func UwUify(message string) string {
	message = UwuifyWords(message)
	message = UwuifyExclamations(message)
	message = UwuifySpaces(message)

	return message
}

// uwuify evewy wowd
func UwuifyWords(message string) string {

	words := strings.SplitAfter(message, " ")

	for i := 0; i < len(words); i++ {
		// t + h bad
		if len(words[i]) > 2 {
			if string(words[i][0])+string(words[i][1]) == "th" {
				words[i] = strings.Replace(words[i], "th", "d", 1)
			}
			if string(words[i][0])+string(words[i][1]) == "Th" {
				words[i] = strings.Replace(words[i], "Th", "D", 1)
			}
			if string(words[i][0])+string(words[i][1]) == "TH" {
				words[i] = strings.Replace(words[i], "TH", "D", 1)
			}
		}

		if !strings.ContainsAny(words[i], "(") {
			// smaww wettews
			words[i] = strings.ReplaceAll(words[i], "r", "w")
			words[i] = strings.ReplaceAll(words[i], "l", "w")
			words[i] = strings.ReplaceAll(words[i], "na", "nya")
			words[i] = strings.ReplaceAll(words[i], "ne", "nye")
			words[i] = strings.ReplaceAll(words[i], "ni", "nyi")
			words[i] = strings.ReplaceAll(words[i], "no", "nyo")
			words[i] = strings.ReplaceAll(words[i], "nu", "nyu")
			words[i] = strings.ReplaceAll(words[i], "ma", "mya")
			words[i] = strings.ReplaceAll(words[i], "me", "mye")
			words[i] = strings.ReplaceAll(words[i], "mi", "myi")
			words[i] = strings.ReplaceAll(words[i], "mo", "myo")
			words[i] = strings.ReplaceAll(words[i], "mu", "myu")
			words[i] = strings.ReplaceAll(words[i], "ove", "uv")

			//BIG WETTEWS
			words[i] = strings.ReplaceAll(words[i], "R", "W")
			words[i] = strings.ReplaceAll(words[i], "L", "W")
			words[i] = strings.ReplaceAll(words[i], "Na", "Nya")
			words[i] = strings.ReplaceAll(words[i], "NA", "NYA")
			words[i] = strings.ReplaceAll(words[i], "Ne", "Nye")
			words[i] = strings.ReplaceAll(words[i], "NE", "NYE")
			words[i] = strings.ReplaceAll(words[i], "Ni", "Nyi")
			words[i] = strings.ReplaceAll(words[i], "NI", "NYI")
			words[i] = strings.ReplaceAll(words[i], "No", "Nyo")
			words[i] = strings.ReplaceAll(words[i], "NO", "NYO")
			words[i] = strings.ReplaceAll(words[i], "Nu", "Nyu")
			words[i] = strings.ReplaceAll(words[i], "NU", "NYU")
			words[i] = strings.ReplaceAll(words[i], "Ma", "Mya")
			words[i] = strings.ReplaceAll(words[i], "MA", "MYA")
			words[i] = strings.ReplaceAll(words[i], "Me", "Mye")
			words[i] = strings.ReplaceAll(words[i], "ME", "MYE")
			words[i] = strings.ReplaceAll(words[i], "Mi", "Myi")
			words[i] = strings.ReplaceAll(words[i], "MI", "MYI")
			words[i] = strings.ReplaceAll(words[i], "Mo", "Myo")
			words[i] = strings.ReplaceAll(words[i], "MO", "MYO")
			words[i] = strings.ReplaceAll(words[i], "Mu", "Myu")
			words[i] = strings.ReplaceAll(words[i], "MU", "MYU")
			words[i] = strings.ReplaceAll(words[i], "OVE", "UV")
			words[i] = strings.ReplaceAll(words[i], "th", "ff")
			words[i] = strings.ReplaceAll(words[i], "TH", "FF")

			// If yowouwu awe hawdcowowe enowouwugh
			//words[i] = strings.ReplaceAll(words[i], "u", "uwu")
			//words[i] = strings.ReplaceAll(words[i], "U", "UwU")
			//words[i] = strings.ReplaceAll(words[i], "o", "owo")
			//words[i] = strings.ReplaceAll(words[i], "O", "OwO")

			// speciaws
			words[i] = strings.ReplaceAll(words[i], " nyo ", " nyo UnU ")
			words[i] = strings.ReplaceAll(words[i], " NYO ", " NYO UnU ")
			words[i] = strings.ReplaceAll(words[i], " nyot ", " nyot UnU ")
			words[i] = strings.ReplaceAll(words[i], " NYOT ", " NYOT UnU ")
			words[i] = strings.ReplaceAll(words[i], "n't ", "nyot UnU ")
			words[i] = strings.ReplaceAll(words[i], "N'T ", "NYOT UnU ")
			words[i] = strings.ReplaceAll(words[i], "nya ", "nya~ ")
			words[i] = strings.ReplaceAll(words[i], "NYA ", "NYA~ ")
		}
	}

	return strings.Join(words, " ")
}

// uwuify ffe excwamation mawks
func UwuifyExclamations(message string) string {
	excl := strings.Split(message, "!")
	for i := 0; i < len(excl)-1; i++ {
		excl[i] += "!"
		randE := rand.Intn(6)
		switch randE {
		case 0:
			excl[i] = strings.ReplaceAll(excl[i], "!", "!?")
		case 1:
			excl[i] = strings.ReplaceAll(excl[i], "!", "?!!")
		case 2:
			excl[i] = strings.ReplaceAll(excl[i], "!", "?!?1")
		case 3:
			excl[i] = strings.ReplaceAll(excl[i], "!", "!!11")
		case 4:
			excl[i] = strings.ReplaceAll(excl[i], "!", "?!?!")
		case 5:
			excl[i] = strings.ReplaceAll(excl[i], "!", "!!?!!")
		}
	}

	return strings.Join(excl, " ")
}

// uwuify wiff wandom faces/actions/stuttews between ffe wowds
func UwuifySpaces(message string) string {
	//Numbers in percent
	faces := 5
	actions := 10
	stutters := 30

	words := strings.Split(message, " ")

	for i := 0; i < len(words); i++ {
		random := rand.Intn(101)
		if 0 <= random && random <= faces {
			//Add random face before the word
			randF := rand.Intn(10)
			switch randF {
			case 0:
				words[i] += " (・`ω´・)"
			case 1:
				words[i] += " ;;w;;"
			case 2:
				words[i] += " OwO"
			case 3:
				words[i] += " UwU"
			case 4:
				words[i] += " >w<"
			case 5:
				words[i] += " ^w^"
			case 6:
				words[i] += " ÚwÚ"
			case 7:
				words[i] += " ^-^"
			case 8:
				words[i] += " :3"
			case 9:
				words[i] += " x3"
			}
		} else if faces < random && random <= actions {
			//Add random action before the word
			randA := rand.Intn(15)
			switch randA {
			case 0:
				words[i] += " *blushes*"
			case 1:
				words[i] += " *whispers to self*"
			case 2:
				words[i] += " *cries*"
			case 3:
				words[i] += " *screams*"
			case 4:
				words[i] += " *sweats*"
			case 5:
				words[i] += " *twerks*"
			case 6:
				words[i] += " *runs away*"
			case 7:
				words[i] += " *screeches*"
			case 8:
				words[i] += " *walks away*"
			case 9:
				words[i] += " *sees bulge*"
			case 10:
				words[i] += " *looks at you*"
			case 11:
				words[i] += " *notices buldge*"
			case 12:
				words[i] += " *starts twerking*"
			case 13:
				words[i] += " *huggles tightly*"
			case 14:
				words[i] += " *boops your nose*"
			}

		} else if actions < random && random <= stutters {
			//Add stutter with a length between 0 and 2
			randS := rand.Intn(3)
			if words[i] != "" && !strings.ContainsAny(words[i], "<[\n") {
				words[i] = strings.Repeat(string(words[i][0])+"-", randS) + words[i]
			}
		}
	}

	return strings.Join(words, " ")
}

// Congrats if you made it until here - Marvin
