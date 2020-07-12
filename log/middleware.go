package log

import (
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/logging"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/labstack/echo/v4"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
)

type ContextLogger struct {
	echo.Context
	logger echo.Logger
}

func (cl *ContextLogger) Logger() echo.Logger {
	return cl.logger
}

type LoggerMiddleware struct {
	client *logging.Client

	moduleID  string
	projectID string
	versionID string
	zone      string
}

func NewLoggerMiddleware(client *logging.Client, moduleID, projectID, versionID, zone string) *LoggerMiddleware {
	return &LoggerMiddleware{client: client, moduleID: moduleID, projectID: projectID, versionID: versionID, zone: zone}
}

func (m *LoggerMiddleware) Logger(next echo.HandlerFunc) echo.HandlerFunc {
	hf := &propagation.HTTPFormat{}

	opt := logging.CommonResource(&mrpb.MonitoredResource{
		Type: "gae_app",
		Labels: map[string]string{
			"module_id":  m.moduleID,
			"project_id": m.projectID,
			"version_id": m.versionID,
			"zone":       m.zone,
		},
	})
	reqLogger := m.client.Logger(fmt.Sprintf("%s_request", m.moduleID), opt)
	appLogger := m.client.Logger(fmt.Sprintf("%s_application", m.moduleID), opt)
	return func(c echo.Context) error {
		req := c.Request()
		sc, _ := hf.SpanContextFromRequest(req)
		trace := fmt.Sprintf("projects/%s/traces/%s", m.projectID, sc.TraceID)
		spanID := sc.SpanID.String()

		logger := New(appLogger, trace, spanID)

		start := time.Now()
		if err := next(&ContextLogger{Context: c, logger: logger}); err != nil {
			c.Error(err) // TODO
		}
		end := time.Now()

		resp := c.Response()
		remoteIP := strings.Split(req.Header.Get("X-Forwarded-For"), ",")[0]
		reqLogger.Log(logging.Entry{
			Timestamp: time.Now(),
			Severity:  logger.maxSeverity,
			HTTPRequest: &logging.HTTPRequest{
				Request:      req,
				Latency:      end.Sub(start),
				Status:       resp.Status,
				RemoteIP:     remoteIP,
				ResponseSize: resp.Size,
			},
			Trace:        trace,
			TraceSampled: true,
			SpanID:       spanID,
		})

		return nil
	}
}
