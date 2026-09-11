package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"scav/config/mqevent"
	"scav/infra"
	"scav/infra/mq"
	"scav/utils"
)

/* ============================================================
   LOGOUT
============================================================ */

func LogoutUser(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CSRF Check
		if r.Header.Get("X-Refresh-Intent") != "1" {
			utils.RespondWithError(w, http.StatusForbidden, "CSRF blocked")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Extract the session ID from cookie
		var sessionID string
		if cookie, err := r.Cookie(SessionCookieName); err == nil {
			sessionID = cookie.Value
		}

		// Service layer handles revoking the session and emitting events
		if sessionID != "" {
			_ = ProcessSingleLogout(ctx, app, sessionID)
		}

		clearSessionCookie(w, r)

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"message": "Logged out",
			"data":    nil,
		})
	}
}

// ProcessSingleLogout wraps session revocation and calls downstream side effects
func ProcessSingleLogout(ctx context.Context, app *infra.Deps, sessionID string) error {
	return RevokeSessionAndEmit(ctx, app, sessionID)
}

// RevokeSessionAndEmit deletes a single session via session ID and publishes the broker event
func RevokeSessionAndEmit(ctx context.Context, app *infra.Deps, sessionID string) error {
	if _, err := LogoutUserBySessionID(ctx, app, sessionID); err != nil {
		return err
	}

	_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.UserLoggedOut, mqevent.UserLoggedOutPayload{})
	return nil
}

/* ============================================================
   LOGOUT ALL SESSIONS
============================================================ */

func LogoutAllSessions(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Extract current session ID from cookie
		cookie, err := r.Cookie(SessionCookieName)
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Retrieve User ID associated with the current active session
		userID, err := GetUserIDBySessionID(ctx, app, cookie.Value)
		if err != nil || userID == "" {
			clearSessionCookie(w, r)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Service layer clears all sessions for the user globally
		if err := ProcessGlobalLogout(ctx, app, userID); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Logout failed")
			return
		}

		clearSessionCookie(w, r)

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"message": "All sessions revoked",
		})
	}
}

// ProcessGlobalLogout handles multi-device session revocation orchestration
func ProcessGlobalLogout(ctx context.Context, app *infra.Deps, userID string) error {
	return RevokeAllSessionsAndEmit(ctx, app, userID)
}

// RevokeAllSessionsAndEmit clears the user's sessions globally from storage and fires a system-wide event
func RevokeAllSessionsAndEmit(ctx context.Context, app *infra.Deps, userID string) error {
	if _, err := LogoutAllUserSessions(ctx, app, userID); err != nil {
		return err
	}

	_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.UserLoggedOutAllSessions, mqevent.UserLoggedOutPayload{})
	return nil
}

/* ============================================================
   HELPERS
============================================================ */

func clearSessionCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https",
		SameSite: http.SameSiteLaxMode,
	})
}

// GetUserIDBySessionID retrieves the UserID corresponding to an active session ID
func GetUserIDBySessionID(ctx context.Context, app *infra.Deps, sessionID string) (string, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	val, err := app.Cache.Get(ctx, key)
	if err != nil || len(val) == 0 {
		return "", fmt.Errorf("session invalid or expired")
	}
	return string(val), nil
}

// LogoutUserBySessionID deletes a single session by session ID from the cache
func LogoutUserBySessionID(ctx context.Context, app *infra.Deps, sessionID string) (any, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	err := app.Cache.Del(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke session: %w", err)
	}
	return true, nil
}

// // GetUserIDBySessionID queries the database for the user ID tied to an active session ID
// func GetUserIDBySessionID(ctx context.Context, app *infra.Deps, sessionID string) (string, error) {
// 	var userID string
// 	query := `SELECT user_id FROM user_sessions WHERE session_id = $1 AND expires_at > NOW()`

// 	err := app.DB.QueryRowContext(ctx, query, sessionID).Scan(&userID)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return "", fmt.Errorf("session not found or expired")
// 		}
// 		return "", fmt.Errorf("failed to fetch user session: %w", err)
// 	}

// 	return userID, nil
// }

// // LogoutUserBySessionID deletes a single session record from the database
// func LogoutUserBySessionID(ctx context.Context, app *infra.Deps, sessionID string) (any, error) {
// 	query := `DELETE FROM user_sessions WHERE session_id = $1`

// 	res, err := app.DB.ExecContext(ctx, query, sessionID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to delete session: %w", err)
// 	}

// 	rowsAffected, _ := res.RowsAffected()
// 	return rowsAffected, nil
// }
