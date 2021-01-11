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

func GetReviewers(a *API) string {
	var ret string

	req, err := a.GetReviewerRequest()
	if err == nil {
		if len(req.Reviewer) == 0 {
			ret = "No Reviewers assigned!"
		} else {
			for _, User := range req.Reviewer {
				ret = ret + User.User.Username + ""
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
