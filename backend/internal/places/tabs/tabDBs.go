// File: internal/places/tabs/tabDBs.go

package places

import (
	"context"
	"time"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
	"scav/internal/events"
	placedb "scav/internal/places/placedb"
)

var (
	eventsTable   = config.Tables.EventsTable
	productsTable = config.Tables.ProductTable
)

func buildPlaceEventQuery(placeID string, now time.Time) (string, []any) {
	return "placeid = $1 AND date >= $2", []any{placeID, now}
}

func buildPlaceEventListOptions(page, limit int) sqldb.FindManyOptions {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	skip := (page - 1) * limit
	return sqldb.FindManyOptions{
		Limit:   int64(limit),
		Offset:  int64(skip),
		OrderBy: "date ASC",
		Columns: []string{
			"eventid",
			"title",
			"description",
			"start_date_time",
			"end_date_time",
			"placename",
			"banner",
			"category",
		},
	}
}

func getPlaceEventsPage(ctx context.Context, app *infra.Deps, placeID string, page, limit int, now time.Time) ([]events.Event, int64, error) {
	query, args := buildPlaceEventQuery(placeID, now)
	opts := buildPlaceEventListOptions(page, limit)

	var placeevents []events.Event
	if err := placedb.FindEventsWithOptions(ctx, app, query, args, opts, &placeevents); err != nil {
		return nil, 0, err
	}

	total, err := placedb.CountEvents(ctx, app, query, args)
	if err != nil {
		return nil, 0, err
	}

	if placeevents == nil {
		placeevents = []events.Event{}
	}

	return placeevents, total, nil
}
