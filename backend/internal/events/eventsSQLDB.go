// File: internal/events/eventsSQLDB.go

package events

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"scav/config"
	"scav/infra"
)

var (
	eventsTable = config.Tables.EventsTable
	ticksTable  = config.Tables.TicketsTable
	mediaTable  = config.Tables.MediaTable
	merchTable  = config.Tables.MerchTable
)

func SQLinsertEvent(ctx context.Context, app *infra.Deps, event Event) error {

}

func SQLensureUniqueEventID(ctx context.Context, app *infra.Deps, event *Event) {

}

func SQLfindEventByID(ctx context.Context, app *infra.Deps, eventID string, event *Event) error {

}

func SQLupdateEvent(ctx context.Context, app *infra.Deps, eventID string, updates map[string]any) (int64, error) {

}

func SQLaggregateEvent(ctx context.Context, app *infra.Deps, eventID string, result *[]Event) error {

}

func SQLlistEvents(ctx context.Context, app *infra.Deps, query string, args []any, opts map[string]any, result *[]Event) error {

}

func SQLcountEvents(ctx context.Context, app *infra.Deps, whereClause string, args []any) (int64, error) {

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

func listEvents(ctx context.Context, app *infra.Deps, filter map[string]any, opts map[string]any, result *[]Event) error {
	where, args := buildEventQuery(filter)
	return SQLlistEvents(ctx, app, where, args, opts, result)
}

func buildEventListOptions(skip, limit int) map[string]any {

}

func getPaginatedEvents(ctx context.Context, app *infra.Deps, filter map[string]any, skip, limit int, result *[]Event) error {
	opts := buildEventListOptions(skip, limit)
	return listEvents(ctx, app, filter, opts, result)
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
