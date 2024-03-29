package handler

import (
	"encoder/app"
	"encoder/helper"
	"encoder/m"
	"encoder/t"
	"encoder/views"
	"fmt"
	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func GetIndex(c echo.Context) error {
	long := c.QueryParam("long")
	var longStats bool
	if long != "" {
		longStats = true
	}
	type historyStats struct {
		SumNewSize float64
		SumOldSize float64
	}
	var historyStatsRes historyStats
	if err := app.DB.
		Model(&m.History{}).
		Select(
			"SUM(histories.new_size) as sum_new_size",
			"SUM(histories.old_size) as sum_old_size",
		).
		Where(&m.History{
			Status: "finished",
		}).
		Scan(&historyStatsRes).
		Error; err != nil {
		log.Error("Failed to get history stats", err)
	}

	savedStorage := historyStatsRes.SumOldSize - historyStatsRes.SumNewSize

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

	var (
		resourcesHistory t.Resources
	)
	shortBreakpoint := 48
	if !longStats && len(app.ResourcesHistory.Cpu) > shortBreakpoint {
		resourcesHistory.Cpu = app.ResourcesHistory.Cpu[len(app.ResourcesHistory.Cpu)-shortBreakpoint:]
	} else {
		resourcesHistory.Cpu = app.ResourcesHistory.Cpu
	}
	if !longStats && len(app.ResourcesHistory.Mem) > shortBreakpoint {
		resourcesHistory.Mem = app.ResourcesHistory.Mem[len(app.ResourcesHistory.Mem)-shortBreakpoint:]
	} else {
		resourcesHistory.Mem = app.ResourcesHistory.Mem
	}
	if !longStats && len(app.ResourcesHistory.NetOut) > shortBreakpoint {
		resourcesHistory.NetOut = app.ResourcesHistory.NetOut[len(app.ResourcesHistory.NetOut)-shortBreakpoint:]
	} else {
		resourcesHistory.NetOut = app.ResourcesHistory.NetOut
	}
	if !longStats && len(app.ResourcesHistory.NetIn) > shortBreakpoint {
		resourcesHistory.NetIn = app.ResourcesHistory.NetIn[len(app.ResourcesHistory.NetIn)-shortBreakpoint:]
	} else {
		resourcesHistory.NetIn = app.ResourcesHistory.NetIn
	}

	return helper.Render(c,
		http.StatusOK,
		views.Index(
			helper.TCtx(c),
			fmt.Sprintf("%s - Home", app.Name),
			resourcesHistory,
			longStats,
			humanize.Bytes(uint64(savedStorage)),
			fmt.Sprint(encodedFiles),
		),
	)
}
