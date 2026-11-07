// File: internal/auth/register.go

package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"scav/config/mqevent"
	"scav/infra"
	log "scav/infra/logger"
	"scav/infra/mq"
	"scav/utils"

	"golang.org/x/crypto/bcrypt"
)

/* ============================================================
   1. HANDLERS (HTTP LAYER)
============================================================ */

func Register(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var input SignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
			return
		}

		user, err := ProcessRegistration(ctx, app, input)
		if err != nil {
			if errors.Is(err, ErrAuthInvalidCredentials) {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid credentials")
				return
			}
			if errors.Is(err, ErrUserAlreadyExists) {
				utils.RespondWithError(w, http.StatusConflict, "User already exists")
				return
			}

			// Log unexpected registration errors with context for debugging. This is
			// especially important for Postgres connectivity issues (missing DB,
			// bad credentials, or failed migrations).
			log.L.Sugar().Errorw(
				"register_error",
				"remote", r.RemoteAddr,
				"method", r.Method,
				"uri", r.RequestURI,
				"error", err,
				"username", input.Username,
				"email", input.Email,
				"db_url", app.Config.DatabaseURL,
				"db_user", app.Config.PostgresUser,
				"db_host", app.Config.PostgresHost,
				"db_port", app.Config.PostgresPort,
				"db_name", app.Config.PostgresDB,
				"env", app.Config.Env,
			)

			status := http.StatusInternalServerError
			message := "Registration failed"
			if err != nil && (strings.Contains(err.Error(), "database not initialized") ||
				strings.Contains(err.Error(), "connection refused") ||
				strings.Contains(err.Error(), "password authentication failed") ||
				strings.Contains(err.Error(), "does not exist") ||
				strings.Contains(err.Error(), "could not connect") ||
				strings.Contains(err.Error(), "dial tcp") ||
				strings.Contains(err.Error(), "pq: ")) {
				status = http.StatusServiceUnavailable
				message = "Registration failed: database unavailable"
			}

			utils.RespondWithError(w, status, message)
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, SignUpResponse{
			Message: "User registered successfully",
			UserID:  user.UserID,
		})
	}
}

/* ============================================================
   2. SERVICES (BUSINESS LAYER)
============================================================ */

func ProcessRegistration(ctx context.Context, app *infra.Deps, input SignUpRequest) (User, error) {
	input.Username = strings.TrimSpace(input.Username)
	input.Password = strings.TrimSpace(input.Password)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	if !validateUsername(input.Username) ||
		!validateEmail(input.Email) ||
		!validatePassword(input.Password) {
		return User{}, ErrAuthInvalidCredentials
	}

	user, err := BuildUser(input)
	if err != nil {
		return User{}, ErrPasswordHashing
	}

	if err := QQQCreateUser(ctx, app, user); err != nil {
		return User{}, err
	}

	_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.UserRegistered, mqevent.UserRegisteredPayload{
		UserID:     user.UserID,
		Username:   user.Username,
		Email:      user.Email,
		OccurredAt: time.Now().UTC(),
	})

	return user, nil
}

func BuildUser(input SignUpRequest) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return User{}, err
	}

	now := time.Now()
	user := User{
		UserID:        "u" + utils.GenerateRandomString(10),
		Username:      input.Username,
		Email:         input.Email,
		Password:      "",
		PasswordHash:  string(hashedPassword),
		Role:          []string{"user"},
		CreatedAt:     now,
		UpdatedAt:     now,
		EmailVerified: false,
		IsVerified:    false,
		Online:        false,
		SocialLinks:   map[string]string{},
	}

	return user, nil
}
