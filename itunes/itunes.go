// Copyright 2013 - Alvaro Vilanova. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package itunes

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"net/http"
	"net/url"
)

const (
	libraryVersion = "0.1"
	defaultBaseURL = "https://itunes.apple.com/"
	userAgent      = "go-itunes/" + libraryVersion
)

// A Client manages communication with the iTunes API. Use NewClient to create
// a new client.
type Client struct {
	BaseURL   *url.URL
	client    *http.Client
	UserAgent string
}

// NewClient returns a new iTunes API client. If a nil httpClient is
// provided, http.DefaultClient will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	c := Client{
		BaseURL:   baseURL,
		client:    httpClient,
		UserAgent: userAgent,
	}
	return &c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// iTunes API calls do not use the body of the request.
func (c Client) NewRequest(method, urlStr string) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

func (c Client) Search(args ...string) ([]simplejson.Json, error) {
	resp, err := c.doAffiliateRequest("search", args...)
	return resp, err
}

func (c Client) LookUp(args ...string) ([]simplejson.Json, error) {
	resp, err := c.doAffiliateRequest("lookup", args...)
	return resp, err
}

func (c Client) doAffiliateRequest(service string, args ...string) ([]simplejson.Json, error) {
	values := argsToValues(args...)
	req, err := c.NewRequest("GET", service+"?"+values.Encode())
	if err != nil {
		return []simplejson.Json{}, err
	}
	rawResp, err := c.client.Do(req)
	if err != nil {
		return []simplejson.Json{}, err
	}
	defer rawResp.Body.Close()
	var resp itunesResponse
	err = json.NewDecoder(rawResp.Body).Decode(&resp)
	if err != nil {
		return []simplejson.Json{}, err
	}
	return resp.Results, nil
}

type itunesResponse struct {
	Results      []simplejson.Json
	ResultsCount int
}

func argsToValues(args ...string) url.Values {
	values := url.Values{}
	currentKey := ""
	for _, arg := range args {
		if len(currentKey) <= 0 {
			currentKey = arg
			continue
		}
		values.Add(currentKey, arg)
		currentKey = ""
	}
	if len(currentKey) > 0 {
		values.Add(currentKey, "")
	}
	return values
}
