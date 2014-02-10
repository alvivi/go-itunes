// Copyright 2013 - Alvaro Vilanova. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package itunes

import (
	"net/http"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if c.client != http.DefaultClient {
		t.Errorf("NewClient default client is no http.DefaultClient")
	}

	if c.BaseURL.String() != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL.String(), defaultBaseURL)
	}

	if c.UserAgent != userAgent {
		t.Errorf("NewClient UserAgent = %v, want %v", c.UserAgent, userAgent)
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)
	inURL, inSlashURL, outURL := "foo", "/foo", defaultBaseURL+"foo"

	// test that relative URL is correct
	req, _ := c.NewRequest("GET", inURL)
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %s, want %v", inURL, req.URL, outURL)
	}

	// test that absolute URL is the same as relative URL
	req, _ = c.NewRequest("GET", inSlashURL)
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %s, want %v", inSlashURL, req.URL, outURL)
	}

	// test that default user-agent is attached to the request
	userAgent := req.Header.Get("User-Agent")
	if c.UserAgent != userAgent {
		t.Errorf("NewRequest() User-Agent = %v, want %v", userAgent, c.UserAgent)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient(nil)
	_, err := c.NewRequest("GET", ":")
	if err == nil {
		t.Errorf("NewRequest(:) Expected an error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("NewRequest(:) Expected a parse error, got %+v", err)
	}
}
