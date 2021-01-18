package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
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
func GetAllPullRequests(a *API) string {
	ret := ">>> "
	lut := ReadLUT()

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			ret = ret + "**There are no active pull requests!**"
		} else {
			for i, val := range req.Values {
				ret = ret + "**" + strconv.Itoa(i+1) + ". " + val.Title + "**" + "\n*Reviewers:*\n"
				for j, rev := range val.Reviewers {
					ret = ret + strconv.Itoa(j+1) + ". " + rev.User.DisplayName
					userid, present := lut[rev.User.Name]
					if present {
						ret = ret + " <@" + userid + ">\n"
					} else {
						ret = ret + "\n"
					}
				}
			}
		}
	} else {
		ret = ret + "**Request returned no data!**"
		fmt.Println(err)
	}

	return ret
}

// GetMyPullRequests returns only the pull requests opened by the requesting user
func GetMyPullRequests(a *API, requesterid string) string {
	ret := ">>> "
	lut := ReadLUT()

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			ret = ret + "**There are no active pull requests!**"
		} else {
			username, present := lut[requesterid]
			if present {
				ret = ret + "**Pull requests by " + username + ":**\n"
				for i, val := range req.Values {
					if val.Author.User.Name == username {
						ret = ret + strconv.Itoa(i+1) + ". " + val.Title + "\n"
					}
				}
			} else {
				ret = ret + "*Couldn't map your Discord ID to a Bitbucket user!*"
			}
		}
	} else {
		ret = ret + "**Request returned no data!**"
		fmt.Println(err)
	}

	return ret
}

// GetMyReviews returns all pull requests that the message requester is a reviewer of
func GetMyReviews(a *API, requesterid string) string {
	ret := ">>> "
	lut := ReadLUT()

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			ret = ret + "**There are no active pull requests!**"
		} else {
			username, present := lut[requesterid]
			if present {
				ret = ret + "**Reviews assigned to " + username + ":**\n"
				for i, val := range req.Values {
					for _, rev := range val.Reviewers {
						if rev.User.Name == username {
							ret = ret + strconv.Itoa(i+1) + ". " + val.Title + "\n"
						}
					}
				}
			} else {
				ret = ret + "*Couldn't map your Discord ID to a Bitbucket user!*"
			}
		}
	} else {
		ret = ret + "**Request returned no data!**"
		fmt.Println(err)
	}

	return ret
}

// TODO This doesn't work yet
func FormatMessage(a *API) *discordgo.MessageEmbed {
	var ret *discordgo.MessageEmbed

	ret.Color = 5

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			ret.Fields[0].Value = "No active pull requests!"
		} else {
			for i, val := range req.Values {
				ret.Fields[i].Name = "Pull-Request:"
				ret.Fields[i].Value = req.Values[i].Title
				ret.Fields[i+1].Name = "Reviewers:"
				for _, rev := range val.Reviewers {
					ret.Fields[i+1].Value = rev.User.DisplayName
					i++
				}
				i++
			}
		}
	} else {
		ret.Fields[0].Value = "Request returned no data!"
		fmt.Println(err)
	}

	return ret
}

// DebugFlag is the global debugging variable
var DebugFlag = false

// Debug outputs debug messages
func Debug(msg interface{}) {
	if DebugFlag {
		fmt.Printf("%+v\n", msg)
	}
}
