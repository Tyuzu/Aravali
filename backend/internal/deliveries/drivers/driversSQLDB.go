// File: internal/deliveries/drivers/driversSQLDB.go

package drivers

import (
	"context"
	"fmt"

	"scav/config"
	"scav/infra"
	"scav/internal/deliveries"
)

var driversTable = config.Tables.DriversTable
var driverJobRejectionsTable = config.Tables.DriverJobRejectionsTable

func getDriverProfileByID(ctx context.Context, app *infra.Deps, driverID, tenantID string) (deliveries.Driver, error) {
	var driver deliveries.Driver

	return driver, nil
}

func updateDriverProfile(ctx context.Context, app *infra.Deps, driverID, tenantID string, updates map[string]any) error {
	return nil
}

func setDriverOnlineState(ctx context.Context, app *infra.Deps, driverID, tenantID string, online bool) error {
	return nil
}

func getDriverStatus(ctx context.Context, app *infra.Deps, driverID, tenantID string) (map[string]any, error) {
	var status map[string]any

	return status, nil
}

func getAvailableJobsForTenant(ctx context.Context, app *infra.Deps, tenantID string) ([]deliveries.Delivery, error) {
	var jobs []deliveries.Delivery

	return jobs, nil
}

func getActiveJobsForDriver(ctx context.Context, app *infra.Deps, driverID, tenantID string) ([]deliveries.Delivery, error) {
	var active []deliveries.Delivery

	return active, nil
}

func findDeliveryForDriver(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) (deliveries.Delivery, error) {
	var current deliveries.Delivery

	return current, nil
}

func saveDriverRejection(ctx context.Context, app *infra.Deps, tenantID, driverID, deliveryID string) error {
	return nil
}

func claimDeliveryAssignment(ctx context.Context, app *infra.Deps, deliveryID, tenantID, driverID string) (deliveries.Delivery, error) {
	var current deliveries.Delivery

	if current.DriverID != nil && *current.DriverID != "" {
		return deliveries.Delivery{}, fmt.Errorf("delivery is already assigned to another driver")
	}
	if err := deliveries.ValidateTransition(current.Status, deliveries.StatusAssigned); err != nil {
		return deliveries.Delivery{}, err
	}

	var updated deliveries.Delivery

	return updated, nil
}

func acceptDeliveryAssignment(ctx context.Context, app *infra.Deps, deliveryID, tenantID, driverID string) (deliveries.Delivery, error) {
	var updated deliveries.Delivery
	return updated, nil
}
