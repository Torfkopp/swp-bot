package main

import (
	"os"
)

var api *API

func main() {

	bt := os.Getenv("BITBUCKET_TOKEN") // Bitbucket
	dt := os.Getenv("DISCORD_TOKEN")   // Bot
	url := os.Getenv("REST_URL")       // Endpoint URL

	// Create API access first
	api, _ = NewAPI(url, bt)

	// Create BOT second
	swpbot := SessionCreate(dt)

	StartBot(api, swpbot)
}
