package delwebhooks

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/internal/deliveries"
)

var webhooksTable = config.Tables.DeliveryWebhooksTable

func sqlcreateWebhook(ctx context.Context, app *infra.Deps, wh deliveries.Webhook) error {
	return app.SQLDB.InsertOne(ctx, webhooksTable, wh)
}

func sqllistWebhooksForTenant(ctx context.Context, app *infra.Deps, tenantID string) ([]deliveries.Webhook, error) {
	var webhooks []deliveries.Webhook
	where := "tenantid = $1"
	args := []any{tenantID}

	if err := app.SQLDB.FindMany(ctx, webhooksTable, where, args, &webhooks); err != nil {
		return nil, err
	}
	if len(webhooks) == 0 {
		return []deliveries.Webhook{}, nil
	}
	return webhooks, nil
}

func sqlgetWebhookByID(ctx context.Context, app *infra.Deps, whID, tenantID string) (deliveries.Webhook, error) {
	var wh deliveries.Webhook
	where := "id = $1 AND tenantid = $2"
	args := []any{whID, tenantID}

	if err := app.SQLDB.FindOne(ctx, webhooksTable, where, args, &wh); err != nil {
		return deliveries.Webhook{}, err
	}
	return wh, nil
}

func sqlupdateWebhookByID(ctx context.Context, app *infra.Deps, whID, tenantID string, updates map[string]any) error {
	where := "id = $1 AND tenantid = $2"
	args := []any{whID, tenantID}

	_, err := app.SQLDB.UpdateOne(ctx, webhooksTable, where, args, updates)
	return err
}

func sqldeleteWebhookByID(ctx context.Context, app *infra.Deps, whID, tenantID string) (int64, error) {
	where := "id = $1 AND tenantid = $2"
	args := []any{whID, tenantID}

	return app.SQLDB.DeleteOne(ctx, webhooksTable, where, args)
}
