package main

import (
	"github.com/bwmarrin/discordgo"
	bitbucketserver "github.com/go-playground/webhooks/bitbucket-server"
	"log"
	"strconv"
	"strings"
)

// Define the color used in the embeds
const color = 4616416

// HelpMessage returns an embed containing an info text about the supported commands
func HelpMessage(vip bool) *discordgo.MessageEmbed {
	if vip {
		title := "__These awe the x3 suppowted intewactive commands:__"
		body := "**!allpullrequests:** Shows the UwU status of aww active puww wequests.\n" +
			"**!mypullrequests:** Shows the x3 status of youw own active puww wequests.\n" +
			"**!myreviews:**  Shows aww puww wequests which you'we a weviewew of, nya.\n" +
			"**!post <something>:** Weways youw text into the x3 bots *huggles tightly* channyew.\n" +
			"**!about:** Some info *sweats* about this bot."
		return MakeEmbed(title, body, nil, nil)
	}
	title := "__These are the supported interactive commands:__"
	body := "**!help:** Shows this help text.\n" +
		"**!allpullrequests:** Shows the status of all active pull requests.\n" +
		"**!mypullrequests:** Shows the status of your own active pull requests.\n" +
		"**!myreviews:** Shows all pull requests which you're a reviewer of.\n" +
		"**!post <something>:** Relays your text into the bots channel.\n" +
		"**!about:** Some info about this bot."
	return MakeEmbed(title, body, nil, nil)
}

// AboutMessage returns an embed containing an "about" text
func AboutMessage(vip bool) *discordgo.MessageEmbed {
	if vip {
		title := "About this owo bot:"
		body := "In case of undesiwed wisks and side effects\n" +
			"pwease wead the x3 [souwce code](https://github.com/MDr164/swp-bot) ow (・`ω´・) ask youw wocaw dev."
		return MakeEmbed(title, body, nil, nil)
	}
	title := "About this bot:"
	body := "In case of undesired risks and side effects\n" +
		"please read the [source code](https://github.com/MDr164/swp-bot) or ask your local dev."
	return MakeEmbed(title, body, nil, nil)
}

// PostMessage strips a string off of its "!post" command
func PostMessage(message string) string {
	return strings.TrimPrefix(message, "!post ")
}

// PostUwUMessage does the same as PostMessage but cursed
func PostUwUMessage(message string) string {
	return UwUify(strings.TrimPrefix(message, "!uwu "))
}

// NewPullRequestCreated returns the latest pull request
func NewPullRequestCreated(event bitbucketserver.PullRequestOpenedPayload) *discordgo.MessageEmbed {
	// Populate title and body from the data extracted of the event
	title := "**New pull request:**"
	body := "[" + event.PullRequest.Title + "](" + event.PullRequest.Links["self"].([]interface{})[0].(map[string]interface{})["href"].(string) + ")"

	// Add title and body to an embed and return it
	return MakeEmbed(title, body, nil, nil)
}

// NewPullRequestPing returns a string containing the reviewers of the latest pull request
func NewPullRequestPing(event bitbucketserver.PullRequestOpenedPayload) string {
	// Populate title and body from the data extracted of the event
	text := "**Pinging Reviewers:**\n"
	for i, reviewer := range event.PullRequest.Reviewers {
		text += strconv.Itoa(i+1) + ". " + reviewer.User.DisplayName
		userid, present := cfg[reviewer.User.Name]
		if present {
			text += " <@" + userid + ">\n"
		} else {
			text += "\n"
		}
	}
	// As pings in Discord don't work in embeds, we need to return a simple string
	return text
}

// PullRequestMerged returns the merged pull request
func PullRequestMerged(event bitbucketserver.PullRequestMergedPayload) *discordgo.MessageEmbed {
	// Populate title and body from the data extracted of the event
	title := "**A pull request was merged:**"
	body := "[" + event.PullRequest.Title + "](" + event.PullRequest.Links["self"].([]interface{})[0].(map[string]interface{})["href"].(string) + ")"

	// Add title and body to an embed and return it
	return MakeEmbed(title, body, nil, nil)
}

// PullRequestApproved returns the approved pull request
func PullRequestApproved(event bitbucketserver.PullRequestReviewerApprovedPayload) *discordgo.MessageEmbed {
	// Populate title and body from the data extracted of the event
	title := "**New review:**"
	body := "Someone approved [" + event.PullRequest.Title + "](" + event.PullRequest.Links["self"].([]interface{})[0].(map[string]interface{})["href"].(string) + ")"

	// Add title and body to an embed and return it
	return MakeEmbed(title, body, nil, nil)
}

// PullRequestNeedsWork returns the pull request that needs work
func PullRequestNeedsWork(event bitbucketserver.PullRequestReviewerNeedsWorkPayload) *discordgo.MessageEmbed {
	// Populate title and body from the data extracted of the event
	title := "**New review:**"
	body := "This PR: [" + event.PullRequest.Title + "](" + event.PullRequest.Links["self"].([]interface{})[0].(map[string]interface{})["href"].(string) + ") " +
		"needs work!"

	// Add title and body to an embed and return it
	return MakeEmbed(title, body, nil, nil)
}

// GetAllPullRequests returns all currently active pull requests from the rest response
func GetAllPullRequests(api *API, vip bool) *discordgo.MessageEmbed {
	var (
		title       string
		body        string
		fieldTitles []string
		fieldBodies []string
	)

	// Make a request to Bitbucket and iterate over it to fill title and field
	request, err := api.GetPullRequests()
	if err == nil {
		if len(request.Values) == 0 {
			title = "**There are no active pull requests!**"
		} else {
			fieldTitles = make([]string, len(request.Values))
			fieldBodies = make([]string, len(request.Values))
			title = "**Active pull requests:**\n"
			for i, values := range request.Values {
				fieldTitles[i] = strconv.Itoa(i+1) + ". " + values.Title
				fieldBodies[i] = "[*Reviewers:*](" + values.Links.Self[0].Href + ")\n"
				for j, reviewer := range values.Reviewers {
					fieldBodies[i] += strconv.Itoa(j+1) + ". [" + reviewer.User.DisplayName + "](" + reviewer.User.Links.Self[0].Href + ") "
					userid, present := cfg[reviewer.User.Name]
					if present {
						fieldBodies[i] += "<@" + userid + "> "
					}
					if reviewer.Approved {
						fieldBodies[i] += "APPROVED!\n"
					} else {
						fieldBodies[i] += "\n"
					}
				}
				fieldBodies[i] += "Comments: " + strconv.Itoa(values.Properties.CommentCount)
				if vip {
					fieldTitles[i] = UwUify(fieldTitles[i])
					fieldBodies[i] = UwUify(fieldBodies[i])
				}
			}
		}
	} else {
		title = "**Request returned no data!**"
		log.Println(err)
	}

	if vip {
		return MakeEmbed(UwUify(title), UwUify(body), fieldTitles, fieldBodies)
	}
	// Add title, body and fields to an embed and return it
	return MakeEmbed(title, body, fieldTitles, fieldBodies)
}

// GetMyPullRequests returns only the pull requests opened by the requesting user
func GetMyPullRequests(api *API, rid string, vip bool) *discordgo.MessageEmbed {
	var (
		title string
		body  string
	)

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
						body += strconv.Itoa(i+1) + ". [" + values.Title + "](" + values.Links.Self[0].Href + ")\n Reviewers:\n"
						for j, reviewer := range values.Reviewers {
							body += strconv.Itoa(j+1) + ". [" + reviewer.User.DisplayName + "](" + reviewer.User.Links.Self[0].Href + ") "
							userid, present := cfg[reviewer.User.Name]
							if present {
								body += "<@" + userid + "> "
							}
							if reviewer.Approved {
								body += "APPROVED!\n"
							} else {
								body += "\n"
							}
						}
						body += "Comments: " + strconv.Itoa(values.Properties.CommentCount)
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

	if vip {
		return MakeEmbed(UwUify(title), UwUify(body), nil, nil)
	}

	// Add title and body to an embed and return it
	return MakeEmbed(title, body, nil, nil)
}

// GetMyReviews returns all pull requests that the message requester is a reviewer of
func GetMyReviews(api *API, rid string, vip bool) *discordgo.MessageEmbed {
	var (
		title string
		body  string
	)

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
							body += strconv.Itoa(i+1) + ". [" + values.Title + "](" + values.Links.Self[0].Href + ")\n"
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

	if vip {
		return MakeEmbed(UwUify(title), UwUify(body), nil, nil)
	}

	// Add title and body to an embed and return it
	return MakeEmbed(title, body, nil, nil)
}

// This is considered WIP, don't use
func GetActiveSprintMessage(api *API, vip bool) *discordgo.MessageEmbed {
	var (
		title string
		body  string
	)

	return MakeEmbed(title, body, nil, nil)
}
