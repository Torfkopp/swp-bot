package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/clinet/discordgo-embed"
	"log"
	"strconv"
)

// HelpMessage returns an embed containing an info text about the supported commands
func HelpMessage() *discordgo.MessageEmbed {
	return embed.NewGenericEmbedAdvanced("__These are the supported interactive commands:__",
		"**!allpullrequests:** Shows the status of all active pull requests.\n"+
			"**!mypullrequests:** Shows the status of your own active pull requests.\n"+
			"**!myreviews:** Shows all pull requests which you're a reviewer of.\n"+
			"**!comments:** Shows the comments under your active pull requests. *(TODO)*\n"+
			"**!about:** Some info about this bot.",
		color)
}

// AboutMessage returns an embed containing an "about" text
func AboutMessage() *discordgo.MessageEmbed {
	return embed.NewGenericEmbedAdvanced("About this bot:",
		"In case of undesired risks and side effects\n"+
			"please read the [source code](https://github.com/MDr164/swp-bot) or ask your local dev.",
		color)
}

// NewPullRequestCreated returns the latest pull request
func NewPullRequestCreated(api *API) *discordgo.MessageEmbed {
	var (
		title string
		body  string
	)

	embedObject := embed.NewEmbed().SetColor(color)

	request, err := api.GetActivePullRequests()
	if err == nil {
		title = "**New pull request:**\n"
		body = "[" + request.Values[0].Title + "](" + request.Values[0].Links.Self[0].Href + ")"
	} else {
		title = "**Request returned no data!**"
		log.Println(err)
	}

	embedObject.SetTitle(title).SetDescription(body)
	return embedObject.MessageEmbed
}

// NewPullRequestPing returns a string containing the reviewers of the latest pull request
func NewPullRequestPing(api *API) string {
	var text string

	request, err := api.GetActivePullRequests()
	if err == nil {
		text = "**Pinging Reviewers:**\n"
		for i, rev := range request.Values[0].Reviewers {
			text = text + strconv.Itoa(i+1) + ". " + rev.User.DisplayName
			userid, present := cfg[rev.User.Name]
			if present {
				text = text + " <@" + userid + ">\n"
			} else {
				text = text + "\n"
			}
		}
	} else {
		text = "**Request returned no data!**"
		log.Println(err)
	}

	return text
}

// GetAllPullRequests returns all currently active pull requests from the rest response
func GetAllPullRequests(api *API) *discordgo.MessageEmbed {
	var (
		title string
		body  string
		field string
	)

	embedObject := embed.NewEmbed().SetColor(color)

	request, err := api.GetActivePullRequests()
	if err == nil {
		if len(request.Values) == 0 {
			title = "**There are no active pull requests!**"
		} else {
			title = "**Active pull requests:**\n"
			for i, val := range request.Values {
				field = ""
				for j, rev := range val.Reviewers {
					field = field + strconv.Itoa(j+1) + ". [" + rev.User.DisplayName + "](" + rev.User.Links.Self[0].Href + ") "
					userid, present := cfg[rev.User.Name]
					if present {
						field = field + "<@" + userid + ">\n"
					} else {
						field = field + "\n"
					}
				}
				embedObject.AddField(strconv.Itoa(i+1)+". "+val.Title, "[*Reviewers:*]("+val.Links.Self[0].Href+")\n"+field)
			}
		}
	} else {
		title = "**Request returned no data!**"
		log.Println(err)
	}

	embedObject.SetTitle(title).SetDescription(body)
	return embedObject.MessageEmbed
}

// GetMyPullRequests returns only the pull requests opened by the requesting user
func GetMyPullRequests(api *API, rid string) *discordgo.MessageEmbed {
	var (
		title string
		body  string
	)

	embedObject := embed.NewEmbed().SetColor(color)

	request, err := api.GetActivePullRequests()
	if err == nil {
		if len(request.Values) == 0 {
			title = "**There are no active pull requests!**"
		} else {
			username, present := cfg[rid]
			if present {
				title = "**Pull requests by " + username + ":**\n"
				i := 0
				for _, val := range request.Values {
					if val.Author.User.Name == username {
						body = body + strconv.Itoa(i+1) + ". [" + val.Title + "](" + val.Links.Self[0].Href + ")\n"
						i++
					}
				}
				if body == "" {
					body = "*None*"
				}
			} else {
				title = "*Couldn'title map your Discord ID to api Bitbucket user!*"
			}
		}
	} else {
		title = "**Request returned no data!**"
		log.Println(err)
	}

	embedObject.SetTitle(title).SetDescription(body)
	return embedObject.MessageEmbed
}

// GetMyReviews returns all pull requests that the message requester is a reviewer of
func GetMyReviews(api *API, rid string) *discordgo.MessageEmbed {
	var (
		title string
		body  string
	)

	embedObject := embed.NewEmbed().SetColor(color)

	request, err := api.GetActivePullRequests()
	if err == nil {
		if len(request.Values) == 0 {
			title = "**There are no active pull requests!**"
		} else {
			username, present := cfg[rid]
			if present {
				title = "**Reviews assigned to " + username + ":**\n"
				i := 0
				for _, val := range request.Values {
					for _, rev := range val.Reviewers {
						if rev.User.Name == username {
							body = body + strconv.Itoa(i+1) + ". [" + val.Title + "](" + val.Links.Self[0].Href + ")\n"
							i++
						}
					}
				}
				if body == "" {
					body = "*None*"
				}
			} else {
				title = "*Couldn'title map your Discord ID to api Bitbucket user!*"
			}
		}
	} else {
		title = "**Request returned no data!**"
		log.Println(err)
	}

	embedObject.SetTitle(title).SetDescription(body)
	return embedObject.MessageEmbed
}
