// File: internal/admin/moderator_applications.go

package admin

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"scav/config"
	"scav/config/mqevent"
	"scav/infra"
	"scav/infra/mq"
	"scav/utils"
)

var moderatorApplicationsCollection = config.Collections.ModeratorApplications

// Sentinel business errors
var (
	ErrMissingRequiredFields = errors.New("missing required fields")
	ErrAlreadyApplied        = errors.New("you have already applied to be a moderator")
	ErrModAppNotFound        = errors.New("application not found or update failed")
)

// Helper Types
type ModeratorApplication struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userid" db:"userid"`
	Reason    string    `json:"reason" db:"reason"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ApplyModeratorPayload struct {
	UserID string `json:"userid"`
	Reason string `json:"reason"`
}

// ============================================================================
// Business Logic / Services
// ============================================================================

func ProcessApplyModerator(ctx context.Context, app *infra.Deps, payload ApplyModeratorPayload) (*ModeratorApplication, error) {
	userID := strings.TrimSpace(payload.UserID)
	reason := strings.TrimSpace(payload.Reason)

	if userID == "" || reason == "" {
		return nil, ErrMissingRequiredFields
	}

	var existing ModeratorApplication
	if err := FindModeratorApplicationByUser(ctx, app, userID, &existing); err == nil {
		return nil, ErrAlreadyApplied
	}

	now := time.Now().UTC()
	appx := ModeratorApplication{
		ID:        "mod_" + utils.GenerateRandomString(16),
		UserID:    userID,
		Reason:    reason,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := InsertModeratorApplication(ctx, app, appx); err != nil {
		return nil, err
	}

	mqpayload, _ := json.Marshal(mqevent.AppliedForModeratorRolePayload{})
	_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.AppliedForModeratorRoleEvent, mqpayload)

	return &appx, nil
}

func FetchModeratorApplications(ctx context.Context, app *infra.Deps, status string) ([]ModeratorApplication, error) {
	filter := map[string]any{}
	if status != "" {
		filter["status"] = status
	}

	query, args := buildFilterQuery(filter)
	var applications []ModeratorApplication
	if err := ListModeratorApplicationsDB(ctx, app, query, args, &applications); err != nil {
		return nil, err
	}

	return applications, nil
}

func ProcessApproveModerator(ctx context.Context, app *infra.Deps, id string) error {
	now := time.Now().UTC()
	if _, err := UpdateModeratorApplicationStatus(ctx, app, id, "approved"); err != nil {
		return ErrModAppNotFound
	}

	mqpayload, _ := json.Marshal(mqevent.ApprovedModeratorRoleRequestPayload{
		ApplicationID: id,
		ApprovedAt:    now,
	})
	_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.ApprovedModeratorRoleRequestEvent, mqpayload)

	return nil
}

func ProcessRejectModerator(ctx context.Context, app *infra.Deps, id string) error {
	now := time.Now().UTC()
	if _, err := UpdateModeratorApplicationStatus(ctx, app, id, "rejected"); err != nil {
		return ErrModAppNotFound
	}

	mqpayload, _ := json.Marshal(mqevent.RejectedModeratorRoleRequestPayload{
		ApplicationID: id,
		RejectedAt:    now,
	})
	_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.RejectedModeratorRoleRequestEvent, mqpayload)

	return nil
}
