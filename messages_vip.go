package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/clinet/discordgo-embed"
	"strconv"
)

// HelpMessageVIP returns an embed containing an info text about the supported commands (VIP version)
func HelpMessageVIP() *discordgo.MessageEmbed {
	return embed.NewGenericEmbedAdvanced("__These awe the x3 suppowted intewactive commands:__",
		"**!allpullrequests:** Shows the UwU status of aww active puww wequests.\n"+
			"**!mypullrequests:** Shows the x3 status of youw own active puww wequests.\n"+
			"**!myreviews:**  Shows aww puww wequests which you'we a weviewew of, nya.\n"+
			"**!comments:** Shows the x3 comments >w< undew youw active puww *boops your nose* wequests. *(TODO)*\n"+
			"**!about:** Some info *sweats* about this bot.",
		color)
}

// AboutMessageVIP returns an embed containing an "about" text (VIP version)
func AboutMessageVIP() *discordgo.MessageEmbed {
	return embed.NewGenericEmbedAdvanced("About this owo bot:",
		"In case of undesiwed wisks and side effects\n"+
			"pwease wead the x3 [souwce code](https://github.com/MDr164/swp-bot) ow (・`ω´・) ask youw wocaw dev.",
		color)
}

// GetAllPullRequestsVIP returns all currently active pull requests from the rest response (VIP version)
func GetAllPullRequestsVIP(a *API) *discordgo.MessageEmbed {
	var (
		t  string
		d  string
		fd string
	)

	e := embed.NewEmbed().SetColor(color)

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			t = "**Thewe awe nyo active puww wequests!** *huggles tightly*"
		} else {
			t = "**Active puww wequests:**\n"
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
				e.AddField(strconv.Itoa(i+1)+". "+val.Title, "[*Wweviewews:*]("+val.Links.Self[0].Href+")\n"+fd)
			}
		}
	} else {
		t = "**Wequest wetuwnyed nyo data?!?!**"
		fmt.Println(err)
	}

	e.SetTitle(t).SetDescription(d)
	return e.MessageEmbed
}

// GetMyPullRequestsVIP returns only the pull requests opened by the requesting user (VIP version)
func GetMyPullRequestsVIP(a *API, rid string) *discordgo.MessageEmbed {
	var (
		t string
		d string
	)

	e := embed.NewEmbed().SetColor(color)

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			t = "**Thewe awe nyo active puww wequests!** *huggles tightly*"
		} else {
			username, _ := cfg[rid]
			t = "**Puww wequests by " + username + ":**\n"
			i := 0
			for _, val := range req.Values {
				if val.Author.User.Name == username {
					d = d + strconv.Itoa(i+1) + ". [" + val.Title + "](" + val.Links.Self[0].Href + ")\n"
					i++
				}
			}
			if d == "" {
				d = "*Nyonye :3*"
			}
		}
	} else {
		t = "**Wequest wetuwnyed nyo data?!?!**"
		fmt.Println(err)
	}

	e.SetTitle(t).SetDescription(d)
	return e.MessageEmbed
}

// GetMyReviewsVIP returns all pull requests that the message requester is a reviewer of (VIP version)
func GetMyReviewsVIP(a *API, rid string) *discordgo.MessageEmbed {
	var (
		t string
		d string
	)

	e := embed.NewEmbed().SetColor(color)

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			t = "**Thewe awe nyo active puww wequests!** *huggles tightly*"
		} else {
			username, _ := cfg[rid]
			t = "**W-W-Weviews assignyed t-to " + username + ":**\n"
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
				d = "*Nyonye :3*"
			}
		}
	} else {
		t = "**Wequest wetuwnyed nyo data?!?!"
		fmt.Println(err)
	}

	e.SetTitle(t).SetDescription(d)
	return e.MessageEmbed
}
