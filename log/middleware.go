package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/logging"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/labstack/echo/v4"
)

type IndexHandler struct {
	requestLogger     *logging.Logger
	applicationLogger *logging.Logger
}

func (h *IndexHandler) Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		if err := next(c); err != nil {
			c.Error(err)
		}

		end := time.Now()

		projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
		req := c.Request()
		remoteIP := strings.Split(req.Header.Get("X-Forwarded-For"), ",")[0]

		resp := c.Response()

		hf := &propagation.HTTPFormat{}
		sc, _ := hf.SpanContextFromRequest(req)
		trace := fmt.Sprintf("projects/%s/traces/%s", projectID, sc.TraceID)
		h.requestLogger.Log(logging.Entry{
			Timestamp: time.Now(),
			Severity:  logging.Error,
			HTTPRequest: &logging.HTTPRequest{
				Request:      req,
				Latency:      end.Sub(start),
				Status:       resp.Status,
				RemoteIP:     remoteIP,
				ResponseSize: resp.Size,
			},
			Trace:        trace,
			TraceSampled: true,
			SpanID:       sc.SpanID.String(),
		})

		return nil
	}
}
