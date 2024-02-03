package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func PostLogout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(-1 * 24 * time.Hour),
	})

	return c.Redirect(http.StatusFound, "/login")
}
