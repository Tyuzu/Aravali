// File: internal/admin/adminSQLDB.go

package admin

import (
	"context"
	"fmt"
	"scav/infra"
	"strings"
	"time"
)

// Role applications
func FindPendingRoleApplication(ctx context.Context, app *infra.Deps, userID, role string, result *RoleApplication) error {
	query := "userid = $1 AND role = $2 AND status = $3"
	args := []any{userID, role, "pending"}
	return app.SQLDB.FindOne(ctx, roleApplicationsTable, query, args, result)
}

func InsertRoleApplication(ctx context.Context, app *infra.Deps, application RoleApplication) error {
	return app.SQLDB.Insert(ctx, roleApplicationsTable, application)
}

func FindRoleApplicationsByUser(ctx context.Context, app *infra.Deps, userID string, result *[]RoleApplication) error {
	query := "userid = $1"
	args := []any{userID}
	return app.SQLDB.FindMany(ctx, roleApplicationsTable, query, args, result)
}

func ListRoleApplicationsDB(ctx context.Context, app *infra.Deps, query string, args []any, result *[]RoleApplication) error {
	return app.SQLDB.FindMany(ctx, roleApplicationsTable, query, args, result)
}

func GetRoleApplicationByID(ctx context.Context, app *infra.Deps, id string, result *RoleApplication) error {
	query := "id = $1"
	args := []any{id}
	return app.SQLDB.FindOne(ctx, roleApplicationsTable, query, args, result)
}

func GetUserRoles(ctx context.Context, app *infra.Deps, userID string, result any) error {
	query := "userid = $1"
	args := []any{userID}
	return app.SQLDB.FindOne(ctx, usersTable, query, args, result)
}

func UpdateUserRoles(ctx context.Context, app *infra.Deps, userID string, roles []string) (int64, error) {
	query := "userid = $1"
	args := []any{userID}

	updateValues := map[string]any{
		"role":       roles,
		"updated_at": time.Now().UTC(),
	}
	return app.SQLDB.UpdateOne(ctx, usersTable, query, args, updateValues)
}

func UpdateRoleApplicationStatus(ctx context.Context, app *infra.Deps, appID, status string) (int64, error) {
	query := "id = $1"
	args := []any{appID}

	updateValues := map[string]any{
		"status":     status,
		"updated_at": time.Now().UTC(),
	}
	return app.SQLDB.UpdateOne(ctx, roleApplicationsTable, query, args, updateValues)
}

// Moderator applications
func FindModeratorApplicationByUser(ctx context.Context, app *infra.Deps, userID string, result *ModeratorApplication) error {
	query := "userid = $1"
	args := []any{userID}
	return app.SQLDB.FindOne(ctx, moderatorApplicationsTable, query, args, result)
}

func InsertModeratorApplication(ctx context.Context, app *infra.Deps, application ModeratorApplication) error {
	return app.SQLDB.Insert(ctx, moderatorApplicationsTable, application)
}

func ListModeratorApplicationsDB(ctx context.Context, app *infra.Deps, query string, args []any, result *[]ModeratorApplication) error {
	return app.SQLDB.FindMany(ctx, moderatorApplicationsTable, query, args, result)
}

func UpdateModeratorApplicationStatus(ctx context.Context, app *infra.Deps, id, status string) (any, error) {
	query := "id = $1"
	args := []any{id}

	updateValues := map[string]any{
		"status":     status,
		"updated_at": time.Now().UTC(),
	}
	return app.SQLDB.UpdateOne(ctx, moderatorApplicationsTable, query, args, updateValues)
}

func buildFilterQuery(filter map[string]any) (string, []any) {
	if len(filter) == 0 {
		return "1 = 1", nil
	}
	clauses := make([]string, 0, len(filter))
	args := make([]any, 0, len(filter))
	for key, value := range filter {
		if value == nil {
			continue
		}
		switch v := value.(type) {
		case map[string]any:
			if inVals, ok := v["$in"]; ok {
				clauses = append(clauses, fmt.Sprintf("%s = ANY($%d)", key, len(args)+1))
				args = append(args, inVals)
				continue
			}
			if ninVals, ok := v["$nin"]; ok {
				clauses = append(clauses, fmt.Sprintf("NOT (%s = ANY($%d))", key, len(args)+1))
				args = append(args, ninVals)
				continue
			}
		}
		if strings.HasSuffix(key, "_ne") {
			clauses = append(clauses, fmt.Sprintf("%s <> $%d", strings.TrimSuffix(key, "_ne"), len(args)+1))
			args = append(args, value)
			continue
		}
		clauses = append(clauses, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value)
	}
	if len(clauses) == 0 {
		return "1 = 1", nil
	}
	return strings.Join(clauses, " AND "), args
}
