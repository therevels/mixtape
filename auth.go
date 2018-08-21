package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/zmb3/spotify"
)

// Login a user by redirecting them to Spotify to initiate the authorization
// code flow. It depends on the following environment variables having been
// set:
//
// SPOTIFY_ID - the oauth2 client ID
// SPOTIFY_SECRET - the oauth2 client secret
func Login(ctx echo.Context) error {
	auth := newAuthenticator(ctx)

	state, err := newState()
	if err != nil {
		return err
	}

	sess, err := session.Get("mixtape-session", ctx)
	if err != nil {
		return err
	}
	sess.Values["auth_state"] = state
	sess.Save(ctx.Request(), ctx.Response())

	authURL := auth.AuthURL(state)
	return ctx.Redirect(http.StatusFound, authURL)
}

// Callback after authentication/authorization is complete and the Spotify
// server redirects back to the redirect URI with an authorization code
func Callback(ctx echo.Context) error {
	sess, err := session.Get("mixtape-session", ctx)
	if err != nil {
		return err
	}
	state := sess.Values["auth_state"].(string)
	delete(sess.Values, "auth_state")

	auth := newAuthenticator(ctx)

	token, err := auth.Token(state, ctx.Request())
	if err != nil {
		return err
	}

	sess.Values["access_token"] = token.AccessToken
	sess.Values["refresh_token"] = token.RefreshToken
	sess.Save(ctx.Request(), ctx.Response())

	fragment := fmt.Sprintf("/#access_token=%s&refresh_token=%s", token.AccessToken, token.RefreshToken)
	return ctx.Redirect(http.StatusFound, fragment)
}

func newState() (string, error) {
	var state string

	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	if err != nil {
		return state, err
	}

	state = base64.StdEncoding.EncodeToString(randBytes)

	return state, nil
}

func newAuthenticator(ctx echo.Context) spotify.Authenticator {
	redirectURI := &url.URL{
		Scheme: ctx.Scheme(),
		Host:   ctx.Request().Host,
		Path:   "auth/callback",
	}

	return spotify.NewAuthenticator(
		redirectURI.String(),
		spotify.ScopeUserReadPrivate,
	)
}
