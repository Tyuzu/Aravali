// File: internal/deliveries/deliverySQLDB.go

package deliveries

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
)

var deliveriesTable = config.Tables.DeliveriesTable

func deliveryMetaMap(d Delivery) map[string]any {
	meta := map[string]any{
		"tenantid":              d.TenantID,
		"userid":                d.UserID,
		"driverid":              d.DriverID,
		"status":                d.Status,
		"pickup_loc":            d.PickupLoc,
		"dropoff_loc":           d.DropoffLoc,
		"current_location":      d.CurrentLocation,
		"proofs":                d.Proofs,
		"status_history":        d.StatusHistory,
		"public_tracking_token": d.PublicTrackingToken,
		"estimated_arrival":     d.EstimatedArrival,
	}
	if d.CreatedAt.IsZero() {
		meta["created_at"] = time.Now().UTC()
	} else {
		meta["created_at"] = d.CreatedAt
	}
	if d.UpdatedAt.IsZero() {
		meta["updated_at"] = time.Now().UTC()
	} else {
		meta["updated_at"] = d.UpdatedAt
	}
	return meta
}

func hydrateDeliveryFromMetadata(metadata []byte, out *Delivery) error {
	if out == nil {
		return nil
	}
	if len(metadata) == 0 || string(metadata) == "null" {
		return nil
	}
	var payload map[string]any
	if err := json.Unmarshal(metadata, &payload); err != nil {
		return err
	}
	if v, ok := payload["tenantid"].(string); ok {
		out.TenantID = v
	}
	if v, ok := payload["userid"].(string); ok {
		out.UserID = v
	}
	if v, ok := payload["driverid"].(string); ok {
		out.DriverID = &v
	}
	if v, ok := payload["status"].(string); ok {
		out.Status = v
	}
	if v, ok := payload["public_tracking_token"].(string); ok {
		out.PublicTrackingToken = v
	}
	if v, ok := payload["pickup_loc"]; ok {
		b, _ := json.Marshal(v)
		_ = json.Unmarshal(b, &out.PickupLoc)
	}
	if v, ok := payload["dropoff_loc"]; ok {
		b, _ := json.Marshal(v)
		_ = json.Unmarshal(b, &out.DropoffLoc)
	}
	if v, ok := payload["current_location"]; ok {
		b, _ := json.Marshal(v)
		var loc *Location
		_ = json.Unmarshal(b, &loc)
		out.CurrentLocation = loc
	}
	if v, ok := payload["status_history"]; ok {
		b, _ := json.Marshal(v)
		_ = json.Unmarshal(b, &out.StatusHistory)
	}
	if v, ok := payload["proofs"]; ok {
		b, _ := json.Marshal(v)
		_ = json.Unmarshal(b, &out.Proofs)
	}
	return nil
}

func findDeliveryByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if strings.TrimSpace(query) == "" {
		query = "1 = 1"
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT deliveryid, orderid, userid, status, driverid, created_at, updated_at, metadata FROM deliveries WHERE `+query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	switch v := out.(type) {
	case *[]Delivery:
		items := make([]Delivery, 0)
		for rows.Next() {
			var d Delivery
			var meta []byte
			var createdAt, updatedAt time.Time
			var orderID, driverID *string
			if err := rows.Scan(&d.DeliveryID, &orderID, &d.UserID, &d.Status, &driverID, &createdAt, &updatedAt, &meta); err != nil {
				return err
			}
			if driverID != nil {
				d.DriverID = driverID
			}
			d.CreatedAt = createdAt
			d.UpdatedAt = updatedAt
			if err := hydrateDeliveryFromMetadata(meta, &d); err != nil {
				return err
			}
			items = append(items, d)
		}
		*v = items
		return rows.Err()
	case *Delivery:
		if !rows.Next() {
			return nil
		}
		var meta []byte
		var createdAt, updatedAt time.Time
		var orderID, driverID *string
		if err := rows.Scan(&v.DeliveryID, &orderID, &v.UserID, &v.Status, &driverID, &createdAt, &updatedAt, &meta); err != nil {
			return err
		}
		if driverID != nil {
			v.DriverID = driverID
		}
		v.CreatedAt = createdAt
		v.UpdatedAt = updatedAt
		return hydrateDeliveryFromMetadata(meta, v)
	default:
		return fmt.Errorf("unsupported delivery output type")
	}
}

func findMyDeliveries(ctx context.Context, app *infra.Deps, userID, tenantID string) ([]Delivery, error) {
	if app == nil || app.SQLDB == nil || userID == "" {
		return []Delivery{}, nil
	}
	query := `SELECT deliveryid, orderid, userid, status, driverid, created_at, updated_at, metadata FROM deliveries WHERE userid = $1`
	args := []any{userID}
	if strings.TrimSpace(tenantID) != "" {
		query += ` AND COALESCE(metadata->>'tenantid', '') = $2`
		args = append(args, tenantID)
	}
	query += ` ORDER BY created_at DESC`
	rows, err := app.SQLDB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Delivery, 0)
	for rows.Next() {
		var d Delivery
		var meta []byte
		var createdAt, updatedAt time.Time
		var orderID, driverID *string
		if err := rows.Scan(&d.DeliveryID, &orderID, &d.UserID, &d.Status, &driverID, &createdAt, &updatedAt, &meta); err != nil {
			return nil, err
		}
		if driverID != nil {
			d.DriverID = driverID
		}
		d.CreatedAt = createdAt
		d.UpdatedAt = updatedAt
		if err := hydrateDeliveryFromMetadata(meta, &d); err != nil {
			return nil, err
		}
		items = append(items, d)
	}
	return items, rows.Err()
}

func saveDelivery(ctx context.Context, app *infra.Deps, delivery Delivery) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	meta, err := json.Marshal(deliveryMetaMap(delivery))
	if err != nil {
		return err
	}
	_, err = app.SQLDB.Exec(
		ctx,
		`INSERT INTO deliveries (deliveryid, orderid, userid, status, driverid, created_at, updated_at, metadata) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (deliveryid) DO UPDATE SET orderid = EXCLUDED.orderid, userid = EXCLUDED.userid, status = EXCLUDED.status, driverid = EXCLUDED.driverid, updated_at = NOW(), metadata = EXCLUDED.metadata`,
		delivery.DeliveryID,
		nil,
		delivery.UserID,
		delivery.Status,
		delivery.DriverID,
		delivery.CreatedAt,
		delivery.UpdatedAt,
		meta,
	)
	return err
}

func findDeliveryAndUpdate(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (Delivery, error) {
	if app == nil || app.SQLDB == nil {
		return Delivery{}, nil
	}
	if len(update) == 0 {
		return Delivery{}, nil
	}
	setParts := make([]string, 0, len(update))
	setArgs := make([]any, 0, len(update)+len(args))
	idx := 1
	for key, value := range update {
		setParts = append(setParts, key+" = $"+fmt.Sprint(idx))
		setArgs = append(setArgs, value)
		idx++
	}
	whereClause := query
	for i := range args {
		whereClause = strings.Replace(whereClause, "$"+fmt.Sprint(i+1), "$"+fmt.Sprint(i+1+len(setArgs)), 1)
	}
	result, err := app.SQLDB.Exec(ctx, `UPDATE deliveries SET `+strings.Join(setParts, ", ")+` WHERE `+whereClause, append(setArgs, args...)...)
	if err != nil {
		return Delivery{}, err
	}
	if result.RowsAffected() == 0 {
		return Delivery{}, fmt.Errorf("delivery not found")
	}
	var d Delivery
	if len(args) == 0 {
		return d, nil
	}
	idValue := args[0]
	if err := findDeliveryByFilter(ctx, app, "deliveryid = $1", []any{idValue}, &d); err != nil {
		return Delivery{}, err
	}
	return d, nil
}

func fetchDeliveryForRead(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) (Delivery, error) {
	if app == nil || app.SQLDB == nil {
		return Delivery{}, nil
	}
	query := `SELECT deliveryid, orderid, userid, status, driverid, created_at, updated_at, metadata FROM deliveries WHERE deliveryid = $1`
	args := []any{deliveryID}
	if strings.TrimSpace(tenantID) != "" {
		query += ` AND COALESCE(metadata->>'tenantid', '') = $2`
		args = append(args, tenantID)
	}
	var d Delivery
	var meta []byte
	var createdAt, updatedAt time.Time
	var orderID, driverID *string
	row := app.SQLDB.QueryRow(ctx, query, args...)
	if err := row.Scan(&d.DeliveryID, &orderID, &d.UserID, &d.Status, &driverID, &createdAt, &updatedAt, &meta); err != nil {
		return Delivery{}, err
	}
	if driverID != nil {
		d.DriverID = driverID
	}
	d.CreatedAt = createdAt
	d.UpdatedAt = updatedAt
	if err := hydrateDeliveryFromMetadata(meta, &d); err != nil {
		return Delivery{}, err
	}
	return d, nil
}

func setDeliveryStatusWithHistory(ctx context.Context, app *infra.Deps, deliveryID, tenantID, userID, newStatus string) (*Delivery, error) {
	current, err := fetchDeliveryForRead(ctx, app, deliveryID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("delivery not found")
	}
	if err := ValidateTransition(current.Status, newStatus); err != nil {
		return nil, err
	}
	newHistoryItem := StatusHistoryItem{Status: newStatus, Timestamp: time.Now(), UpdatedBy: userID}
	current.StatusHistory = append(current.StatusHistory, newHistoryItem)
	current.Status = newStatus
	current.UpdatedAt = time.Now()
	if err := saveDelivery(ctx, app, current); err != nil {
		return nil, err
	}
	return &current, nil
}

func findDeliveryByIDTenant(ctx context.Context, app *infra.Deps, deliveryID, tenantID string) (Delivery, error) {
	return fetchDeliveryForRead(ctx, app, deliveryID, tenantID)
}

func listDeliveryEvents(ctx context.Context, app *infra.Deps, deliveryID, _ string) ([]map[string]any, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT deliveryeventid, deliveryid, event_type, message, created_at, metadata FROM deliveryEvents WHERE deliveryid = $1 ORDER BY created_at DESC`, deliveryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, eventType, message string
		var createdAt time.Time
		var metadata []byte
		if err := rows.Scan(&id, &deliveryID, &eventType, &message, &createdAt, &metadata); err != nil {
			return nil, err
		}
		items = append(items, map[string]any{"id": id, "deliveryid": deliveryID, "event_type": eventType, "message": message, "created_at": createdAt, "metadata": string(metadata)})
	}
	return items, rows.Err()
}

func upsertDeliveryAssignment(ctx context.Context, app *infra.Deps, deliveryID, tenantID, driverID, userID string) (Delivery, error) {
	current, err := fetchDeliveryForRead(ctx, app, deliveryID, tenantID)
	if err != nil {
		return Delivery{}, fmt.Errorf("delivery not found")
	}
	if err := ValidateTransition(current.Status, StatusAssigned); err != nil {
		return Delivery{}, err
	}
	current.DriverID = &driverID
	current.Status = StatusAssigned
	current.UpdatedAt = time.Now()
	current.StatusHistory = append(current.StatusHistory, StatusHistoryItem{Status: StatusAssigned, Timestamp: time.Now(), UpdatedBy: userID})
	if err := saveDelivery(ctx, app, current); err != nil {
		return Delivery{}, err
	}
	return current, nil
}

func updateDeliveryStatusRecord(ctx context.Context, app *infra.Deps, deliveryID, tenantID, userID, newStatus string) (*Delivery, error) {
	currentDelivery, err := fetchDeliveryForRead(ctx, app, deliveryID, tenantID)
	if err != nil {
		return nil, err
	}
	if err := ValidateTransition(currentDelivery.Status, newStatus); err != nil {
		return nil, err
	}
	currentDelivery.StatusHistory = append(currentDelivery.StatusHistory, StatusHistoryItem{Status: newStatus, Timestamp: time.Now(), UpdatedBy: userID})
	currentDelivery.Status = newStatus
	currentDelivery.UpdatedAt = time.Now()
	if err := saveDelivery(ctx, app, currentDelivery); err != nil {
		return nil, err
	}
	return &currentDelivery, nil
}
