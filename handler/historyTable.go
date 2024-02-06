package handler

import (
	"encoder/app"
	"encoder/components"
	"encoder/helper"
	"encoder/m"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetHistoryTable(c echo.Context) error {
	var histories []m.History
	if err := app.DB.
		Order("id DESC").
		Limit(50).
		Find(&histories).Error; err != nil {
		c.Logger().Error("failed to list histories", err)
		return c.String(http.StatusInternalServerError, "failed to list histories")
	}

	return helper.Render(c,
		http.StatusOK,
		components.HistoryTable(
			helper.TCtx(c),
			histories,
		),
	)
}
