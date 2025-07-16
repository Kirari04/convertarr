package mware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func Caching(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.HasPrefix(c.Request().URL.Path, "/resources") {
			c.Response().Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(time.Hour*24)))
			c.Response().Header().Set("Expires", time.Now().Add(time.Hour*24).Format(http.TimeFormat))
		}

		return next(c)
	}
}
