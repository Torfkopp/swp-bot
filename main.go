package main

import (
	"log"
	"os"
)

func main() {

	bt := os.Getenv("BITBUCKET_TOKEN") // Bitbucket
	dt := os.Getenv("DISCORD_TOKEN")   // Bot
	url := os.Getenv("REST_URL")       // Endpoint URL

	// Create API access first
	api, err := NewAPI(url, bt)
	if err != nil {
		log.Fatal(err)
	}

	// Create BOT second
	swpbot := SessionCreate(dt)

	StartBot(api, swpbot)
}
