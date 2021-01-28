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

// ReadLUT reads the provided look up table and turns it into a map
func ReadLUT() map[string]string {
	file, err := os.Open(UserLUT)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var lut map[string]string
	err = json.NewDecoder(file).Decode(&lut)
	if err != nil {
		log.Fatal(err)
	}

	return lut
}

// TODO Implement me
func ListenForNewRequests(a *API) string {
	ret := ">>> "

	return ret
}

// GetAllPullRequests returns all currently active pull requests from the rest response
func GetAllPullRequests(a *API) *discordgo.MessageEmbed {
	var (
		t   string
		d   string
		u   string
		lut = ReadLUT()
	)

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			t = "**There are no active pull requests!**"
		} else {
			for i, val := range req.Values {
				d = d + "**" + strconv.Itoa(i+1) + ". " + val.Title + "**" + "\n*Reviewers:*\n"
				for j, rev := range val.Reviewers {
					d = d + strconv.Itoa(j+1) + ". " + rev.User.DisplayName
					userid, present := lut[rev.User.Name]
					if present {
						d = d + " <@" + userid + ">\n"
					} else {
						d = d + "\n"
					}
				}
			}
		}
	} else {
		t = "**Request returned no data!**"
		fmt.Println(err)
	}

	e := embed.NewEmbed().SetTitle(t).SetColor(color).SetDescription(d).SetURL(u).MessageEmbed
	return e
}

// GetMyPullRequests returns only the pull requests opened by the requesting user
func GetMyPullRequests(a *API, rid string) *discordgo.MessageEmbed {
	var (
		t   string
		d   string
		u   string
		lut = ReadLUT()
	)

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			t = "**There are no active pull requests!**"
		} else {
			username, present := lut[rid]
			if present {
				d = "**Pull requests by " + username + ":**\n"
				for i, val := range req.Values {
					if val.Author.User.Name == username {
						d = d + strconv.Itoa(i+1) + ". " + val.Title + "\n"
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

	e := embed.NewEmbed().SetTitle(t).SetColor(color).SetDescription(d).SetURL(u).MessageEmbed
	return e
}

// GetMyReviews returns all pull requests that the message requester is a reviewer of
func GetMyReviews(a *API, rid string) *discordgo.MessageEmbed {
	var (
		t   string
		d   string
		u   string
		lut = ReadLUT()
	)

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			t = "**There are no active pull requests!**"
		} else {
			username, present := lut[rid]
			if present {
				d = "**Reviews assigned to " + username + ":**\n"
				for i, val := range req.Values {
					for _, rev := range val.Reviewers {
						if rev.User.Name == username {
							d = d + strconv.Itoa(i+1) + ". " + val.Title + "\n"
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

	e := embed.NewEmbed().SetTitle(t).SetColor(color).SetDescription(d).SetURL(u).MessageEmbed
	return e
}

// GetAllPullRequestsVIP returns all currently active pull requests from the rest response (VIP version)
func GetAllPullRequestsVIP(a *API) *discordgo.MessageEmbed {
	var (
		t   string
		d   string
		u   string
		lut = ReadLUT()
	)

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			t = "**Thewe awe nyo active puww wequests!** *huggles tightly*"
		} else {
			for i, val := range req.Values {
				d = d + "**" + strconv.Itoa(i+1) + ". " + val.Title + "**" + "\n*Wweviewews:*\n"
				for j, rev := range val.Reviewers {
					d = d + strconv.Itoa(j+1) + ". " + rev.User.DisplayName
					userid, present := lut[rev.User.Name]
					if present {
						d = d + " <@" + userid + ">\n"
					} else {
						d = d + "\n"
					}
				}
			}
		}
	} else {
		t = "**Wequest wetuwnyed nyo data?!?!**"
		fmt.Println(err)
	}

	e := embed.NewEmbed().SetTitle(t).SetColor(color).SetDescription(d).SetURL(u).MessageEmbed
	return e
}

// GetMyPullRequestsVIP returns only the pull requests opened by the requesting user (VIP version)
func GetMyPullRequestsVIP(a *API, rid string) *discordgo.MessageEmbed {
	var (
		t   string
		d   string
		u   string
		lut = ReadLUT()
	)

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			t = "**Thewe awe nyo active puww wequests!** *huggles tightly*"
		} else {
			username, _ := lut[rid]
			d = "**Puww wequests by " + username + ":**\n"
			for i, val := range req.Values {
				if val.Author.User.Name == username {
					d = d + strconv.Itoa(i+1) + ". " + val.Title + "\n"
				}
			}
		}
	} else {
		t = "**Wequest wetuwnyed nyo data?!?!**"
		fmt.Println(err)
	}

	e := embed.NewEmbed().SetTitle(t).SetColor(color).SetDescription(d).SetURL(u).MessageEmbed
	return e
}

// GetMyReviewsVIP returns all pull requests that the message requester is a reviewer of (VIP version)
func GetMyReviewsVIP(a *API, rid string) *discordgo.MessageEmbed {
	var (
		t   string
		d   string
		u   string
		lut = ReadLUT()
	)

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			t = "**Thewe awe nyo active puww wequests!** *huggles tightly*"
		} else {
			username, _ := lut[rid]
			d = "**W-W-Weviews assignyed t-to " + username + ":**\n"
			for i, val := range req.Values {
				for _, rev := range val.Reviewers {
					if rev.User.Name == username {
						d = d + strconv.Itoa(i+1) + ". " + val.Title + "\n"
					}
				}
			}
		}
	} else {
		t = "**Wequest wetuwnyed nyo data?!?!**"
		fmt.Println(err)
	}

	e := embed.NewEmbed().SetTitle(t).SetColor(color).SetDescription(d).SetURL(u).MessageEmbed
	return e
}

// DebugFlag is the global debugging variable
var DebugFlag = false

// Debug outputs debug messages
func Debug(msg interface{}) {
	if DebugFlag {
		fmt.Printf("%+v\n", msg)
	}
}
