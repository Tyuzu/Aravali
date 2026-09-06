package vendors

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"scav/config"
	"scav/config/mqevent"
	"scav/infra"
	"scav/infra/mq"
	"scav/utils"
)

func ListAvailabilityHandler(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vendorID := utils.GetParam(r, "vendorID")
		if vendorID == "" {
			http.Error(w, "vendorID required", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		slots, err := FindAvailabilitySlots(ctx, app, vendorID)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"slots": slots})
	}
}

// Create availability slot (vendor sets unavailable dates or recurring availability)
func CreateAvailabilityHandler(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(config.UserIDKey).(string)

		vendorID := utils.GetParam(r, "vendorID")
		if vendorID == "" {
			http.Error(w, "vendorID required", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var slot AvailabilitySlot
		if err := json.NewDecoder(r.Body).Decode(&slot); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		// Basic validation
		slot.VendorID = vendorID
		if slot.StartDate == "" || slot.EndDate == "" {
			http.Error(w, "start_date and end_date required", http.StatusBadRequest)
			return
		}

		// Verify caller owns the vendor profile
		vendor, err := FindVendorByID(ctx, app, vendorID)
		if err != nil || vendor == nil {
			http.Error(w, "vendor not found", http.StatusNotFound)
			return
		}
		if userID == "" || vendor.UserID != userID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		existing, _ := FindAvailabilitySlots(ctx, app, vendorID)
		// simple in-app overlap check
		for _, ex := range existing {
			if !(ex.EndDate < slot.StartDate || ex.StartDate > slot.EndDate) {
				http.Error(w, "overlap", http.StatusConflict)
				return
			}
		}

		slot.SlotID = genSlotID()
		slot.CreatedAt = time.Now().UTC()

		if err := InsertAvailabilitySlotDB(ctx, app, slot); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		if err := mq.PublishWithMeta(ctx, app.MQ, mqevent.SlotCreatedEvent, mqevent.SlotCreatedPayload{}); err != nil {
			log.Printf("failed to publish slot created event: %v", err)
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"ok": true, "slot": slot})
	}
}

// Delete availability slot
func DeleteAvailabilityHandler(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(config.UserIDKey).(string)

		vendorID := utils.GetParam(r, "vendorID")
		slotID := utils.GetParam(r, "slotID")
		if vendorID == "" || slotID == "" {
			http.Error(w, "missing params", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		slot, err := FindAvailabilitySlotByID(ctx, app, slotID, vendorID)
		if err != nil || slot == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		// verify ownership
		vendor, err := FindVendorByID(ctx, app, vendorID)
		if err != nil || vendor == nil {
			http.Error(w, "vendor not found", http.StatusNotFound)
			return
		}
		if userID == "" || vendor.UserID != userID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		if _, err := DeleteAvailabilitySlotDB(ctx, app, slotID); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		if err := mq.PublishWithMeta(ctx, app.MQ, mqevent.SlotDeletedEvent, mqevent.SlotDeletedPayload{}); err != nil {
			log.Printf("failed to publish slot deleted event: %v", err)
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"ok": true})
	}
}

// helper to generate simple slot id
func genSlotID() string {
	return time.Now().UTC().Format("20060102T150405")
}
