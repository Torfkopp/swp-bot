package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

// Request implements the basic Request function
func (a *API) Request(req *http.Request) ([]byte, error) {
	req.Header.Add("Accept", "application/json, */*")

	a.Auth(req)

	Debug("====== Request ======")
	Debug(req)
	Debug("====== Request Body ======")
	if DebugFlag {
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(requestDump))
	}
	Debug("====== /Request Body ======")
	Debug("====== /Request ======")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	Debug(fmt.Sprintf("====== Response Status Code: %d ======", resp.StatusCode))

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	Debug("====== Response Body ======")
	Debug(string(res))
	Debug("====== /Response Body ======")

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusPartialContent:
		return res, nil
	case http.StatusNoContent, http.StatusResetContent:
		return nil, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("authentication failed")
	case http.StatusServiceUnavailable:
		return nil, fmt.Errorf("service is not available: %s", resp.Status)
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("internal server error: %s", resp.Status)
	case http.StatusConflict:
		return nil, fmt.Errorf("conflict: %s", resp.Status)
	}

	return nil, fmt.Errorf("unknown response status: %s", resp.Status)
}

// GetReviewerRequest sends reviewer related requests
func (a *API) GetReviewerRequest() (*Reviewers, error) {

	req, err := http.NewRequest("GET", a.endPoint.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := a.Request(req)
	if err != nil {
		return nil, err
	}

	var reviewers Reviewers

	err = json.Unmarshal(res, &reviewers)
	if err != nil {
		return nil, err
	}

	return &reviewers, nil
}
