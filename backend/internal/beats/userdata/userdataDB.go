package userdata

import (
	"context"
	"scav/config"
	"scav/infra"
	"time"
)

var userdataCollection = config.Collections.UserDataCollection

// InsertUserData inserts a single user data document.
func InsertUserData(ctx context.Context, app *infra.Deps, content UserData) error {
	return app.DB.InsertOne(ctx, userdataCollection, content)
}

// DeleteUserData removes user data matching filter.
func DeleteUserData(ctx context.Context, app *infra.Deps, filter map[string]any) error {
	return app.DB.DeleteMany(ctx, userdataCollection, filter)
}

// InsertUserDataMany inserts many user data documents.
func InsertUserDataMany(ctx context.Context, app *infra.Deps, docs []any) error {
	return app.DB.InsertMany(ctx, userdataCollection, docs)
}

// FindUserData finds user data for a given filter.
func FindUserData(ctx context.Context, app *infra.Deps, filter map[string]any, out *[]UserData) error {
	return app.DB.FindMany(ctx, userdataCollection, filter, out)
}

// Data Models
type postDoc struct {
	PostID    string    `bson:"postid"`
	Title     string    `bson:"title"`
	Thumb     string    `bson:"thumb"`
	CreatedBy string    `bson:"createdBy"`
	Username  string    `bson:"username"`
	CreatedAt time.Time `bson:"createdAt"`
	Blocks    []struct {
		Type    string `bson:"type"`
		URL     string `bson:"url"`
		Caption string `bson:"caption"`
	} `bson:"blocks"`
}

// Database Helpers

// FetchUserDataByEntity queries user data based on entity type and user ID.
func FetchUserDataByEntity(ctx context.Context, app *infra.Deps, entityType, username string) ([]UserData, error) {
	filter := map[string]any{
		"entity_type": entityType,
		"userid":      username,
	}

	var results []UserData
	if err := FindUserData(ctx, app, filter, &results); err != nil {
		return nil, err
	}

	if results == nil {
		return []UserData{}, nil
	}

	return results, nil
}

// FetchOtherUserFeedPosts retrieves posts for a given user from the database.
func FetchOtherUserFeedPosts(ctx context.Context, app *infra.Deps, username string) ([]postDoc, error) {
	filter := map[string]any{
		"$or": []map[string]any{
			{"createdBy": username},
			{"username": username},
		},
	}

	var posts []postDoc
	if err := app.DB.FindMany(ctx, config.Collections.FeedPostsCollection, filter, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}
