package main

import (
	"golang.org/x/oauth2"
	"context"
	"fmt"
	"log"
	"os"
	"io"
	"net/http"
)

func main2() {
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     "bab74605b741c3ba9aa8",
		ClientSecret: "39412d02bbbca3170efb83311a0a11b60cbcf176",
		Scopes:       []string{"user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}

func main() {
	conf := &oauth2.Config{
		ClientID:     "bab74605b741c3ba9aa8",
		ClientSecret: "39412d02bbbca3170efb83311a0a11b60cbcf176",
		Scopes:       []string{"user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	server := oauthServer(conf)
	log.Fatal(http.ListenAndServe(":8080", server))
}

func oauthServer(conf *oauth2.Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, geturl(conf), http.StatusSeeOther)
	})

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request){
		code := r.URL.Query().Get("code")
		tok, err := conf.Exchange(r.Context(), code)
		if err != nil {
			log.Fatal(err)
		}

		client := conf.Client(r.Context(), tok)
		resp, err := client.Get("https://api.github.com/user")
		if err != nil {
			log.Fatal(err)
		}
	
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
	})

	return mux
}

func geturl(conf *oauth2.Config) string{
	return conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

