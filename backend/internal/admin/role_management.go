package admin

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/utils"
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
	UserID    string    `json:"userid" bson:"userid"`
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
		if err := FindPendingRoleApplication(ctx, app.DB, userID, payload.Role, &existing); err == nil {
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
		if err := InsertRoleApplication(ctx, app.DB, application); err != nil {
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
		if err := FindRoleApplicationsByUser(ctx, app.DB, userID, &applications); err != nil {
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
		filter := map[string]any{}
		if status := normalizeRoleRequestStatus(r.URL.Query().Get("status")); status != "" {
			filter["status"] = status
		}
		if err := ListRoleApplicationsDB(ctx, app.DB, filter, &applications); err != nil {
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
		if err := GetRoleApplicationByID(ctx, app.DB, appID, &application); err != nil {
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
		if err := GetUserRoles(ctx, app.DB, application.UserID, &user); err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		user.Role = MergeRoleList(user.Role, application.Role)
		if _, err := UpdateUserRoles(ctx, app.DB, application.UserID, user.Role); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update user role")
			return
		}

		if _, err := UpdateRoleApplicationStatus(ctx, app.DB, appID, "approved"); err != nil {
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
		if err := GetRoleApplicationByID(ctx, app.DB, appID, &application); err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Application not found")
			return
		}
		if isFinalRoleRequestStatus(application.Status) {
			utils.RespondWithError(w, http.StatusConflict, "This role request has already been resolved")
			return
		}

		if _, err := UpdateRoleApplicationStatus(ctx, app.DB, appID, "rejected"); err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Application not found or update failed")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Role request rejected"})
	}
}
