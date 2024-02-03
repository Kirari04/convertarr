package handler

import (
	"encoder/app"
	"encoder/helper"
	"encoder/views"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetIndex(c echo.Context) error {
	return helper.Render(c, http.StatusOK, views.Index(nil, fmt.Sprintf("%s - Home", app.Name)))
}
