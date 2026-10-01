// File: internal/events/eventsSQLDB.go

package events

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
	"scav/utils"
)

var (
	eventsTable = config.Tables.EventsTable
	ticksTable  = config.Tables.TicketsTable
	mediaTable  = config.Tables.MediaTable
	merchTable  = config.Tables.MerchTable
)

func SQLinsertEvent(ctx context.Context, app *infra.Deps, event Event) error {
	return app.SQLDB.InsertOne(ctx, eventsTable, event)
}

func SQLensureUniqueEventID(ctx context.Context, app *infra.Deps, event *Event) {
	if event == nil {
		return
	}

	event.EventID = utils.GenerateRandomString(14)
	var existingEvent Event
	query := "eventid = $1"
	args := []any{event.EventID}

	if err := app.SQLDB.FindOne(ctx, eventsTable, query, args, &existingEvent); err == nil {
		event.EventID = utils.GenerateRandomString(14)
	}
}

func SQLfindEventByID(ctx context.Context, app *infra.Deps, eventID string, event *Event) error {
	query := "eventid = $1"
	args := []any{eventID}

	return app.SQLDB.FindOne(ctx, eventsTable, query, args, event)
}

func SQLupdateEvent(ctx context.Context, app *infra.Deps, eventID string, updates map[string]any) (int64, error) {
	query := "eventid = $1"
	args := []any{eventID}

	return app.SQLDB.UpdateOne(ctx, eventsTable, query, args, updates)
}

func SQLaggregateEvent(ctx context.Context, app *infra.Deps, eventID string, result *[]Event) error {
	rawQuery := `
		SELECT 
			e.*,
			COALESCE(
				(SELECT json_agg(t.*) FROM ` + ticksTable + ` t WHERE t.eventid = e.eventid),
				'[]'
			) AS tickets,
			COALESCE(
				(SELECT json_agg(m.*) FROM ` + mediaTable + ` m WHERE m.entityid = e.eventid AND m.entitytype = 'event'),
				'[]'
			) AS media,
			COALESCE(
				(SELECT json_agg(mc.*) FROM ` + merchTable + ` mc WHERE mc.entity_id = e.eventid AND mc.entity_type = 'event'),
				'[]'
			) AS merch
		FROM ` + eventsTable + ` e
		WHERE e.eventid = $1
	`

	return app.SQLDB.QueryRaw(ctx, rawQuery, []any{eventID}, result)
}

func SQLlistEvents(ctx context.Context, app *infra.Deps, query string, args []any, opts sqldb.FindManyOptions, result *[]Event) error {
	return app.SQLDB.FindManyWithOptions(ctx, eventsTable, query, args, opts, result)
}

func SQLcountEvents(ctx context.Context, app *infra.Deps, whereClause string, args []any) (int64, error) {
	return app.SQLDB.Count(ctx, eventsTable, whereClause, args)
}

func insertEvent(ctx context.Context, app *infra.Deps, event Event) error {
	return SQLinsertEvent(ctx, app, event)
}

func ensureUniqueEventID(ctx context.Context, app *infra.Deps, event *Event) {
	SQLensureUniqueEventID(ctx, app, event)
}

func findEventByID(ctx context.Context, app *infra.Deps, eventID string, event *Event) error {
	return SQLfindEventByID(ctx, app, eventID, event)
}

func updateEvent(ctx context.Context, app *infra.Deps, eventID string, updates map[string]any) (int64, error) {
	return SQLupdateEvent(ctx, app, eventID, updates)
}

func aggregateEvent(ctx context.Context, app *infra.Deps, eventID string, result *[]Event) error {
	return SQLaggregateEvent(ctx, app, eventID, result)
}

func countEvents(ctx context.Context, app *infra.Deps, filter map[string]any) (int64, error) {
	where, args := buildEventQuery(filter)
	return SQLcountEvents(ctx, app, where, args)
}

func listEvents(ctx context.Context, app *infra.Deps, filter map[string]any, opts sqldb.FindManyOptions, result *[]Event) error {
	where, args := buildEventQuery(filter)
	return SQLlistEvents(ctx, app, where, args, opts, result)
}

func buildEventQuery(filter map[string]any) (string, []any) {
	if len(filter) == 0 {
		return "1 = 1", nil
	}

	clauses := make([]string, 0, len(filter))
	args := make([]any, 0, len(filter))
	keys := make([]string, 0, len(filter))
	for k := range filter {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := filter[key]
		if value == nil {
			continue
		}
		clauses = append(clauses, fmt.Sprintf("%s = $%d", normalizeEventColumn(key), len(args)+1))
		args = append(args, value)
	}

	if len(clauses) == 0 {
		return "1 = 1", nil
	}
	return strings.Join(clauses, " AND "), args
}

func normalizeEventColumn(key string) string {
	return strings.ToLower(strings.ReplaceAll(key, "-", "_"))
}
