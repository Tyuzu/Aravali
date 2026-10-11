// File: internal/workers/workerSQLDB.go

package workers

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
)

var BaitoWorkersTable = config.Tables.BaitoWorkerTable
var UsersTable = config.Tables.UserTable

func ensureWorkerDB(app *infra.Deps) error {
	if app == nil || app.SQLDB == nil {
		return errors.New("database not initialized")
	}
	return nil
}

func decodeWorkerMetadata(metadata []byte, dst any) error {
	if len(metadata) == 0 {
		return nil
	}
	if dst == nil {
		return errors.New("nil destination")
	}
	return json.Unmarshal(metadata, dst)
}

func metadataToWorker(metadata []byte, existing *BaitoWorker) error {
	if existing == nil {
		return errors.New("nil worker")
	}
	if len(metadata) == 0 {
		return nil
	}
	if err := decodeWorkerMetadata(metadata, existing); err != nil {
		return err
	}
	if existing.BaitoWorkerId == "" {
		existing.BaitoWorkerId = existing.UserID
	}
	return nil
}

func metadataToResponse(metadata []byte) (BaitoWorkersResponse, error) {
	var worker BaitoWorker
	if err := metadataToWorker(metadata, &worker); err != nil {
		return BaitoWorkersResponse{}, err
	}

	return BaitoWorkersResponse{
		UserID:        worker.UserID,
		BaitoWorkerId: worker.BaitoWorkerId,
		Name:          worker.Name,
		Age:           worker.Age,
		Phone:         worker.Phone,
		Location:      worker.Location,
		Preferred:     worker.Preferred,
		Bio:           worker.Bio,
		ProfilePic:    worker.Avatar,
		CreatedAt:     worker.CreatedAt,
	}, nil
}

func findWorkerByIDFromDB(ctx context.Context, app *infra.Deps, workerID string) (BaitoWorker, error) {
	var worker BaitoWorker
	if err := ensureWorkerDB(app); err != nil {
		return worker, err
	}
	if strings.TrimSpace(workerID) == "" {
		return worker, nil
	}

	query := `SELECT baitoworkerid, userid, created_at, updated_at, metadata
		FROM baitoworkers
		WHERE baitoworkerid = $1
		LIMIT 1`

	var workerIDRow, userID string
	var createdAt, updatedAt time.Time
	var metadata []byte
	if err := app.SQLDB.QueryRow(ctx, query, workerID).Scan(&workerIDRow, &userID, &createdAt, &updatedAt, &metadata); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return worker, nil
		}
		return worker, err
	}

	worker.BaitoWorkerId = workerIDRow
	worker.UserID = userID
	worker.CreatedAt = createdAt.Unix()
	worker.UpdatedAt = updatedAt.Unix()
	if err := metadataToWorker(metadata, &worker); err != nil {
		return worker, err
	}
	if worker.BaitoWorkerId == "" {
		worker.BaitoWorkerId = workerIDRow
	}
	return worker, nil
}

func getUniqueWorkerSkillsFromDB(ctx context.Context, app *infra.Deps) ([]string, error) {
	if err := ensureWorkerDB(app); err != nil {
		return nil, err
	}

	query := `SELECT DISTINCT value
		FROM (
			SELECT jsonb_array_elements_text(COALESCE(metadata->'preferredRoles', '[]'::jsonb)) AS value
			FROM baitoworkers
			WHERE metadata ? 'preferredRoles'
		) t
		WHERE trim(value) <> ''
		ORDER BY value`

	rows, err := app.SQLDB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []string
	for rows.Next() {
		var skill string
		if err := rows.Scan(&skill); err != nil {
			return nil, err
		}
		skill = strings.TrimSpace(skill)
		if skill != "" {
			skills = append(skills, skill)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return skills, nil
}

func findWorkersFromDB(ctx context.Context, app *infra.Deps, query string, args []any, opts map[string]any) ([]BaitoWorkersResponse, error) {
	if err := ensureWorkerDB(app); err != nil {
		return nil, err
	}
	if query == "" {
		query = "TRUE"
	}
	if len(args) == 0 && query == "TRUE" {
		query = "1 = 1"
	}

	offset := 0
	limit := 10
	if v, ok := opts["skip"]; ok {
		if n, ok := v.(int); ok {
			offset = n
		}
		if n, ok := v.(int64); ok {
			offset = int(n)
		}
		if n, ok := v.(float64); ok {
			offset = int(n)
		}
	}
	if v, ok := opts["limit"]; ok {
		if n, ok := v.(int); ok {
			limit = n
		}
		if n, ok := v.(int64); ok {
			limit = int(n)
		}
		if n, ok := v.(float64); ok {
			limit = int(n)
		}
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	limitIndex := len(args) + 1
	offsetIndex := len(args) + 2
	statement := fmt.Sprintf(`SELECT baitoworkerid, userid, created_at, updated_at, metadata
		FROM baitoworkers
		WHERE %s
		ORDER BY updated_at DESC NULLS LAST, created_at DESC
		LIMIT $%d OFFSET $%d`, query, limitIndex, offsetIndex)
	queryArgs := append(append([]any{}, args...), limit, offset)

	rows, err := app.SQLDB.Query(ctx, statement, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workers := make([]BaitoWorkersResponse, 0, limit)
	for rows.Next() {
		var workerID, userID string
		var createdAt, updatedAt time.Time
		var metadata []byte
		if err := rows.Scan(&workerID, &userID, &createdAt, &updatedAt, &metadata); err != nil {
			return nil, err
		}

		resp, err := metadataToResponse(metadata)
		if err != nil {
			return nil, err
		}
		if workerID != "" {
			resp.BaitoWorkerId = workerID
		}
		if userID != "" {
			resp.UserID = userID
		}
		if resp.CreatedAt == 0 {
			resp.CreatedAt = createdAt.Unix()
		}
		if resp.BaitoWorkerId == "" {
			resp.BaitoWorkerId = workerID
		}
		if resp.UserID == "" {
			resp.UserID = userID
		}
		if resp.ProfilePic == "" {
			resp.ProfilePic = ""
		}
		workers = append(workers, resp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return workers, nil
}

func countWorkersFromDB(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	if err := ensureWorkerDB(app); err != nil {
		return 0, err
	}
	if query == "" {
		query = "TRUE"
	}
	if len(args) == 0 && query == "TRUE" {
		query = "1 = 1"
	}

	statement := fmt.Sprintf(`SELECT COUNT(*) FROM baitoworkers WHERE %s`, query)
	var total int64
	if err := app.SQLDB.QueryRow(ctx, statement, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func findExistingWorkerProfile(ctx context.Context, app *infra.Deps, userID string, result any) error {
	if err := ensureWorkerDB(app); err != nil {
		return err
	}
	if strings.TrimSpace(userID) == "" {
		return nil
	}
	if result == nil {
		return errors.New("nil result")
	}

	query := `SELECT baitoworkerid, userid, created_at, updated_at, metadata
		FROM baitoworkers
		WHERE userid = $1
		ORDER BY created_at DESC
		LIMIT 1`

	var workerID, matchedUserID string
	var createdAt, updatedAt time.Time
	var metadata []byte
	if err := app.SQLDB.QueryRow(ctx, query, userID).Scan(&workerID, &matchedUserID, &createdAt, &updatedAt, &metadata); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	switch v := result.(type) {
	case *BaitoWorker:
		v.BaitoWorkerId = workerID
		v.UserID = matchedUserID
		v.CreatedAt = createdAt.Unix()
		v.UpdatedAt = updatedAt.Unix()
		if err := metadataToWorker(metadata, v); err != nil {
			return err
		}
		if v.BaitoWorkerId == "" {
			v.BaitoWorkerId = workerID
		}
		if v.UserID == "" {
			v.UserID = matchedUserID
		}
		return nil
	default:
		if len(metadata) == 0 {
			return nil
		}
		return json.Unmarshal(metadata, result)
	}
}

func createWorkerProfileRecord(ctx context.Context, app *infra.Deps, worker BaitoWorker) error {
	if err := ensureWorkerDB(app); err != nil {
		return err
	}
	if worker.BaitoWorkerId == "" {
		worker.BaitoWorkerId = worker.UserID
	}
	if worker.UserID == "" {
		worker.UserID = worker.BaitoWorkerId
	}
	if worker.CreatedAt == 0 {
		worker.CreatedAt = time.Now().Unix()
	}
	if worker.UpdatedAt == 0 {
		worker.UpdatedAt = worker.CreatedAt
	}
	payload, err := json.Marshal(worker)
	if err != nil {
		return err
	}

	query := `INSERT INTO baitoworkers (baitoworkerid, userid, status, created_at, updated_at, metadata)
		VALUES ($1, $2, 'active', NOW(), NOW(), $3)
		ON CONFLICT (baitoworkerid) DO UPDATE SET
			userid = EXCLUDED.userid,
			status = EXCLUDED.status,
			updated_at = NOW(),
			metadata = EXCLUDED.metadata`

	_, err = app.SQLDB.Exec(ctx, query, worker.BaitoWorkerId, worker.UserID, payload)
	return err
}

func updateWorkerProfileRecord(ctx context.Context, app *infra.Deps, workerID, userID string, update map[string]any) error {
	if err := ensureWorkerDB(app); err != nil {
		return err
	}
	if strings.TrimSpace(workerID) == "" && strings.TrimSpace(userID) == "" {
		return errors.New("worker id or user id required")
	}
	if len(update) == 0 {
		return nil
	}

	incoming := map[string]any{}
	if set, ok := update["$set"]; ok {
		if m, ok := set.(map[string]any); ok {
			incoming = m
		} else {
			incoming = map[string]any{"data": set}
		}
	} else {
		incoming = update
	}
	if len(incoming) == 0 {
		return nil
	}

	var metadata []byte
	if strings.TrimSpace(workerID) != "" && strings.TrimSpace(userID) != "" {
		query := `SELECT metadata FROM baitoworkers WHERE baitoworkerid = $1 AND userid = $2 LIMIT 1`
		if err := app.SQLDB.QueryRow(ctx, query, workerID, userID).Scan(&metadata); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}
		}
	}
	if len(metadata) == 0 {
		query := `SELECT metadata FROM baitoworkers WHERE baitoworkerid = $1 OR userid = $2 ORDER BY created_at DESC LIMIT 1`
		if err := app.SQLDB.QueryRow(ctx, query, workerID, userID).Scan(&metadata); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return err
		}
	}

	base := map[string]any{}
	if len(metadata) > 0 {
		if err := json.Unmarshal(metadata, &base); err != nil {
			return err
		}
	}
	for key, value := range incoming {
		base[key] = value
	}

	payload, err := json.Marshal(base)
	if err != nil {
		return err
	}

	statement := `UPDATE baitoworkers SET metadata = $1::jsonb, updated_at = NOW() WHERE baitoworkerid = $2 AND userid = $3`
	if strings.TrimSpace(workerID) == "" || strings.TrimSpace(userID) == "" {
		statement = `UPDATE baitoworkers SET metadata = $1::jsonb, updated_at = NOW() WHERE baitoworkerid = $2 OR userid = $3`
	}
	_, err = app.SQLDB.Exec(ctx, statement, payload, workerID, userID)
	return err
}

func addWorkerRoleToUser(ctx context.Context, app *infra.Deps, userID string) error {
	if err := ensureWorkerDB(app); err != nil {
		return err
	}
	if strings.TrimSpace(userID) == "" {
		return nil
	}

	query := `UPDATE users
		SET role = CASE
			WHEN role IS NULL THEN ARRAY['worker']::TEXT[]
			WHEN NOT (role @> ARRAY['worker']) THEN array_append(role, 'worker')
			ELSE role
		END,
		updated_at = NOW()
		WHERE userid = $1`
	_, err := app.SQLDB.Exec(ctx, query, userID)
	return err
}

func touchUserUpdatedAt(ctx context.Context, app *infra.Deps, userID string) error {
	if err := ensureWorkerDB(app); err != nil {
		return err
	}
	if strings.TrimSpace(userID) == "" {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `UPDATE users SET updated_at = NOW() WHERE userid = $1`, userID)
	return err
}

func buildWorkerQuery(search string, skill string) (string, []any) {
	clauses := make([]string, 0, 2)
	args := make([]any, 0, 4)

	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		namePH := fmt.Sprintf("$%d", len(args)+1)
		args = append(args, searchTerm)
		locPH := fmt.Sprintf("$%d", len(args)+1)
		args = append(args, searchTerm)
		bioPH := fmt.Sprintf("$%d", len(args)+1)
		args = append(args, searchTerm)
		clauses = append(clauses, fmt.Sprintf("(LOWER(COALESCE(metadata->>'name', '')) LIKE LOWER(%s) OR LOWER(COALESCE(metadata->>'location', '')) LIKE LOWER(%s) OR LOWER(COALESCE(metadata->>'bio', '')) LIKE LOWER(%s))", namePH, locPH, bioPH))
	}

	if skill != "" {
		skillTerm := "%" + strings.ToLower(skill) + "%"
		skillPH := fmt.Sprintf("$%d", len(args)+1)
		args = append(args, skillTerm)
		clauses = append(clauses, fmt.Sprintf("(LOWER(COALESCE(metadata->>'preferredRoles', '')) LIKE LOWER(%s) OR LOWER(COALESCE(metadata->>'skills', '')) LIKE LOWER(%s))", skillPH, skillPH))
	}

	if len(clauses) == 0 {
		return "TRUE", nil
	}

	return strings.Join(clauses, " AND "), args
}

func buildWorkerListOptions(skip, limit int) map[string]any {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if skip < 0 {
		skip = 0
	}

	return map[string]any{
		"skip":  skip,
		"limit": limit,
	}
}

func getWorkersPage(ctx context.Context, app *infra.Deps, search, skill string, skip, limit int) ([]BaitoWorkersResponse, int64, error) {
	whereClause, args := buildWorkerQuery(search, skill)
	opts := buildWorkerListOptions(skip, limit)

	workers, err := findWorkersFromDB(ctx, app, whereClause, args, opts)
	if err != nil {
		return nil, 0, err
	}

	total, err := countWorkersFromDB(ctx, app, whereClause, args)
	if err != nil {
		return nil, 0, err
	}

	return workers, total, nil
}
