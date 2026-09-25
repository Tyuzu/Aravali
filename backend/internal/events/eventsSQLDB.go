package events

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/db"
	"scav/utils"
)

var eventsTable = config.Tables.EventsTable

func sqlSQLinsertEvent(ctx context.Context, app *infra.Deps, event Event) error {
	return app.DB.Insert(ctx, eventsTable, event)
}

func sqlSQLensureUniqueEventID(ctx context.Context, app *infra.Deps, event *Event) {
	if event == nil {
		return
	}

	event.EventID = utils.GenerateRandomString(14)
	var existingEvent Event
	if err := app.DB.FindOne(ctx, eventsTable, map[string]string{"eventid": event.EventID}, &existingEvent); err == nil {
		event.EventID = utils.GenerateRandomString(14)
	}
}

func sqlSQLfindEventByID(ctx context.Context, app *infra.Deps, eventID string, event *Event) error {
	return app.DB.FindOne(ctx, eventsTable, map[string]string{"eventid": eventID}, event)
}

func sqlSQLupdateEvent(ctx context.Context, app *infra.Deps, eventID string, updates map[string]any) (any, error) {
	return app.DB.UpdateOne(ctx, eventsTable, map[string]string{"eventid": eventID}, map[string]any{"$set": updates})
}

func sqlSQLaggregateEvent(ctx context.Context, app *infra.Deps, eventID string, result *[]Event) error {
	pipeline := []any{
		map[string]any{"$match": map[string]any{"eventid": eventID}},
		map[string]any{"$lookup": map[string]any{
			"from":         "ticks",
			"localField":   "eventid",
			"foreignField": "eventid",
			"as":           "tickets",
		}},
		map[string]any{"$lookup": map[string]any{
			"from": "media",
			"let":  map[string]any{"eid": "$eventid"},
			"pipeline": []any{
				map[string]any{"$match": map[string]any{
					"$expr": map[string]any{
						"$and": []any{
							map[string]any{"$eq": []any{"$entityid", "$$eid"}},
							map[string]any{"$eq": []any{"$entitytype", "event"}},
						},
					},
				}},
			},
			"as": "media",
		}},
		map[string]any{"$lookup": map[string]any{
			"from": "merch",
			"let":  map[string]any{"eid": "$eventid"},
			"pipeline": []any{
				map[string]any{"$match": map[string]any{
					"$expr": map[string]any{
						"$and": []any{
							map[string]any{"$eq": []any{"$entity_id", "$$eid"}},
							map[string]any{"$eq": []any{"$entity_type", "event"}},
						},
					},
				}},
			},
			"as": "merch",
		}},
	}

	return app.DB.Aggregate(ctx, eventsTable, pipeline, result)
}

func sqlSQLlistEvents(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions, result *[]Event) error {
	return app.DB.FindManyWithOptions(ctx, eventsTable, filter, opts, result)
}

func sqlSQLcountEvents(ctx context.Context, app *infra.Deps, filter map[string]any) (int64, error) {
	return app.DB.CountDocuments(ctx, eventsTable, filter)
}
