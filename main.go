package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/logging"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/labstack/gommon/log"
)

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	log.SetFlags(0)

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
	e.Logger.SetLevel(echolog.INFO)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/_ah/warmup", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })
	e.GET("/", hello)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		e.Logger.Debugf("Defaulting to port %s.", port)
	}

	e.Logger.Fatalf("%v.", e.Start(fmt.Sprintf(":%s", port)))
}

func hello(c echo.Context) error {
	c.Logger().Info("Start.")
	defer c.Logger().Info("End.")

	var trace string
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID != "" {
		traceHeader := c.Request().Header.Get("X-Cloud-Trace-Context")
		fmt.Println(traceHeader)
		traceParts := strings.Split(traceHeader, "/")
		if len(traceParts) > 0 && len(traceParts[0]) > 0 {
			trace = fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
		}
	}

	log.Println(Entry{Severity: "DEFAULT", Message: "Default", Trace: trace})
	log.Println(Entry{Severity: "DEBUG", Message: "Debug", Trace: trace})
	log.Println(Entry{Severity: "INFO", Message: "Info", Trace: trace})
	log.Println(Entry{Severity: "NOTICE", Message: "Notice", Trace: trace})
	log.Println(Entry{Severity: "WARNING", Message: "Warning", Trace: trace})
	log.Println(Entry{Severity: "ERROR", Message: "Error", Trace: trace})
	log.Println(Entry{Severity: "CRITICAL", Message: "Critical", Trace: trace})
	log.Println(Entry{Severity: "ALERT", Message: "Alert", Trace: trace})
	log.Println(Entry{Severity: "EMERGENCY", Message: "Emergency", Trace: trace})

	return c.String(http.StatusOK, "Hello, World!")
}

// Entry defines a log entry.
type Entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`
}

// String renders an entry structure to the JSON format expected by Stackdriver.
func (e Entry) String() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}
