// File: internal/workers/workerSQLDB.go

package workers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

var BaitoWorkersTable = config.Tables.BaitoWorkerTable
var UsersTable = config.Tables.UserTable

func findWorkerByIDFromDB(ctx context.Context, app *infra.Deps, workerID string) (BaitoWorker, error) {
	var worker BaitoWorker
	query := "baitoWorkerId = $1"
	args := []any{workerID}

	err := app.SQLDB.FindOne(ctx, BaitoWorkersTable, query, args, &worker)
	if errors.Is(err, sql.ErrNoRows) {
		return BaitoWorker{}, nil
	}
	return worker, err
}

func getUniqueWorkerSkillsFromDB(ctx context.Context, app *infra.Deps) ([]string, error) {
	var skills []string
	err := app.SQLDB.Distinct(ctx, BaitoWorkersTable, "preferredRoles", "", nil, &skills)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []string{}, nil
		}
		return nil, err
	}

	return skills, nil
}

func findWorkersFromDB(ctx context.Context, app *infra.Deps, query string, args []any, opts sqldb.FindManyOptions) ([]BaitoWorkersResponse, error) {
	var workers []BaitoWorkersResponse
	err := app.SQLDB.FindManyWithOptions(ctx, BaitoWorkersTable, query, args, opts, &workers)
	if errors.Is(err, sql.ErrNoRows) {
		return []BaitoWorkersResponse{}, nil
	}
	return workers, err
}

func countWorkersFromDB(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.CountDocuments(ctx, BaitoWorkersTable, query, args)
}

func findExistingWorkerProfile(ctx context.Context, app *infra.Deps, userID string, result any) error {
	query := "userid = $1"
	args := []any{userID}

	err := app.SQLDB.FindOne(ctx, BaitoWorkersTable, query, args, result)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

func createWorkerProfileRecord(ctx context.Context, app *infra.Deps, worker BaitoWorker) error {
	return app.SQLDB.Insert(ctx, BaitoWorkersTable, worker)
}

func updateWorkerProfileRecord(ctx context.Context, app *infra.Deps, workerID, userID string, update map[string]any) error {
	query := "baitoWorkerId = $1 AND userid = $2"
	args := []any{workerID, userID}

	_, err := app.SQLDB.UpdateOne(ctx, BaitoWorkersTable, query, args, update)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

func addWorkerRoleToUser(ctx context.Context, app *infra.Deps, userID string) error {
	query := "userid = $1"
	args := []any{userID}

	err := app.SQLDB.AddToSet(ctx, UsersTable, query, args, "role", "worker")
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

func touchUserUpdatedAt(ctx context.Context, app *infra.Deps, userID string) error {
	query := "userid = $1"
	args := []any{userID}

	update := map[string]any{"updated_at": time.Now()}
	_, err := app.SQLDB.UpdateOne(ctx, UsersTable, query, args, update)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
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
		clauses = append(clauses, fmt.Sprintf("(LOWER(name) LIKE LOWER(%s) OR LOWER(location) LIKE LOWER(%s) OR LOWER(bio) LIKE LOWER(%s))", namePH, locPH, bioPH))
	}

	if skill != "" {
		skillTerm := "%" + strings.ToLower(skill) + "%"
		skillPH := fmt.Sprintf("$%d", len(args)+1)
		args = append(args, skillTerm)
		clauses = append(clauses, fmt.Sprintf("(LOWER(CAST(preferred AS TEXT)) LIKE LOWER(%s))", skillPH))
	}

	if len(clauses) == 0 {
		return "TRUE", nil
	}

	return strings.Join(clauses, " AND "), args
}

func buildWorkerListOptions(skip, limit int) sqldb.FindManyOptions {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if skip < 0 {
		skip = 0
	}

	return sqldb.FindManyOptions{
		Limit:  int64(limit),
		Offset: int64(skip),
		Sort: []sqldb.OrderBy{{
			Column:     "createdat",
			Descending: true,
		}},
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
