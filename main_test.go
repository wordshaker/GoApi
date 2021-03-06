package main

import (
	"golang.org/x/oauth2"
	"testing"
	"net/http/httptest"
	"net/http"
)

func TestThing(t *testing.T) {
	conf := &oauth2.Config{
		ClientID:     "id",
		ClientSecret: "secret",
		Scopes:       []string{"user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	s := httptest.NewServer(oauthServer(conf))
	defer s.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, _	:= http.NewRequest("GET", s.URL + "/login", nil)
	resp, _ := client.Do(req)

	if resp.StatusCode != http.StatusSeeOther {
		t.Fail()
	}
}