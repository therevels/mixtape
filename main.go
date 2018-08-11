package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.File("/", "static/index.html")
	e.GET("/auth/login", Login)

	e.Logger.Fatal(e.StartTLS(":8088", "cert.pem", "key.pem"))
}
