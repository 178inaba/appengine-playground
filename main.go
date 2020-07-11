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
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/labstack/gommon/log"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
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
	opt := logging.CommonResource(&mrpb.MonitoredResource{
		Type: "gae_app",
		Labels: map[string]string{
			"module_id":  service,
			"project_id": projectID,
			"version_id": os.Getenv("GAE_VERSION"),
			"zone":       zone,
		},
	})
	reqLogger := loggingClient.Logger(fmt.Sprintf("%s_request", service), opt)
	appLogger := loggingClient.Logger(fmt.Sprintf("%s_application", service), opt)

	ih := &IndexHandler{requestLogger: reqLogger, applicationLogger: appLogger}

	e := echo.New()
	e.Logger.SetLevel(echolog.INFO)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", ih.index)
	e.GET("/hello", hello)
	e.GET("/_ah/warmup", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		e.Logger.Debugf("Defaulting to port %s.", port)
	}

	e.Logger.Fatalf("%v.", e.Start(fmt.Sprintf(":%s", port)))
}

type IndexHandler struct {
	requestLogger     *logging.Logger
	applicationLogger *logging.Logger
}

func (h *IndexHandler) index(c echo.Context) error {
	c.Logger().Info("Start.")
	defer c.Logger().Info("End.")

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	req := c.Request()
	hf := &propagation.HTTPFormat{}
	sc, _ := hf.SpanContextFromRequest(req)
	trace := fmt.Sprintf("projects/%s/traces/%s", projectID, sc.TraceID)
	h.requestLogger.Log(logging.Entry{
		Timestamp: time.Now(),
		Severity:  logging.Error,
		HTTPRequest: &logging.HTTPRequest{
			Request: req,
			Latency: 10 * time.Second,

			// RequestSize is the size of the HTTP request message in bytes, including
			// the request headers and the request body.
			//RequestSize int64

			// Status is the response code indicating the status of the response.
			// Examples: 200, 404.
			//Status int

			// ResponseSize is the size of the HTTP response message sent back to the client, in bytes,
			// including the response headers and the response body.
			//ResponseSize int64

			// Latency is the request processing latency on the server, from the time the request was
			// received until the response was sent.
			//Latency time.Duration

			// LocalIP is the IP address (IPv4 or IPv6) of the origin server that the request
			// was sent to.
			//LocalIP string

			// RemoteIP is the IP address (IPv4 or IPv6) of the client that issued the
			// HTTP request. Examples: "192.168.1.1", "FE80::0202:B3FF:FE1E:8329".
			//RemoteIP string

			// CacheHit reports whether an entity was served from cache (with or without
			// validation).
			//CacheHit bool

			// CacheValidatedWithOriginServer reports whether the response was
			// validated with the origin server before being served from cache. This
			// field is only meaningful if CacheHit is true.
			//CacheValidatedWithOriginServer bool
		},
		Trace:        trace,
		TraceSampled: true,
		SpanID:       sc.SpanID.String(),
		//Payload interface{}
		//Labels map[string]string
		//InsertID string
		//Operation *logpb.LogEntryOperation
	})

	h.applicationLogger.Log(logging.Entry{
		Timestamp:    time.Now(),
		Severity:     logging.Critical,
		Trace:        trace,
		TraceSampled: true,
		SourceLocation: &logpb.LogEntrySourceLocation{
			File:     "main.go",
			Line:     101,
			Function: "index",
		},
		SpanID:  sc.SpanID.String(),
		Payload: "hello log!!",
		//Labels map[string]string
		//InsertID string
		//Operation *logpb.LogEntryOperation
	})

	return c.String(http.StatusOK, "Index!")
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
