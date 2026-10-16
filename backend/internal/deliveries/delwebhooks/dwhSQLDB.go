// File: internal/deliveries/delwebhooks/dwhSQLDB.go

package delwebhooks

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/internal/deliveries"
)

var webhooksTable = config.Tables.DeliveryWebhooksTable

func createWebhook(_ context.Context, _ *infra.Deps, _ deliveries.Webhook) error {
	return nil
}

func listWebhooksForTenant(_ context.Context, _ *infra.Deps, _ string) ([]deliveries.Webhook, error) {
	var webhooks []deliveries.Webhook
	return webhooks, nil
}

func getWebhookByID(_ context.Context, _ *infra.Deps, _, _ string) (deliveries.Webhook, error) {
	var wh deliveries.Webhook
	return wh, nil
}

func updateWebhookByID(_ context.Context, _ *infra.Deps, _, _ string, _ map[string]any) error {
	return nil
}

func deleteWebhookByID(_ context.Context, _ *infra.Deps, _, _ string) (int64, error) {
	return 0, nil
}
