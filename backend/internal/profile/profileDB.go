package profile

import (
	"context"
	"net/http"
	"time"

	"scav/config"
	"scav/infra"
	"scav/infra/db"
	"scav/internal/auth"
	"scav/utils"
)

var usersCollection = config.Collections.UserCollection

// Wrapper helpers that accept infra.Deps to simplify call sites.
func FindUserByFilter(ctx context.Context, app *infra.Deps, filter map[string]any) (*auth.User, error) {
	return findUser(ctx, filter, app.DB)
}

func ApplyProfileUpdatesDeps(ctx context.Context, app *infra.Deps, userID string, updates map[string]any) (any, error) {
	return ApplyProfileUpdates(ctx, app.DB, userID, updates)
}

func DeleteUserByIDDeps(ctx context.Context, app *infra.Deps, userID string) (int64, error) {
	return DeleteUserByID(ctx, app.DB, userID)
}

func RespondWithUserProfileDeps(w http.ResponseWriter, userid string, app *infra.Deps) {
	RespondWithUserProfile(w, userid, app.DB)
}

// findUser returns a user by filter, or nil if not found
func findUser(ctx context.Context, filter map[string]any, database db.Database) (*auth.User, error) {
	var user auth.User
	_ = database.FindOne(ctx, usersCollection, filter, &user)
	// ignore errors; return nil if user not found
	if user.UserID == "" {
		return nil, nil
	}
	return &user, nil
}

// RespondWithUserProfile writes user profile as JSON to the response
func RespondWithUserProfile(w http.ResponseWriter, userid string, database db.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userProfile auth.User
	_ = database.FindOne(ctx, usersCollection, map[string]any{"userid": userid}, &userProfile)

	if userProfile.UserID == "" {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")

	}

	utils.RespondWithJSON(w, http.StatusOK, userProfile)
}
