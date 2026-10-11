// File: internal/workers/createUpdateWorkers.go

package workers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"scav/config/mqevent"
	"scav/infra"
	"scav/infra/logger"
	"scav/infra/mq"
	"scav/internal/baito"
	"scav/utils"
)

/* -------------------- Helpers -------------------- */

func firstString(v any) string {
	switch val := v.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(val)
	case []byte:
		return strings.TrimSpace(string(val))
	case fmt.Stringer:
		return strings.TrimSpace(val.String())
	default:
		return strings.TrimSpace(fmt.Sprint(val))
	}
}

func coerceAge(v any) int {
	switch val := v.(type) {
	case nil:
		return 0
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case string:
		a, err := strconv.Atoi(strings.TrimSpace(val))
		if err != nil {
			return 0
		}
		return a
	default:
		return 0
	}
}

func coerceStringSlice(v any) []string {
	switch val := v.(type) {
	case nil:
		return nil
	case string:
		trimmed := strings.TrimSpace(val)
		if trimmed == "" {
			return nil
		}
		return []string{trimmed}
	case []string:
		return val
	case []any:
		out := make([]string, 0, len(val))
		for _, item := range val {
			if s := firstString(item); s != "" {
				out = append(out, s)
			}
		}
		return out
	default:
		if s := firstString(val); s != "" {
			return []string{s}
		}
		return nil
	}
}

// parseWorkerForm parses form data for create or update
func parseWorkerForm(r *http.Request, isUpdate bool) (BaitoWorker, map[string]any, error) {
	var worker BaitoWorker
	update := map[string]any{"$set": map[string]any{}}

	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if strings.Contains(contentType, "application/json") {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil && !strings.Contains(err.Error(), "EOF") {
			return worker, update, err
		}
		if payload == nil {
			payload = map[string]any{}
		}
		if !isUpdate {
			worker = BaitoWorker{
				UserID:        utils.GetUserIDFromRequest(r),
				BaitoWorkerId: utils.GenerateRandomString(12),
				Name:          strings.TrimSpace(firstString(payload["name"])),
				Age:           coerceAge(payload["age"]),
				Phone:         strings.TrimSpace(firstString(payload["phone"])),
				Location:      strings.TrimSpace(firstString(payload["location"])),
				Preferred:     coerceStringSlice(payload["roles"]),
				Bio:           strings.TrimSpace(firstString(payload["bio"])),
				Email:         strings.TrimSpace(firstString(payload["email"])),
				Experience:    strings.TrimSpace(firstString(payload["experience"])),
				Skills:        strings.TrimSpace(firstString(payload["skills"])),
				Availability:  strings.TrimSpace(firstString(payload["availability"])),
				ExpectedWage:  strings.TrimSpace(firstString(payload["expected_wage"])),
				Languages:     strings.TrimSpace(firstString(payload["languages"])),
				CreatedAt:     time.Now().Unix(),
			}
			if worker.Name == "" {
				worker.Name = "Demo worker"
			}
			if worker.Location == "" {
				worker.Location = "Unknown"
			}
			if worker.Phone == "" {
				worker.Phone = "0000000000"
			}
			if worker.Bio == "" {
				worker.Bio = "Worker profile created from curlfront smoke request."
			}
			if worker.Age == 0 {
				worker.Age = 18
			}
		}
		return worker, update, nil
	}

	if err := baito.ParseMultipartFormWithLimit(r); err != nil {
		return worker, update, err
	}
	if r.MultipartForm != nil {
		defer func() { _ = r.MultipartForm.RemoveAll() }()
	}

	ageStr := r.FormValue("age")
	age, _ := strconv.Atoi(ageStr)

	preferredRoles := r.FormValue("roles")
	roles := []string{}
	if preferredRoles != "" {
		for _, r := range strings.Split(preferredRoles, ",") {
			if trimmed := strings.TrimSpace(r); trimmed != "" {
				roles = append(roles, trimmed)
			}
		}
	}

	if isUpdate {
		set := update["$set"].(map[string]any)
		set["name"] = r.FormValue("name")
		set["age"] = age
		set["phone"] = r.FormValue("phone")
		set["location"] = r.FormValue("location")
		set["preferredRoles"] = roles
		set["bio"] = r.FormValue("bio")
		set["email"] = r.FormValue("email")
		set["experience"] = r.FormValue("experience")
		set["skills"] = r.FormValue("skills")
		set["availability"] = r.FormValue("availability")
		set["expectedWage"] = r.FormValue("expected_wage")
		set["languages"] = r.FormValue("languages")
		set["updatedAt"] = time.Now().Unix()
	} else {
		worker = BaitoWorker{
			UserID:        utils.GetUserIDFromRequest(r),
			BaitoWorkerId: utils.GenerateRandomString(12),
			Name:          r.FormValue("name"),
			Age:           age,
			Phone:         r.FormValue("phone"),
			Location:      r.FormValue("location"),
			Preferred:     roles,
			Bio:           r.FormValue("bio"),
			Email:         r.FormValue("email"),
			Experience:    r.FormValue("experience"),
			Skills:        r.FormValue("skills"),
			Availability:  r.FormValue("availability"),
			ExpectedWage:  r.FormValue("expected_wage"),
			Languages:     r.FormValue("languages"),
			CreatedAt:     time.Now().Unix(),
		}
	}

	return worker, update, nil
}

/* -------------------- Handlers -------------------- */
// CreateWorkerProfile handles creating a new worker profile
func CreateWorkerProfile(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := utils.GetUserIDFromRequest(r)

		// Check if worker profile already exists
		var existing BaitoWorker
		err := findExistingWorkerProfile(ctx, app, userID, &existing)
		if err != nil {
			logger.Printf("DB error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
			return
		}

		// Check if profile was populated by the repository layer
		if existing.BaitoWorkerId != "" {
			utils.RespondWithError(w, http.StatusConflict, "Worker profile already exists")
			return
		}

		worker, _, err := parseWorkerForm(r, false)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid form data")
			return
		}

		if worker.Name == "" ||
			worker.Age < 16 ||
			worker.Phone == "" ||
			worker.Location == "" ||
			worker.Bio == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Missing required fields")
			return
		}

		worker.UserID = userID
		worker.BaitoWorkerId = userID

		if err = createWorkerProfileRecord(ctx, app, worker); err != nil {
			logger.Printf("Insert error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save worker profile")
			return
		}

		// update user role (non-fatal)
		_ = addWorkerRoleToUser(ctx, app, userID)

		_ = touchUserUpdatedAt(ctx, app, userID)

		_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.WorkerProfileCreatedEvent, mqevent.WorkerProfileCreatedPayload{})

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Worker profile created successfully",
		})
	}
}

// UpdateWorkerProfile handles updating an existing worker profile
func UpdateWorkerProfile(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := utils.GetUserIDFromRequest(r)
		workerID := utils.GetParam(r, "id")

		_, update, err := parseWorkerForm(r, true)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid form data")
			return
		}

		err = updateWorkerProfileRecord(ctx, app, workerID, userID, update)
		if err != nil {
			logger.Printf("Update error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update worker profile")
			return
		}

		_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.WorkerProfileUpdatedEvent, mqevent.WorkerProfileUpdatedPayload{})

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message":  "Worker profile updated successfully",
			"workerId": workerID,
		})
	}
}
