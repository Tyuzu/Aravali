// File: internal/deliveries/delwebhooks/dwhSQLDB.go

package delwebhooks

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/internal/deliveries"
)

var webhooksTable = config.Tables.DeliveryWebhooksTable

func createWebhook(ctx context.Context, app *infra.Deps, wh deliveries.Webhook) error {

}

func listWebhooksForTenant(ctx context.Context, app *infra.Deps, tenantID string) ([]deliveries.Webhook, error) {
	var webhooks []deliveries.Webhook
	return webhooks, nil
}

func getWebhookByID(ctx context.Context, app *infra.Deps, whID, tenantID string) (deliveries.Webhook, error) {
	var wh deliveries.Webhook
	return wh, nil
}

func updateWebhookByID(ctx context.Context, app *infra.Deps, whID, tenantID string, updates map[string]any) error {
}

func deleteWebhookByID(ctx context.Context, app *infra.Deps, whID, tenantID string) (int64, error) {

}
