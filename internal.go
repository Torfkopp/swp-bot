package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/bwmarrin/discordgo"
	bitbucketserver "github.com/go-playground/webhooks/bitbucket-server"
	"log"
	"net/http"
	"net/url"
	"os"
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

// Debug outputs debug messages
func Debug(msg interface{}) {
	if debugFlag {
		log.Printf("%+v\n", msg)
	}
}
