// File: internal/workers/workerSQLDB.go

package workers

import (
	"context"
	"fmt"
	"strings"

	"scav/config"
	"scav/infra"
)

var BaitoWorkersTable = config.Tables.BaitoWorkerTable
var UsersTable = config.Tables.UserTable

func findWorkerByIDFromDB(ctx context.Context, app *infra.Deps, workerID string) (BaitoWorker, error) {
	var worker BaitoWorker

	return worker, err
}

func getUniqueWorkerSkillsFromDB(ctx context.Context, app *infra.Deps) ([]string, error) {
	var skills []string

	return skills, nil
}

func findWorkersFromDB(ctx context.Context, app *infra.Deps, query string, args []any, opts map[string]any) ([]BaitoWorkersResponse, error) {
	var workers []BaitoWorkersResponse

	return workers, err
}

func countWorkersFromDB(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
}

func findExistingWorkerProfile(ctx context.Context, app *infra.Deps, userID string, result any) error {

	return err
}

func createWorkerProfileRecord(ctx context.Context, app *infra.Deps, worker BaitoWorker) error {
}

func updateWorkerProfileRecord(ctx context.Context, app *infra.Deps, workerID, userID string, update map[string]any) error {

	return err
}

func addWorkerRoleToUser(ctx context.Context, app *infra.Deps, userID string) error {

	return err
}

func touchUserUpdatedAt(ctx context.Context, app *infra.Deps, userID string) error {
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
