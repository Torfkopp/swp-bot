package main

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

// Request implements the basic Request function
func (api *API) Request(request *http.Request) ([]byte, error) {
	// Add necessary header
	request.Header.Add("Accept", "application/json, */*")

	// Authenticate with endpoint
	api.Auth(request)

	Debug("====== Request ======")
	Debug(request)
	Debug("====== Request Body ======")
	if debugFlag {
		requestDump, err := httputil.DumpRequest(request, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(requestDump))
	}
	Debug("====== /Request Body ======")
	Debug("====== /Request ======")

	// Actually post the request
	response, err := api.client.Do(request)
	if err != nil {
		return nil, err
	}
	Debug(fmt.Sprintf("====== Response Status Code: %d ======", response.StatusCode))

	// Extract the response body
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Close resource now
	err = response.Body.Close()
	if err != nil {
		return nil, err
	}

	Debug("====== Response Body ======")
	Debug(string(responseBody))
	Debug("====== /Response Body ======")

	// Check which http status code we got and act on it
	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusPartialContent:
		return responseBody, nil
	case http.StatusNoContent, http.StatusResetContent:
		return nil, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("authentication failed")
	case http.StatusServiceUnavailable:
		return nil, fmt.Errorf("service is not available: %s", response.Status)
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("internal server error: %s", response.Status)
	case http.StatusConflict:
		return nil, fmt.Errorf("conflict: %s", response.Status)
	default:
		return nil, fmt.Errorf("unknown response status: %s", response.Status)
	}
}

// GetPullRequests sends Bitbucket GET requests
func (api *API) GetPullRequests() (*Response, error) {

	// Craft a GET request here
	request, err := http.NewRequest("GET", api.endPoint.String(), nil)
	if err != nil {
		return nil, err
	}

	// Send the request and save the response
	response, err := api.Request(request)
	if err != nil {
		return nil, err
	}

	// Parse the response from json into structs
	var parsed Response
	err = json.Unmarshal(response, &parsed)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

// This is considered WIP, don't use
func (api *API) GetActiveSprint() (*Response, error) {

	// Craft a GET request here
	request, err := http.NewRequest("GET", api.endPoint.String(), nil)
	if err != nil {
		return nil, err
	}

	// Send the request and save the response
	response, err := api.Request(request)
	if err != nil {
		return nil, err
	}

	var parser fastjson.Parser

	values, err := parser.ParseBytes(response)
	if err != nil {
		log.Println(err)
	}

	sprintID := values.GetInt("values", string(rune(len(values.GetArray("values"))-1)), "id")
	fmt.Println(sprintID)

	var parsed Response
	return &parsed, nil
}
