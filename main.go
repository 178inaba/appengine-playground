package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", hello)
	e.GET("/_ah/start", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })
	e.GET("/_ah/stop", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })
	e.GET("/_ah/warmup", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	port := os.Getenv("PORT")
	if port != "" {
		port = "8080"
		e.Logger.Debugf("Defaulting to port %s.", port)
	}

	e.Logger.Fatalf("%v.", e.Start(fmt.Sprintf(":%s", port)))
}

func hello(c echo.Context) error {
	c.Logger().Info("Start.")
	defer c.Logger().Info("End.")

	return c.String(http.StatusOK, "Hello, World!")
}
