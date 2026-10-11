// File: internal/auth/authSQLDB.go

package auth

import (
	"context"
	"fmt"

	"scav/config"
	"scav/infra"
)

var UsersTable = config.Tables.UserTable

/* ============================================================
   3. REPOSITORIES (DATA ACCESS LAYER)
============================================================ */

func CreateUser(ctx context.Context, app *infra.Deps, user User) error {

	return nil
}

func FindUserByUsername(ctx context.Context, app *infra.Deps, username string) (User, error) {
	var user User

	return user, nil
}

func UpdateUserSession(ctx context.Context, app *infra.Deps, userID, refreshTokenHash, ua, ip string) (int64, error) {

}

func LogoutUserByRefreshToken(ctx context.Context, app *infra.Deps, hashedToken string) (int64, error) {

}

func LogoutAllUserSessions(ctx context.Context, app *infra.Deps, userID string) (int64, error) {

}

func FindValidRefreshSession(ctx context.Context, app *infra.Deps, hashedToken string) (User, error) {

	return User{}, fmt.Errorf("no valid refresh session")
}

// InvalidateUserSession clears all refresh token fields for a user.
func InvalidateUserSession(ctx context.Context, app *infra.Deps, userID string) (int64, error) {

}

func RotateRefreshTokenForUser(ctx context.Context, app *infra.Deps, userID, newRefreshHash, prevRefreshHash, ua string) (int64, error) {

}

func VerifyUserEmail(ctx context.Context, app *infra.Deps, email string) (int64, error) {

}

func MigrateUserPasswordHash(ctx context.Context, app *infra.Deps, userID, passwordHash string) (int64, error) {

}
