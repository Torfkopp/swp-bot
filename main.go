package main

import (
	"flag"
	"log"
)

var (
	configFile string
	listenPort string
	cfg        map[string]string
	debugFlag  bool
)

// init runs before main and parses the cli arguments
func init() {
	flag.StringVar(&configFile, "config", "/home/user/swp_bot.config", "SWP-Bot configuration file")
	flag.StringVar(&listenPort, "port", "50015", "Port on which to receive Webhooks")
	flag.BoolVar(&debugFlag, "debug", false, "Run bot in foreground and enable debugging output")
	flag.Parse()
}

// main is the main function to run obviously
func main() {
	var err error

	// Read config file
	cfg = ReadConfig()

	// Create API access first
	bAPI1, err := NewAPI(cfg["BITBUCKET_URL_1"], cfg["BITBUCKET_TOKEN"], 1)
	bAPI2, err := NewAPI(cfg["BITBUCKET_URL_2"], cfg["BITBUCKET_TOKEN"], 1)
	//jAPI1, err := NewAPI(cfg["JIRA_URL_1"], cfg["JIRA_TOKEN"], 2)
	if err != nil {
		log.Fatal(err)
	}

	// Create BOT and start it
	bot := SessionCreate(cfg["DISCORD_TOKEN"])
	StartBot(bot, bAPI1, bAPI2)
}
