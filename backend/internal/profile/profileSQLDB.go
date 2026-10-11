// File: internal/profile/profileSQLDB.go

package profile

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
	if app == nil || app.SQLDB == nil {
		return &auth.User{}, nil
	}
	if query == "" {
		query = "1 = 1"
	}
	var user auth.User
	row := app.SQLDB.QueryRow(ctx, `SELECT * FROM `+usersTable+` WHERE `+query+` LIMIT 1`, args...)
	if row == nil {
		return &user, nil
	}
	if err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Name); err != nil {
		return &user, err
	}
	return &user, nil
}

func SQLRespondWithUserProfile(w http.ResponseWriter, userid string, app *infra.Deps) {
	if w == nil {
		return
	}
	profile := map[string]any{"userid": userid}
	if app != nil && app.SQLDB != nil && userid != "" {
		user, err := SQLfindUser(context.Background(), app, "userid = $1", []any{userid})
		if err == nil && user != nil {
			profile = map[string]any{"userid": user.UserID, "username": user.Username, "email": user.Email, "displayname": user.Name}
		}
	}
	utils.RespondWithJSON(w, http.StatusOK, profile)
}

// SQLRespondWithUserProfile writes user profile as JSON to the response
func RespondWithUserProfile(w http.ResponseWriter, userid string, app *infra.Deps) {
	SQLRespondWithUserProfile(w, userid, app)
}

/*
	-------------------------------------------------------
	  Apply updates / Delete user

-------------------------------------------------------
*/
func SQLApplyProfileUpdates(ctx context.Context, app *infra.Deps, userID string, updates map[string]any) (int64, error) {
	if app == nil || app.SQLDB == nil || userID == "" || len(updates) == 0 {
		return 0, nil
	}
	parts := make([]string, 0, len(updates))
	args := make([]any, 0, len(updates)+1)
	for key, value := range updates {
		parts = append(parts, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value)
	}
	args = append(args, userID)
	setClause := strings.Join(parts, ", ")
	result, err := app.SQLDB.Exec(ctx, `UPDATE `+usersTable+` SET `+setClause+` WHERE userid = $`+fmt.Sprintf("%d", len(args)), args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func SQLDeleteUserByID(ctx context.Context, app *infra.Deps, userID string) (int64, error) {
	if app == nil || app.SQLDB == nil || userID == "" {
		return 0, nil
	}
	result, err := app.SQLDB.Exec(ctx, `DELETE FROM `+usersTable+` WHERE userid = $1`, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
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
