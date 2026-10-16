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

func getDriverProfileByID(_ context.Context, _ *infra.Deps, _, _ string) (deliveries.Driver, error) {
	var driver deliveries.Driver

	return driver, nil
}

func updateDriverProfile(_ context.Context, _ *infra.Deps, _, _ string, _ map[string]any) error {
	return nil
}

func setDriverOnlineState(_ context.Context, _ *infra.Deps, _, _ string, _ bool) error {
	return nil
}

func getDriverStatus(_ context.Context, _ *infra.Deps, _, _ string) (map[string]any, error) {
	var status map[string]any

	return status, nil
}

func getAvailableJobsForTenant(_ context.Context, _ *infra.Deps, _ string) ([]deliveries.Delivery, error) {
	var jobs []deliveries.Delivery

	return jobs, nil
}

func getActiveJobsForDriver(_ context.Context, _ *infra.Deps, _, _ string) ([]deliveries.Delivery, error) {
	var active []deliveries.Delivery

	return active, nil
}

func findDeliveryForDriver(_ context.Context, _ *infra.Deps, _, _ string) (deliveries.Delivery, error) {
	var current deliveries.Delivery

	return current, nil
}

func saveDriverRejection(_ context.Context, _ *infra.Deps, _, _, _ string) error {
	return nil
}

func claimDeliveryAssignment(_ context.Context, _ *infra.Deps, _, _, _ string) (deliveries.Delivery, error) {
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

func acceptDeliveryAssignment(_ context.Context, _ *infra.Deps, _, _, _ string) (deliveries.Delivery, error) {
	var updated deliveries.Delivery
	return updated, nil
}
