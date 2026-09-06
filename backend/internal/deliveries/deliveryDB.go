package deliveries

import (
	"context"
	"fmt"
	"time"

	"scav/infra"
)

const deliveriesCollection = "deliveries"

func findDeliveryByFilter(ctx context.Context, app *infra.Deps, filter map[string]any, out any) error {
	return app.DB.FindOne(ctx, deliveriesCollection, filter, out)
}

func findMyDeliveries(ctx context.Context, app *infra.Deps, userID, tenantID string) ([]Delivery, error) {
	filter := map[string]any{"userid": userID, "tenantid": tenantID}
	var deliveries []Delivery
	if err := app.DB.FindMany(ctx, deliveriesCollection, filter, &deliveries); err != nil {
		return nil, err
	}
	if len(deliveries) == 0 {
		return []Delivery{}, nil
	}
	return deliveries, nil
}

func saveDelivery(ctx context.Context, app *infra.Deps, delivery Delivery) error {
	return app.DB.InsertOne(ctx, deliveriesCollection, delivery)
}

func findDeliveryAndUpdate(ctx context.Context, app *infra.Deps, filter map[string]any, update map[string]any) (Delivery, error) {
	var updated Delivery
	if err := app.DB.FindOneAndUpdate(ctx, deliveriesCollection, filter, update, &updated); err != nil {
		return Delivery{}, err
	}
	return updated, nil
}

func fetchDeliveryForRead(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) (Delivery, error) {
	var delivery Delivery
	filter := map[string]any{"id": deliveryID, "tenantid": tenantID}
	if err := app.DB.FindOne(ctx, deliveriesCollection, filter, &delivery); err != nil {
		return Delivery{}, err
	}
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
	update := map[string]any{
		"$set": map[string]any{
			"status":     newStatus,
			"updated_at": now,
		},
		"$push": map[string]any{
			"status_history": StatusHistoryItem{
				Status:    newStatus,
				Timestamp: now,
				UpdatedBy: userID,
			},
		},
	}
	updated, err := findDeliveryAndUpdate(ctx, app, map[string]any{"id": deliveryID, "tenantid": tenantID}, update)
	if err != nil {
		return nil, err
	}
	_ = app.Cache.Del(ctx, fmt.Sprintf("delivery:%s", deliveryID))
	_ = app.NatsConn.Publish(fmt.Sprintf("deliveries.status.%s", newStatus), []byte(deliveryID))
	return &updated, nil
}

func findDeliveryByIDTenant(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) (Delivery, error) {
	return fetchDeliveryForRead(ctx, app, deliveryID, tenantID)
}

func listDeliveryEvents(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) ([]map[string]any, error) {
	var events []map[string]any
	filter := map[string]any{"deliveryid": deliveryID, "tenantid": tenantID}
	if err := app.DB.FindMany(ctx, "delivery_events", filter, &events); err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return []map[string]any{}, nil
	}
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
	update := map[string]any{
		"$set": map[string]any{
			"driverid":   driverID,
			"status":     StatusAssigned,
			"updated_at": now,
		},
		"$push": map[string]any{
			"status_history": StatusHistoryItem{
				Status:    StatusAssigned,
				Timestamp: now,
				UpdatedBy: userID,
			},
		},
	}
	return findDeliveryAndUpdate(ctx, app, map[string]any{"id": deliveryID, "tenantid": tenantID}, update)
}

func updateDeliveryStatusRecord(ctx context.Context, app *infra.Deps, deliveryID, tenantID, userID, newStatus string) (*Delivery, error) {
	var currentDelivery Delivery
	filter := map[string]any{"id": deliveryID, "tenantid": tenantID}
	if err := app.DB.FindOne(ctx, deliveriesCollection, filter, &currentDelivery); err != nil {
		return nil, fmt.Errorf("delivery not found")
	}
	if err := ValidateTransition(currentDelivery.Status, newStatus); err != nil {
		return nil, err
	}

	now := time.Now()
	update := map[string]any{
		"$set": map[string]any{
			"status":     newStatus,
			"updated_at": now,
		},
		"$push": map[string]any{
			"status_history": StatusHistoryItem{
				Status:    newStatus,
				Timestamp: now,
				UpdatedBy: userID,
			},
		},
	}

	var updatedDelivery Delivery
	if err := app.DB.FindOneAndUpdate(ctx, deliveriesCollection, filter, update, &updatedDelivery); err != nil {
		return nil, err
	}
	_ = app.Cache.Del(ctx, fmt.Sprintf("delivery:%s", deliveryID))
	_ = app.NatsConn.Publish(fmt.Sprintf("deliveries.status.%s", newStatus), []byte(deliveryID))
	return &updatedDelivery, nil
}
