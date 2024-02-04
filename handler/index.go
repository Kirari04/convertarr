package handler

import (
	"encoder/app"
	"encoder/helper"
	"encoder/m"
	"encoder/views"
	"fmt"
	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func GetIndex(c echo.Context) error {
	type historyStats struct {
		AvgNewSize float64
		AvgOldSize float64
	}
	var historyStatsRes historyStats
	if err := app.DB.
		Model(&m.History{}).
		Select(
			"AVG(histories.new_size) as avg_new_size",
			"AVG(histories.old_size) as avg_old_size",
		).
		Scan(&historyStatsRes).
		Error; err != nil {
		log.Error("Failed to get history stats", err)
	}

	savedStorage := historyStatsRes.AvgOldSize - historyStatsRes.AvgNewSize

	var encodedFiles int64
	if err := app.DB.
		Model(&m.History{}).
		Where(&m.History{
			Status: "finished",
		}).
		Count(&encodedFiles).
		Error; err != nil {
		log.Error("Failed to get history stats", err)
	}

	return helper.Render(c,
		http.StatusOK,
		views.Index(
			helper.TCtx(c),
			fmt.Sprintf("%s - Home", app.Name),
			app.ResourcesHistory,
			humanize.Bytes(uint64(savedStorage)),
			fmt.Sprint(encodedFiles),
		),
	)
}
