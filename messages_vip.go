package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/clinet/discordgo-embed"
	"log"
	"strconv"
)

// HelpMessageVIP returns an embed containing an info text about the supported commands (VIP version)
func HelpMessageVIP() *discordgo.MessageEmbed {
	return embed.NewGenericEmbedAdvanced("__These awe the x3 suppowted intewactive commands:__",
		"**!allpullrequests:** Shows the UwU status of aww active puww wequests.\n"+
			"**!mypullrequests:** Shows the x3 status of youw own active puww wequests.\n"+
			"**!myreviews:**  Shows aww puww wequests which you'we a weviewew of, nya.\n"+
			"**!post <something>:** Weways youw text into the x3 bots *huggles tightly* channyew.\n"+
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
func GetAllPullRequestsVIP(api *API) *discordgo.MessageEmbed {
	var (
		title string
		body  string
		field string
	)

	embedObject := embed.NewEmbed().SetColor(color)

	request, err := api.GetPullRequests()
	if err == nil {
		if len(request.Values) == 0 {
			title = "**Thewe awe nyo active puww wequests!** *huggles tightly*"
		} else {
			title = "**Active puww wequests:**\n"
			for i, values := range request.Values {
				field = "[*Wweviewews:*](" + values.Links.Self[0].Href + ")\n"
				for j, reviewer := range values.Reviewers {
					field = field + strconv.Itoa(j+1) + ". [" + reviewer.User.DisplayName + "](" + reviewer.User.Links.Self[0].Href + ") "
					userid, present := cfg[reviewer.User.Name]
					if present {
						field = field + "<@" + userid + ">\n"
						if reviewer.Approved {
							field = field + "NYAAA!\n"
						} else {
							field = field + "\n"
						}
					}
				}
				field = field + "Tweets: " + strconv.Itoa(values.Properties.CommentCount)
				embedObject.AddField(strconv.Itoa(i+1)+". "+values.Title, field)
			}
		}
	} else {
		title = "**Wequest wetuwnyed nyo data?!?!**"
		log.Println(err)
	}

	embedObject.SetTitle(title).SetDescription(body)
	return embedObject.MessageEmbed
}

// GetMyPullRequestsVIP returns only the pull requests opened by the requesting user (VIP version)
func GetMyPullRequestsVIP(api *API, rid string) *discordgo.MessageEmbed {
	var (
		title string
		body  string
	)

	embedObject := embed.NewEmbed().SetColor(color)

	request, err := api.GetPullRequests()
	if err == nil {
		if len(request.Values) == 0 {
			title = "**Thewe awe nyo active puww wequests!** *huggles tightly*"
		} else {
			username, _ := cfg[rid]
			title = "**Puww wequests by " + username + ":**\n"
			i := 0
			for _, values := range request.Values {
				if values.Author.User.Name == username {
					body = body + strconv.Itoa(i+1) + ". [" + values.Title + "](" + values.Links.Self[0].Href + ")\n"
					i++
				}
			}
			if body == "" {
				body = "*Nyonye :3*"
			}
		}
	} else {
		title = "**Wequest wetuwnyed nyo data?!?!**"
		log.Println(err)
	}

	embedObject.SetTitle(title).SetDescription(body)
	return embedObject.MessageEmbed
}

// GetMyReviewsVIP returns all pull requests that the message requester is a reviewer of (VIP version)
func GetMyReviewsVIP(api *API, rid string) *discordgo.MessageEmbed {
	var (
		title string
		body  string
	)

	embedObject := embed.NewEmbed().SetColor(color)

	request, err := api.GetPullRequests()
	if err == nil {
		if len(request.Values) == 0 {
			title = "**Thewe awe nyo active puww wequests!** *huggles tightly*"
		} else {
			username, _ := cfg[rid]
			title = "**W-W-Weviews assignyed t-to " + username + ":**\n"
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
				body = "*Nyonye :3*"
			}
		}
	} else {
		title = "**Wequest wetuwnyed nyo data?!?!"
		log.Println(err)
	}

	embedObject.SetTitle(title).SetDescription(body)
	return embedObject.MessageEmbed
}
