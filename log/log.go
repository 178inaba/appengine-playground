package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"cloud.google.com/go/logging"
	"github.com/labstack/gommon/log"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
)

var severityLogLevel = map[logging.Severity]log.Lvl{
	logging.Default:  0,
	logging.Debug:    log.DEBUG,
	logging.Info:     log.INFO,
	logging.Warning:  log.WARN,
	logging.Error:    log.ERROR,
	logging.Critical: 6,
	logging.Alert:    7,
}

type Logger struct {
	logger *logging.Logger

	trace  string
	spanID string
	level  log.Lvl
}

func New(logger *logging.Logger, trace, spanID string) *Logger {
	return &Logger{
		logger: logger,
		trace:  trace,
		spanID: spanID,
	}
}

func (l *Logger) Output() io.Writer     { return nil }
func (l *Logger) SetOutput(w io.Writer) {}
func (l *Logger) Prefix() string        { return "" }
func (l *Logger) SetPrefix(p string)    {}
func (l *Logger) SetHeader(h string)    {}

func (l *Logger) Level() log.Lvl {
	return l.level
}

func (l *Logger) SetLevel(v log.Lvl) {
	l.level = v
}

func (l *Logger) Print(i ...interface{}) {
	l.log(logging.Default, fmt.Sprint(i...))
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.log(logging.Default, fmt.Sprintf(format, args...))
}

func (l *Logger) Printj(j log.JSON) {
	l.log(logging.Default, j)
}

func (l *Logger) Debug(i ...interface{}) {
	l.log(logging.Debug, fmt.Sprint(i...))
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(logging.Debug, fmt.Sprintf(format, args...))
}

func (l *Logger) Debugj(j log.JSON) {
	l.log(logging.Debug, j)
}

func (l *Logger) Info(i ...interface{}) {
	l.log(logging.Info, fmt.Sprint(i...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(logging.Info, fmt.Sprintf(format, args...))
}

func (l *Logger) Infoj(j log.JSON) {
	l.log(logging.Info, j)
}

func (l *Logger) Warn(i ...interface{}) {
	l.log(logging.Warning, fmt.Sprint(i...))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(logging.Warning, fmt.Sprintf(format, args...))
}

func (l *Logger) Warnj(j log.JSON) {
	l.log(logging.Warning, j)
}

func (l *Logger) Error(i ...interface{}) {
	l.log(logging.Error, fmt.Sprint(i...))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(logging.Error, fmt.Sprintf(format, args...))
}

func (l *Logger) Errorj(j log.JSON) {
	l.log(logging.Error, j)
}

func (l *Logger) Fatal(i ...interface{}) {
	l.log(logging.Critical, fmt.Sprint(i...))
	l.logger.Flush()
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(logging.Critical, fmt.Sprintf(format, args...))
	l.logger.Flush()
	os.Exit(1)
}

func (l *Logger) Fatalj(j log.JSON) {
	l.log(logging.Critical, j)
	l.logger.Flush()
	os.Exit(1)
}

func (l *Logger) Panic(i ...interface{}) {
	s := fmt.Sprint(i...)
	l.log(logging.Alert, s)
	panic(s)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	l.log(logging.Alert, s)
	panic(s)
}

func (l *Logger) Panicj(j log.JSON) {
	l.log(logging.Alert, j)
	panic(j)
}

func (l *Logger) log(severity logging.Severity, payload interface{}) {
	if l.level >= severityLogLevel[severity] {
		return
	}

	pc, file, line, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	l.logger.Log(logging.Entry{
		Timestamp:    time.Now(),
		Severity:     severity,
		Trace:        l.trace,
		TraceSampled: true,
		SourceLocation: &logpb.LogEntrySourceLocation{
			File:     file,
			Line:     int64(line),
			Function: f.Name(),
		},
		SpanID:  l.spanID,
		Payload: payload,
	})
}
