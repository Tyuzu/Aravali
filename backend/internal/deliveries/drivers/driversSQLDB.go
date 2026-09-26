package drivers

import (
	"context"
	"fmt"
	"time"

	"scav/config"
	"scav/infra"
	"scav/internal/deliveries"
)

var driversTable = config.Tables.DriversTable
var driverJobRejectionsTable = config.Tables.DriverJobRejectionsTable

func sqlgetDriverProfileByID(ctx context.Context, app *infra.Deps, driverID, tenantID string) (deliveries.Driver, error) {
	var driver deliveries.Driver
	where := "id = $1 AND tenantid = $2"
	args := []any{driverID, tenantID}

	if err := app.SQLDB.FindOne(ctx, driversTable, where, args, &driver); err != nil {
		return deliveries.Driver{}, err
	}
	return driver, nil
}

func sqlupdateDriverProfile(ctx context.Context, app *infra.Deps, driverID, tenantID string, updates map[string]any) error {
	where := "id = $1 AND tenantid = $2"
	args := []any{driverID, tenantID}

	_, err := app.SQLDB.UpdateOne(ctx, driversTable, where, args, updates)
	return err
}

func sqlsetDriverOnlineState(ctx context.Context, app *infra.Deps, driverID, tenantID string, online bool) error {
	where := "id = $1 AND tenantid = $2"
	args := []any{driverID, tenantID}
	updates := map[string]any{"is_online": online}

	_, err := app.SQLDB.UpdateOne(ctx, driversTable, where, args, updates)
	return err
}

func sqlgetDriverStatus(ctx context.Context, app *infra.Deps, driverID, tenantID string) (map[string]any, error) {
	var status map[string]any
	where := "id = $1 AND tenantid = $2"
	args := []any{driverID, tenantID}
	columns := []string{"is_online", "current_state"}

	if err := app.SQLDB.FindOneWithProjection(ctx, driversTable, columns, where, args, &status); err != nil {
		return nil, err
	}
	return status, nil
}

func sqlgetAvailableJobsForTenant(ctx context.Context, app *infra.Deps, tenantID string) ([]deliveries.Delivery, error) {
	var jobs []deliveries.Delivery
	where := "status = $1 AND driverid IS NULL AND tenantid = $2"
	args := []any{deliveries.StatusCreated, tenantID}

	if err := app.SQLDB.FindMany(ctx, "deliveries", where, args, &jobs); err != nil {
		return nil, err
	}
	if len(jobs) == 0 {
		return []deliveries.Delivery{}, nil
	}
	return jobs, nil
}

func sqlgetActiveJobsForDriver(ctx context.Context, app *infra.Deps, driverID, tenantID string) ([]deliveries.Delivery, error) {
	var active []deliveries.Delivery
	where := "driverid = $1 AND tenantid = $2 AND status = ANY($3)"
	statuses := []string{
		deliveries.StatusAssigned,
		deliveries.StatusAccepted,
		deliveries.StatusPickedUp,
		deliveries.StatusInTransit,
	}
	args := []any{driverID, tenantID, statuses}

	if err := app.SQLDB.FindMany(ctx, "deliveries", where, args, &active); err != nil {
		return nil, err
	}
	if len(active) == 0 {
		return []deliveries.Delivery{}, nil
	}
	return active, nil
}

func sqlfindDeliveryForDriver(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) (deliveries.Delivery, error) {
	var current deliveries.Delivery
	where := "id = $1 AND tenantid = $2"
	args := []any{deliveryID, tenantID}

	if err := app.SQLDB.FindOne(ctx, "deliveries", where, args, &current); err != nil {
		return deliveries.Delivery{}, err
	}
	return current, nil
}

func sqlsaveDriverRejection(ctx context.Context, app *infra.Deps, tenantID, driverID, deliveryID string) error {
	record := map[string]any{
		"rejectionid": time.Now().UnixNano(),
		"tenantid":    tenantID,
		"driverid":    driverID,
		"deliveryid":  deliveryID,
		"rejected_at": time.Now(),
	}
	return app.SQLDB.InsertOne(ctx, driverJobRejectionsTable, record)
}

func sqlclaimDeliveryAssignment(ctx context.Context, app *infra.Deps, deliveryID, tenantID, driverID string) (deliveries.Delivery, error) {
	var current deliveries.Delivery
	where := "id = $1 AND tenantid = $2"
	args := []any{deliveryID, tenantID}

	if err := app.SQLDB.FindOne(ctx, "deliveries", where, args, &current); err != nil {
		return deliveries.Delivery{}, err
	}
	if current.DriverID != nil && *current.DriverID != "" {
		return deliveries.Delivery{}, fmt.Errorf("delivery is already assigned to another driver")
	}
	if err := deliveries.ValidateTransition(current.Status, deliveries.StatusAssigned); err != nil {
		return deliveries.Delivery{}, err
	}

	now := time.Now()
	newHistoryItem := deliveries.StatusHistoryItem{
		Status:    deliveries.StatusAssigned,
		Timestamp: now,
		UpdatedBy: driverID,
	}

	updatedHistory := append(current.StatusHistory, newHistoryItem)
	updates := map[string]any{
		"status":         deliveries.StatusAssigned,
		"driverid":       driverID,
		"updated_at":     now,
		"status_history": updatedHistory,
	}

	if _, err := app.SQLDB.UpdateOne(ctx, "deliveries", where, args, updates); err != nil {
		return deliveries.Delivery{}, err
	}

	var updated deliveries.Delivery
	if err := app.SQLDB.FindOne(ctx, "deliveries", where, args, &updated); err != nil {
		return deliveries.Delivery{}, err
	}
	return updated, nil
}

func sqlacceptDeliveryAssignment(ctx context.Context, app *infra.Deps, deliveryID, tenantID, driverID string) (deliveries.Delivery, error) {
	var current deliveries.Delivery
	where := "id = $1 AND tenantid = $2"
	args := []any{deliveryID, tenantID}

	if err := app.SQLDB.FindOne(ctx, "deliveries", where, args, &current); err != nil {
		return deliveries.Delivery{}, err
	}
	if err := deliveries.ValidateTransition(current.Status, deliveries.StatusAccepted); err != nil {
		return deliveries.Delivery{}, err
	}

	now := time.Now()
	newHistoryItem := deliveries.StatusHistoryItem{
		Status:    deliveries.StatusAccepted,
		Timestamp: now,
		UpdatedBy: driverID,
	}

	updatedHistory := append(current.StatusHistory, newHistoryItem)
	updates := map[string]any{
		"status":         deliveries.StatusAccepted,
		"driverid":       driverID,
		"updated_at":     now,
		"status_history": updatedHistory,
	}

	if _, err := app.SQLDB.UpdateOne(ctx, "deliveries", where, args, updates); err != nil {
		return deliveries.Delivery{}, err
	}

	var updated deliveries.Delivery
	if err := app.SQLDB.FindOne(ctx, "deliveries", where, args, &updated); err != nil {
		return deliveries.Delivery{}, err
	}
	return updated, nil
}
