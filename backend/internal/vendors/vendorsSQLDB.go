// File: internal/vendors/vendorsSQLDB.go

package vendors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"

	"github.com/jackc/pgx/v5"
)

var (
	vendorTable = config.Tables.VendorTable
	hiringTable = config.Tables.HiringTable
)

func ensureVendorDB(app *infra.Deps) error {
	if app == nil || app.SQLDB == nil {
		return errors.New("database not initialized")
	}
	return nil
}

func encodeMetadata(value any) ([]byte, error) {
	if value == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(value)
}

func decodeVendorMetadata(metadata []byte) (map[string]any, error) {
	if len(metadata) == 0 {
		return map[string]any{}, nil
	}
	var payload map[string]any
	if err := json.Unmarshal(metadata, &payload); err != nil {
		return nil, err
	}
	if payload == nil {
		payload = map[string]any{}
	}
	return payload, nil
}

func vendorFromRow(vendorID, userID, name, description, status string, createdAt, updatedAt time.Time, metadata []byte) (*Vendor, error) {
	vendor := &Vendor{
		VendorID:    vendorID,
		UserID:      userID,
		Name:        name,
		Description: description,
		Available:   strings.EqualFold(status, "available") || status == "active",
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	if len(metadata) > 0 {
		payload, err := decodeVendorMetadata(metadata)
		if err != nil {
			return nil, err
		}
		if v, ok := payload["description"]; ok && vendor.Description == "" {
			if s, ok := v.(string); ok {
				vendor.Description = s
			}
		}
		if v, ok := payload["category"]; ok {
			if s, ok := v.(string); ok {
				vendor.Category = s
			}
		}
		if v, ok := payload["email"]; ok {
			if s, ok := v.(string); ok {
				vendor.Email = s
			}
		}
		if v, ok := payload["phone"]; ok {
			if s, ok := v.(string); ok {
				vendor.Phone = s
			}
		}
		if v, ok := payload["location"]; ok {
			if s, ok := v.(string); ok {
				vendor.Location = s
			}
		}
		if v, ok := payload["rating"]; ok {
			switch val := v.(type) {
			case float64:
				vendor.Rating = val
			case int:
				vendor.Rating = float64(val)
			}
		}
		if v, ok := payload["rating_count"]; ok {
			switch val := v.(type) {
			case float64:
				vendor.RatingCount = int(val)
			case int:
				vendor.RatingCount = val
			}
		}
		if v, ok := payload["profile_image"]; ok {
			if s, ok := v.(string); ok {
				vendor.ProfileImage = s
			}
		}
		if v, ok := payload["portfolio"]; ok {
			if arr, ok := v.([]any); ok {
				vendor.Portfolio = make([]string, 0, len(arr))
				for _, item := range arr {
					if s, ok := item.(string); ok {
						vendor.Portfolio = append(vendor.Portfolio, s)
					}
				}
			}
		}
		if v, ok := payload["verified"]; ok {
			if b, ok := v.(bool); ok {
				vendor.Verified = b
			}
		}
		if v, ok := payload["available"]; ok {
			if b, ok := v.(bool); ok {
				vendor.Available = b
			}
		}
		if v, ok := payload["name"]; ok {
			if s, ok := v.(string); ok && vendor.Name == "" {
				vendor.Name = s
			}
		}
		if v, ok := payload["status"]; ok {
			if s, ok := v.(string); ok {
				vendor.Available = strings.EqualFold(s, "available") || s == "active"
			}
		}
	}
	if vendor.Name == "" {
		vendor.Name = name
	}
	if vendor.Description == "" {
		vendor.Description = description
	}
	if vendor.UserID == "" {
		vendor.UserID = userID
	}
	if vendor.VendorID == "" {
		vendor.VendorID = vendorID
	}
	return vendor, nil
}

func hiringFromRow(hiringID, userID, vendorID, status string, createdAt, updatedAt time.Time, metadata []byte) (*VendorHiring, error) {
	result := &VendorHiring{
		HiringID:  hiringID,
		VendorID:  vendorID,
		HiredBy:   userID,
		Status:    status,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	if len(metadata) == 0 {
		return result, nil
	}
	payload, err := decodeVendorMetadata(metadata)
	if err != nil {
		return nil, err
	}
	if v, ok := payload["eventid"]; ok {
		if s, ok := v.(string); ok {
			result.EventID = s
		}
	}
	if v, ok := payload["vendor_name"]; ok {
		if s, ok := v.(string); ok {
			result.VendorName = s
		}
	}
	if v, ok := payload["vendor_category"]; ok {
		if s, ok := v.(string); ok {
			result.VendorCategory = s
		}
	}
	if v, ok := payload["hired_at"]; ok {
		if s, ok := v.(string); ok {
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				result.HiredAt = t
			}
		}
	}
	if v, ok := payload["hired_by"]; ok {
		if s, ok := v.(string); ok {
			result.HiredBy = s
		}
	}
	if v, ok := payload["notes"]; ok {
		if s, ok := v.(string); ok {
			result.Notes = s
		}
	}
	if v, ok := payload["status"]; ok {
		if s, ok := v.(string); ok {
			result.Status = s
		}
	}
	return result, nil
}

func rewriteVendorQuery(query string, args []any) string {
	if query == "" {
		return "1 = 1"
	}
	q := query
	for _, token := range []string{"available", "category", "name", "description", "location", "vendorid", "userid", "status", "email", "phone", "verified", "rating"} {
		q = strings.ReplaceAll(q, token, "metadata->>'"+token+"'")
	}
	if strings.Contains(q, "LOWER(") {
		q = strings.ReplaceAll(q, "LOWER(metadata->>'name')", "LOWER(COALESCE(metadata->>'name',''))")
		q = strings.ReplaceAll(q, "LOWER(metadata->>'category')", "LOWER(COALESCE(metadata->>'category',''))")
		q = strings.ReplaceAll(q, "LOWER(metadata->>'description')", "LOWER(COALESCE(metadata->>'description',''))")
		q = strings.ReplaceAll(q, "LOWER(metadata->>'location')", "LOWER(COALESCE(metadata->>'location',''))")
	}
	q = strings.ReplaceAll(q, "metadata->>'available'", "COALESCE((metadata->>'available')::boolean, false)")
	q = strings.ReplaceAll(q, "metadata->>'category'", "LOWER(COALESCE(metadata->>'category',''))")
	q = strings.ReplaceAll(q, "metadata->>'status'", "COALESCE(metadata->>'status','active')")
	q = strings.ReplaceAll(q, "metadata->>'vendorid'", "COALESCE(metadata->>'vendorid','')")
	q = strings.ReplaceAll(q, "metadata->>'userid'", "COALESCE(metadata->>'userid','')")
	return q
}

func rewriteHiringQuery(query string) string {
	if query == "" {
		return "1 = 1"
	}
	q := query
	q = strings.ReplaceAll(q, "eventid", "metadata->>'eventid'")
	q = strings.ReplaceAll(q, "vendorid", "vendorid")
	q = strings.ReplaceAll(q, "hiringid", "hiringid")
	q = strings.ReplaceAll(q, "status", "status")
	return q
}

// InsertVendor inserts a vendor document into the vendor table.
func InsertVendor(ctx context.Context, app *infra.Deps, vendor *Vendor) error {
	if err := ensureVendorDB(app); err != nil {
		return err
	}
	if vendor == nil {
		return errors.New("nil vendor")
	}
	vendor.CreatedAt = time.Now().UTC()
	vendor.UpdatedAt = vendor.CreatedAt
	payload, err := encodeMetadata(map[string]any{
		"vendorid":      vendor.VendorID,
		"userid":        vendor.UserID,
		"name":          vendor.Name,
		"category":      vendor.Category,
		"description":   vendor.Description,
		"email":         vendor.Email,
		"phone":         vendor.Phone,
		"location":      vendor.Location,
		"rating":        vendor.Rating,
		"rating_count":  vendor.RatingCount,
		"profile_image": vendor.ProfileImage,
		"portfolio":     vendor.Portfolio,
		"verified":      vendor.Verified,
		"available":     vendor.Available,
		"status":        mapStatus(vendor.Available),
	})
	if err != nil {
		return err
	}
	_, err = app.SQLDB.Exec(ctx,
		`INSERT INTO vendors (vendorid, name, description, userid, status, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (vendorid) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			userid = EXCLUDED.userid,
			status = EXCLUDED.status,
			updated_at = NOW(),
			metadata = EXCLUDED.metadata`,
		vendor.VendorID,
		vendor.Name,
		vendor.Description,
		vendor.UserID,
		mapStatus(vendor.Available),
		vendor.CreatedAt,
		vendor.UpdatedAt,
		payload,
	)
	return err
}

// FindVendorByID returns a vendor by vendorID. Returns ErrVendorNotFound if not found.
func FindVendorByID(ctx context.Context, app *infra.Deps, vendorID string) (*Vendor, error) {
	if err := ensureVendorDB(app); err != nil {
		return nil, err
	}
	var vendorIDRow, userID, name, description, status string
	var createdAt, updatedAt time.Time
	var metadata []byte
	if err := app.SQLDB.QueryRow(ctx, `SELECT vendorid, name, description, userid, status, created_at, updated_at, metadata FROM vendors WHERE vendorid = $1 LIMIT 1`, vendorID).Scan(&vendorIDRow, &name, &description, &userID, &status, &createdAt, &updatedAt, &metadata); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVendorNotFound
		}
		return nil, err
	}
	return vendorFromRow(vendorIDRow, userID, name, description, status, createdAt, updatedAt, metadata)
}

// FindVendorByUserID returns a vendor by userID. Returns nil,err when FindOne fails.
func FindVendorByUserID(ctx context.Context, app *infra.Deps, userID string) (*Vendor, error) {
	if err := ensureVendorDB(app); err != nil {
		return nil, err
	}
	var vendorIDRow, userIDRow, name, description, status string
	var createdAt, updatedAt time.Time
	var metadata []byte
	if err := app.SQLDB.QueryRow(ctx, `SELECT vendorid, name, description, userid, status, created_at, updated_at, metadata FROM vendors WHERE userid = $1 ORDER BY created_at DESC LIMIT 1`, userID).Scan(&vendorIDRow, &name, &description, &userIDRow, &status, &createdAt, &updatedAt, &metadata); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVendorNotFound
		}
		return nil, err
	}
	return vendorFromRow(vendorIDRow, userIDRow, name, description, status, createdAt, updatedAt, metadata)
}

// FindVendors finds many vendors using provided query and args.
func FindVendors(ctx context.Context, app *infra.Deps, query string, args []any, out *[]Vendor) error {
	if err := ensureVendorDB(app); err != nil {
		return err
	}
	if out == nil {
		return errors.New("nil result")
	}
	if query == "" {
		query = "1 = 1"
	}
	converted := rewriteVendorQuery(query, args)
	statement := fmt.Sprintf(`SELECT vendorid, name, description, userid, status, created_at, updated_at, metadata FROM vendors WHERE %s ORDER BY created_at DESC`, converted)
	rows, err := app.SQLDB.Query(ctx, statement, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	items := make([]Vendor, 0)
	for rows.Next() {
		var vendorIDRow, userID, name, description, status string
		var createdAt, updatedAt time.Time
		var metadata []byte
		if err := rows.Scan(&vendorIDRow, &name, &description, &userID, &status, &createdAt, &updatedAt, &metadata); err != nil {
			return err
		}
		v, err := vendorFromRow(vendorIDRow, userID, name, description, status, createdAt, updatedAt, metadata)
		if err != nil {
			return err
		}
		items = append(items, *v)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*out = items
	return nil
}

// UpdateVendorDB updates vendor documents matching query with update map.
func UpdateVendorDB(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	if err := ensureVendorDB(app); err != nil {
		return 0, err
	}
	if len(update) == 0 {
		return 0, nil
	}
	set := map[string]any{}
	if v, ok := update["$set"]; ok {
		if nested, ok := v.(map[string]any); ok {
			set = nested
		} else {
			set = map[string]any{"value": v}
		}
	} else {
		set = update
	}
	if len(set) == 0 {
		return 0, nil
	}
	var metadata []byte
	if whereQuery := query; whereQuery != "" {
		// fetch current metadata for the matching vendor row
		baseQuery := fmt.Sprintf(`SELECT metadata FROM vendors WHERE %s LIMIT 1`, rewriteVendorQuery(whereQuery, args))
		if err := app.SQLDB.QueryRow(ctx, baseQuery, args...).Scan(&metadata); err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return 0, err
		}
	}
	base := map[string]any{}
	if len(metadata) > 0 {
		if err := json.Unmarshal(metadata, &base); err != nil {
			return 0, err
		}
	}
	for k, v := range set {
		base[k] = v
	}
	if value, ok := base["available"]; ok {
		if b, ok2 := value.(bool); ok2 {
			base["status"] = mapStatus(b)
		}
	}
	payload, err := encodeMetadata(base)
	if err != nil {
		return 0, err
	}
	updated, err := app.SQLDB.Exec(ctx,
		`UPDATE vendors SET metadata = $1::jsonb, status = COALESCE((metadata->>'status'), 'active'), updated_at = NOW() WHERE `+rewriteVendorQuery(query, args),
		append([]any{payload}, args...)...,
	)
	if err != nil {
		return 0, err
	}
	return updated.RowsAffected(), nil
}

// DeleteVendorDB marks a vendor as unavailable.
func DeleteVendorDB(ctx context.Context, app *infra.Deps, vendorID string) (int64, error) {
	if err := ensureVendorDB(app); err != nil {
		return 0, err
	}
	result, err := app.SQLDB.Exec(ctx,
		`UPDATE vendors SET status = 'inactive', metadata = COALESCE(metadata, '{}'::jsonb) || jsonb_build_object('available', false, 'status', 'inactive'), updated_at = NOW() WHERE vendorid = $1`,
		vendorID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// --- Hiring related DB helpers ---

func FindHiringByID(ctx context.Context, app *infra.Deps, hiringID string) (*VendorHiring, error) {
	if err := ensureVendorDB(app); err != nil {
		return nil, err
	}
	var hiringIDRow, vendorID, userID, status string
	var createdAt, updatedAt time.Time
	var metadata []byte
	if err := app.SQLDB.QueryRow(ctx, `SELECT hiringid, userid, vendorid, status, created_at, updated_at, metadata FROM hirings WHERE hiringid = $1 LIMIT 1`, hiringID).Scan(&hiringIDRow, &userID, &vendorID, &status, &createdAt, &updatedAt, &metadata); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVendorNotInEvent
		}
		return nil, err
	}
	return hiringFromRow(hiringIDRow, userID, vendorID, status, createdAt, updatedAt, metadata)
}

func FindHiringByEventAndVendor(ctx context.Context, app *infra.Deps, eventID, vendorID string) (*VendorHiring, error) {
	if err := ensureVendorDB(app); err != nil {
		return nil, err
	}
	var hiringIDRow, userID, rowVendorID, status string
	var createdAt, updatedAt time.Time
	var metadata []byte
	if err := app.SQLDB.QueryRow(ctx, `SELECT hiringid, userid, vendorid, status, created_at, updated_at, metadata FROM hirings WHERE metadata->>'eventid' = $1 AND vendorid = $2 ORDER BY created_at DESC LIMIT 1`, eventID, vendorID).Scan(&hiringIDRow, &userID, &rowVendorID, &status, &createdAt, &updatedAt, &metadata); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVendorNotInEvent
		}
		return nil, err
	}
	return hiringFromRow(hiringIDRow, userID, rowVendorID, status, createdAt, updatedAt, metadata)
}

func InsertHiring(ctx context.Context, app *infra.Deps, hiring *VendorHiring) error {
	if err := ensureVendorDB(app); err != nil {
		return err
	}
	if hiring == nil {
		return errors.New("nil hiring")
	}
	if hiring.HiredAt.IsZero() {
		hiring.HiredAt = time.Now().UTC()
	}
	if hiring.CreatedAt.IsZero() {
		hiring.CreatedAt = hiring.HiredAt
	}
	if hiring.UpdatedAt.IsZero() {
		hiring.UpdatedAt = hiring.CreatedAt
	}
	payload, err := encodeMetadata(map[string]any{
		"eventid":         hiring.EventID,
		"vendor_name":     hiring.VendorName,
		"vendor_category": hiring.VendorCategory,
		"vendorid":        hiring.VendorID,
		"hired_at":        hiring.HiredAt.Format(time.RFC3339),
		"hired_by":        hiring.HiredBy,
		"notes":           hiring.Notes,
		"status":          hiring.Status,
	})
	if err != nil {
		return err
	}
	_, err = app.SQLDB.Exec(ctx,
		`INSERT INTO hirings (hiringid, userid, vendorid, status, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (hiringid) DO UPDATE SET
			userid = EXCLUDED.userid,
			vendorid = EXCLUDED.vendorid,
			status = EXCLUDED.status,
			updated_at = NOW(),
			metadata = EXCLUDED.metadata`,
		hiring.HiringID,
		hiring.HiredBy,
		hiring.VendorID,
		hiring.Status,
		hiring.CreatedAt,
		hiring.UpdatedAt,
		payload,
	)
	return err
}

func FindHiringsByEvent(ctx context.Context, app *infra.Deps, eventID string, out *[]VendorHiring) error {
	if err := ensureVendorDB(app); err != nil {
		return err
	}
	if out == nil {
		return errors.New("nil result")
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT hiringid, userid, vendorid, status, created_at, updated_at, metadata FROM hirings WHERE metadata->>'eventid' = $1 ORDER BY created_at DESC`, eventID)
	if err != nil {
		return err
	}
	defer rows.Close()
	result := make([]VendorHiring, 0)
	for rows.Next() {
		var hiringIDRow, userID, vendorID, status string
		var createdAt, updatedAt time.Time
		var metadata []byte
		if err := rows.Scan(&hiringIDRow, &userID, &vendorID, &status, &createdAt, &updatedAt, &metadata); err != nil {
			return err
		}
		h, err := hiringFromRow(hiringIDRow, userID, vendorID, status, createdAt, updatedAt, metadata)
		if err != nil {
			return err
		}
		result = append(result, *h)
	}
	*out = result
	return rows.Err()
}

func FindHiringsByVendorID(ctx context.Context, app *infra.Deps, vendorID string, out *[]VendorHiring) error {
	if err := ensureVendorDB(app); err != nil {
		return err
	}
	if out == nil {
		return errors.New("nil result")
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT hiringid, userid, vendorid, status, created_at, updated_at, metadata FROM hirings WHERE vendorid = $1 ORDER BY created_at DESC`, vendorID)
	if err != nil {
		return err
	}
	defer rows.Close()
	result := make([]VendorHiring, 0)
	for rows.Next() {
		var hiringIDRow, userID, rowVendorID, status string
		var createdAt, updatedAt time.Time
		var metadata []byte
		if err := rows.Scan(&hiringIDRow, &userID, &rowVendorID, &status, &createdAt, &updatedAt, &metadata); err != nil {
			return err
		}
		h, err := hiringFromRow(hiringIDRow, userID, rowVendorID, status, createdAt, updatedAt, metadata)
		if err != nil {
			return err
		}
		result = append(result, *h)
	}
	*out = result
	return rows.Err()
}

func UpdateHiringDB(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	if err := ensureVendorDB(app); err != nil {
		return 0, err
	}
	if len(update) == 0 {
		return 0, nil
	}
	base := map[string]any{}
	if query != "" {
		statement := fmt.Sprintf(`SELECT metadata FROM hirings WHERE %s LIMIT 1`, rewriteHiringQuery(query))
		var metadata []byte
		if err := app.SQLDB.QueryRow(ctx, statement, args...).Scan(&metadata); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return 0, err
			}
		} else if len(metadata) > 0 {
			if err := json.Unmarshal(metadata, &base); err != nil {
				return 0, err
			}
		}
	}
	for k, v := range update {
		if k == "updated_at" {
			continue
		}
		base[k] = v
	}
	if v, ok := update["status"]; ok {
		if s, ok2 := v.(string); ok2 && s != "" {
			base["status"] = s
		}
	}
	base["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	payload, err := encodeMetadata(base)
	if err != nil {
		return 0, err
	}
	result, err := app.SQLDB.Exec(ctx,
		`UPDATE hirings SET metadata = $1::jsonb, status = COALESCE($2::text, status), updated_at = NOW() WHERE `+rewriteHiringQuery(query),
		append([]any{payload, update["status"]}, args...)...,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// --- Availability related DB helpers ---

func FindAvailabilitySlots(ctx context.Context, app *infra.Deps, vendorID string) ([]AvailabilitySlot, error) {
	if err := ensureVendorDB(app); err != nil {
		return nil, err
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT vendoravailabilityid, vendorid, day_of_week, starts_at, ends_at, created_at, updated_at, metadata FROM vendor_availability WHERE vendorid = $1 ORDER BY created_at DESC`, vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]AvailabilitySlot, 0)
	for rows.Next() {
		var slotID, rowVendorID, startDate, endDate string
		var dayOfWeek int
		var createdAt, updatedAt time.Time
		var metadata []byte
		if err := rows.Scan(&slotID, &rowVendorID, &dayOfWeek, &startDate, &endDate, &createdAt, &updatedAt, &metadata); err != nil {
			return nil, err
		}
		slot := AvailabilitySlot{
			SlotID:    slotID,
			VendorID:  rowVendorID,
			StartDate: startDate,
			EndDate:   endDate,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
		if len(metadata) > 0 {
			var payload map[string]any
			if err := json.Unmarshal(metadata, &payload); err == nil {
				if v, ok := payload["start_date"]; ok {
					if s, ok := v.(string); ok {
						slot.StartDate = s
					}
				}
				if v, ok := payload["end_date"]; ok {
					if s, ok := v.(string); ok {
						slot.EndDate = s
					}
				}
				if v, ok := payload["notes"]; ok {
					if s, ok := v.(string); ok {
						slot.Notes = s
					}
				}
				if v, ok := payload["recurring"]; ok {
					if b, ok := v.(bool); ok {
						slot.Recurring = b
					}
				}
				if v, ok := payload["recurrence_rule"]; ok {
					if s, ok := v.(string); ok {
						slot.RecurrenceRule = s
					}
				}
			}
		}
		result = append(result, slot)
	}
	return result, rows.Err()
}

func InsertAvailabilitySlotDB(ctx context.Context, app *infra.Deps, slot AvailabilitySlot) error {
	if err := ensureVendorDB(app); err != nil {
		return err
	}
	if slot.SlotID == "" {
		slot.SlotID = time.Now().UTC().Format("20060102T150405")
	}
	if slot.CreatedAt.IsZero() {
		slot.CreatedAt = time.Now().UTC()
	}
	if slot.UpdatedAt.IsZero() {
		slot.UpdatedAt = slot.CreatedAt
	}
	payload, err := encodeMetadata(map[string]any{
		"slotid":          slot.SlotID,
		"vendorid":        slot.VendorID,
		"start_date":      slot.StartDate,
		"end_date":        slot.EndDate,
		"recurring":       slot.Recurring,
		"recurrence_rule": slot.RecurrenceRule,
		"notes":           slot.Notes,
	})
	if err != nil {
		return err
	}
	_, err = app.SQLDB.Exec(ctx,
		`INSERT INTO vendor_availability (vendoravailabilityid, vendorid, day_of_week, starts_at, ends_at, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (vendoravailabilityid) DO UPDATE SET
			vendorid = EXCLUDED.vendorid,
			updated_at = NOW(),
			metadata = EXCLUDED.metadata`,
		slot.SlotID,
		slot.VendorID,
		0,
		slot.StartDate,
		slot.EndDate,
		slot.CreatedAt,
		slot.UpdatedAt,
		payload,
	)
	return err
}

func FindAvailabilitySlotByID(ctx context.Context, app *infra.Deps, slotID, vendorID string) (*AvailabilitySlot, error) {
	if err := ensureVendorDB(app); err != nil {
		return nil, err
	}
	var slotRowID, rowVendorID, startDate, endDate string
	var createdAt, updatedAt time.Time
	var metadata []byte
	if err := app.SQLDB.QueryRow(ctx, `SELECT vendoravailabilityid, vendorid, starts_at, ends_at, created_at, updated_at, metadata FROM vendor_availability WHERE vendoravailabilityid = $1 AND vendorid = $2 LIMIT 1`, slotID, vendorID).Scan(&slotRowID, &rowVendorID, &startDate, &endDate, &createdAt, &updatedAt, &metadata); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVendorNotFound
		}
		return nil, err
	}
	slot := &AvailabilitySlot{SlotID: slotRowID, VendorID: rowVendorID, StartDate: startDate, EndDate: endDate, CreatedAt: createdAt, UpdatedAt: updatedAt}
	if len(metadata) > 0 {
		var payload map[string]any
		if err := json.Unmarshal(metadata, &payload); err == nil {
			if v, ok := payload["start_date"]; ok {
				if s, ok := v.(string); ok {
					slot.StartDate = s
				}
			}
			if v, ok := payload["end_date"]; ok {
				if s, ok := v.(string); ok {
					slot.EndDate = s
				}
			}
			if v, ok := payload["notes"]; ok {
				if s, ok := v.(string); ok {
					slot.Notes = s
				}
			}
		}
	}
	return slot, nil
}

func DeleteAvailabilitySlotDB(ctx context.Context, app *infra.Deps, slotID string) (int64, error) {
	if err := ensureVendorDB(app); err != nil {
		return 0, err
	}
	result, err := app.SQLDB.Exec(ctx, `DELETE FROM vendor_availability WHERE vendoravailabilityid = $1`, slotID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func mapStatus(available bool) string {
	if available {
		return "available"
	}
	return "inactive"
}
