// File: internal/auth/authSQLDB.go

package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
	if dst == nil {
		return errors.New("nil user result")
	}

	var (
		userID         string
		username       string
		email          string
		password       sql.NullString
		passwordHash   sql.NullString
		role           []string
		name           sql.NullString
		createdAt      sql.NullTime
		updatedAt      sql.NullTime
		bio            sql.NullString
		online         bool
		lastLogin      sql.NullTime
		avatar         sql.NullString
		banner         sql.NullString
		profileViews   int
		phoneNumber    sql.NullString
		address        sql.NullString
		socialLinks    []byte
		isVerified     bool
		emailVerified  bool
		followersCount int
		followingCount int
		walletBalance  float64
		refreshToken   sql.NullString
		refreshExpiry  sql.NullTime
		refreshUA      sql.NullString
		refreshIP      sql.NullString
		refreshPrev    sql.NullString
	)

	if err := row.Scan(
		&userID,
		&username,
		&email,
		&password,
		&passwordHash,
		&role,
		&name,
		&createdAt,
		&updatedAt,
		&bio,
		&online,
		&lastLogin,
		&avatar,
		&banner,
		&profileViews,
		&phoneNumber,
		&address,
		&socialLinks,
		&isVerified,
		&emailVerified,
		&followersCount,
		&followingCount,
		&walletBalance,
		&refreshToken,
		&refreshExpiry,
		&refreshUA,
		&refreshIP,
		&refreshPrev,
	); err != nil {
		return err
	}

	*dst = User{
		UserID:         userID,
		Username:       username,
		Email:          email,
		Role:           role,
		Name:           name.String,
		CreatedAt:      toTime(createdAt),
		UpdatedAt:      toTime(updatedAt),
		Bio:            bio.String,
		Online:         online,
		LastLogin:      toTime(lastLogin),
		Avatar:         avatar.String,
		Banner:         banner.String,
		ProfileViews:   profileViews,
		PhoneNumber:    phoneNumber.String,
		Address:        address.String,
		IsVerified:     isVerified,
		EmailVerified:  emailVerified,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
		WalletBalance:  walletBalance,
		RefreshToken:   refreshToken.String,
		RefreshExpiry:  toTime(refreshExpiry),
		RefreshUA:      refreshUA.String,
		RefreshIP:      refreshIP.String,
		RefreshPrev:    refreshPrev.String,
	}

	if password.Valid {
		dst.Password = password.String
	}
	if passwordHash.Valid {
		dst.PasswordHash = passwordHash.String
	}
	if len(socialLinks) > 0 && string(socialLinks) != "null" {
		if err := json.Unmarshal(socialLinks, &dst.SocialLinks); err != nil {
			dst.SocialLinks = map[string]string{}
		}
	}
	if dst.SocialLinks == nil {
		dst.SocialLinks = map[string]string{}
	}
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

func CreateUser(ctx context.Context, app *infra.Deps, user User) error {
	if err := ensureAuthDB(app); err != nil {
		return err
	}

	socialLinks := map[string]string{}
	if user.SocialLinks != nil {
		socialLinks = user.SocialLinks
	}
	payload, err := json.Marshal(socialLinks)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (
		userid, username, email, password, password_hash, role, name, created_at, updated_at,
		bio, online, last_login, avatar, banner, profile_views, phone_number, address,
		social_links, is_verified, email_verified, followerscount, followscount,
		wallet_balance, metadata
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24)`
	_, err = app.SQLDB.Exec(ctx, query,
		user.UserID,
		user.Username,
		user.Email,
		user.Password,
		user.PasswordHash,
		user.Role,
		user.Name,
		user.CreatedAt,
		user.UpdatedAt,
		user.Bio,
		user.Online,
		user.LastLogin,
		user.Avatar,
		user.Banner,
		user.ProfileViews,
		user.PhoneNumber,
		user.Address,
		payload,
		user.IsVerified,
		user.EmailVerified,
		user.FollowersCount,
		user.FollowingCount,
		user.WalletBalance,
		map[string]any{},
	)
	return err
}

func FindUserByUsername(ctx context.Context, app *infra.Deps, username string) (User, error) {
	if err := ensureAuthDB(app); err != nil {
		return User{}, err
	}

	candidate := strings.TrimSpace(username)
	query := `SELECT
		userid, username, email, password, password_hash, role, name, created_at, updated_at,
		bio, online, last_login, avatar, banner, profile_views, phone_number, address,
		social_links, is_verified, email_verified, followerscount, followscount,
		wallet_balance, refresh_token, refresh_expiry, refresh_ua, refresh_ip, refresh_prev
	FROM users
	WHERE username = $1 OR email = $1
	LIMIT 1`
	row := app.SQLDB.QueryRow(ctx, query, candidate)
	var user User
	if err := scanUserRow(row, &user); err != nil {
		return User{}, err
	}
	if user.UserID == "" {
		return User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func UpdateUserSession(ctx context.Context, app *infra.Deps, userID, refreshTokenHash, ua, ip string) (int64, error) {
	if err := ensureAuthDB(app); err != nil {
		return 0, err
	}
	query := `UPDATE users
		SET refresh_token = $2,
		    refresh_ua = $3,
		    refresh_ip = $4,
		    refresh_expiry = NOW() + INTERVAL '7 days',
		    updated_at = NOW()
		WHERE userid = $1`
	tag, err := app.SQLDB.Exec(ctx, query, userID, refreshTokenHash, ua, ip)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func LogoutUserByRefreshToken(ctx context.Context, app *infra.Deps, hashedToken string) (int64, error) {
	if err := ensureAuthDB(app); err != nil {
		return 0, err
	}
	query := `UPDATE users
		SET refresh_token = NULL,
		    refresh_prev = NULL,
		    refresh_ua = NULL,
		    refresh_ip = NULL,
		    refresh_expiry = NULL,
		    updated_at = NOW()
		WHERE refresh_token = $1`
	tag, err := app.SQLDB.Exec(ctx, query, hashedToken)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func LogoutAllUserSessions(ctx context.Context, app *infra.Deps, userID string) (int64, error) {
	if err := ensureAuthDB(app); err != nil {
		return 0, err
	}
	query := `UPDATE users
		SET refresh_token = NULL,
		    refresh_prev = NULL,
		    refresh_ua = NULL,
		    refresh_ip = NULL,
		    refresh_expiry = NULL,
		    updated_at = NOW()
		WHERE userid = $1`
	tag, err := app.SQLDB.Exec(ctx, query, userID)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func FindValidRefreshSession(ctx context.Context, app *infra.Deps, hashedToken string) (User, error) {
	if err := ensureAuthDB(app); err != nil {
		return User{}, err
	}

	query := `SELECT
		userid, username, email, password, password_hash, role, name, created_at, updated_at,
		bio, online, last_login, avatar, banner, profile_views, phone_number, address,
		social_links, is_verified, email_verified, followerscount, followscount,
		wallet_balance, refresh_token, refresh_expiry, refresh_ua, refresh_ip, refresh_prev
	FROM users
	WHERE refresh_token = $1 AND refresh_expiry > NOW()
	LIMIT 1`
	row := app.SQLDB.QueryRow(ctx, query, hashedToken)
	var user User
	if err := scanUserRow(row, &user); err != nil {
		return User{}, err
	}
	if user.UserID == "" {
		return User{}, fmt.Errorf("no valid refresh session")
	}
	return user, nil
}

// InvalidateUserSession clears all refresh token fields for a user.
func InvalidateUserSession(ctx context.Context, app *infra.Deps, userID string) (int64, error) {
	return LogoutAllUserSessions(ctx, app, userID)
}

func RotateRefreshTokenForUser(ctx context.Context, app *infra.Deps, userID, newRefreshHash, prevRefreshHash, ua string) (int64, error) {
	if err := ensureAuthDB(app); err != nil {
		return 0, err
	}
	query := `UPDATE users
		SET refresh_prev = $3,
		    refresh_token = $2,
		    refresh_ua = $4,
		    refresh_expiry = NOW() + INTERVAL '7 days',
		    updated_at = NOW()
		WHERE userid = $1`
	tag, err := app.SQLDB.Exec(ctx, query, userID, newRefreshHash, prevRefreshHash, ua)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func VerifyUserEmail(ctx context.Context, app *infra.Deps, email string) (int64, error) {
	if err := ensureAuthDB(app); err != nil {
		return 0, err
	}
	query := `UPDATE users SET email_verified = TRUE, updated_at = NOW() WHERE email = $1`
	tag, err := app.SQLDB.Exec(ctx, query, strings.TrimSpace(email))
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func MigrateUserPasswordHash(ctx context.Context, app *infra.Deps, userID, passwordHash string) (int64, error) {
	if err := ensureAuthDB(app); err != nil {
		return 0, err
	}
	query := `UPDATE users SET password_hash = $2, password = NULL, updated_at = NOW() WHERE userid = $1`
	tag, err := app.SQLDB.Exec(ctx, query, userID, passwordHash)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}
