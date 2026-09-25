package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
)

var UsersTable = config.Tables.UserTable

/* ============================================================
   3. REPOSITORIES (DATA ACCESS LAYER)
============================================================ */

func SQLCreateUser(ctx context.Context, app *infra.Deps, user User) error {
	err := app.DB.Insert(ctx, UsersTable, user)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func SQLFindUserByUsername(ctx context.Context, app *infra.Deps, username string) (User, error) {
	var user User
	if err := app.DB.FindOne(ctx, UsersTable, map[string]any{"username": username}, &user); err != nil {
		return User{}, err
	}
	return user, nil
}

func SQLUpdateUserSession(ctx context.Context, app *infra.Deps, userID, refreshTokenHash, ua, ip string) (any, error) {
	update := map[string]any{
		"refresh_token":  refreshTokenHash,
		"refresh_expiry": time.Now().Add(RefreshTokenTTL),
		"refresh_ua":     ua,
		"refresh_ip":     ip,
		"last_login":     time.Now(),
		"online":         true,
		"updated_at":     time.Now(),
	}
	return app.DB.Update(ctx, UsersTable, map[string]any{"userid": userID}, update)
}

func SQLLogoutUserByRefreshToken(ctx context.Context, app *infra.Deps, hashedToken string) (any, error) {
	update := map[string]any{
		"refresh_token":  nil,
		"refresh_expiry": nil,
		"online":         false,
		"updated_at":     time.Now(),
	}
	return app.DB.Update(ctx, UsersTable, map[string]any{"refresh_token": hashedToken}, update)
}

func SQLLogoutAllUserSessions(ctx context.Context, app *infra.Deps, userID string) (any, error) {
	update := map[string]any{
		"refresh_token":  nil,
		"refresh_prev":   nil,
		"refresh_expiry": nil,
		"refresh_ua":     nil,
		"refresh_ip":     nil,
		"online":         false,
		"updated_at":     time.Now(),
	}
	return app.DB.Update(ctx, UsersTable, map[string]any{"userid": userID}, update)
}

func SQLFindValidRefreshSession(ctx context.Context, app *infra.Deps, hashedToken string) (User, error) {
	now := time.Now()
	var user User

	// Try direct match against current refresh token
	if err := app.DB.FindOne(ctx, UsersTable, map[string]any{"refresh_token": hashedToken}, &user); err == nil {
		if user.RefreshExpiry.After(now) {
			return user, nil
		}
		return User{}, fmt.Errorf("refresh token expired")
	}

	// Try match against previous refresh token
	if err := app.DB.FindOne(ctx, UsersTable, map[string]any{"refresh_prev": hashedToken}, &user); err == nil {
		if user.RefreshExpiry.After(now) {
			return user, nil
		}
		return User{}, fmt.Errorf("refresh token expired")
	}

	return User{}, fmt.Errorf("no valid refresh session")
}

// InvalidateUserSession clears all refresh token fields for a user.
func SQLInvalidateUserSession(ctx context.Context, app *infra.Deps, userID string) (any, error) {
	update := map[string]any{
		"refresh_token":  nil,
		"refresh_prev":   nil,
		"refresh_expiry": nil,
		"refresh_ua":     nil,
		"updated_at":     time.Now(),
	}
	return app.DB.Update(ctx, UsersTable, map[string]any{"userid": userID}, update)
}

func SQLRotateRefreshTokenForUser(ctx context.Context, app *infra.Deps, userID, newRefreshHash, prevRefreshHash, ua string) (any, error) {
	now := time.Now()
	update := map[string]any{
		"refresh_prev":   prevRefreshHash,
		"refresh_token":  newRefreshHash,
		"refresh_expiry": now.Add(RefreshTokenTTL),
		"refresh_ua":     ua,
		"updated_at":     now,
	}
	return app.DB.Update(ctx, UsersTable, map[string]any{"userid": userID}, update)
}

func SQLVerifyUserEmail(ctx context.Context, app *infra.Deps, email string) (any, error) {
	update := map[string]any{"email_verified": true}
	return app.DB.Update(ctx, UsersTable, map[string]any{"email": email}, update)
}

func SQLMigrateUserPasswordHash(ctx context.Context, app *infra.Deps, userID, passwordHash string) (any, error) {
	update := map[string]any{"$set": map[string]any{"password_hash": passwordHash}}
	return app.DB.Update(ctx, UsersTable, map[string]any{"userid": userID}, update)
}
