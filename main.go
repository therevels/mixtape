package main

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	secret := os.Getenv("SECRET_TOKEN")
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(secret))))

	e.File("/", "static/index.html")
	e.GET("/auth/login", Login)
	e.GET("/auth/callback", Callback)

	e.Logger.Fatal(e.StartTLS(":8088", "cert.pem", "key.pem"))
}
