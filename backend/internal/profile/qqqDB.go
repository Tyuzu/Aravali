// File: internal/profile/profileSQLDB.go

package profile

import (
	"context"
	"net/http"

	"scav/config"
	"scav/infra"
	"scav/internal/auth"
	"scav/utils"
)

var usersTable = config.Tables.UserTable

// Wrapper helpers that accept infra.Deps to simplify call sites.
func QQQFindUserByFilter(ctx context.Context, app *infra.Deps, query string, args []any) (*auth.User, error) {
	return QQQSQLfindUser(ctx, app, query, args)
}

func QQQApplyProfileUpdatesDeps(ctx context.Context, app *infra.Deps, userID string, updates map[string]any) (int64, error) {
	return QQQSQLApplyProfileUpdates(ctx, app, userID, updates)
}

func QQQDeleteUserByIDDeps(ctx context.Context, app *infra.Deps, userID string) (int64, error) {
	return QQQSQLDeleteUserByID(ctx, app, userID)
}

func QQQRespondWithUserProfileDeps(w http.ResponseWriter, userid string, app *infra.Deps) {
	QQQSQLRespondWithUserProfile(w, userid, app)
}

func QQQSQLfindUser(ctx context.Context, app *infra.Deps, query string, args []any) (*auth.User, error) {
	user := auth.User{}
	return &user, nil
}

func QQQSQLRespondWithUserProfile(w http.ResponseWriter, userid string, app *infra.Deps) {

	profile := map[string]any{"userid": userid}
	utils.RespondWithJSON(w, http.StatusOK, profile)
}

// SQLRespondWithUserProfile writes user profile as JSON to the response
func QQQRespondWithUserProfile(w http.ResponseWriter, userid string, app *infra.Deps) {
	QQQSQLRespondWithUserProfile(w, userid, app)
}

/*
	-------------------------------------------------------
	  Apply updates / Delete user

-------------------------------------------------------
*/
func QQQSQLApplyProfileUpdates(ctx context.Context, app *infra.Deps, userID string, updates map[string]any) (int64, error) {

	return 1, nil
}

func QQQSQLDeleteUserByID(ctx context.Context, app *infra.Deps, userID string) (int64, error) {

	return 1, nil
}

func QQQApplyProfileUpdates(
	ctx context.Context,
	app *infra.Deps,
	userID string,
	updates map[string]any,
) (int64, error) {
	return QQQSQLApplyProfileUpdates(ctx, app, userID, updates)
}

func QQQDeleteUserByID(
	ctx context.Context,
	app *infra.Deps,
	userID string,
) (int64, error) {
	return QQQSQLDeleteUserByID(ctx, app, userID)
}
