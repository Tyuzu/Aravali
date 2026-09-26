package deliveries

import (
	"context"
	"fmt"
	"time"

	"scav/config"
	"scav/infra"
)

var deliveriesTable = config.Tables.DeliveriesTable

func SQLfindDeliveryByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindOne(ctx, deliveriesTable, query, args, out)
}

func SQLfindMyDeliveries(ctx context.Context, app *infra.Deps, userID, tenantID string) ([]Delivery, error) {
	query := "userid = $1 AND tenantid = $2"
	args := []any{userID, tenantID}

	var deliveries []Delivery
	if err := app.SQLDB.FindMany(ctx, deliveriesTable, query, args, &deliveries); err != nil {
		return nil, err
	}
	if len(deliveries) == 0 {
		return []Delivery{}, nil
	}
	return deliveries, nil
}

func SQLsaveDelivery(ctx context.Context, app *infra.Deps, delivery Delivery) error {
	return app.SQLDB.InsertOne(ctx, deliveriesTable, delivery)
}

func SQLfindDeliveryAndUpdate(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (Delivery, error) {
	if _, err := app.SQLDB.UpdateOne(ctx, deliveriesTable, query, args, update); err != nil {
		return Delivery{}, err
	}

	var updated Delivery
	if err := app.SQLDB.FindOne(ctx, deliveriesTable, query, args, &updated); err != nil {
		return Delivery{}, err
	}
	return updated, nil
}

func SQLfetchDeliveryForRead(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) (Delivery, error) {
	var delivery Delivery
	query := "id = $1 AND tenantid = $2"
	args := []any{deliveryID, tenantID}

	if err := app.SQLDB.FindOne(ctx, deliveriesTable, query, args, &delivery); err != nil {
		return Delivery{}, err
	}
	return delivery, nil
}

func SQLsetDeliveryStatusWithHistory(ctx context.Context, app *infra.Deps, deliveryID, tenantID, userID, newStatus string) (*Delivery, error) {
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

	updated, err := SQLfindDeliveryAndUpdate(ctx, app, query, args, update)
	if err != nil {
		return nil, err
	}

	_ = app.Cache.Del(ctx, fmt.Sprintf("delivery:%s", deliveryID))
	_ = app.NatsConn.Publish(fmt.Sprintf("deliveries.status.%s", newStatus), []byte(deliveryID))
	return &updated, nil
}

func SQLfindDeliveryByIDTenant(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) (Delivery, error) {
	return fetchDeliveryForRead(ctx, app, deliveryID, tenantID)
}

func SQLlistDeliveryEvents(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) ([]map[string]any, error) {
	var events []map[string]any
	query := "deliveryid = $1 AND tenantid = $2"
	args := []any{deliveryID, tenantID}

	if err := app.SQLDB.FindMany(ctx, "delivery_events", query, args, &events); err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return []map[string]any{}, nil
	}
	return events, nil
}

func SQLupsertDeliveryAssignment(ctx context.Context, app *infra.Deps, deliveryID, tenantID, driverID, userID string) (Delivery, error) {
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

	return SQLfindDeliveryAndUpdate(ctx, app, query, args, update)
}

func SQLupdateDeliveryStatusRecord(ctx context.Context, app *infra.Deps, deliveryID, tenantID, userID, newStatus string) (*Delivery, error) {
	var currentDelivery Delivery
	query := "id = $1 AND tenantid = $2"
	args := []any{deliveryID, tenantID}

	if err := app.SQLDB.FindOne(ctx, deliveriesTable, query, args, &currentDelivery); err != nil {
		return nil, fmt.Errorf("delivery not found")
	}
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

	updatedDelivery, err := SQLfindDeliveryAndUpdate(ctx, app, query, args, update)
	if err != nil {
		return nil, err
	}

	_ = app.Cache.Del(ctx, fmt.Sprintf("delivery:%s", deliveryID))
	_ = app.NatsConn.Publish(fmt.Sprintf("deliveries.status.%s", newStatus), []byte(deliveryID))
	return &updatedDelivery, nil
}
