package log

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var _ echo.Logger = (*Logger)(nil)

type Logger struct {
	logger *log.Logger
}

func New(prefix string) *Logger {
	return &Logger{
		logger: log.New(prefix),
	}
}

func (l *Logger) Output() io.Writer {
	return l.logger.Output()
}

func (l *Logger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

func (l *Logger) Prefix() string {
	return l.logger.Prefix()
}

func (l *Logger) SetPrefix(p string) {
	l.logger.SetPrefix(p)
}

func (l *Logger) Level() log.Lvl {
	return l.logger.Level()
}

func (l *Logger) SetLevel(v log.Lvl) {
	l.SetLevel(v)
}

func (l *Logger) SetHeader(h string) {
	l.logger.SetHeader(h)
}

func (l *Logger) Print(i ...interface{}) {
	l.logger.Print(i...)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *Logger) Printj(j log.JSON) {
	l.logger.Panicj(j)
}

func (l *Logger) Debug(i ...interface{}) {
	l.logger.Debug(i...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *Logger) Debugj(j log.JSON) {
	l.logger.Debugj(j)
}

func (l *Logger) Info(i ...interface{}) {
	l.logger.Info(i...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *Logger) Infoj(j log.JSON) {
	l.logger.Infoj(j)
}

func (l *Logger) Warn(i ...interface{}) {
	l.logger.Warn(i...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *Logger) Warnj(j log.JSON) {
	l.logger.Warnj(j)
}

func (l *Logger) Error(i ...interface{}) {
	l.logger.Error(i...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *Logger) Errorj(j log.JSON) {
	l.logger.Errorj(j)
}

func (l *Logger) Fatal(i ...interface{}) {
	l.logger.Fatal(i...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *Logger) Fatalj(j log.JSON) {
	l.logger.Fatalj(j)
}

func (l *Logger) Panic(i ...interface{}) {
	l.logger.Panic(i...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l *Logger) Panicj(j log.JSON) {
	l.logger.Panicj(j)
}
