// File: internal/profile/editProfile.go

package profile

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"scav/config/mqevents"
	"scav/infra"
	"scav/infra/cache"
	"scav/infra/mq"
	"scav/middleware"
	"scav/utils"

	log "scav/infra/logger"
)

/* -------------------------------------------------------
   Edit Profile
------------------------------------------------------- */

func EditProfile(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// 1. Validate JWT before trusting request context.
		claims, err := validateJWT(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 2. Parse form data (~10 MB)
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// 3. Invalidate cached profile
		_ = InvalidateCachedProfile(ctx, app.Cache, claims.Username)
		_ = UpdateCachedUsername(ctx, app.Cache, claims.UserID)

		// 4. Build updates map
		updates, err := BuildProfileUpdates(ctx, r, claims, app.Cache)
		if err != nil {
			http.Error(w, "Failed to update profile fields", http.StatusInternalServerError)
			return
		}

		// 5. Apply updates in DB
		if _, err := ApplyProfileUpdatesDeps(ctx, app, claims.UserID, updates); err != nil {
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}

		if app != nil && app.MQ != nil {
			evtCtx := mq.WithServiceName(ctx, "profile")
			if err := mq.PublishWithMeta(evtCtx, app.MQ, mqevents.SubjectUserProfileUpdated, mqevents.UserProfileUpdatedEvent{
				UserID:    claims.UserID,
				Username:  claims.Username,
				UpdatedAt: time.Now().UTC(),
			}); err != nil {
				log.L.Sugar().Warnw("profile: publish update event failed", "userid", claims.UserID, "error", err)
			}
		}

		// 6. Respond with updated profile
		RespondWithUserProfileDeps(w, claims.UserID, app)
	}
}

/* -------------------------------------------------------
   Delete Profile
------------------------------------------------------- */

func DeleteProfile(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		claims, err := validateJWT(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Invalidate cached profile
		_ = InvalidateCachedProfile(ctx, app.Cache, claims.Username)
		_ = UpdateCachedUsername(ctx, app.Cache, claims.UserID)

		// Delete user in DB
		if _, err := DeleteUserByIDDeps(ctx, app, claims.UserID); err != nil {
			http.Error(w, "Failed to delete profile", http.StatusInternalServerError)
			return
		}

		if app != nil && app.MQ != nil {
			evtCtx := mq.WithServiceName(ctx, "profile")
			if err := mq.PublishWithMeta(evtCtx, app.MQ, mqevents.SubjectUserProfileDeleted, mqevents.UserProfileDeletedEvent{
				UserID:    claims.UserID,
				Username:  claims.Username,
				DeletedAt: time.Now().UTC(),
			}); err != nil {
				log.L.Sugar().Warnw("profile: publish delete event failed", "userid", claims.UserID, "error", err)
			}
		}

		// Success response
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Profile deleted successfully",
		})
	}
}

/* -------------------------------------------------------
   Helper: Build profile updates
------------------------------------------------------- */

func BuildProfileUpdates(
	ctx context.Context,
	r *http.Request,
	claims *middleware.Claims,
	c cache.Cache,
) (map[string]any, error) {

	updates := map[string]any{}

	// Username
	if newUsername := r.FormValue("username"); newUsername != "" && newUsername != claims.Username {
		updates["username"] = newUsername
		_ = c.HSet(ctx, "users", claims.UserID, []byte(newUsername))
	}

	// Email, Bio, Name, Phone
	if val := r.FormValue("email"); val != "" {
		updates["email"] = val
	}
	if val := r.FormValue("bio"); val != "" {
		updates["bio"] = val
	}
	if val := r.FormValue("name"); val != "" {
		updates["name"] = val
	}
	if val := r.FormValue("phone_number"); val != "" {
		updates["phone_number"] = val
	}

	// Password (hashed)
	if val := r.FormValue("password"); val != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(val), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		// store in both fields for compatibility
		// store bcrypt hash in canonical `password_hash` field
		updates["password_hash"] = string(hashed)
	}

	return updates, nil
}
