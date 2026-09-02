package admin

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/utils"

	"go.mongodb.org/mongo-driver/bson"
)

const roleApplicationsCollection = "role_applications"

var usersCollection = config.Collections.UserCollection

func NormalizeRoleName(role string) string {
	return strings.ToLower(strings.TrimSpace(role))
}

func normalizeRoleRequestStatus(status string) string {
	switch NormalizeRoleName(status) {
	case "pending", "approved", "rejected":
		return NormalizeRoleName(status)
	default:
		return ""
	}
}

func isFinalRoleRequestStatus(status string) bool {
	normalized := normalizeRoleRequestStatus(status)
	return normalized == "approved" || normalized == "rejected"
}

func MergeRoleList(existing []string, roles ...string) []string {
	seen := map[string]struct{}{}
	merged := make([]string, 0, len(existing)+len(roles))
	for _, role := range existing {
		normalized := NormalizeRoleName(role)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; !ok {
			seen[normalized] = struct{}{}
			merged = append(merged, normalized)
		}
	}
	for _, role := range roles {
		normalized := NormalizeRoleName(role)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; !ok {
			seen[normalized] = struct{}{}
			merged = append(merged, normalized)
		}
	}
	return merged
}

type RoleApplication struct {
	ID        string    `json:"id" bson:"id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Role      string    `json:"role" bson:"role"`
	Reason    string    `json:"reason" bson:"reason"`
	Status    string    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func ApplyForRole(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		var payload struct {
			Role   string `json:"role"`
			Reason string `json:"reason"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		payload.Role = NormalizeRoleName(payload.Role)
		payload.Reason = strings.TrimSpace(payload.Reason)
		if payload.Role == "" || payload.Reason == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Role and reason are required")
			return
		}

		if payload.Role == "user" {
			utils.RespondWithError(w, http.StatusBadRequest, "User role is already assigned to every account")
			return
		}

		var existing RoleApplication
		if err := app.DB.FindOne(ctx, roleApplicationsCollection, bson.M{"user_id": userID, "role": payload.Role, "status": "pending"}, &existing); err == nil {
			utils.RespondWithError(w, http.StatusConflict, "You already submitted a pending request for this role")
			return
		}

		application := RoleApplication{
			ID:        "role_" + utils.GenerateRandomString(16),
			UserID:    userID,
			Role:      payload.Role,
			Reason:    payload.Reason,
			Status:    "pending",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		if err := app.DB.Insert(ctx, roleApplicationsCollection, application); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save role request")
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, map[string]any{
			"message": "Role request submitted",
			"id":      application.ID,
			"status":  application.Status,
		})
	}
}

func GetMyRoleRequests(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		var applications []RoleApplication
		if err := app.DB.FindMany(ctx, roleApplicationsCollection, bson.M{"user_id": userID}, &applications); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to load your role requests")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, applications)
	}
}

func ListRoleRequests(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var applications []RoleApplication
		filter := bson.M{}
		if status := normalizeRoleRequestStatus(r.URL.Query().Get("status")); status != "" {
			filter["status"] = status
		}
		if err := app.DB.FindMany(ctx, roleApplicationsCollection, filter, &applications); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to load role applications")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, applications)
	}
}

func ApproveRoleRequest(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		appID := utils.GetParam(r, "id")
		if appID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Missing application id")
			return
		}

		var application RoleApplication
		if err := app.DB.FindOne(ctx, roleApplicationsCollection, bson.M{"id": appID}, &application); err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Application not found")
			return
		}
		if isFinalRoleRequestStatus(application.Status) {
			utils.RespondWithError(w, http.StatusConflict, "This role request has already been resolved")
			return
		}

		var user struct {
			Role []string `json:"role" bson:"role"`
		}
		if err := app.DB.FindOne(ctx, usersCollection, bson.M{"userid": application.UserID}, &user); err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		user.Role = MergeRoleList(user.Role, application.Role)
		if _, err := app.DB.UpdateOne(ctx, usersCollection, bson.M{"userid": application.UserID}, bson.M{"$set": bson.M{"role": user.Role, "updated_at": time.Now().UTC()}}); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update user role")
			return
		}

		if _, err := app.DB.UpdateOne(ctx, roleApplicationsCollection, bson.M{"id": appID}, bson.M{"$set": bson.M{"status": "approved", "updated_at": time.Now().UTC()}}); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update role application")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"message": "Role approved", "role": application.Role})
	}
}

func RejectRoleRequest(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		appID := utils.GetParam(r, "id")
		if appID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Missing application id")
			return
		}

		var application RoleApplication
		if err := app.DB.FindOne(ctx, roleApplicationsCollection, bson.M{"id": appID}, &application); err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Application not found")
			return
		}
		if isFinalRoleRequestStatus(application.Status) {
			utils.RespondWithError(w, http.StatusConflict, "This role request has already been resolved")
			return
		}

		if _, err := app.DB.UpdateOne(ctx, roleApplicationsCollection, bson.M{"id": appID}, bson.M{"$set": bson.M{"status": "rejected", "updated_at": time.Now().UTC()}}); err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Application not found or update failed")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Role request rejected"})
	}
}
