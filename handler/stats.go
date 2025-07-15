package handler

import (
	"encoder/app"
	"encoder/t"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetStatsData handles the request for resource statistics.
// It checks for a 'long' query parameter to determine whether to return
// the full history or a shorter, more recent snapshot of the data.
func GetStatsData(c echo.Context) error {
	long := c.QueryParam("long")
	var longStats bool
	if long != "" {
		longStats = true
	}

	var resourcesHistory t.Resources

	// The number of data points to show for the "live" view.
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

	return c.JSON(http.StatusOK, resourcesHistory)
}
