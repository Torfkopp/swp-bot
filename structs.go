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

// User defines user information
type User struct {
	Type        string `json:"type"`
	Username    string `json:"name"`
	Mail        string `json:"emailAddress"`
	DisplayName string `json:"displayName"`
}

type PullRequest struct {
	Request []Request `json:"values"`
	Start   int       `json:"start,omitempty"`
	Limit   int       `json:"limit,omitempty"`
	Size    int       `json:"size,omitempty"`
}

type Request struct {
	Title string `json:"title"`
}

type Reviewer struct {
	User     []User `json:"user"`
	Approved bool   `json:"approved"`
	Status   string `json:"status"`
}

type Reviewers struct {
	Reviewer []Reviewer
}
