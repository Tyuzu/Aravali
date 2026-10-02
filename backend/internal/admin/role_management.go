// File: internal/admin/role_management.go

package admin

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/utils"
)

var roleApplicationsTable = config.Tables.RoleApplicationsTable

var usersTable = config.Tables.UserTable

// Sentinel business errors
var (
	ErrInvalidPayload       = errors.New("role and reason are required")
	ErrUserRoleDefault      = errors.New("user role is already assigned to every account")
	ErrPendingRequestExists = errors.New("you already submitted a pending request for this role")
	ErrApplicationNotFound  = errors.New("application not found")
	ErrAlreadyResolved      = errors.New("this role request has already been resolved")
	ErrUserNotFound         = errors.New("user not found")
)

// Helper Types
type RoleApplication struct {
	ID        string    `json:"id" db:"roleapplicationid"`
	UserID    string    `json:"userid" db:"userid"`
	Role      string    `json:"role" db:"role"`
	Reason    string    `json:"reason" db:"-"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Metadata  []byte    `json:"metadata" db:"metadata"`
}

type ApplyPayload struct {
	Role   string `json:"role"`
	Reason string `json:"reason"`
}

// Pure Utility Functions

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

// ============================================================================
// Business Logic / Services
// ============================================================================

func ProcessApplyForRole(ctx context.Context, app *infra.Deps, userID string, payload ApplyPayload) (*RoleApplication, error) {
	role := NormalizeRoleName(payload.Role)
	reason := strings.TrimSpace(payload.Reason)

	if role == "" || reason == "" {
		return nil, ErrInvalidPayload
	}

	if role == "user" {
		return nil, ErrUserRoleDefault
	}

	var existing RoleApplication
	if err := FindPendingRoleApplication(ctx, app, userID, role, &existing); err == nil {
		return nil, ErrPendingRequestExists
	}

	application := RoleApplication{
		ID:        "role_" + utils.GenerateRandomString(16),
		UserID:    userID,
		Role:      role,
		Reason:    reason,
		Status:    "pending",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	meta := map[string]any{"reason": reason}
	if mb, err := json.Marshal(meta); err == nil {
		application.Metadata = mb
	}

	if err := InsertRoleApplication(ctx, app, application); err != nil {
		return nil, err
	}

	return &application, nil
}

func FetchUserRoleRequests(ctx context.Context, app *infra.Deps, userID string) ([]RoleApplication, error) {
	var applications []RoleApplication
	if err := FindRoleApplicationsByUser(ctx, app, userID, &applications); err != nil {
		return nil, err
	}
	return applications, nil
}

func FetchAllRoleRequests(ctx context.Context, app *infra.Deps, rawStatus string) ([]RoleApplication, error) {
	var applications []RoleApplication
	filter := map[string]any{}

	if status := normalizeRoleRequestStatus(rawStatus); status != "" {
		filter["status"] = status
	}

	query, args := buildFilterQuery(filter)
	if err := ListRoleApplicationsDB(ctx, app, query, args, &applications); err != nil {
		return nil, err
	}
	return applications, nil
}

func ProcessApproveRoleRequest(ctx context.Context, app *infra.Deps, appID string) (string, error) {
	var application RoleApplication
	if err := GetRoleApplicationByID(ctx, app, appID, &application); err != nil {
		return "", ErrApplicationNotFound
	}

	if isFinalRoleRequestStatus(application.Status) {
		return "", ErrAlreadyResolved
	}

	var user struct {
		Role []string `json:"role" db:"role"`
	}
	if err := GetUserRoles(ctx, app, application.UserID, &user); err != nil {
		return "", ErrUserNotFound
	}

	user.Role = MergeRoleList(user.Role, application.Role)
	if _, err := UpdateUserRoles(ctx, app, application.UserID, user.Role); err != nil {
		return "", err
	}

	if _, err := UpdateRoleApplicationStatus(ctx, app, appID, "approved"); err != nil {
		return "", err
	}

	return application.Role, nil
}

func ProcessRejectRoleRequest(ctx context.Context, app *infra.Deps, appID string) error {
	var application RoleApplication
	if err := GetRoleApplicationByID(ctx, app, appID, &application); err != nil {
		return ErrApplicationNotFound
	}

	if isFinalRoleRequestStatus(application.Status) {
		return ErrAlreadyResolved
	}

	if _, err := UpdateRoleApplicationStatus(ctx, app, appID, "rejected"); err != nil {
		return err
	}

	return nil
}
