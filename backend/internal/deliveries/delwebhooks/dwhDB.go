package delwebhooks

import (
	"context"

	"scav/infra"
	"scav/internal/deliveries"
)

const webhooksCollection = "webhooks"

func createWebhook(ctx context.Context, app *infra.Deps, wh deliveries.Webhook) error {
	return app.DB.InsertOne(ctx, webhooksCollection, wh)
}

func listWebhooksForTenant(ctx context.Context, app *infra.Deps, tenantID string) ([]deliveries.Webhook, error) {
	var webhooks []deliveries.Webhook
	if err := app.DB.FindMany(ctx, webhooksCollection, map[string]any{"tenantid": tenantID}, &webhooks); err != nil {
		return nil, err
	}
	if len(webhooks) == 0 {
		return []deliveries.Webhook{}, nil
	}
	return webhooks, nil
}

func getWebhookByID(ctx context.Context, app *infra.Deps, whID, tenantID string) (deliveries.Webhook, error) {
	var wh deliveries.Webhook
	filter := map[string]any{"id": whID, "tenantid": tenantID}
	if err := app.DB.FindOne(ctx, webhooksCollection, filter, &wh); err != nil {
		return deliveries.Webhook{}, err
	}
	return wh, nil
}

func updateWebhookByID(ctx context.Context, app *infra.Deps, whID, tenantID string, updates map[string]any) error {
	filter := map[string]any{"_id": whID, "tenantid": tenantID}
	_, err := app.DB.UpdateOne(ctx, webhooksCollection, filter, map[string]any{"$set": updates})
	return err
}

func deleteWebhookByID(ctx context.Context, app *infra.Deps, whID, tenantID string) (int64, error) {
	filter := map[string]any{"_id": whID, "tenantid": tenantID}
	return app.DB.DeleteOne(ctx, webhooksCollection, filter)
}
