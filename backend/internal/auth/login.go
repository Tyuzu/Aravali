package auth

import (
	"context"
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
	"scav/middleware"
	"scav/utils"

	log "scav/utils/logger"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

		accessToken, refreshToken, userID, err := AuthenticateAndCreateSession(ctx, app, creds, uaHashStr, ipPrefixStr)
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

		setRefreshCookie(w, refreshToken)
		_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.UserLoggedIn, mqevent.UserLoggedInPayload{})

		utils.RespondWithJSON(w, http.StatusOK, LoginResponse{
			Message: "Login successful",
			Status:  http.StatusOK,
			Token:   accessToken,
			UserID:  userID,
		})
	}
}

/* ============================================================
   2. SERVICES (BUSINESS LAYER)
============================================================ */

func AuthenticateAndCreateSession(ctx context.Context, app *infra.Deps, creds LoginRequest, uaHash string, ipPrefix string) (string, string, string, error) {
	user, err := GetUserByUsername(ctx, app, creds.Username)
	if err != nil {
		log.Printf("auth: user lookup failed for username=%s: %v", creds.Username, err)
		return "", "", "", ErrAuthInvalidCredentials
	}
	// Ensure incoming password is trimmed
	creds.Password = strings.TrimSpace(creds.Password)

	isBcrypt := func(s string) bool {
		return strings.HasPrefix(s, "$2a$") || strings.HasPrefix(s, "$2b$") || strings.HasPrefix(s, "$2y$")
	}

	// Sanity info (do not log full hashes)
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

	// Fallback to legacy `Password` field (migrate on success)
	if !matched && len(user.Password) > 0 && isBcrypt(user.Password) {
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) == nil {
			matched = true
			// migrate existing hash into canonical field
			_, _ = app.DB.Update(ctx, UsersCollection, map[string]string{"userid": user.UserID}, map[string]any{"$set": map[string]any{"password_hash": user.Password}})
			log.Printf("auth: migrated password -> password_hash for userid=%s", user.UserID)
		} else {
			log.Printf("auth: password mismatch for username=%s userid=%s using field=password", creds.Username, user.UserID)
		}
	}

	if !matched {
		return "", "", "", ErrAuthInvalidCredentials
	}

	claims := &middleware.Claims{
		UserID:   user.UserID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := createAccessToken(claims)
	if err != nil {
		return "", "", "", ErrTokenGeneration
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return "", "", "", ErrTokenGeneration
	}

	_, err = PersistUserSession(
		ctx,
		app,
		user.UserID,
		hashRefreshToken(refreshToken),
		uaHash,
		ipPrefix,
	)
	if err != nil {
		return "", "", "", ErrSessionPersistence
	}

	return accessToken, refreshToken, user.UserID, nil
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

func PersistUserSession(ctx context.Context, app *infra.Deps, userID, hashedRefreshToken, uaHash, ipPrefix string) (any, error) {
	return UpdateUserSession(ctx, app, userID, hashedRefreshToken, uaHash, ipPrefix)
}
