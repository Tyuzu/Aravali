package tickets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"scav/config/mqevent"
	"scav/infra"
	"scav/infra/mq"
	"scav/internal/beats/userdata"
	"scav/utils"
	log "scav/utils/logger"
)

func TransferTicket(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := utils.GetParam(r, "eventid")
		requestingUserId := utils.GetUserIDFromRequest(r)

		type TransferPayload struct {
			UniqueCode string `json:"uniqueCode"`
			Recipient  string `json:"recipient"`
		}

		var payload TransferPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		if payload.UniqueCode == "" || payload.Recipient == "" {
			http.Error(w, "uniqueCode and recipient are required", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		ticket, err := FindPurchasedTicketByUnique(ctx, app, eventID, payload.UniqueCode)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ticket not found: %v", err), http.StatusNotFound)
			return
		}

		// Ownership check
		if ticket.UserID != requestingUserId {
			http.Error(w, "You are not authorized to transfer this ticket", http.StatusForbidden)
			return
		}

		// Update ownership
		if _, err := UpdatePurchasedTicket(ctx, app, map[string]any{"eventid": eventID, "uniquecode": payload.UniqueCode}, map[string]any{"$set": map[string]any{"userid": payload.Recipient, "transferred": true, "transferredto": payload.Recipient}}); err != nil {
			http.Error(w, fmt.Sprintf("Failed to transfer ticket: %v", err), http.StatusInternalServerError)
			return
		}

		// Update user data indexes
		userdata.DelUserData("ticket", payload.UniqueCode, requestingUserId, app)
		userdata.SetUserData("ticket", payload.UniqueCode, payload.Recipient, "event", eventID, app)

		if err := mq.PublishWithMeta(ctx, app.MQ, mqevent.TicketTransferredEvent, mqevent.TicketTransferredPayload{}); err != nil {
			log.Printf("failed to publish ticket transferred event: %v", err)
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"message": "Ticket transferred successfully",
		})
	}
}
