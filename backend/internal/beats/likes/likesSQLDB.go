package likes

import (
	"context"
	"errors"
	"scav/config"
	"scav/infra"

	"go.mongodb.org/mongo-driver/mongo"
)

var usersTable = config.Tables.UserTable
var likesTable = config.Tables.LikesTable

type SQLRepository interface {
	Insert(ctx context.Context, like Like) error
	Delete(ctx context.Context, userID, entityType, entityID string) (int64, error)
	FindOne(ctx context.Context, userID, entityType, entityID string) (bool, error)
	Count(ctx context.Context, entityType, entityID string) (int64, error)
}

type pgxRepository struct {
	app *infra.Deps
}

func SQLNewPgxRepository(app *infra.Deps) Repository {
	return &pgxRepository{app: app}
}

func (r *pgxRepository) Insert(ctx context.Context, like Like) error {
	return r.app.DB.Insert(ctx, likesTable, like)
}

func (r *pgxRepository) Delete(ctx context.Context, userID, entityType, entityID string) (int64, error) {
	filter := map[string]any{
		"userid":      userID,
		"entity_type": entityType,
		"entity_id":   entityID,
	}
	return r.app.DB.DeleteOne(ctx, likesTable, filter)
}

func (r *pgxRepository) FindOne(ctx context.Context, userID, entityType, entityID string) (bool, error) {
	filter := map[string]any{
		"userid":      userID,
		"entity_type": entityType,
		"entity_id":   entityID,
	}

	var like Like
	err := r.app.DB.FindOne(ctx, likesTable, filter, &like)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}

	return false, err
}

func (r *pgxRepository) Count(ctx context.Context, entityType, entityID string) (int64, error) {
	return r.app.DB.CountDocuments(
		ctx,
		likesTable,
		map[string]any{
			"entity_type": entityType,
			"entity_id":   entityID,
		},
	)
}

func SQLFindLikesByEntity(ctx context.Context, app *infra.Deps, entityType, entityID string) ([]Like, error) {
	var likes []Like
	err := app.DB.FindMany(
		ctx,
		likesTable,
		map[string]any{
			"entity_type": entityType,
			"entity_id":   entityID,
		},
		&likes,
	)
	return likes, err
}

func SQLFindUserLikesByEntityIDs(ctx context.Context, app *infra.Deps, userID, entityType string, entityIDs []string) ([]Like, error) {
	var likes []Like
	err := app.DB.FindMany(
		ctx,
		likesTable,
		map[string]any{
			"userid":      userID,
			"entity_type": entityType,
			"entity_id":   map[string]any{"$in": entityIDs},
		},
		&likes,
	)
	return likes, err
}

func SQLFindUsersByIDs(ctx context.Context, app *infra.Deps, userIDs []string) ([]map[string]string, error) {
	var users []struct {
		UserID   string `bson:"userid"`
		Username string `bson:"username"`
		Avatar   string `bson:"avatar,omitempty"`
	}

	err := app.DB.FindMany(
		ctx,
		usersTable,
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
