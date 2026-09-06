package menu

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"scav/config"
	"scav/config/mqevent"
	"scav/infra"
	"scav/infra/mq"
	"scav/internal/userdata"
	"scav/utils"
	log "scav/utils/logger"
)

// BuyMenu validates the request, atomically decreases stock, records the purchase,
// and emits the purchase event. This is the single purchase handler kept in the route layer.
func BuyMenu(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		placeID := utils.GetParam(r, "placeid")
		menuID := utils.GetParam(r, "menuid")

		var body struct {
			Quantity int `json:"quantity"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Quantity <= 0 {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}

		requestingUserID, ok := r.Context().Value(config.UserIDKey).(string)
		if !ok || requestingUserID == "" {
			http.Error(w, "Invalid user", http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		menu, err := findMenuByPlaceAndID(ctx, app, placeID, menuID)
		if err != nil {
			http.Error(w, "Menu not found", http.StatusNotFound)
			return
		}

		if menu.Stock < body.Quantity {
			http.Error(w, "Not enough menu available", http.StatusBadRequest)
			return
		}

		if err := decrementMenuStock(ctx, app, placeID, menuID, body.Quantity); err != nil {
			http.Error(w, "Failed to update menu stock", http.StatusInternalServerError)
			return
		}

		userdata.SetUserData("menu", menuID, requestingUserID, "place", placeID, app)

		if err := mq.PublishWithMeta(ctx, app.MQ, mqevent.MenuBoughtEvent, mqevent.MenuBoughtPayload{}); err != nil {
			log.Printf("failed to publish menu bought event: %v", err)
		}

		resp := map[string]any{
			"success":        true,
			"remainingStock": menu.Stock - body.Quantity,
			"message":        "Payment successfully processed. Menu purchased.",
		}
		utils.RespondWithJSON(w, http.StatusOK, resp)
	}
}
