package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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

	a.client = &http.Client{Transport: tr, Timeout: time.Minute}

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

	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

// CheckNewPullRequest compares the date of the latest pull request with an internal variable
func CheckNewPullRequest(a *API) bool {
	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if req != nil && len(req.Values) == 0 {
			return false
		} else {
			ts, err := ioutil.ReadFile(timestamp)
			if err != nil {
				fmt.Println(err)
			}

			tss := strings.TrimSuffix(string(ts), "\n")

			n, err := strconv.ParseInt(tss, 10, 64)
			if err != nil {
				fmt.Println(err)
			}

			if req.Values[0].CreatedDate > n {
				n = req.Values[0].CreatedDate

				f, err := os.OpenFile(timestamp, os.O_WRONLY, 0600)
				if err != nil {
					fmt.Println(err)
				}
				defer f.Close()

				_, err = f.WriteString(strconv.FormatInt(n, 10))
				if err != nil {
					fmt.Println(err)
				}

				return true
			} else {
				return false
			}
		}
	} else {
		fmt.Println(err)
		return false
	}
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

// DebugFlag is the global debugging variable
var DebugFlag = false

// Debug outputs debug messages
func Debug(msg interface{}) {
	if DebugFlag {
		fmt.Printf("%+v\n", msg)
	}
}
