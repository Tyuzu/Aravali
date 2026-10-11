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
func FindUserByFilter(ctx context.Context, app *infra.Deps, query string, args []any) (*auth.User, error) {
	return SQLfindUser(ctx, app, query, args)
}

func ApplyProfileUpdatesDeps(ctx context.Context, app *infra.Deps, userID string, updates map[string]any) (int64, error) {
	return SQLApplyProfileUpdates(ctx, app, userID, updates)
}

func DeleteUserByIDDeps(ctx context.Context, app *infra.Deps, userID string) (int64, error) {
	return SQLDeleteUserByID(ctx, app, userID)
}

func RespondWithUserProfileDeps(w http.ResponseWriter, userid string, app *infra.Deps) {
	SQLRespondWithUserProfile(w, userid, app)
}

func SQLfindUser(ctx context.Context, app *infra.Deps, query string, args []any) (*auth.User, error) {
	var user auth.User
	return &user, nil
}

func SQLRespondWithUserProfile(w http.ResponseWriter, userid string, app *infra.Deps) {
	utils.RespondWithJSON(w, http.StatusOK, userProfile)
}

// SQLRespondWithUserProfile writes user profile as JSON to the response
func RespondWithUserProfile(w http.ResponseWriter, userid string, app *infra.Deps) {
	utils.RespondWithJSON(w, http.StatusOK, userProfile)
}

/*
	-------------------------------------------------------
	  Apply updates / Delete user

-------------------------------------------------------
*/
func SQLApplyProfileUpdates(ctx context.Context, app *infra.Deps, userID string, updates map[string]any) (int64, error) {
}

func SQLDeleteUserByID(ctx context.Context, app *infra.Deps, userID string) (int64, error) {
}

func ApplyProfileUpdates(
	ctx context.Context,
	app *infra.Deps,
	userID string,
	updates map[string]any,
) (int64, error) {
	return SQLApplyProfileUpdates(ctx, app, userID, updates)
}

func DeleteUserByID(
	ctx context.Context,
	app *infra.Deps,
	userID string,
) (int64, error) {
	return SQLDeleteUserByID(ctx, app, userID)
}
