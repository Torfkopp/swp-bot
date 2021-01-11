package main

import (
	"log"
	"os"
)

func main() {

	t1 := os.Getenv("discord_token") // Bot
	t2 := os.Getenv("bitbucket_token") // Bitbucket
	url := os.Getenv("url") // Endpoint URL
	
	// Create API access first
	api, err := NewAPI(url, t2)
	if err != nil {
		log.Fatal(err)
	}

	// Create BOT second
	swpbot := SessionCreate(t1)


	StartBot(api, swpbot)
}