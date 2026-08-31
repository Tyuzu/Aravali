package handlers

import (
	"nae/models"
	"nae/utils"
	"net/http"
)

func AnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	utils.SendJSON(w, http.StatusOK, models.AnalyticsData{
		DailyActiveUsers: 14850,
		TotalRevenue:     32450.75,
		ServerUptime:     "99.98%",
		ConversionRate:   "4.2%",
	})
}
