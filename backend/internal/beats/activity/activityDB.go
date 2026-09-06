package activity

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"scav/config"
	"scav/infra"
	"scav/infra/db"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	ActivitiesCollection = config.Collections.ActivitiesCollection
	AnalyticsCollection  = config.Collections.AnalyticsCollection

	analyticsIdemTTL = 24 * time.Hour
	defaultPageSize  = 20
	maxPageSize      = 100
)

func insertActivities(ctx context.Context, app *infra.Deps, activities []Activity) error {
	docs := make([]any, len(activities))
	for i := range activities {
		docs[i] = activities[i]
	}

	return app.DB.WithDB(ctx, func(ctx context.Context) error {
		return app.DB.InsertMany(ctx, ActivitiesCollection, docs)
	})
}

func getActivities(ctx context.Context, app *infra.Deps, userID string, cursor time.Time, limit int) ([]Activity, error) {
	filter := map[string]any{
		"userid": userID,
	}

	if !cursor.IsZero() {
		filter["timestamp"] = map[string]any{"$lt": cursor}
	}

	opts := db.FindManyOptions{
		Limit: limit,
		Sort:  []bson.E{{Key: "timestamp", Value: -1}},
	}

	var activities []Activity
	err := app.DB.FindManyWithOptions(ctx, ActivitiesCollection, filter, opts, &activities)
	return activities, err
}

func insertAnalyticsEvents(ctx context.Context, app *infra.Deps, payload AnalyticsPayload, remoteAddr string) (int, error) {
	var docsToInsert []any
	meta := payload.Meta
	user, _ := meta["user"].(string)
	session, _ := meta["session"].(string)
	url, _ := meta["url"].(string)

	err := app.DB.WithDB(ctx, func(ctx context.Context) error {
		for _, ev := range payload.Events {
			key := analyticsIdempotencyKey(ev)

			ok, err := app.Cache.SetNX(ctx, key, []byte("1"), analyticsIdemTTL)
			if err != nil || !ok {
				continue
			}

			doc := map[string]any{
				"type":      ev["type"],
				"data":      ev["data"],
				"url":       url,
				"user":      user,
				"session":   session,
				"timestamp": time.Now(),
				"ip":        remoteAddr,
			}

			docsToInsert = append(docsToInsert, doc)
		}

		if len(docsToInsert) == 0 {
			return nil
		}

		return app.DB.InsertMany(ctx, AnalyticsCollection, docsToInsert)
	})

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
