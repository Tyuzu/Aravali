package events

import (
	"context"

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

func sqlSQLinsertEvent(ctx context.Context, app *infra.Deps, event Event) error {
	return app.SQLDB.InsertOne(ctx, eventsTable, event)
}

func sqlSQLensureUniqueEventID(ctx context.Context, app *infra.Deps, event *Event) {
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

func sqlSQLfindEventByID(ctx context.Context, app *infra.Deps, eventID string, event *Event) error {
	query := "eventid = $1"
	args := []any{eventID}

	return app.SQLDB.FindOne(ctx, eventsTable, query, args, event)
}

func sqlSQLupdateEvent(ctx context.Context, app *infra.Deps, eventID string, updates map[string]any) (int64, error) {
	query := "eventid = $1"
	args := []any{eventID}

	return app.SQLDB.UpdateOne(ctx, eventsTable, query, args, updates)
}

func sqlSQLaggregateEvent(ctx context.Context, app *infra.Deps, eventID string, result *[]Event) error {
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

func sqlSQLlistEvents(ctx context.Context, app *infra.Deps, query string, args []any, opts sqldb.FindManyOptions, result *[]Event) error {
	return app.SQLDB.FindManyWithOptions(ctx, eventsTable, query, args, opts, result)
}

func sqlSQLcountEvents(ctx context.Context, app *infra.Deps, whereClause string, args []any) (int64, error) {
	return app.SQLDB.Count(ctx, eventsTable, whereClause, args)
}
