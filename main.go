package main

import (
	"flag"
)

var (
	api    *API
	config string
	cfg    map[string]string
	color  = 4616416
	n      int64
)

func init() {
	flag.StringVar(&config, "config", "/home/user/swp_bot.config", "SWP-Bot configuration file")
	flag.BoolVar(&DebugFlag, "debug", false, "Run bot in foreground and enable debugging output")
	flag.Parse()
}

func main() {
	// Read config file
	cfg = ReadConfig()

	// Create API access first
	api, _ = NewAPI(cfg["REST_URL"], cfg["BITBUCKET_TOKEN"])

	// Create BOT second
	swpbot := SessionCreate(cfg["DISCORD_TOKEN"])

	StartBot(swpbot)
}
