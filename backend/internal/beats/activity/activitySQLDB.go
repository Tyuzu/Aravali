// File: internal/beats/activity/activitySQLDB.go

package activity

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
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
	docs := make([]any, len(activities))
	for i := range activities {
		docs[i] = activities[i]
	}

	return app.SQLDB.InsertMany(ctx, ActivitiesTable, docs)
}

func getActivities(ctx context.Context, app *infra.Deps, userID string, cursor time.Time, limit int) ([]Activity, error) {
	where := "userid = $1"
	args := []any{userID}

	if !cursor.IsZero() {
		where += " AND timestamp < $2"
		args = append(args, cursor)
	}

	opts := sqldb.FindManyOptions{
		Limit:   int64(limit),
		OrderBy: "timestamp DESC",
	}

	var activities []Activity
	err := app.SQLDB.FindManyWithOptions(ctx, ActivitiesTable, where, args, opts, &activities)
	return activities, err
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

		if app.Cache != nil {
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

	err := app.SQLDB.InsertMany(ctx, AnalyticsTable, docsToInsert)
	return len(docsToInsert), err
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
