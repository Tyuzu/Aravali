package delwebhooks

import (
	"encoding/json"
	"net/http"
	"time"

	"scav/infra"
	"scav/internal/deliveries"
	"scav/utils"
)

func CreateWebhook(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var wh deliveries.Webhook
		if err := json.NewDecoder(r.Body).Decode(&wh); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid body")
			return
		}

		wh.WebhookID = utils.GenerateRandomString(18)
		wh.TenantID = deliveries.GetTenantIDFromContext(r.Context())
		wh.CreatedAt = time.Now()

		ctx := r.Context()
		if err := createWebhook(ctx, app, wh); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create webhook")
			return
		}
		utils.RespondWithJSON(w, http.StatusCreated, wh)
	}
}

func ListWebhooks(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		ctx := r.Context()

		webhooks, err := listWebhooksForTenant(ctx, app, tenantID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to list webhooks")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, webhooks)
	}
}

func GetWebhook(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		whID := utils.GetParam(r, "webhookid")
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		ctx := r.Context()

		wh, err := getWebhookByID(ctx, app, whID, tenantID)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Webhook not found")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, wh)
	}
}

func UpdateWebhook(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		whID := utils.GetParam(r, "webhookid")
		tenantID := deliveries.GetTenantIDFromContext(r.Context())

		var updates map[string]any
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid payload")
			return
		}

		ctx := r.Context()
		delete(updates, "_id")
		delete(updates, "tenantid")

		if err := updateWebhookByID(ctx, app, whID, tenantID, updates); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update webhook")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

func DeleteWebhook(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		whID := utils.GetParam(r, "webhookid")
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		ctx := r.Context()

		count, err := deleteWebhookByID(ctx, app, whID, tenantID)
		if err != nil || count == 0 {
			utils.RespondWithError(w, http.StatusNotFound, "Webhook not found or failed to delete")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

func TestWebhook(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		whID := utils.GetParam(r, "webhookid")

		payload := map[string]string{"event": "ping", "webhookid": whID}
		data, _ := json.Marshal(payload)
		_ = app.NatsConn.Publish("webhooks.test", data)

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "test_triggered"})
	}
}
