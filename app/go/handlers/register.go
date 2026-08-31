package handlers

import (
	"encoding/json"
	"nae/models"
	"nae/utils"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "Method not allowed"})
		return
	}

	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "Invalid request payload"})
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.SendJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "Email and password required"})
		return
	}

	token, err := utils.GenerateToken(req.Email)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "Token generation failed"})
		return
	}

	utils.SendJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data: models.AuthData{
			Token: token,
			User:  req.Email,
		},
	})
}
