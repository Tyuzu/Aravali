package home

import (
	"context"

	"scav/infra"
	"scav/infra/db"

	"go.mongodb.org/mongo-driver/bson"
)

func fetchHomeCardsFromDB(ctx context.Context, app *infra.Deps, category string, skip, limit int) ([]HomeCard, error) {
	collection, projector := categoryProjection(category)
	if collection == "" || projector == nil {
		return []HomeCard{}, nil
	}

	opts := db.FindManyOptions{
		Skip:  skip,
		Limit: limit,
		Sort:  []bson.E{{Key: "createdAt", Value: -1}},
	}

	var docs []map[string]any
	if err := app.DB.FindManyWithOptions(ctx, collection, map[string]any{}, opts, &docs); err != nil {
		return nil, err
	}

	cards := make([]HomeCard, 0, len(docs))
	for _, doc := range docs {
		cards = append(cards, projector(doc))
	}

	return cards, nil
}
