package main

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zmb3/spotify"
)

// Login a user by redirecting them to Spotify to initiate the authorization
// code flow. It depends on the following environment variables having been
// set:
//
// SPOTIFY_ID - the oauth2 client ID
// SPOTIFY_SECRET - the oauth2 client secret
func Login(w http.ResponseWriter, r *http.Request) {
	redirectURI := &url.URL{
		Scheme: r.URL.Scheme,
		Host:   r.URL.Host,
		Path:   "auth/callback",
	}

	auth := spotify.NewAuthenticator(redirectURI.String(), spotify.ScopeUserReadPrivate)

	// TODO this is supposed to be random and unique
	state := strconv.Itoa(int(time.Now().UnixNano()))

	authURL := auth.AuthURL(state)
	http.Redirect(w, r, authURL, 302)
}
