package handler

import (
	"encoder/app"
	"encoder/helper"
	"encoder/m"
	"encoder/views"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetHistory(c echo.Context) error {
	var histories []m.History
	if err := app.DB.Find(&histories).Error; err != nil {
		c.Logger().Error("failed to list histories", err)
		return helper.Render(c,
			http.StatusBadRequest,
			views.History(
				helper.TCtxWError(c, errors.New("failed to list histories")),
				fmt.Sprintf("%s - History", app.Name),
				histories,
			),
		)
	}

	return helper.Render(c,
		http.StatusOK,
		views.History(
			helper.TCtx(c),
			fmt.Sprintf("%s - History", app.Name),
			histories,
		),
	)
}
