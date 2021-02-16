package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
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

// DebugFlag is the global debugging variable
var DebugFlag = false

// Debug outputs debug messages
func Debug(msg interface{}) {
	if DebugFlag {
		fmt.Printf("%+v\n", msg)
	}
}
