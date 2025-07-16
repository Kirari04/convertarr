package handler

import (
	"encoder/app"
	"encoder/m"
	"encoder/t"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetConversionStats(c echo.Context) error {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	var histories []m.History
	if err := app.DB.
		Model(&m.History{}).
		Where("created_at >= ? AND (status = 'finished' OR status = 'failed')", thirtyDaysAgo).
		Find(&histories).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	stats := t.ConversionStats{
		Labels:     make([]string, 30),
		Successful: make([]int, 30),
		Failed:     make([]int, 30),
	}

	resultMap := make(map[string]map[string]int)
	for _, h := range histories {
		dateStr := h.CreatedAt.Format("2006-01-02")
		if _, ok := resultMap[dateStr]; !ok {
			resultMap[dateStr] = make(map[string]int)
		}
		resultMap[dateStr][h.Status]++
	}

	for i := 0; i < 30; i++ {
		day := time.Now().AddDate(0, 0, -29+i)
		dateStr := day.Format("2006-01-02")
		stats.Labels[i] = day.Format("Jan 2")

		if dayData, ok := resultMap[dateStr]; ok {
			if val, ok := dayData["finished"]; ok {
				stats.Successful[i] = val
			}
			if val, ok := dayData["failed"]; ok {
				stats.Failed[i] = val
			}
		}
	}

	return c.JSON(http.StatusOK, stats)
}