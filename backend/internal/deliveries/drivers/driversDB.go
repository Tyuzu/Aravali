package drivers

import (
	"context"
	"fmt"
	"time"

	"scav/infra"
	"scav/internal/deliveries"
)

const driversCollection = "drivers"
const driverJobRejectionsCollection = "driver_job_rejections"

func getDriverProfileByID(ctx context.Context, app *infra.Deps, driverID, tenantID string) (deliveries.Driver, error) {
	var driver deliveries.Driver
	filter := map[string]any{"id": driverID, "tenantid": tenantID}
	if err := app.DB.FindOne(ctx, driversCollection, filter, &driver); err != nil {
		return deliveries.Driver{}, err
	}
	return driver, nil
}

func updateDriverProfile(ctx context.Context, app *infra.Deps, driverID, tenantID string, updates map[string]any) error {
	filter := map[string]any{"id": driverID, "tenantid": tenantID}
	_, err := app.DB.UpdateOne(ctx, driversCollection, filter, map[string]any{"$set": updates})
	return err
}

func setDriverOnlineState(ctx context.Context, app *infra.Deps, driverID, tenantID string, online bool) error {
	filter := map[string]any{"id": driverID, "tenantid": tenantID}
	_, err := app.DB.UpdateOne(ctx, driversCollection, filter, map[string]any{"$set": map[string]any{"is_online": online}})
	return err
}

func getDriverStatus(ctx context.Context, app *infra.Deps, driverID, tenantID string) (map[string]any, error) {
	var status map[string]any
	filter := map[string]any{"id": driverID, "tenantid": tenantID}
	if err := app.DB.FindOneWithProjection(ctx, driversCollection, filter, []string{"is_online", "current_state"}, &status); err != nil {
		return nil, err
	}
	return status, nil
}

func getAvailableJobsForTenant(ctx context.Context, app *infra.Deps, tenantID string) ([]deliveries.Delivery, error) {
	var jobs []deliveries.Delivery
	filter := map[string]any{
		"status":   deliveries.StatusCreated,
		"driverid": nil,
		"tenantid": tenantID,
	}
	if err := app.DB.FindMany(ctx, "deliveries", filter, &jobs); err != nil {
		return nil, err
	}
	if len(jobs) == 0 {
		return []deliveries.Delivery{}, nil
	}
	return jobs, nil
}

func getActiveJobsForDriver(ctx context.Context, app *infra.Deps, driverID, tenantID string) ([]deliveries.Delivery, error) {
	var active []deliveries.Delivery
	filter := map[string]any{
		"driverid": driverID,
		"tenantid": tenantID,
		"status":   map[string]any{"$in": []string{deliveries.StatusAssigned, deliveries.StatusAccepted, deliveries.StatusPickedUp, deliveries.StatusInTransit}},
	}
	if err := app.DB.FindMany(ctx, "deliveries", filter, &active); err != nil {
		return nil, err
	}
	if len(active) == 0 {
		return []deliveries.Delivery{}, nil
	}
	return active, nil
}

func findDeliveryForDriver(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) (deliveries.Delivery, error) {
	var current deliveries.Delivery
	filter := map[string]any{"id": deliveryID, "tenantid": tenantID}
	if err := app.DB.FindOne(ctx, "deliveries", filter, &current); err != nil {
		return deliveries.Delivery{}, err
	}
	return current, nil
}

func saveDriverRejection(ctx context.Context, app *infra.Deps, tenantID, driverID, deliveryID string) error {
	return app.DB.InsertOne(ctx, driverJobRejectionsCollection, map[string]any{
		"rejectionid": time.Now().UnixNano(),
		"tenantid":    tenantID,
		"driverid":    driverID,
		"deliveryid":  deliveryID,
		"rejected_at": time.Now(),
	})
}

func claimDeliveryAssignment(ctx context.Context, app *infra.Deps, deliveryID, tenantID, driverID string) (deliveries.Delivery, error) {
	var current deliveries.Delivery
	filter := map[string]any{"id": deliveryID, "tenantid": tenantID}
	if err := app.DB.FindOne(ctx, "deliveries", filter, &current); err != nil {
		return deliveries.Delivery{}, err
	}
	if current.DriverID != nil && *current.DriverID != "" {
		return deliveries.Delivery{}, fmt.Errorf("delivery is already assigned to another driver")
	}
	if err := deliveries.ValidateTransition(current.Status, deliveries.StatusAssigned); err != nil {
		return deliveries.Delivery{}, err
	}

	now := time.Now()
	update := map[string]any{
		"$set": map[string]any{
			"status":     deliveries.StatusAssigned,
			"driverid":   driverID,
			"updated_at": now,
		},
		"$push": map[string]any{
			"status_history": deliveries.StatusHistoryItem{
				Status:    deliveries.StatusAssigned,
				Timestamp: now,
				UpdatedBy: driverID,
			},
		},
	}
	var updated deliveries.Delivery
	if err := app.DB.FindOneAndUpdate(ctx, "deliveries", filter, update, &updated); err != nil {
		return deliveries.Delivery{}, err
	}
	return updated, nil
}

func acceptDeliveryAssignment(ctx context.Context, app *infra.Deps, deliveryID, tenantID, driverID string) (deliveries.Delivery, error) {
	var current deliveries.Delivery
	filter := map[string]any{"id": deliveryID, "tenantid": tenantID}
	if err := app.DB.FindOne(ctx, "deliveries", filter, &current); err != nil {
		return deliveries.Delivery{}, err
	}
	if err := deliveries.ValidateTransition(current.Status, deliveries.StatusAccepted); err != nil {
		return deliveries.Delivery{}, err
	}

	now := time.Now()
	update := map[string]any{
		"$set": map[string]any{
			"status":     deliveries.StatusAccepted,
			"driverid":   driverID,
			"updated_at": now,
		},
		"$push": map[string]any{
			"status_history": deliveries.StatusHistoryItem{
				Status:    deliveries.StatusAccepted,
				Timestamp: now,
				UpdatedBy: driverID,
			},
		},
	}
	var updated deliveries.Delivery
	if err := app.DB.FindOneAndUpdate(ctx, "deliveries", filter, update, &updated); err != nil {
		return deliveries.Delivery{}, err
	}
	return updated, nil
}
