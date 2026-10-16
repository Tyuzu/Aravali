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

func getTrackingDetails(_ context.Context, _ *infra.Deps, _, _ string) (map[string]any, error) {
	var result map[string]any

	return result, nil
}

func getDeliveryEvents(_ context.Context, _ *infra.Deps, _, _ string) ([]map[string]any, error) {
	var events []map[string]any

	return events, nil
}

func getStatusHistory(_ context.Context, _ *infra.Deps, _, _ string) ([]deliveries.StatusHistoryItem, error) {
	var res struct {
		StatusHistory []deliveries.StatusHistoryItem `db:"status_history" json:"status_history"`
	}

	return res.StatusHistory, nil
}

func addProofToDelivery(_ context.Context, _ *infra.Deps, _, _ string, _ deliveries.Proof) error {
	return nil
}

func getProofs(_ context.Context, _ *infra.Deps, _, _ string) ([]deliveries.Proof, error) {
	var res struct {
		Proofs []deliveries.Proof `db:"proofs" json:"proofs"`
	}

	return res.Proofs, nil
}

func getPublicTrackingInfo(_ context.Context, _ *infra.Deps, _ string) (map[string]any, error) {
	var res map[string]any

	return res, nil
}
