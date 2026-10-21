// File: internal/auth/logout.go

package auth

import (
	"context"
	"net/http"
	"time"

	"scav/config/mqevents"
	"scav/infra"
	"scav/infra/mq"
	"scav/middleware"
	"scav/utils"

	log "scav/infra/logger"
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

		// Safely extract the refresh token from cookie
		var tokenStr string
		if cookie, err := r.Cookie("refresh_token"); err == nil {
			tokenStr = cookie.Value
		}

		// Handoff logic & event publishing down to the service layer
		if tokenStr != "" {
			_ = ProcessSingleLogout(ctx, app, tokenStr)
		}

		clearRefreshCookie(w, r)

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"message": "Logged out",
			"data":    nil,
		})
	}
}

// ProcessSingleLogout wraps token transformation and calls downstream side effects
func ProcessSingleLogout(ctx context.Context, app *infra.Deps, rawRefreshToken string) error {
	hashedToken := hashRefreshToken(rawRefreshToken)

	return RevokeSessionAndEmit(ctx, app, hashedToken)
}

// RevokeSessionAndEmit deletes a single session via token hash and publishes the broker event
func RevokeSessionAndEmit(ctx context.Context, app *infra.Deps, hashedToken string) error {
	if _, err := LogoutUserByRefreshToken(ctx, app, hashedToken); err != nil {
		return err
	}

	if app != nil && app.MQ != nil {
		evtCtx := mq.WithServiceName(ctx, "auth")
		if err := mq.PublishWithMeta(evtCtx, app.MQ, mqevents.SubjectUserSessionRevoked, mqevents.UserSessionRevokedEvent{
			Reason:    "single_session_revoked",
			Timestamp: time.Now().UTC(),
		}); err != nil {
			log.L.Sugar().Warnw("auth: publish single logout event failed", "error", err)
		}
	}

	return nil
}

/* ============================================================
   LOGOUT ALL SESSIONS
============================================================ */

func LogoutAllSessions(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract and validate token
		claims, err := middleware.ValidateJWT(authHeader)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Handoff to service layer to clear DB sessions and fire MQ events
		if err := ProcessGlobalLogout(ctx, app, claims.UserID); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Logout failed")
			return
		}

		clearRefreshCookie(w, r)

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

	if app != nil && app.MQ != nil {
		evtCtx := mq.WithServiceName(ctx, "auth")
		if err := mq.PublishWithMeta(evtCtx, app.MQ, mqevents.SubjectUserLoggedOut, mqevents.UserLoggedOutEvent{
			UserID:    userID,
			Reason:    "all_sessions_revoked",
			Timestamp: time.Now().UTC(),
		}); err != nil {
			log.L.Sugar().Warnw("auth: publish logout event failed", "userid", userID, "error", err)
		}
	}

	return nil
}
