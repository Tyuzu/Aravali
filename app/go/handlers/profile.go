package handlers

import (
	"nae/models"
	"nae/utils"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	utils.SendJSON(w, http.StatusOK, models.ProfileData{
		Name:   "Jane Doe",
		Email:  "jane.doe@example.com",
		Role:   "Lead Systems Architect",
		Status: "Active",
	})
}
