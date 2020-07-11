package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/logging"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	loggingClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("New logging client: %s.", err)
	}
	defer loggingClient.Close()

	// ---
	logName := "my-log"
	logger := loggingClient.Logger(logName).StandardLogger(logging.Info)
	logger.Println("hello world")
	// ---

	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		c.Logger().Info("Start.")
		defer c.Logger().Info("End.")
		logger.Println("hello world")

		return c.String(http.StatusOK, "Hello, World!")
	})

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
