package main

import (
	"fmt"
	"net/http"
	"os"
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
// REDIRECT_URI - the exact URI (including protocol) for spotify to redirect
//                back to after the user has authenticated and authorized
//                access
func Login(w http.ResponseWriter, r *http.Request) {
	redirectURI := os.Getenv("REDIRECT_URI")
	auth := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate)

	// TODO this is supposed to be random and unique
	state := strconv.Itoa(int(time.Now().UnixNano()))

	url := auth.AuthURL(state)
	fmt.Fprintf(w, "url: %v", url)

	// TODO redirect to the url
}
