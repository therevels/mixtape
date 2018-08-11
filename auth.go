package main

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/zmb3/spotify"
)

// Login a user by redirecting them to Spotify to initiate the authorization
// code flow. It depends on the following environment variables having been
// set:
//
// SPOTIFY_ID - the oauth2 client ID
func Login(ctx echo.Context) error {
	req := ctx.Request()

	redirectURI := &url.URL{
		Scheme: "https",
		Host:   req.Host,
		Path:   "auth/callback",
	}

	auth := spotify.NewAuthenticator(redirectURI.String(), spotify.ScopeUserReadPrivate)

	// TODO this is supposed to be random and unique, and stored somewhere (like a session
	// cookie so we can implement CSRF protection for auth by checking the state value when
	// Spotify redirects back to us)
	state := strconv.Itoa(int(time.Now().UnixNano()))

	authURL := auth.AuthURL(state)
	return ctx.Redirect(http.StatusFound, authURL)
}
