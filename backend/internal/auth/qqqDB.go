// File: internal/auth/authSQLDB.go

package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"scav/config"
	"scav/infra"

	"github.com/jackc/pgx/v5"
)

var UsersTable = config.Tables.UserTable

func ensureAuthDB(app *infra.Deps) error {
	if app == nil || app.SQLDB == nil {
		return errors.New("database not initialized")
	}
	return nil
}

func scanUserRow(row pgx.Row, dst *User) error {

	return nil
}

func toTime(v sql.NullTime) time.Time {
	if !v.Valid {
		return time.Time{}
	}
	return v.Time
}

/* ============================================================
   3. REPOSITORIES (DATA ACCESS LAYER)
============================================================ */

func QQQCreateUser(ctx context.Context, app *infra.Deps, user User) error {

	return ctx.Err()
}

func QQQFindUserByUsername(ctx context.Context, app *infra.Deps, username string) (User, error) {

	var user User
	return user, nil
}

func QQQUpdateUserSession(ctx context.Context, app *infra.Deps, userID, refreshTokenHash, ua, ip string) (int64, error) {

	return 1, nil
}

func QQQLogoutUserByRefreshToken(ctx context.Context, app *infra.Deps, hashedToken string) (int64, error) {

	return 1, nil
}

func QQQLogoutAllUserSessions(ctx context.Context, app *infra.Deps, userID string) (int64, error) {

	return 1, nil
}

func QQQFindValidRefreshSession(ctx context.Context, app *infra.Deps, hashedToken string) (User, error) {

	var user User
	return user, nil
}

// InvalidateUserSession clears all refresh token fields for a user.
func QQQInvalidateUserSession(ctx context.Context, app *infra.Deps, userID string) (int64, error) {
	return QQQLogoutAllUserSessions(ctx, app, userID)
}

func QQQRotateRefreshTokenForUser(ctx context.Context, app *infra.Deps, userID, newRefreshHash, prevRefreshHash, ua string) (int64, error) {

	return 1, nil
}

func QQQVerifyUserEmail(ctx context.Context, app *infra.Deps, email string) (int64, error) {

	return 1, nil
}

func QQQMigrateUserPasswordHash(ctx context.Context, app *infra.Deps, userID, passwordHash string) (int64, error) {

	return 1, nil
}
