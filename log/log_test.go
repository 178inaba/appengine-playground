package log

import (
	"testing"

	"github.com/labstack/echo/v4"
)

func TestLogger_ImplementEchoLogger(t *testing.T) {
	var _ echo.Logger = (*Logger)(nil)
}
