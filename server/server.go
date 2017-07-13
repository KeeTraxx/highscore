package main

import (
	"io/ioutil"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
)

func main() {
	InitLogging(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	initDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionName))))

	initAPI(e.Group("/api"))

	e.Static("/", "public")
	setupAuth(e)
	PORT := getenv("PORT", "1323")
	e.Logger.Fatal(e.Start(":" + PORT))
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
