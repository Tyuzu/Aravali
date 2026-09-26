package tracking

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/internal/deliveries"
)

var deliveriesTable = config.Tables.DeliveriesTable
var deliveryEventsTable = config.Tables.DeliveryEventsTable

func sqlgetTrackingDetails(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) (map[string]any, error) {
	var result map[string]any
	where := "id = $1 AND tenantid = $2"
	args := []any{deliveryID, tenantID}
	proj := []string{"status", "status_history", "current_location"}

	if err := app.SQLDB.FindOneWithProjection(ctx, deliveriesTable, proj, where, args, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func sqlgetDeliveryEvents(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) ([]map[string]any, error) {
	var events []map[string]any
	where := "deliveryid = $1 AND tenantid = $2"
	args := []any{deliveryID, tenantID}

	if err := app.SQLDB.FindMany(ctx, deliveryEventsTable, where, args, &events); err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return []map[string]any{}, nil
	}
	return events, nil
}

func sqlgetStatusHistory(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) ([]deliveries.StatusHistoryItem, error) {
	var res struct {
		StatusHistory []deliveries.StatusHistoryItem `db:"status_history" json:"status_history"`
	}
	where := "id = $1 AND tenantid = $2"
	args := []any{deliveryID, tenantID}

	if err := app.SQLDB.FindOneWithProjection(ctx, deliveriesTable, []string{"status_history"}, where, args, &res); err != nil {
		return nil, err
	}
	return res.StatusHistory, nil
}

func sqladdProofToDelivery(ctx context.Context, app *infra.Deps, deliveryID, tenantID string, proof deliveries.Proof) error {
	where := "id = $1 AND tenantid = $2"
	args := []any{deliveryID, tenantID}

	var res struct {
		Proofs []deliveries.Proof `db:"proofs" json:"proofs"`
	}
	if err := app.SQLDB.FindOneWithProjection(ctx, deliveriesTable, []string{"proofs"}, where, args, &res); err != nil {
		return err
	}

	updatedProofs := append(res.Proofs, proof)
	update := map[string]any{
		"proofs": updatedProofs,
	}

	_, err := app.SQLDB.UpdateOne(ctx, deliveriesTable, where, args, update)
	return err
}

func sqlgetProofs(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) ([]deliveries.Proof, error) {
	var res struct {
		Proofs []deliveries.Proof `db:"proofs" json:"proofs"`
	}
	where := "id = $1 AND tenantid = $2"
	args := []any{deliveryID, tenantID}

	if err := app.SQLDB.FindOneWithProjection(ctx, deliveriesTable, []string{"proofs"}, where, args, &res); err != nil {
		return nil, err
	}
	return res.Proofs, nil
}

func sqlgetPublicTrackingInfo(ctx context.Context, app *infra.Deps, token string) (map[string]any, error) {
	var res map[string]any
	where := "public_tracking_token = $1"
	args := []any{token}
	proj := []string{"status", "pickup_loc", "dropoff_loc", "estimated_arrival"}

	if err := app.SQLDB.FindOneWithProjection(ctx, deliveriesTable, proj, where, args, &res); err != nil {
		return nil, err
	}
	return res, nil
}
