package menu

import (
	"encoding/json"
	"net/http"

	"scav/config/mqevent"
	"scav/infra"
	"scav/infra/mq"
	"scav/internal/pay/stripe"
	"scav/utils"
	log "scav/utils/logger"
)

// POST /menu/event/:placeId/:menuId/payment-session
func CreateMenuPaymentSession(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		menuId := utils.GetParam(r, "menuid")
		placeId := utils.GetParam(r, "placeid")

		var body struct {
			Stock int `json:"quantity"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Stock < 1 {
			http.Error(w, "Invalid request or stock", http.StatusBadRequest)
			return
		}

		session, err := stripe.CreateMenuSession(menuId, placeId, body.Stock)
		if err != nil {
			log.Printf("Error creating payment session: %v", err)
			http.Error(w, "Failed to create payment session", http.StatusInternalServerError)
			return
		}

		dataResponse := map[string]any{
			"paymentUrl": session.URL,
			"placeid":    session.PlaceID,
			"menuid":     session.MenuID,
			"stock":      session.Stock,
		}

		if err := mq.PublishWithMeta(ctx, app.MQ, mqevent.MenuPaymentSessionInitiatedEvent, mqevent.MenuPaymentSessionInitiatedPayload{}); err != nil {
			log.Printf("failed to publish menu payment session initiated event: %v", err)
		}

		response := map[string]any{
			"success": true,
			"data":    dataResponse,
		}
		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}
