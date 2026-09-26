package home

import (
	"context"

	"scav/infra"
	"scav/infra/sqldb"
)

func SQLfetchHomeCardsFromDB(ctx context.Context, app *infra.Deps, category string, offset, limit int) ([]HomeCard, error) {
	table, projector := categoryProjection(category)
	if table == "" || projector == nil {
		return []HomeCard{}, nil
	}

	opts := sqldb.FindManyOptions{
		Offset:  offset,
		Limit:   limit,
		OrderBy: "created_at DESC",
	}

	var docs []map[string]any
	// Passing empty query and empty args to select all records
	if err := app.SQLDB.FindManyWithOptions(ctx, table, "", nil, opts, &docs); err != nil {
		return nil, err
	}

	cards := make([]HomeCard, 0, len(docs))
	for _, doc := range docs {
		cards = append(cards, projector(doc))
	}

	return cards, nil
}
