package main_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/therevels/mixtape"
)

var _ = Describe("Mixtape", func() {
	var origClientID, clientID string

	BeforeEach(func() {
		origClientID = os.Getenv("SPOTIFY_ID")
		clientID = "my-test-client-id"
		os.Setenv("SPOTIFY_ID", clientID)
	})

	AfterEach(func() {
		os.Setenv("SPOTIFY_ID", origClientID)
	})

	Describe("Login handler", func() {
		var req *http.Request
		var res *httptest.ResponseRecorder
		var handler http.HandlerFunc

		BeforeEach(func() {
			var err error

			req, err = http.NewRequest("GET", "/auth/login", nil)
			Expect(err).ToNot(HaveOccurred())

			res = httptest.NewRecorder()
			handler = http.HandlerFunc(Login)
		})

		Context("when not logged in", func() {
			It("redirects with a 302", func() {
				handler.ServeHTTP(res, req)
				Expect(res.Code).To(Equal(http.StatusFound))
			})

			It("redirects to Spotify authorization endpoint", func() {
				handler.ServeHTTP(res, req)
				loc, err := res.Result().Location()
				Expect(err).ToNot(HaveOccurred())
				Expect(loc.Scheme).To(Equal("https"))
				Expect(loc.Host).To(Equal("accounts.spotify.com"))
				Expect(loc.Path).To(Equal("/authorize"))
			})

			It("includes client ID in redirect", func() {
				handler.ServeHTTP(res, req)
				loc, _ := res.Result().Location()
				Expect(loc.Query().Get("client_id")).To(Equal(clientID))
			})

			It("includes redirect URI", func() {
				loginURL, _ := url.Parse("http://my.testserver.com:8088/auth/login")
				req.URL = loginURL
				handler.ServeHTTP(res, req)
				loc, _ := res.Result().Location()
				Expect(loc.Query().Get("redirect_uri")).To(Equal("http://my.testserver.com:8088/auth/callback"))
			})

			It("includes the auth scopes", func() {
				handler.ServeHTTP(res, req)
				loc, _ := res.Result().Location()
				Expect(loc.Query().Get("scope")).To(Equal("user-read-private"))
			})

			PIt("includes unique state", func() {})
		})

		PContext("when logged in", func() {
			// TODO: have to figure out what this means - a cookie? redis? are we using a session manager?
		})
	})

	Describe("Callback handler", func() {
		// var req *http.Request
		// var res *httptest.ResponseRecorder
		// var handler http.HandlerFunc

		// BeforeEach(func() {
		// 	var err error

		// 	req, err = http.NewRequest("GET", "/auth/callback", nil)
		// 	Expect(err).ToNot(HaveOccurred())

		// 	res = httptest.NewRecorder()
		// 	handler = http.HandlerFunc(Login)
		// })

		PIt("exchanges the authorization code for an access token", func() {})
		PIt("validates the session state", func() {})
		PIt("stores the access token in the session", func() {})
		PIt("stores the refresh token in the session", func() {})
	})
})
