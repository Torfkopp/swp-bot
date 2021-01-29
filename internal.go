package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// NewAPI implements API constructor
func NewAPI(location string, token string) (*API, error) {
	if len(location) == 0 {
		return nil, errors.New("url empty")
	}

	u, err := url.ParseRequestURI(location)

	if err != nil {
		return nil, err
	}

	a := new(API)
	a.endPoint = u
	a.token = token

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}

	a.client = &http.Client{Transport: tr, Timeout: 10 * time.Second}

	return a, nil
}

// Auth implements token auth
func (a *API) Auth(req *http.Request) {
	// Supports unauthenticated access as well:
	// If token is not set, no authorization header is added
	if a.token != "" {
		req.Header.Set("Authorization", "Bearer "+a.token)
	}
}

// ReadConfig reads the provided config file and turns it into a map
func ReadConfig() map[string]string {
	file, err := os.Open(config)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var cfg map[string]string
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

// TODO Implement me
func NewPullRequestCreated(a *API) *discordgo.MessageEmbed {
	var (
		t string
		d string
	)

	e := embed.NewEmbed().SetColor(color).SetTitle(t).SetDescription(d)
	return e.MessageEmbed
}

// GetAllPullRequests returns all currently active pull requests from the rest response
func GetAllPullRequests(a *API) *discordgo.MessageEmbed {
	var (
		t   string
		d   string
		fd  string
		cfg = ReadConfig()
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
		t   string
		d   string
		cfg = ReadConfig()
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
		t   string
		d   string
		cfg = ReadConfig()
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

// GetAllPullRequestsVIP returns all currently active pull requests from the rest response (VIP version)
func GetAllPullRequestsVIP(a *API) *discordgo.MessageEmbed {
	var (
		t   string
		d   string
		fd  string
		cfg = ReadConfig()
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
		t   string
		d   string
		cfg = ReadConfig()
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
		t   string
		d   string
		cfg = ReadConfig()
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
		}
	} else {
		t = "**Wequest wetuwnyed nyo data?!?!"
		fmt.Println(err)
	}

	e.SetTitle(t).SetDescription(d)
	return e.MessageEmbed
}

// DebugFlag is the global debugging variable
var DebugFlag = false

// Debug outputs debug messages
func Debug(msg interface{}) {
	if DebugFlag {
		fmt.Printf("%+v\n", msg)
	}
}
