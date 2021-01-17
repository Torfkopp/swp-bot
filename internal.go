package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

func GetPullRequests(a *API) string {
	var ret string

	req, err := a.GetPullRequestsRequest()
	if err == nil {
		if len(req.Values) == 0 {
			ret = "No active pull requests!"
		} else {
			for _, value := range req.Values {
				ret = ret + value.Title + ""
			}
		}
	} else {
		ret = "Request returned no data!"
		panic(err)
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
