package main_test

import (
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/quasoft/memstore"

	. "github.com/therevels/mixtape"
)

var _ = Describe("Auth", func() {
	var origClientID, clientID string

	BeforeEach(func() {
		origClientID = os.Getenv("SPOTIFY_ID")
		clientID = "my-test-client-id"
		os.Setenv("SPOTIFY_ID", clientID)
	})

	AfterEach(func() {
		os.Setenv("SPOTIFY_ID", origClientID)
	})

	Describe("Login", func() {
		var req *http.Request
		var rec *httptest.ResponseRecorder
		var ctx echo.Context

		BeforeEach(func() {
			e := echo.New()
			e.Use(session.Middleware(
				memstore.NewMemStore(
					[]byte("authkey123"),
					[]byte("enckey12341234567890123456789012"),
				),
			))

			req = httptest.NewRequest(echo.GET, "https://example.com:443/auth/login", nil)
			rec = httptest.NewRecorder()
			ctx = e.NewContext(req, rec)
		})

		Context("when not logged in", func() {
			It("redirects with a 302", func() {
				Login(ctx)
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

			PIt("includes unique state", func() {})
		})

		PContext("when logged in", func() {
			// TODO: have to figure out what this means - a cookie? redis? are we using a session manager?
		})
	})

	Describe("Callback", func() {
		// TODO: ideally this would have full test coverage, but between
		// the spotify and oauth2 libraries, the abstractions do not make it
		// easily testable (no way to inject server URL, etc)
		PIt("validates the session state", func() {})
		PIt("exchanges the authorization code for an access token", func() {})
		PIt("stores the tokens in the session", func() {})
	})
})
