package main

import (
	"flag"
	"log"
)

// TODO some of these shouldn't be global
var (
	bAPI1         *API
	bAPI2         *API
	jAPI1         *API
	configFile    string
	timestampFile string
	cfg           map[string]string
	color         = 4616416
	debugFlag     bool
)

// init runs before main and parses the cli arguments
func init() {
	flag.StringVar(&configFile, "config", "/home/user/swp_bot.config", "SWP-Bot configuration file")
	flag.StringVar(&timestampFile, "timestamp", "/tmp/swp_bot.timestamp", "Timestamp of latest PR in UNIX time")
	flag.BoolVar(&debugFlag, "debug", false, "Run bot in foreground and enable debugging output")
	flag.Parse()
}

// main is the main function to run obviously
func main() {
	var err error

	// Read config file
	cfg = ReadConfig()

	// Create API access first
	bAPI1, err = NewAPI(cfg["BITBUCKET_URL_1"], cfg["BITBUCKET_TOKEN"])
	bAPI2, err = NewAPI(cfg["BITBUCKET_URL_2"], cfg["BITBUCKET_TOKEN"])
	//jAPI1, err = NewAPI(cfg["JIRA_URL_1"], cfg["JIRA_TOKEN"])
	if err != nil {
		log.Fatal(err)
	}

	// Create BOT and start it
	bot := SessionCreate(cfg["DISCORD_TOKEN"])
	StartBot(bot)
}
