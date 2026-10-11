// File: internal/deliveries/deliverySQLDB.go

package deliveries

import (
	"context"
	"fmt"
	"time"

	"scav/config"
	"scav/infra"
)

var deliveriesTable = config.Tables.DeliveriesTable

func findDeliveryByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {

}

func findMyDeliveries(ctx context.Context, app *infra.Deps, userID, tenantID string) ([]Delivery, error) {

	var deliveries []Delivery

	return deliveries, nil
}

func saveDelivery(ctx context.Context, app *infra.Deps, delivery Delivery) error {

}

func findDeliveryAndUpdate(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (Delivery, error) {

	var updated Delivery
	return updated, nil
}

func fetchDeliveryForRead(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) (Delivery, error) {
	var delivery Delivery
	return delivery, nil
}

func setDeliveryStatusWithHistory(ctx context.Context, app *infra.Deps, deliveryID, tenantID, userID, newStatus string) (*Delivery, error) {
	current, err := fetchDeliveryForRead(ctx, app, deliveryID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("delivery not found")
	}
	if err := ValidateTransition(current.Status, newStatus); err != nil {
		return nil, err
	}

	now := time.Now()
	newHistoryItem := StatusHistoryItem{
		Status:    newStatus,
		Timestamp: now,
		UpdatedBy: userID,
	}

	updatedHistory := append(current.StatusHistory, newHistoryItem)

	update := map[string]any{
		"status":         newStatus,
		"updated_at":     now,
		"status_history": updatedHistory,
	}

	query := "id = $1 AND tenantid = $2"
	args := []any{deliveryID, tenantID}

	updated, err := findDeliveryAndUpdate(ctx, app, query, args, update)
	if err != nil {
		return nil, err
	}

	_ = app.Cache.Del(ctx, fmt.Sprintf("delivery:%s", deliveryID))
	_ = app.MQ.Publish(ctx, fmt.Sprintf("deliveries.status.%s", newStatus), []byte(deliveryID))
	return &updated, nil
}

func findDeliveryByIDTenant(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) (Delivery, error) {
	return fetchDeliveryForRead(ctx, app, deliveryID, tenantID)
}

func listDeliveryEvents(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) ([]map[string]any, error) {
	var events []map[string]any
	return events, nil
}

func upsertDeliveryAssignment(ctx context.Context, app *infra.Deps, deliveryID, tenantID, driverID, userID string) (Delivery, error) {
	current, err := fetchDeliveryForRead(ctx, app, deliveryID, tenantID)
	if err != nil {
		return Delivery{}, fmt.Errorf("delivery not found")
	}
	if err := ValidateTransition(current.Status, StatusAssigned); err != nil {
		return Delivery{}, err
	}

	now := time.Now()
	newHistoryItem := StatusHistoryItem{
		Status:    StatusAssigned,
		Timestamp: now,
		UpdatedBy: userID,
	}

	updatedHistory := append(current.StatusHistory, newHistoryItem)

	update := map[string]any{
		"driverid":       driverID,
		"status":         StatusAssigned,
		"updated_at":     now,
		"status_history": updatedHistory,
	}

	query := "id = $1 AND tenantid = $2"
	args := []any{deliveryID, tenantID}

	return findDeliveryAndUpdate(ctx, app, query, args, update)
}

func updateDeliveryStatusRecord(ctx context.Context, app *infra.Deps, deliveryID, tenantID, userID, newStatus string) (*Delivery, error) {
	var currentDelivery Delivery
	if err := ValidateTransition(currentDelivery.Status, newStatus); err != nil {
		return nil, err
	}

	now := time.Now()
	newHistoryItem := StatusHistoryItem{
		Status:    newStatus,
		Timestamp: now,
		UpdatedBy: userID,
	}

	updatedHistory := append(currentDelivery.StatusHistory, newHistoryItem)

	update := map[string]any{
		"status":         newStatus,
		"updated_at":     now,
		"status_history": updatedHistory,
	}

	updatedDelivery, err := findDeliveryAndUpdate(ctx, app, query, args, update)
	if err != nil {
		return nil, err
	}

	_ = app.Cache.Del(ctx, fmt.Sprintf("delivery:%s", deliveryID))
	_ = app.MQ.Publish(ctx, fmt.Sprintf("deliveries.status.%s", newStatus), []byte(deliveryID))
	return &updatedDelivery, nil
}
