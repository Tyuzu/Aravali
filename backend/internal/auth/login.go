package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"scav/config/mqevent"
	"scav/infra"
	"scav/infra/mq"
	"scav/utils"

	log "scav/utils/logger"

	"golang.org/x/crypto/bcrypt"
)

const (
	SessionCookieName = "session_id"
	SessionTTL        = 24 * time.Hour
)

/* ============================================================
   1. HANDLERS (HTTP LAYER)
============================================================ */

func Login(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var creds LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
			return
		}
		creds.Username = strings.TrimSpace(creds.Username)
		creds.Password = strings.TrimSpace(creds.Password)

		ip := clientIP(r)
		failKey := fmt.Sprintf("auth:fail:%s:%s", creds.Username, ipPrefix(ip))
		uaHashStr := uaHash(r)
		ipPrefixStr := ipPrefix(ip)

		if isLocked := CheckRateLimitLockout(ctx, app, failKey); isLocked {
			utils.RespondWithError(w, http.StatusTooManyRequests, "Too many attempts")
			return
		}

		sessionID, userID, err := AuthenticateAndCreateSession(ctx, app, creds, uaHashStr, ipPrefixStr)
		if err != nil {
			IncrementRateLimitCounter(ctx, app, failKey)

			if errors.Is(err, ErrAuthInvalidCredentials) {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = ClearRateLimitCounter(ctx, app, failKey)

		setSessionCookie(w, r, sessionID)
		_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.UserLoggedIn, mqevent.UserLoggedInPayload{})

		utils.RespondWithJSON(w, http.StatusOK, LoginResponse{
			Message: "Login successful",
			Status:  http.StatusOK,
			UserID:  userID,
		})
	}
}

/* ============================================================
   2. SERVICES (BUSINESS LAYER)
============================================================ */

func AuthenticateAndCreateSession(ctx context.Context, app *infra.Deps, creds LoginRequest, uaHash string, ipPrefix string) (string, string, error) {
	user, err := GetUserByUsername(ctx, app, creds.Username)
	if err != nil {
		log.Printf("auth: user lookup failed for username=%s: %v", creds.Username, err)
		return "", "", ErrAuthInvalidCredentials
	}

	creds.Password = strings.TrimSpace(creds.Password)

	isBcrypt := func(s string) bool {
		return strings.HasPrefix(s, "$2a$") || strings.HasPrefix(s, "$2b$") || strings.HasPrefix(s, "$2y$")
	}

	log.Printf("auth: user=%s userid=%s password_hash_present=%v password_present=%v", creds.Username, user.UserID, len(user.PasswordHash) > 0, len(user.Password) > 0)

	var matched bool

	// Prefer canonical `PasswordHash` field
	if len(user.PasswordHash) > 0 && isBcrypt(user.PasswordHash) {
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)) == nil {
			matched = true
		} else {
			log.Printf("auth: password mismatch for username=%s userid=%s using field=password_hash", creds.Username, user.UserID)
		}
	}

	// Fallback to legacy `Password` field
	if !matched && len(user.Password) > 0 && isBcrypt(user.Password) {
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) == nil {
			matched = true
			if _, err := MigrateUserPasswordHash(ctx, app, user.UserID, user.Password); err != nil {
				log.Printf("auth: failed to migrate password hash for userid=%s: %v", user.UserID, err)
			}
		}
	}

	if !matched {
		return "", "", ErrAuthInvalidCredentials
	}

	// Generate secure, cryptographically random session ID
	sessionID, err := generateSessionID()
	if err != nil {
		return "", "", ErrTokenGeneration
	}

	// Persist the active session ID into the database/cache
	_, err = PersistUserSession(
		ctx,
		app,
		user.UserID,
		sessionID,
		uaHash,
		ipPrefix,
	)
	if err != nil {
		return "", "", ErrSessionPersistence
	}

	return sessionID, user.UserID, nil
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func setSessionCookie(w http.ResponseWriter, r *http.Request, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(SessionTTL),
		MaxAge:   int(SessionTTL.Seconds()),
		HttpOnly: true,                                                         // Prevents XSS script access
		Secure:   r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https", // Transmit over HTTPS only
		SameSite: http.SameSiteLaxMode,                                         // Protection against CSRF
	})
}

/* ============================================================
   3. REPOSITORIES (DATA ACCESS / CACHE LAYER)
============================================================ */

func CheckRateLimitLockout(ctx context.Context, app *infra.Deps, failKey string) bool {
	val, err := app.Cache.Get(ctx, failKey)
	var cnt int64
	if err == nil && len(val) > 0 {
		cnt, _ = strconv.ParseInt(string(val), 10, 64)
	}
	return cnt >= maxFailedAttempts
}

func IncrementRateLimitCounter(ctx context.Context, app *infra.Deps, failKey string) {
	cnt, err := app.Cache.Incr(ctx, failKey)
	if err != nil {
		cnt = 0
	}
	_ = app.Cache.Set(ctx, failKey, []byte(strconv.FormatInt(cnt, 10)), lockoutDuration)
}

func ClearRateLimitCounter(ctx context.Context, app *infra.Deps, failKey string) error {
	return app.Cache.Del(ctx, failKey)
}

func GetUserByUsername(ctx context.Context, app *infra.Deps, username string) (User, error) {
	return FindUserByUsername(ctx, app, username)
}

func PersistUserSession(ctx context.Context, app *infra.Deps, userID, sessionID, uaHash, ipPrefix string) (any, error) {
	return UpdateUserSession(ctx, app, userID, sessionID, uaHash, ipPrefix)
}
