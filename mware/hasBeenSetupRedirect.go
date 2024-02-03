package mware

import (
	"encoder/app"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HasBeenSetupRedirect(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !app.Setting.HasBeenSetup {
			if c.Request().URL.Path != "/setup" && c.Request().URL.Path != "/favicon.ico" {
				return c.Redirect(http.StatusFound, "/setup")
			}
		}
		return next(c)
	}
}
