// File: internal/baito/baitoSQLDB.go

package baito

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var UsersTable = config.Tables.UserTable
var BaitoTable = config.Tables.BaitoTable
var BaitoAppTable = config.Tables.BaitoApplicationsTable

func ensureBaitoDB(app *infra.Deps) error {
	if app == nil || app.SQLDB == nil {
		return errors.New("database not initialized")
	}
	return nil
}

func decodeBaitoMetadata(metadata []byte, dst any) error {
	if len(metadata) == 0 {
		return nil
	}
	if dst == nil {
		return errors.New("nil destination")
	}
	return json.Unmarshal(metadata, dst)
}

func deleteBaitoRecord(ctx context.Context, app *infra.Deps, baitoID, userID string) (int64, error) {
	if err := ensureBaitoDB(app); err != nil {
		return 0, err
	}
	if strings.TrimSpace(baitoID) == "" {
		return 0, nil
	}
	statement := `DELETE FROM baitos WHERE baitoid = $1`
	if strings.TrimSpace(userID) != "" {
		statement += ` AND userid = $2`
	}
	var tag pgconn.CommandTag
	var err error
	if strings.TrimSpace(userID) != "" {
		tag, err = app.SQLDB.Exec(ctx, statement, baitoID, userID)
	} else {
		tag, err = app.SQLDB.Exec(ctx, statement, baitoID)
	}
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func saveBaitoApplication(ctx context.Context, app *infra.Deps, application BaitoApplication) error {
	if err := ensureBaitoDB(app); err != nil {
		return err
	}
	if application.BaitoAppId == "" {
		application.BaitoAppId = application.UserID + "-" + application.BaitoID
	}
	payload, err := json.Marshal(application)
	if err != nil {
		return err
	}
	query := `INSERT INTO baitoapply (baitoappid, userid, baitoid, status, notes, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, 'pending', '', NOW(), NOW(), $4)
		ON CONFLICT (baitoappid) DO UPDATE SET userid = EXCLUDED.userid, baitoid = EXCLUDED.baitoid, updated_at = NOW(), metadata = EXCLUDED.metadata`
	_, err = app.SQLDB.Exec(ctx, query, application.BaitoAppId, application.UserID, application.BaitoID, payload)
	return err
}

func incrementBaitoApplicationCount(ctx context.Context, app *infra.Deps, baitoID string) error {
	if err := ensureBaitoDB(app); err != nil {
		return err
	}
	if strings.TrimSpace(baitoID) == "" {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `UPDATE baitos SET metadata = COALESCE(metadata, '{}'::jsonb) || jsonb_build_object('applicationcount', COALESCE((metadata->>'applicationcount')::int, 0) + 1), updated_at = NOW() WHERE baitoid = $1`, baitoID)
	return err
}

func buildMyApplicationsResult(applications []map[string]any, jobs []Baito) []map[string]any {
	jobByID := make(map[string]Baito, len(jobs))
	for _, job := range jobs {
		jobByID[job.BaitoId] = job
	}

	results := make([]map[string]any, 0, len(applications))
	for _, application := range applications {
		result := map[string]any{}
		for k, v := range application {
			result[k] = v
		}

		jobID, _ := application["baitoid"].(string)
		if job, ok := jobByID[jobID]; ok {
			result["jobId"] = job.BaitoId
			result["title"] = job.Title
			result["location"] = job.Location
			result["wage"] = job.Wage
		}

		if _, ok := result["id"]; !ok {
			if id, ok := application["_id"]; ok {
				result["id"] = id
			}
		}
		if _, ok := result["jobId"]; !ok {
			result["jobId"] = jobID
		}
		results = append(results, result)
	}
	return results
}

func createBaitoRecord(ctx context.Context, app *infra.Deps, baito Baito) error {
	if err := ensureBaitoDB(app); err != nil {
		return err
	}
	if baito.BaitoId == "" {
		baito.BaitoId = strings.TrimSpace(baito.Title)
	}
	payload, err := json.Marshal(baito)
	if err != nil {
		return err
	}
	_, err = app.SQLDB.Exec(ctx, `INSERT INTO baitos (baitoid, userid, status, created_at, updated_at, metadata) VALUES ($1, $2, 'active', NOW(), NOW(), $3) ON CONFLICT (baitoid) DO UPDATE SET userid = EXCLUDED.userid, updated_at = NOW(), metadata = EXCLUDED.metadata`, baito.BaitoId, baito.OwnerID, payload)
	return err
}

func updateBaitoRecord(ctx context.Context, app *infra.Deps, baitoID, userID string, update map[string]any) (int64, error) {
	if err := ensureBaitoDB(app); err != nil {
		return 0, err
	}
	if strings.TrimSpace(baitoID) == "" {
		return 0, nil
	}
	var payload []byte
	if len(update) > 0 {
		base := map[string]any{}
		query := `SELECT metadata FROM baitos WHERE baitoid = $1 AND userid = $2 LIMIT 1`
		if err := app.SQLDB.QueryRow(ctx, query, baitoID, userID).Scan(&payload); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return 0, err
			}
		}
		if len(payload) > 0 {
			if err := json.Unmarshal(payload, &base); err != nil {
				return 0, err
			}
		}
		for key, value := range update {
			if key == "$set" {
				if nested, ok := value.(map[string]any); ok {
					for k, v := range nested {
						base[k] = v
					}
				}
				continue
			}
			base[key] = value
		}
		payload, err := json.Marshal(base)
		if err != nil {
			return 0, err
		}
		result, err := app.SQLDB.Exec(ctx, `UPDATE baitos SET metadata = $1::jsonb, updated_at = NOW() WHERE baitoid = $2 AND userid = $3`, payload, baitoID, userID)
		if err != nil {
			return 0, err
		}
		return result.RowsAffected(), nil
	}
	return 0, nil
}

func findLatestBaitosFromDB(ctx context.Context, app *infra.Deps, query string, args []any, limit int) ([]BaitosResponse, error) {
	if err := ensureBaitoDB(app); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if query == "" {
		query = "1 = 1"
	}
	statement := fmt.Sprintf(`SELECT baitoid, userid, metadata, created_at, updated_at FROM baitos WHERE %s ORDER BY created_at DESC LIMIT $%d`, query, len(args)+1)
	queryArgs := append(append([]any{}, args...), limit)
	rows, err := app.SQLDB.Query(ctx, statement, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]BaitosResponse, 0, limit)
	for rows.Next() {
		var baitoID, userID string
		var createdAt, updatedAt time.Time
		var metadata []byte
		if err := rows.Scan(&baitoID, &userID, &metadata, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		var baito Baito
		if err := decodeBaitoMetadata(metadata, &baito); err != nil {
			return nil, err
		}
		baito.BaitoId = baitoID
		baito.OwnerID = userID
		baito.CreatedAt = createdAt
		baito.UpdatedAt = updatedAt
		items = append(items, BaitosResponse{
			BaitoId:         baito.BaitoId,
			Title:           baito.Title,
			Description:     baito.Description,
			Category:        baito.Category,
			SubCategory:     baito.SubCategory,
			Location:        baito.Location,
			Wage:            baito.Wage,
			Requirements:    baito.Requirements,
			BannerURL:       baito.Banner,
			WorkHours:       baito.WorkHours,
			Duration:        baito.Duration,
			LastDateToApply: baito.LastDateToApply,
			CreatedAt:       baito.CreatedAt,
			OwnerID:         baito.OwnerID,
		})
	}
	return items, rows.Err()
}

func findRelatedBaitosFromDB(ctx context.Context, app *infra.Deps, query string, args []any, limit int) ([]BaitosResponse, error) {
	return findLatestBaitosFromDB(ctx, app, query, args, limit)
}

func findBaitoByIDFromDB(ctx context.Context, app *infra.Deps, baitoID string) (Baito, error) {
	var baito Baito
	if err := ensureBaitoDB(app); err != nil {
		return baito, err
	}
	var userID string
	var createdAt, updatedAt time.Time
	var metadata []byte
	if err := app.SQLDB.QueryRow(ctx, `SELECT baitoid, userid, created_at, updated_at, metadata FROM baitos WHERE baitoid = $1 LIMIT 1`, baitoID).Scan(&baito.BaitoId, &userID, &createdAt, &updatedAt, &metadata); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return baito, nil
		}
		return baito, err
	}
	baito.OwnerID = userID
	baito.CreatedAt = createdAt
	baito.UpdatedAt = updatedAt
	if err := decodeBaitoMetadata(metadata, &baito); err != nil {
		return baito, err
	}
	if baito.BaitoId == "" {
		baito.BaitoId = baitoID
	}
	return baito, nil
}

func findMyBaitosFromDB(ctx context.Context, app *infra.Deps, userID string) ([]BaitosResponse, error) {
	return findLatestBaitosFromDB(ctx, app, "userid = $1", []any{userID}, 50)
}

func countApplicationsForBaito(ctx context.Context, app *infra.Deps, baitoID string) (int64, error) {
	if err := ensureBaitoDB(app); err != nil {
		return 0, err
	}
	var count int64
	if err := app.SQLDB.QueryRow(ctx, `SELECT COUNT(*) FROM baitoapply WHERE baitoid = $1`, baitoID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func findBaitoApplicantsFromDB(ctx context.Context, app *infra.Deps, baitoID string) ([]map[string]any, error) {
	if err := ensureBaitoDB(app); err != nil {
		return nil, err
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT baitoappid, userid, baitoid, status, metadata FROM baitoapply WHERE baitoid = $1 ORDER BY created_at DESC`, baitoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]any, 0)
	for rows.Next() {
		var appID, userID, jobID, status string
		var metadata []byte
		if err := rows.Scan(&appID, &userID, &jobID, &status, &metadata); err != nil {
			return nil, err
		}
		entry := map[string]any{"baitoappid": appID, "userid": userID, "baitoid": jobID, "status": status}
		if len(metadata) > 0 {
			var payload map[string]any
			if err := json.Unmarshal(metadata, &payload); err == nil {
				for k, v := range payload {
					entry[k] = v
				}
			}
		}
		results = append(results, entry)
	}
	return results, rows.Err()
}

func findMyApplicationsFromDB(ctx context.Context, app *infra.Deps, userID string) ([]map[string]any, error) {
	var applications []map[string]any
	if err := ensureBaitoDB(app); err != nil {
		return nil, err
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT baitoappid, userid, baitoid, status, metadata FROM baitoapply WHERE userid = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var appID, rowUserID, jobID, status string
		var metadata []byte
		if err := rows.Scan(&appID, &rowUserID, &jobID, &status, &metadata); err != nil {
			return nil, err
		}
		entry := map[string]any{"baitoappid": appID, "userid": rowUserID, "baitoid": jobID, "status": status}
		if len(metadata) > 0 {
			var payload map[string]any
			if err := json.Unmarshal(metadata, &payload); err == nil {
				for k, v := range payload {
					entry[k] = v
				}
			}
		}
		applications = append(applications, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(applications) == 0 {
		return []map[string]any{}, nil
	}
	jobIDs := make([]string, 0, len(applications))
	seen := make(map[string]struct{}, len(applications))
	for _, application := range applications {
		jobID, _ := application["baitoid"].(string)
		if jobID == "" {
			continue
		}
		if _, ok := seen[jobID]; ok {
			continue
		}
		seen[jobID] = struct{}{}
		jobIDs = append(jobIDs, jobID)
	}
	jobs := make([]Baito, 0, len(jobIDs))
	for _, id := range jobIDs {
		job, err := findBaitoByIDFromDB(ctx, app, id)
		if err != nil {
			continue
		}
		jobs = append(jobs, job)
	}
	return buildMyApplicationsResult(applications, jobs), nil
}
