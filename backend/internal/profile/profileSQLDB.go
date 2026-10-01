// File: internal/profile/profileSQLDB.go

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
var usersCollection = usersTable

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
	if err := app.SQLDB.FindOne(ctx, usersTable, query, args, &user); err != nil {
		return nil, err
	}
	if user.UserID == "" {
		return nil, nil
	}
	return &user, nil
}

func SQLRespondWithUserProfile(w http.ResponseWriter, userid string, app *infra.Deps) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "userid = $1"
	args := []any{userid}
	var userProfile auth.User
	if err := app.SQLDB.FindOne(ctx, usersTable, query, args, &userProfile); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	if userProfile.UserID == "" {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, userProfile)
}

// SQLRespondWithUserProfile writes user profile as JSON to the response
func RespondWithUserProfile(w http.ResponseWriter, userid string, app *infra.Deps) {
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
func SQLApplyProfileUpdates(ctx context.Context, app *infra.Deps, userID string, updates map[string]any) (int64, error) {
	query := "userid = $1"
	args := []any{userID}

	rowsAffected, err := app.SQLDB.UpdateOne(ctx, usersTable, query, args, updates)
	return rowsAffected, err
}

func SQLDeleteUserByID(ctx context.Context, app *infra.Deps, userID string) (int64, error) {
	query := "userid = $1"
	args := []any{userID}

	return app.SQLDB.DeleteOne(ctx, usersTable, query, args)
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
