package mware

import (
	"encoder/app"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HasBeenSetupRedirect(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !app.Setting.HasBeenSetup {
			if c.Path() != "/setup" {
				c.Redirect(http.StatusTemporaryRedirect, "/setup")
			}
		}
		return next(c)
	}
}
