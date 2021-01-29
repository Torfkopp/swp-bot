package main

import (
	"flag"
)

var (
	api    *API
	config string
	color  = 4616416
)

func init() {
	flag.StringVar(&config, "config", "/home/user/swp_bot.config", "SWP-Bot configuration file")
	flag.BoolVar(&DebugFlag, "debug", false, "Run bot in foreground and enable debugging output")
	flag.Parse()
}

func main() {
	cfg := ReadConfig()

	// Create API access first
	api, _ = NewAPI(cfg["REST_URL"], cfg["BITBUCKET_TOKEN"])

	// Create BOT second
	swpbot := SessionCreate(cfg["DISCORD_TOKEN"])

	StartBot(api, swpbot)
}
