package main

import (
	"net/http"
	"net/url"
)

// API is the main api data structure
type API struct {
	endPoint *url.URL
	client   *http.Client
	token    string
}

// Response defines the information sent via the REST-API
type Response struct {
	Size       int      `json:"size"`
	Limit      int      `json:"limit"`
	Start      int      `json:"start"`
	IsLastPage bool     `json:"isLastPage"`
	Values     []Values `json:"values"`
}

// Values defines all relevant values inside the response
type Values struct {
	Open        bool       `json:"open"`
	CreatedDate int64      `json:"createdDate"`
	Author      Author     `json:"author"`
	UpdatedDate int64      `json:"updatedDate"`
	Description string     `json:"description"`
	Reviewers   []Reviewer `json:"reviewers"`
	Title       string     `json:"title"`
	//Links       []Links    `json:"links"`
	ID int `json:"id"`
}

// Author defines author information
type Author struct {
	User     User `json:"user"`
	Approved bool `json:"approved"`
}

// Reviewer defines reviewer information
type Reviewer struct {
	User     User `json:"user"`
	Approved bool `json:"approved"`
}

// User defines user information
type User struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Mail        string `json:"emailAddress"`
	DisplayName string `json:"displayName"`
	//Links       []Links `json:"links"`
}

// TODO This struct is probably incorrect as in includes an array inside the array
// Links defines hyperlink information
type Links struct {
	Self []string `json:"self"`
}
