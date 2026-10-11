// File: internal/deliveries/tracking/trackingSQLDB.go

package tracking

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/internal/deliveries"
)

var deliveriesTable = config.Tables.DeliveriesTable
var deliveryEventsTable = config.Tables.DeliveryEventsTable

func getTrackingDetails(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) (map[string]any, error) {
	var result map[string]any

	return result, nil
}

func getDeliveryEvents(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) ([]map[string]any, error) {
	var events []map[string]any

	return events, nil
}

func getStatusHistory(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) ([]deliveries.StatusHistoryItem, error) {
	var res struct {
		StatusHistory []deliveries.StatusHistoryItem `db:"status_history" json:"status_history"`
	}

	return res.StatusHistory, nil
}

func addProofToDelivery(ctx context.Context, app *infra.Deps, deliveryID, tenantID string, proof deliveries.Proof) error {
	return nil
}

func getProofs(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) ([]deliveries.Proof, error) {
	var res struct {
		Proofs []deliveries.Proof `db:"proofs" json:"proofs"`
	}

	return res.Proofs, nil
}

func getPublicTrackingInfo(ctx context.Context, app *infra.Deps, token string) (map[string]any, error) {
	var res map[string]any

	return res, nil
}
