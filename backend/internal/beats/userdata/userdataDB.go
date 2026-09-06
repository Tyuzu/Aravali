package userdata

import (
	"context"
	"scav/config"
	"scav/infra"
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
