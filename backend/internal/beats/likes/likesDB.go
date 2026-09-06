package likes

import (
	"context"
	"errors"
	"scav/config"
	"scav/infra"

	"go.mongodb.org/mongo-driver/mongo"
)

var usersCollection = config.Collections.UserCollection
var likesCollection = config.Collections.LikesCollection

type Repository interface {
	Insert(ctx context.Context, like Like) error
	Delete(ctx context.Context, userID, entityType, entityID string) (int64, error)
	FindOne(ctx context.Context, userID, entityType, entityID string) (bool, error)
	Count(ctx context.Context, entityType, entityID string) (int64, error)
}

type mongoRepository struct {
	app *infra.Deps
}

func NewMongoRepository(app *infra.Deps) Repository {
	return &mongoRepository{app: app}
}

func (r *mongoRepository) Insert(ctx context.Context, like Like) error {
	return r.app.DB.Insert(ctx, likesCollection, like)
}

func (r *mongoRepository) Delete(ctx context.Context, userID, entityType, entityID string) (int64, error) {
	filter := map[string]any{
		"userid":      userID,
		"entity_type": entityType,
		"entity_id":   entityID,
	}
	return r.app.DB.DeleteOne(ctx, likesCollection, filter)
}

func (r *mongoRepository) FindOne(ctx context.Context, userID, entityType, entityID string) (bool, error) {
	filter := map[string]any{
		"userid":      userID,
		"entity_type": entityType,
		"entity_id":   entityID,
	}

	var like Like
	err := r.app.DB.FindOne(ctx, likesCollection, filter, &like)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}

	return false, err
}

func (r *mongoRepository) Count(ctx context.Context, entityType, entityID string) (int64, error) {
	return r.app.DB.CountDocuments(
		ctx,
		likesCollection,
		map[string]any{
			"entity_type": entityType,
			"entity_id":   entityID,
		},
	)
}

func FindLikesByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Like, error) {
	var likes []Like
	err := app.DB.FindMany(
		ctx,
		likesCollection,
		map[string]any{
			"entity_type": entityType,
			"entity_id":   entityID,
		},
		&likes,
	)
	return likes, err
}

func FindUserLikesByEntityIDs(ctx context.Context, app *infra.Deps, userID, entityType string, entityIDs []string) ([]Like, error) {
	var likes []Like
	err := app.DB.FindMany(
		ctx,
		likesCollection,
		map[string]any{
			"userid":      userID,
			"entity_type": entityType,
			"entity_id": map[string]any{"$in": entityIDs},
		},
		&likes,
	)
	return likes, err
}

func FindUsersByIDs(ctx context.Context, app *infra.Deps, userIDs []string) ([]map[string]string, error) {
	var users []struct {
		UserID   string `bson:"userid"`
		Username string `bson:"username"`
		Avatar   string `bson:"avatar,omitempty"`
	}

	err := app.DB.FindMany(
		ctx,
		usersCollection,
		map[string]any{"userid": map[string]any{"$in": userIDs}},
		&users,
	)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]string, 0, len(users))
	for _, user := range users {
		result = append(result, map[string]string{
			"userid":   user.UserID,
			"username": user.Username,
			"avatar":   user.Avatar,
		})
	}
	return result, nil
}
