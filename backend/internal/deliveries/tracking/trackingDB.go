package tracking

import (
	"context"

	"scav/infra"
	"scav/internal/deliveries"
)

const deliveriesCollection = "deliveries"
const deliveryEventsCollection = "delivery_events"

func getTrackingDetails(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) (map[string]any, error) {
	var result map[string]any
	filter := map[string]any{"id": deliveryID, "tenantid": tenantID}
	proj := []string{"status", "status_history", "current_location"}
	if err := app.DB.FindOneWithProjection(ctx, deliveriesCollection, filter, proj, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func getDeliveryEvents(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) ([]map[string]any, error) {
	var events []map[string]any
	filter := map[string]any{"deliveryid": deliveryID, "tenantid": tenantID}
	if err := app.DB.FindMany(ctx, deliveryEventsCollection, filter, &events); err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return []map[string]any{}, nil
	}
	return events, nil
}

func getStatusHistory(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) ([]deliveries.StatusHistoryItem, error) {
	var res struct {
		StatusHistory []deliveries.StatusHistoryItem `bson:"status_history" json:"status_history"`
	}
	filter := map[string]any{"id": deliveryID, "tenantid": tenantID}
	if err := app.DB.FindOneWithProjection(ctx, deliveriesCollection, filter, []string{"status_history"}, &res); err != nil {
		return nil, err
	}
	return res.StatusHistory, nil
}

func addProofToDelivery(ctx context.Context, app *infra.Deps, deliveryID, tenantID string, proof deliveries.Proof) error {
	filter := map[string]any{"id": deliveryID, "tenantid": tenantID}
	return app.DB.AddToSet(ctx, deliveriesCollection, filter, "proofs", proof)
}

func getProofs(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) ([]deliveries.Proof, error) {
	var res struct {
		Proofs []deliveries.Proof `bson:"proofs" json:"proofs"`
	}
	filter := map[string]any{"id": deliveryID, "tenantid": tenantID}
	if err := app.DB.FindOneWithProjection(ctx, deliveriesCollection, filter, []string{"proofs"}, &res); err != nil {
		return nil, err
	}
	return res.Proofs, nil
}

func getPublicTrackingInfo(ctx context.Context, app *infra.Deps, token string) (map[string]any, error) {
	var res map[string]any
	filter := map[string]any{"public_tracking_token": token}
	proj := []string{"status", "pickup_loc", "dropoff_loc", "estimated_arrival"}
	if err := app.DB.FindOneWithProjection(ctx, deliveriesCollection, filter, proj, &res); err != nil {
		return nil, err
	}
	return res, nil
}
