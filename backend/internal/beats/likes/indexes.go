package likes

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureIndexes(
	ctx context.Context,
	collection *mongo.Collection,
) error {

	_, err := collection.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys: bson.M{
					"userid":      1,
					"entity_type": 1,
					"entity_id":   1,
				},
				Options: options.Index().
					SetUnique(true).
					SetName("unique_user_entity_like"),
			},
			{
				Keys: bson.M{
					"entity_type": 1,
					"entity_id":   1,
					"created_at":  -1,
				},
				Options: options.Index().
					SetName("entity_likes_created_at"),
			},
		},
	)

	return err
}
