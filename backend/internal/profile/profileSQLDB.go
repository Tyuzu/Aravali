package profile

import (
	"context"
	"net/http"
	"time"

	"scav/config"
	"scav/infra"
	"scav/internal/auth"
	"scav/utils"
)

var usersTable = config.Tables.UserTable

// Wrapper helpers that accept infra.Deps to simplify call sites.
func SQLFindUserByFilter(ctx context.Context, app *infra.Deps, query string, args []any) (*auth.User, error) {
	return SQLfindUser(ctx, query, args, app)
}

func SQLApplyProfileUpdatesDeps(ctx context.Context, app *infra.Deps, userID string, updates map[string]any) (int64, error) {
	return SQLApplyProfileUpdates(ctx, app, userID, updates)
}

func SQLDeleteUserByIDDeps(ctx context.Context, app *infra.Deps, userID string) (int64, error) {
	return SQLDeleteUserByID(ctx, app, userID)
}

func SQLRespondWithUserProfileDeps(w http.ResponseWriter, userid string, app *infra.Deps) {
	SQLRespondWithUserProfile(w, userid, app)
}

// SQLfindUser returns a user by SQL query and args, or nil if not found
func SQLfindUser(ctx context.Context, query string, args []any, app *infra.Deps) (*auth.User, error) {
	var user auth.User
	_ = app.SQLDB.FindOne(ctx, usersTable, query, args, &user)
	// ignore errors; return nil if user not found
	if user.UserID == "" {
		return nil, nil
	}
	return &user, nil
}

// SQLRespondWithUserProfile writes user profile as JSON to the response
func SQLRespondWithUserProfile(w http.ResponseWriter, userid string, app *infra.Deps) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "userid = $1"
	args := []any{userid}

	var userProfile auth.User
	_ = app.SQLDB.FindOne(ctx, usersTable, query, args, &userProfile)

	if userProfile.UserID == "" {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, userProfile)
}

/*
	-------------------------------------------------------
	  Apply updates / Delete user

-------------------------------------------------------
*/
func SQLApplyProfileUpdates(
	ctx context.Context,
	app *infra.Deps,
	userID string,
	updates map[string]any,
) (int64, error) {
	query := "userid = $1"
	args := []any{userID}

	rowsAffected, err := app.SQLDB.UpdateOne(ctx, usersCollection, query, args, updates)
	return rowsAffected, err
}

func SQLDeleteUserByID(
	ctx context.Context,
	app *infra.Deps,
	userID string,
) (int64, error) {
	query := "userid = $1"
	args := []any{userID}

	return app.SQLDB.DeleteOne(ctx, usersCollection, query, args)
}
