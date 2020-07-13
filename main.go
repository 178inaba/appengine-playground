package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/logging"
	aelogging "github.com/178inaba/appengine-echo-logging"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/labstack/gommon/log"
)

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	log.SetFlags(0)

	zone, err := metadataZone(ctx)
	if err != nil {
		log.Fatalf("metadata zone: %s.", err)
	}

	loggingClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("New logging client: %s.", err)
	}
	defer loggingClient.Close()

	service := os.Getenv("GAE_SERVICE")
	version := os.Getenv("GAE_VERSION")
	lm := aelogging.NewLoggerMiddleware(loggingClient, service, projectID, version, zone)

	e := echo.New()
	e.Logger.SetLevel(echolog.INFO)
	e.Use(lm.Logger)
	e.Use(middleware.Recover())

	e.GET("/", index)
	e.GET("/sleep", sleep)
	e.GET("/hello", hello)
	e.GET("/_ah/warmup", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		e.Logger.Debugf("Defaulting to port %s.", port)
	}

	e.Logger.Fatalf("%v.", e.Start(fmt.Sprintf(":%s", port)))
}

func index(c echo.Context) error {
	c.Logger().Info("Start index.")
	defer c.Logger().Info("End index.")

	c.Logger().Print("Print.")
	c.Logger().Debug("Debug.")
	c.Logger().Info("Info.")
	c.Logger().Warn("Warning.")
	c.Logger().Error("Error.")
	c.Logger().Error("c.Request().ContentLength: %d.", c.Request().ContentLength)

	return c.String(http.StatusTeapot, "Index!")
}

func sleep(c echo.Context) error {
	c.Logger().Info("Start sleep.")
	defer c.Logger().Info("End sleep.")

	time.Sleep(10 * time.Second)
	c.Logger().Info("sleeped!!!")

	return c.String(http.StatusCreated, "Sleep!")
}

func hello(c echo.Context) error {
	c.Logger().Info("Start.")
	defer c.Logger().Info("End.")

	var trace string
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID != "" {
		traceHeader := c.Request().Header.Get("X-Cloud-Trace-Context")
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

func metadataZone(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://metadata/computeMetadata/v1/instance/zone", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ss := strings.Split(string(bs), "/")

	return ss[len(ss)-1], nil
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
