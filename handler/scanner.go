package handler

import (
	"encoder/setup"
	"net/http"

	"github.com/labstack/echo/v4"
)

func PostScanner(c echo.Context) error {
	go setup.ScannFolders()
	return c.Redirect(http.StatusFound, "/history")
}
