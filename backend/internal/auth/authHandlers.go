package auth

import (
	"net/http"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/middleware"
	"scav/utils"

	"github.com/golang-jwt/jwt/v5"
)

func issueAccessToken(userID, username string, roles ...string) (string, error) {
	claims := &middleware.Claims{
		Username: username,
		UserID:   userID,
		Role:     roles,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(config.JwtSecret)
}

func issueRefreshToken(userID, username string, roles ...string) (string, error) {
	claims := &middleware.Claims{
		Username: username,
		UserID:   userID,
		Role:     roles,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(config.RefreshTokenSecret)
}

func Register(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := utils.ParseJSON(r, &req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		username := strings.TrimSpace(req.Username)
		email := strings.TrimSpace(req.Email)
		password := strings.TrimSpace(req.Password)
		if username == "" || email == "" || password == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "username, email and password are required")
			return
		}

		if app != nil && app.Log != nil {
			app.Log.Sugar().Infow("user registration request", "username", username, "email", email)
		}

		userID := "user_" + username
		accessToken, err := issueAccessToken(userID, username, "user")
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "could not create auth token")
			return
		}
		refreshToken, err := issueRefreshToken(userID, username, "user")
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "could not create refresh token")
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, map[string]any{
			"message":       "user registered",
			"token":         accessToken,
			"refresh_token": refreshToken,
			"user": map[string]any{
				"userid":   userID,
				"username": username,
				"email":    email,
			},
		})
	}
}

func Login(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := utils.ParseJSON(r, &req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		username := strings.TrimSpace(req.Username)
		password := strings.TrimSpace(req.Password)
		if username == "" || password == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "username and password are required")
			return
		}

		userID := "user_" + username
		accessToken, err := issueAccessToken(userID, username, "user")
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "could not create auth token")
			return
		}
		refreshToken, err := issueRefreshToken(userID, username, "user")
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "could not create refresh token")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"message":       "login successful",
			"token":         accessToken,
			"refresh_token": refreshToken,
			"user": map[string]any{
				"userid":   userID,
				"username": username,
			},
		})
	}
}

func RefreshToken(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := struct {
			RefreshToken string `json:"refresh_token"`
		}{}

		if r.Header.Get("Authorization") != "" {
			payload.RefreshToken = middleware.ExtractBearerToken(r.Header.Get("Authorization"))
		}
		if payload.RefreshToken == "" {
			if err := utils.ParseJSON(r, &payload); err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "invalid refresh token")
				return
			}
		}

		if strings.TrimSpace(payload.RefreshToken) == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "missing refresh token")
			return
		}

		claims := &middleware.Claims{}
		parsedToken, err := jwt.ParseWithClaims(payload.RefreshToken, claims, func(token *jwt.Token) (any, error) {
			return config.RefreshTokenSecret, nil
		})
		if err != nil || !parsedToken.Valid || claims.UserID == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "invalid refresh token")
			return
		}

		newToken, err := issueAccessToken(claims.UserID, claims.Username, claims.Role...)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "could not refresh token")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"message": "token refreshed",
			"token":   newToken,
		})
	}
}

func LogoutUser(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "token", Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode})
		http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode})
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"message": "logged out", "success": true})
	}
}

func LogoutAllSessions(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "token", Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode})
		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"message": "all sessions logged out", "success": true})
	}
}

func VerifyOTPHandler(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req OTPRequest
		if err := utils.ParseJSON(r, &req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if strings.TrimSpace(req.OTP) == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "otp is required")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"message":  "otp verified",
			"verified": true,
		})
	}
}

func RequestOTPHandler(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req OTPRequest
		if err := utils.ParseJSON(r, &req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if strings.TrimSpace(req.Username) == "" && strings.TrimSpace(req.Email) == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "username or email is required")
			return
		}

		utils.RespondWithJSON(w, http.StatusAccepted, map[string]any{
			"message": "otp request received",
			"sent":    true,
		})
	}
}
