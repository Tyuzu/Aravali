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
)

var (
	ActivitiesTable = config.Tables.ActivitiesTable
	AnalyticsTable  = config.Tables.AnalyticsTable
)

func SQLinsertActivities(ctx context.Context, app *infra.Deps, activities []Activity) error {
	docs := make([]any, len(activities))
	for i := range activities {
		docs[i] = activities[i]
	}

	return app.SQLDB.InsertMany(ctx, ActivitiesTable, docs)
}

func SQLgetActivities(ctx context.Context, app *infra.Deps, userID string, cursor time.Time, limit int) ([]Activity, error) {
	where := "userid = $1"
	args := []any{userID}

	if !cursor.IsZero() {
		where += " AND timestamp < $2"
		args = append(args, cursor)
	}

	opts := sqldb.FindManyOptions{
		Limit:   limit,
		OrderBy: "timestamp DESC",
	}

	var activities []Activity
	err := app.SQLDB.FindManyWithOptions(ctx, ActivitiesTable, where, args, opts, &activities)
	return activities, err
}

func SQLinsertAnalyticsEvents(ctx context.Context, app *infra.Deps, payload AnalyticsPayload, remoteAddr string) (int, error) {
	var docsToInsert []any
	meta := payload.Meta
	user, _ := meta["user"].(string)
	session, _ := meta["session"].(string)
	url, _ := meta["url"].(string)

	for _, ev := range payload.Events {
		key := analyticsIdempotencyKey(ev)

		ok, err := app.Cache.SetNX(ctx, key, []byte("1"), analyticsIdemTTL)
		if err != nil || !ok {
			continue
		}

		// Marshal nested map/interface data into JSON for SQL JSONB/TEXT columns if required
		eventDataRaw, _ := json.Marshal(ev["data"])

		doc := map[string]any{
			"type":      ev["type"],
			"data":      string(eventDataRaw),
			"url":       url,
			"user":      user,
			"session":   session,
			"timestamp": time.Now(),
			"ip":        remoteAddr,
		}

		docsToInsert = append(docsToInsert, doc)
	}

	if len(docsToInsert) == 0 {
		return 0, nil
	}

	err := app.SQLDB.InsertMany(ctx, AnalyticsTable, docsToInsert)
	return len(docsToInsert), err
}

func SQLparseCursor(r *http.Request) (time.Time, int) {
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

func SQLanalyticsIdempotencyKey(ev map[string]any) string {
	raw, _ := json.Marshal(ev)
	sum := sha256.Sum256(raw)
	return "analytics:idemp:" + hex.EncodeToString(sum[:])
}
