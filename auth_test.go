package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/quasoft/memstore"

	. "github.com/therevels/mixtape"
)

var _ = Describe("Auth", func() {
	var origClientID, clientID string
	var req *http.Request
	var rec *httptest.ResponseRecorder
	var e *echo.Echo
	var ctx echo.Context
	var store sessions.Store
	var sess *sessions.Session

	BeforeEach(func() {
		origClientID = os.Getenv("SPOTIFY_ID")
		clientID = "my-test-client-id"
		os.Setenv("SPOTIFY_ID", clientID)

		e = echo.New()
		rec = httptest.NewRecorder()
		store = memstore.NewMemStore(
			[]byte("authkey123"),
			[]byte("enckey12341234567890123456789012"),
		)
	})

	AfterEach(func() {
		os.Setenv("SPOTIFY_ID", origClientID)
	})

	Describe("Login", func() {
		BeforeEach(func() {
			req = httptest.NewRequest(echo.GET, "https://example.com:443/auth/login", nil)
			ctx = e.NewContext(req, rec)
			ctx.Set("_session_store", store)
		})

		Context("when not logged in", func() {
			It("redirects with a 302", func() {
				Expect(Login(ctx)).To(Succeed())
				Expect(rec.Code).To(Equal(http.StatusFound))
			})

			It("redirects to Spotify authorization endpoint", func() {
				Login(ctx)
				loc, err := rec.Result().Location()
				Expect(err).ToNot(HaveOccurred())
				Expect(loc.Scheme).To(Equal("https"))
				Expect(loc.Host).To(Equal("accounts.spotify.com"))
				Expect(loc.Path).To(Equal("/authorize"))
			})

			It("includes client ID in redirect", func() {
				Login(ctx)
				loc, _ := rec.Result().Location()
				Expect(loc.Query().Get("client_id")).To(Equal(clientID))
			})

			It("includes redirect URI", func() {
				Login(ctx)
				loc, _ := rec.Result().Location()
				Expect(loc.Query().Get("redirect_uri")).To(Equal("https://example.com:443/auth/callback"))
			})

			It("includes the auth scopes", func() {
				Login(ctx)
				loc, _ := rec.Result().Location()
				Expect(loc.Query().Get("scope")).To(Equal("user-read-private"))
			})

			It("includes unique session state", func() {
				Expect(Login(ctx)).To(Succeed())
				loc, _ := rec.Result().Location()
				state := loc.Query().Get("state")

				sess, err := store.Get(req, "mixtape-session")
				Expect(err).ToNot(HaveOccurred())
				sessState := sess.Values["auth_state"].(string)
				Expect(state).To(Equal(sessState))
			})
		})

		Context("when logged in", func() {
			BeforeEach(func() {
				sess, _ = store.Get(req, "mixtape-session")
				sess.Values["access_token"] = "my-existing-access-token"
				sess.Values["refresh_token"] = "my-existing-refresh-token"
			})

			AfterEach(func() {
				delete(sess.Values, "access_token")
				delete(sess.Values, "refresh_token")
			})

			It("redirects to root", func() {
				Expect(Login(ctx)).To(Succeed())
				Expect(rec.Code).To(Equal(http.StatusFound))
				loc, _ := rec.Result().Location()
				Expect(loc.String()).To(Equal("/"))
			})
		})
	})

	Describe("Callback", func() {
		var code, state string

		Context("with authorization code", func() {
			BeforeEach(func() {
				code = "my-authorization-code"
				state = "my-redirect-state"
				callbackURL := fmt.Sprintf("https://example.com/auth/callback?code=%s&state=%s", code, state)
				req = httptest.NewRequest(echo.GET, callbackURL, nil)
				ctx = e.NewContext(req, rec)
				ctx.Set("_session_store", store)
				sess, _ = store.Get(req, "mixtape-session")
			})

			Context("when there is no session state", func() {
				It("returns an error", func() {
					err := Callback(ctx)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError("spotify: redirect state parameter doesn't match"))
				})
			})

			Context("when session state is invalid", func() {
				BeforeEach(func() {
					sess.Values["state"] = "some-completely-different-value"
				})

				AfterEach(func() {
					delete(sess.Values, "state")
				})

				It("returns an error", func() {
					err := Callback(ctx)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError("spotify: redirect state parameter doesn't match"))
				})
			})

			PContext("when session state is valid", func() {
				// TODO: ideally this would have full test coverage, but between
				// the spotify and oauth2 libraries, the abstractions do not make it
				// easily testable (no way to inject server URL, etc)
				It("exchanges the authorization code for an access token", func() {})
				It("stores the tokens in the session", func() {})
			})
		})

		Context("without authorization code", func() {
			BeforeEach(func() {
				callbackURL := "https://example.com/auth/callback"
				req = httptest.NewRequest(echo.GET, callbackURL, nil)
				ctx = e.NewContext(req, rec)
				ctx.Set("_session_store", store)
			})

			It("returns an error", func() {
				err := Callback(ctx)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("spotify: didn't get access code"))
			})
		})
	})
})
