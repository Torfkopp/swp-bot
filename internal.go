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
func NewAPI(location string, token string, tokenType int) (*API, error) {
	// Check if url isn't empty
	if len(location) == 0 {
		return nil, errors.New("url empty")
	}

	// Parse URL
	endPoint, err := url.ParseRequestURI(location)
	if err != nil {
		return nil, err
	}

	// Create new API object
	api := new(API)
	api.endPoint = endPoint
	api.tokenType = tokenType
	api.token = token

	// Make sure we use a valid and secure connection
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}

	// Set up http connection with a reasonable timeout
	api.client = &http.Client{Transport: transport, Timeout: time.Minute}

	return api, nil
}

// Auth implements token auth
func (api *API) Auth(request *http.Request) {
	// Supports unauthenticated access as well:
	// If token is not set, no authorization header is added
	if api.tokenType == 1 && api.token != "" {
		request.Header.Set("Authorization", "Bearer "+api.token)
	}
	if api.tokenType == 2 && api.token != "" {
		request.Header.Set("Authorization", "Basic "+api.token)
		request.Header.Set("Content-Type", "application/json")
	}
}

// ReadConfig reads the provided config file and turns it into a map
func ReadConfig() map[string]string {
	// Open config file
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Defer file closure so it runs after the return
	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Decode config file from json
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

// CheckNewPullRequest compares the date of the latest pull request with an internal variable
func CheckNewPullRequest(api *API) bool {
	// Craft a GET request
	request, err := api.GetActivePullRequests()
	// Only run if no error occurred yet
	if err == nil {
		// Return false if there are open PRs
		if request != nil && len(request.Values) == 0 {
			return false
		} else {
			// Read contents of the timestamp file
			timestamp, err := ioutil.ReadFile(timestampFile)
			if err != nil {
				log.Println(err)
			}

			// Trim the newline character to avoid trouble later
			timestampCleaned := strings.TrimSuffix(string(timestamp), "\n")

			// Parse the timestamp into a variable
			n, err := strconv.ParseInt(timestampCleaned, 10, 64)
			if err != nil {
				log.Println(err)
			}

			// Compare the timestamp of the latest PR with our previous timestamp
			if request.Values[0].CreatedDate > n {
				n = request.Values[0].CreatedDate

				// Open the timestamp file for writing here
				file, err := os.OpenFile(timestampFile, os.O_WRONLY, 0600)
				if err != nil {
					fmt.Println(err)
				}

				// Defer file closure so it runs after the return
				defer func() {
					err = file.Close()
					if err != nil {
						log.Fatal(err)
					}
				}()

				// Overwrite the timestamp in our file with the latest recorded one
				_, err = file.WriteString(strconv.FormatInt(n, 10))
				if err != nil {
					log.Println(err)
				}

				return true
			} else {
				return false
			}
		}
	} else {
		log.Println(err)
		return false
	}
}

// Debug outputs debug messages
func Debug(msg interface{}) {
	if debugFlag {
		log.Printf("%+v\n", msg)
	}
}

func JiraTest(api *API) {
	request, err := api.GetActivePullRequests()
	if err != nil {
		fmt.Println(err)
	}

	if request.Size == 1 {
		fmt.Println("wow")
	}
}
