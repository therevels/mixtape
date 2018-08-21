package main

import (
	"encoding/gob"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"golang.org/x/oauth2"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// TODO: research config for encryption in gorilla sessions
	secret := os.Getenv("SECRET_TOKEN")
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(secret))))

	// Permit serializing the oauth2.Token type to the session
	gob.Register(&oauth2.Token{})

	// Routing
	e.File("/", "static/index.html")
	e.GET("/auth/login", Login)
	e.GET("/auth/callback", Callback)
	e.GET("/auth/logout", Logout)

	e.Logger.Fatal(e.StartTLS(":8088", "cert.pem", "key.pem"))
}
