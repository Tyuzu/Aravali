package handlers

import (
	"nae/models"
	"nae/utils"
	"net/http"
)

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	utils.SendJSON(w, http.StatusOK, models.SettingsData{
		NotificationsEnabled: true,
		DarkMode:             true,
		ApiVersion:           "v1.4.0",
		Theme:                "Midnight Blue",
	})
}
