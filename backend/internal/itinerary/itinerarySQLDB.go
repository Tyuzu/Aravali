// File: internal/itinerary/itinerarySQLDB.go

package itinerary

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"scav/config"
	"scav/infra"
)

var ItineraryTable = config.Tables.ItineraryTable

func insertItinerary(ctx context.Context, app *infra.Deps, itinerary Itinerary) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	payload, err := json.Marshal(itinerary)
	if err != nil {
		return err
	}
	_, err = app.SQLDB.Exec(
		ctx,
		`INSERT INTO itinerary (itineraryid, userid, name, description, start_date, end_date, status, published, forked_from, deleted, days, created_at, updated_at, metadata) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NOW(),NOW(),$12) ON CONFLICT (itineraryid) DO UPDATE SET userid = EXCLUDED.userid, name = EXCLUDED.name, description = EXCLUDED.description, start_date = EXCLUDED.start_date, end_date = EXCLUDED.end_date, status = EXCLUDED.status, published = EXCLUDED.published, forked_from = EXCLUDED.forked_from, deleted = EXCLUDED.deleted, days = EXCLUDED.days, updated_at = NOW(), metadata = EXCLUDED.metadata`,
		itinerary.ItineraryID,
		itinerary.UserID,
		itinerary.Name,
		itinerary.Description,
		itinerary.StartDate,
		itinerary.EndDate,
		itinerary.Status,
		itinerary.Published,
		itinerary.ForkedFrom,
		itinerary.Deleted,
		itinerary.Days,
		payload,
	)
	return err
}

func findItineraryByID(ctx context.Context, app *infra.Deps, itineraryID string) (Itinerary, error) {
	if app == nil || app.SQLDB == nil {
		return Itinerary{}, nil
	}
	var itinerary Itinerary
	var metadata []byte
	var createdAt, updatedAt interface{}
	row := app.SQLDB.QueryRow(ctx, `SELECT itineraryid, userid, name, description, start_date, end_date, status, published, forked_from, deleted, days, created_at, updated_at, metadata FROM itinerary WHERE itineraryid = $1 LIMIT 1`, itineraryID)
	if err := row.Scan(&itinerary.ItineraryID, &itinerary.UserID, &itinerary.Name, &itinerary.Description, &itinerary.StartDate, &itinerary.EndDate, &itinerary.Status, &itinerary.Published, &itinerary.ForkedFrom, &itinerary.Deleted, &itinerary.Days, &createdAt, &updatedAt, &metadata); err != nil {
		return Itinerary{}, err
	}
	_ = metadata
	_ = createdAt
	_ = updatedAt
	return itinerary, nil
}

func findItineraries(ctx context.Context, app *infra.Deps, query string, args []any) ([]Itinerary, error) {
	if app == nil || app.SQLDB == nil {
		return []Itinerary{}, nil
	}
	if query == "" {
		query = "1 = 1"
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT itineraryid, userid, name, description, start_date, end_date, status, published, forked_from, deleted, days, created_at, updated_at, metadata FROM itinerary WHERE `+query+` ORDER BY created_at DESC`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Itinerary, 0)
	for rows.Next() {
		var itinerary Itinerary
		var metadata []byte
		var createdAt, updatedAt interface{}
		if err := rows.Scan(&itinerary.ItineraryID, &itinerary.UserID, &itinerary.Name, &itinerary.Description, &itinerary.StartDate, &itinerary.EndDate, &itinerary.Status, &itinerary.Published, &itinerary.ForkedFrom, &itinerary.Deleted, &itinerary.Days, &createdAt, &updatedAt, &metadata); err != nil {
			return nil, err
		}
		_ = metadata
		_ = createdAt
		_ = updatedAt
		items = append(items, itinerary)
	}
	return items, rows.Err()
}

func updateItineraryFields(ctx context.Context, app *infra.Deps, itineraryID string, update map[string]any) (int64, error) {
	if app == nil || app.SQLDB == nil || len(update) == 0 {
		return 0, nil
	}
	parts := make([]string, 0, len(update))
	args := make([]any, 0, len(update)+1)
	idx := 1
	for key, value := range update {
		parts = append(parts, key+" = $"+fmt.Sprint(idx))
		args = append(args, value)
		idx++
	}
	args = append(args, itineraryID)
	result, err := app.SQLDB.Exec(ctx, `UPDATE itinerary SET `+strings.Join(parts, ", ")+`, updated_at = NOW() WHERE itineraryid = $`+fmt.Sprint(len(args)), args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func softDeleteItinerary(ctx context.Context, app *infra.Deps, itineraryID, userID string) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	result, err := app.SQLDB.Exec(ctx, `UPDATE itinerary SET deleted = true, updated_at = NOW() WHERE itineraryid = $1 AND userid = $2`, itineraryID, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func publishItinerary(ctx context.Context, app *infra.Deps, itineraryID, userID string) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	result, err := app.SQLDB.Exec(ctx, `UPDATE itinerary SET published = true, updated_at = NOW() WHERE itineraryid = $1 AND userid = $2`, itineraryID, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
