package httpserver

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func mySkipper(c echo.Context) bool {
	if (strings.HasPrefix(c.Path(), "/env")) ||
		(strings.HasPrefix(c.Path(), "/info")) ||
		(strings.HasPrefix(c.Path(), "/health")) ||
		(strings.HasPrefix(c.Path(), "/favicon.ico")) {
		return true
	}
	return false
}
