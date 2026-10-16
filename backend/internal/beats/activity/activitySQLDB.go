// File: internal/beats/activity/activitySQLDB.go

package activity

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/utils"
)

var (
	ActivitiesTable = config.Tables.ActivitiesTable
	AnalyticsTable  = config.Tables.AnalyticsTable
)

const (
	analyticsIdemTTL = time.Minute
	defaultPageSize  = 20
	maxPageSize      = 100
)

func insertActivities(ctx context.Context, app *infra.Deps, activities []Activity) error {
	if app == nil || app.SQLDB == nil || len(activities) == 0 {
		return nil
	}

	values := make([]any, 0, len(activities)*14)
	clauses := make([]string, 0, len(activities))
	for _, activity := range activities {
		if strings.TrimSpace(activity.ActivityID) == "" {
			activity.ActivityID = "act_" + utils.GenerateRandomString(18)
		}
		if activity.Timestamp.IsZero() {
			activity.Timestamp = time.Now().UTC()
		}
		values = append(values,
			activity.ActivityID,
			activity.PlaceID,
			activity.Action,
			activity.PerformedBy,
			activity.Timestamp,
			activity.Details,
			activity.IPAddress,
			activity.DeviceInfo,
			activity.ActivityType,
			activity.EntityID,
			activity.EntityType,
			activity.UserID,
			time.Now().UTC(),
			time.Now().UTC(),
		)
		clauses = append(clauses, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", len(values)-13, len(values)-12, len(values)-11, len(values)-10, len(values)-9, len(values)-8, len(values)-7, len(values)-6, len(values)-5, len(values)-4, len(values)-3, len(values)-2, len(values)-1, len(values)))
	}

	query := `INSERT INTO activities (activityid, placeid, action, performedby, timestamp, details, ipaddress, deviceinfo, activity_type, entity_id, entity_type, userid, created_at, updated_at) VALUES ` + strings.Join(clauses, ", ") + ` ON CONFLICT (activityid) DO NOTHING`
	_, err := app.SQLDB.Exec(ctx, query, values...)
	return err
}

func getActivities(ctx context.Context, app *infra.Deps, userID string, cursor time.Time, limit int) ([]Activity, error) {
	if app == nil || app.SQLDB == nil || strings.TrimSpace(userID) == "" {
		return []Activity{}, nil
	}
	if limit <= 0 {
		limit = defaultPageSize
	}
	if limit > maxPageSize {
		limit = maxPageSize
	}
	if cursor.IsZero() {
		cursor = time.Now().UTC()
	}

	rows, err := app.SQLDB.Query(
		ctx,
		`SELECT activityid, placeid, action, performedby, timestamp, details, ipaddress, deviceinfo, activity_type, entity_id, entity_type, userid FROM activities WHERE userid = $1 AND timestamp < $2 ORDER BY timestamp DESC LIMIT $3`,
		userID,
		cursor,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Activity, 0, limit)
	for rows.Next() {
		var item Activity
		if err := rows.Scan(&item.ActivityID, &item.PlaceID, &item.Action, &item.PerformedBy, &item.Timestamp, &item.Details, &item.IPAddress, &item.DeviceInfo, &item.ActivityType, &item.EntityID, &item.EntityType, &item.UserID); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return []Activity{}, nil
	}
	return items, nil
}

type analyticsEventRow struct {
	ID         string          `db:"id"`
	EntityType string          `db:"entity_type"`
	EntityID   string          `db:"entity_id"`
	EventName  string          `db:"event_name"`
	Payload    json.RawMessage `db:"payload"`
}

func analyticsEventRowFromPayload(ev map[string]any, user, session, url, remoteAddr string) (analyticsEventRow, bool) {
	if ev == nil {
		return analyticsEventRow{}, false
	}

	eventType, _ := ev["type"].(string)
	if eventType == "" {
		return analyticsEventRow{}, false
	}

	payload := map[string]any{
		"data":    ev["data"],
		"url":     url,
		"user":    user,
		"session": session,
		"ip":      remoteAddr,
	}
	if ts, exists := ev["ts"]; exists {
		payload["ts"] = ts
	}

	payloadRaw, err := json.Marshal(payload)
	if err != nil {
		return analyticsEventRow{}, false
	}

	return analyticsEventRow{
		ID:         "ana_" + utils.GenerateRandomString(16),
		EntityType: "event",
		EntityID:   session,
		EventName:  eventType,
		Payload:    payloadRaw,
	}, true
}

func insertAnalyticsEvents(ctx context.Context, app *infra.Deps, payload AnalyticsPayload, remoteAddr string) (int, error) {
	var docsToInsert []any
	meta := payload.Meta
	user, _ := meta["user"].(string)
	session, _ := meta["session"].(string)
	url, _ := meta["url"].(string)

	for _, ev := range payload.Events {
		key := analyticsIdempotencyKey(ev)

		if app != nil && app.Cache != nil {
			ok, err := app.Cache.SetNX(ctx, key, []byte("1"), analyticsIdemTTL)
			if err != nil || !ok {
				continue
			}
		}

		row, ok := analyticsEventRowFromPayload(ev, user, session, url, remoteAddr)
		if !ok {
			continue
		}

		docsToInsert = append(docsToInsert, row)
	}

	if len(docsToInsert) == 0 {
		return 0, nil
	}
	if app == nil || app.SQLDB == nil {
		return len(docsToInsert), nil
	}

	inserted := 0
	for _, doc := range docsToInsert {
		row, ok := doc.(analyticsEventRow)
		if !ok {
			continue
		}
		_, err := app.SQLDB.Exec(ctx, `INSERT INTO analytics (id, entity_type, entity_id, event_name, payload, created_at) VALUES ($1,$2,$3,$4,$5,NOW()) ON CONFLICT (id) DO NOTHING`, row.ID, row.EntityType, row.EntityID, row.EventName, row.Payload)
		if err == nil {
			inserted++
		}
	}
	return inserted, nil
}

func parseCursor(r *http.Request) (time.Time, int) {
	q := r.URL.Query()

	limit := defaultPageSize
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			if n > maxPageSize {
				n = maxPageSize
			}
			limit = n
		}
	}

	var cursor time.Time
	if v := q.Get("cursor"); v != "" {
		if ts, err := strconv.ParseInt(v, 10, 64); err == nil {
			cursor = time.UnixMilli(ts)
		}
	}

	return cursor, limit
}

func analyticsIdempotencyKey(ev map[string]any) string {
	raw, _ := json.Marshal(ev)
	sum := sha256.Sum256(raw)
	return "analytics:idemp:" + hex.EncodeToString(sum[:])
}
