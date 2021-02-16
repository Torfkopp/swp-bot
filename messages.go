package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/clinet/discordgo-embed"
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
func NewPullRequestCreated(a *API) *discordgo.MessageEmbed {
	var (
		t string
		d string
	)

	e := embed.NewEmbed().SetColor(color)

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		t = "**New pull request:**\n"
		d = "[" + req.Values[0].Title + "](" + req.Values[0].Links.Self[0].Href + ")"
	} else {
		t = "**Request returned no data!**"
		fmt.Println(err)
	}

	e.SetTitle(t).SetDescription(d)
	return e.MessageEmbed
}

// NewPullRequestPing returns a string containing the reviewers of the latest pull request
func NewPullRequestPing(a *API) string {
	var t string

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		t = "**Pinging Reviewers:**\n"
		for i, rev := range req.Values[0].Reviewers {
			t = t + strconv.Itoa(i+1) + ". " + rev.User.DisplayName
			userid, present := cfg[rev.User.Name]
			if present {
				t = t + " <@" + userid + ">\n"
			} else {
				t = t + "\n"
			}
		}
	} else {
		t = "**Request returned no data!**"
		fmt.Println(err)
	}

	return t
}

// GetAllPullRequests returns all currently active pull requests from the rest response
func GetAllPullRequests(a *API) *discordgo.MessageEmbed {
	var (
		t  string
		d  string
		fd string
	)

	e := embed.NewEmbed().SetColor(color)

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			t = "**There are no active pull requests!**"
		} else {
			t = "**Active pull requests:**\n"
			for i, val := range req.Values {
				fd = ""
				for j, rev := range val.Reviewers {
					fd = fd + strconv.Itoa(j+1) + ". [" + rev.User.DisplayName + "](" + rev.User.Links.Self[0].Href + ") "
					userid, present := cfg[rev.User.Name]
					if present {
						fd = fd + "<@" + userid + ">\n"
					} else {
						fd = fd + "\n"
					}
				}
				e.AddField(strconv.Itoa(i+1)+". "+val.Title, "[*Reviewers:*]("+val.Links.Self[0].Href+")\n"+fd)
			}
		}
	} else {
		t = "**Request returned no data!**"
		fmt.Println(err)
	}

	e.SetTitle(t).SetDescription(d)
	return e.MessageEmbed
}

// GetMyPullRequests returns only the pull requests opened by the requesting user
func GetMyPullRequests(a *API, rid string) *discordgo.MessageEmbed {
	var (
		t string
		d string
	)

	e := embed.NewEmbed().SetColor(color)

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			t = "**There are no active pull requests!**"
		} else {
			username, present := cfg[rid]
			if present {
				t = "**Pull requests by " + username + ":**\n"
				i := 0
				for _, val := range req.Values {
					if val.Author.User.Name == username {
						d = d + strconv.Itoa(i+1) + ". [" + val.Title + "](" + val.Links.Self[0].Href + ")\n"
						i++
					}
				}
				if d == "" {
					d = "*None*"
				}
			} else {
				t = "*Couldn't map your Discord ID to a Bitbucket user!*"
			}
		}
	} else {
		t = "**Request returned no data!**"
		fmt.Println(err)
	}

	e.SetTitle(t).SetDescription(d)
	return e.MessageEmbed
}

// GetMyReviews returns all pull requests that the message requester is a reviewer of
func GetMyReviews(a *API, rid string) *discordgo.MessageEmbed {
	var (
		t string
		d string
	)

	e := embed.NewEmbed().SetColor(color)

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			t = "**There are no active pull requests!**"
		} else {
			username, present := cfg[rid]
			if present {
				t = "**Reviews assigned to " + username + ":**\n"
				i := 0
				for _, val := range req.Values {
					for _, rev := range val.Reviewers {
						if rev.User.Name == username {
							d = d + strconv.Itoa(i+1) + ". [" + val.Title + "](" + val.Links.Self[0].Href + ")\n"
							i++
						}
					}
				}
				if d == "" {
					d = "*None*"
				}
			} else {
				t = "*Couldn't map your Discord ID to a Bitbucket user!*"
			}
		}
	} else {
		t = "**Request returned no data!**"
		fmt.Println(err)
	}

	e.SetTitle(t).SetDescription(d)
	return e.MessageEmbed
}
