package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/clinet/discordgo-embed"
	bitbucketserver "github.com/go-playground/webhooks/bitbucket-server"
	"log"
	"strconv"
	"strings"
)

// Define the color used in the embeds
const color = 4616416

// HelpMessage returns an embed containing an info text about the supported commands
func HelpMessage() *discordgo.MessageEmbed {
	return embed.NewGenericEmbedAdvanced("__These are the supported interactive commands:__",
		"**!help:** Shows this help text.\n"+
			"**!allpullrequests:** Shows the status of all active pull requests.\n"+
			"**!mypullrequests:** Shows the status of your own active pull requests.\n"+
			"**!myreviews:** Shows all pull requests which you're a reviewer of.\n"+
			"**!post <something>:** Relays your text into the bots channel.\n"+
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

// PostMessage strips a string off of its "!post" command
func PostMessage(message string) string {
	return strings.TrimPrefix(message, "!post ")
}

// NewPullRequestCreated returns the latest pull request
func NewPullRequestCreated(event bitbucketserver.PullRequestOpenedPayload) *discordgo.MessageEmbed {

	// Create an empty embed with a predefined color
	embedObject := embed.NewEmbed().SetColor(color)

	// Populate title and body from the data extracted of the event
	title := "**New pull request:**"
	body := "[" + event.PullRequest.Title + "](" + event.PullRequest.Links["self"].([]interface{})[0].(map[string]interface{})["href"].(string) + ")"

	// Add title and body previously populated to the embed and return it
	embedObject.SetTitle(title).SetDescription(body)
	return embedObject.MessageEmbed
}

// NewPullRequestPing returns a string containing the reviewers of the latest pull request
func NewPullRequestPing(event bitbucketserver.PullRequestOpenedPayload) string {
	// Populate title and body from the data extracted of the event
	text := "**Pinging Reviewers:**\n"
	for i, reviewer := range event.PullRequest.Reviewers {
		text = text + strconv.Itoa(i+1) + ". " + reviewer.User.DisplayName
		userid, present := cfg[reviewer.User.Name]
		if present {
			text = text + " <@" + userid + ">\n"
		} else {
			text = text + "\n"
		}
	}
	// As pings in Discord don't work in embeds, we need to return a simple string
	return text
}

// PullRequestMerged returns the merged pull request
func PullRequestMerged(event bitbucketserver.PullRequestMergedPayload) *discordgo.MessageEmbed {
	// Create an empty embed with a predefined color
	embedObject := embed.NewEmbed().SetColor(color)

	// Populate title and body from the data extracted of the event
	title := "**A pull request was merged:**"
	body := "[" + event.PullRequest.Title + "](" + event.PullRequest.Links["self"].([]interface{})[0].(map[string]interface{})["href"].(string) + ")"

	// Add title and body previously populated to the embed and return it
	embedObject.SetTitle(title).SetDescription(body)
	return embedObject.MessageEmbed
}

// PullRequestApproved returns the approved pull request
func PullRequestApproved(event bitbucketserver.PullRequestReviewerApprovedPayload) *discordgo.MessageEmbed {
	// Create an empty embed with a predefined color
	embedObject := embed.NewEmbed().SetColor(color)

	// Populate title and body from the data extracted of the event
	title := "**New review:**"
	body := "Someone approved [" + event.PullRequest.Title + "](" + event.PullRequest.Links["self"].([]interface{})[0].(map[string]interface{})["href"].(string) + ")"

	// Add title and body previously populated to the embed and return it
	embedObject.SetTitle(title).SetDescription(body)
	return embedObject.MessageEmbed
}

// PullRequestNeedsWork returns the pull request that needs work
func PullRequestNeedsWork(event bitbucketserver.PullRequestReviewerNeedsWorkPayload) *discordgo.MessageEmbed {
	// Create an empty embed with a predefined color
	embedObject := embed.NewEmbed().SetColor(color)

	// Populate title and body from the data extracted of the event
	title := "**New review:**"
	body := "This PR: [" + event.PullRequest.Title + "](" + event.PullRequest.Links["self"].([]interface{})[0].(map[string]interface{})["href"].(string) + ") " +
		"needs work!"

	// Add title and body previously populated to the embed and return it
	embedObject.SetTitle(title).SetDescription(body)
	return embedObject.MessageEmbed
}

// GetAllPullRequests returns all currently active pull requests from the rest response
func GetAllPullRequests(api *API) *discordgo.MessageEmbed {
	var (
		title string
		body  string
		field string
	)

	// Create an empty embed with a predefined color
	embedObject := embed.NewEmbed().SetColor(color)

	// Make a request to Bitbucket and iterate over it to fill title and field
	request, err := api.GetPullRequests()
	if err == nil {
		if len(request.Values) == 0 {
			title = "**There are no active pull requests!**"
		} else {
			title = "**Active pull requests:**\n"
			for i, values := range request.Values {
				field = "[*Reviewers:*](" + values.Links.Self[0].Href + ")\n"
				for j, reviewer := range values.Reviewers {
					field = field + strconv.Itoa(j+1) + ". [" + reviewer.User.DisplayName + "](" + reviewer.User.Links.Self[0].Href + ") "
					userid, present := cfg[reviewer.User.Name]
					if present {
						field = field + "<@" + userid + "> "
					}
					if reviewer.Approved {
						field = field + "APPROVED!\n"
					} else {
						field = field + "\n"
					}
				}
				field = field + "Comments: " + strconv.Itoa(values.Properties.CommentCount)
				embedObject.AddField(strconv.Itoa(i+1)+". "+values.Title, field)
			}
		}
	} else {
		title = "**Request returned no data!**"
		log.Println(err)
	}

	// Add title and body previously populated to the embed and return it
	embedObject.SetTitle(title).SetDescription(body)
	return embedObject.MessageEmbed
}

// GetMyPullRequests returns only the pull requests opened by the requesting user
func GetMyPullRequests(api *API, rid string) *discordgo.MessageEmbed {
	var (
		title string
		body  string
	)

	// Create an empty embed with a predefined color
	embedObject := embed.NewEmbed().SetColor(color)

	// Make a request to Bitbucket and iterate over it to fill title and body
	request, err := api.GetPullRequests()
	if err == nil {
		if len(request.Values) == 0 {
			title = "**There are no active pull requests!**"
		} else {
			username, present := cfg[rid]
			if present {
				title = "**Pull requests by " + username + ":**\n"
				i := 0
				for _, values := range request.Values {
					if values.Author.User.Name == username {
						body = body + strconv.Itoa(i+1) + ". [" + values.Title + "](" + values.Links.Self[0].Href + ")\n Reviewers:\n"
						for j, reviewer := range values.Reviewers {
							body = body + strconv.Itoa(j+1) + ". [" + reviewer.User.DisplayName + "](" + reviewer.User.Links.Self[0].Href + ") "
							userid, present := cfg[reviewer.User.Name]
							if present {
								body = body + "<@" + userid + "> "
							}
							if reviewer.Approved {
								body = body + "APPROVED!\n"
							} else {
								body = body + "\n"
							}
						}
						body = body + "Comments: " + strconv.Itoa(values.Properties.CommentCount)
						i++
					}
				}
				if body == "" {
					body = "*None*"
				}
			} else {
				title = "*Couldn't map your Discord ID to a Bitbucket user!*"
			}
		}
	} else {
		title = "**Request returned no data!**"
		log.Println(err)
	}

	// Add title and body previously populated to the embed and return it
	embedObject.SetTitle(title).SetDescription(body)
	return embedObject.MessageEmbed
}

// GetMyReviews returns all pull requests that the message requester is a reviewer of
func GetMyReviews(api *API, rid string) *discordgo.MessageEmbed {
	var (
		title string
		body  string
	)

	// Create an empty embed with a predefined color
	embedObject := embed.NewEmbed().SetColor(color)

	// Make a request to Bitbucket and iterate over it to fill title and body
	request, err := api.GetPullRequests()
	if err == nil {
		if len(request.Values) == 0 {
			title = "**There are no active pull requests!**"
		} else {
			username, present := cfg[rid]
			if present {
				title = "**Reviews assigned to " + username + ":**\n"
				i := 0
				for _, values := range request.Values {
					for _, reviewer := range values.Reviewers {
						if reviewer.User.Name == username {
							body = body + strconv.Itoa(i+1) + ". [" + values.Title + "](" + values.Links.Self[0].Href + ")\n"
							i++
						}
					}
				}
				if body == "" {
					body = "*None*"
				}
			} else {
				title = "*Couldn't map your Discord ID to a Bitbucket user!*"
			}
		}
	} else {
		title = "**Request returned no data!**"
		log.Println(err)
	}

	// Add title and body previously populated to the embed and return it
	embedObject.SetTitle(title).SetDescription(body)
	return embedObject.MessageEmbed
}

// This is considered WIP, don't use
func GetActiveSprintMessage(api *API) *discordgo.MessageEmbed {
	var (
		title string
		body  string
	)

	// Create an empty embed with a predefined color
	embedObject := embed.NewEmbed().SetColor(color)

	// New code here

	// Add title and body previously populated to the embed and return it
	embedObject.SetTitle(title).SetDescription(body)
	return embedObject.MessageEmbed
}
