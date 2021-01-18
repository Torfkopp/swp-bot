package main

import (
	"flag"
	"os"
)

var (
	api     *API
	UserLUT string
)

func init() {
	flag.StringVar(&UserLUT, "u", "/home/user/swp_users.lut", "File containing the UserLUT")
	flag.BoolVar(&DebugFlag, "d", false, "Run Bot in foreground and enable debugging output")
	flag.Parse()
}

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
