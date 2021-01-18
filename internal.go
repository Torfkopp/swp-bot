package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"net/url"
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

// TODO Implement me
func ListenForNewRequests(a *API) string {
	ret := ">>> "

	return ret
}

func GetAllPullRequests(a *API) string {
	ret := ">>> "

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			ret = ret + "**No active pull requests!**"
		} else {
			for i, val := range req.Values {
				ret = ret + "**" + strconv.Itoa(i+1) + ". " + val.Title + "**" + "\n*Reviewers:*\n"
				for j, rev := range val.Reviewers {
					ret = ret + strconv.Itoa(j+1) + ". " + rev.User.DisplayName + "\n"
				}
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

// SetDebug enables debug output
func SetDebug(state bool) {
	DebugFlag = state
}

// Debug outputs debug messages
func Debug(msg interface{}) {
	if DebugFlag {
		fmt.Printf("%+v\n", msg)
	}
}
