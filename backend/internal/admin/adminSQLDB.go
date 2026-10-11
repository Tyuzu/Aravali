// File: internal/admin/adminSQLDB.go

package admin

import (
	"context"
	"fmt"
	"scav/infra"
	"strings"
)

// Role applications
func FindPendingRoleApplication(ctx context.Context, app *infra.Deps, userID, role string, result *RoleApplication) error {

}

func InsertRoleApplication(ctx context.Context, app *infra.Deps, application RoleApplication) error {

}

func FindRoleApplicationsByUser(ctx context.Context, app *infra.Deps, userID string, result *[]RoleApplication) error {

}

func ListRoleApplicationsDB(ctx context.Context, app *infra.Deps, query string, args []any, result *[]RoleApplication) error {

}

func GetRoleApplicationByID(ctx context.Context, app *infra.Deps, id string, result *RoleApplication) error {

}

func GetUserRoles(ctx context.Context, app *infra.Deps, userID string, result any) error {

}

func UpdateUserRoles(ctx context.Context, app *infra.Deps, userID string, roles []string) (int64, error) {

}

func UpdateRoleApplicationStatus(ctx context.Context, app *infra.Deps, appID, status string) (int64, error) {

}

// Moderator applications
func FindModeratorApplicationByUser(ctx context.Context, app *infra.Deps, userID string, result *ModeratorApplication) error {

}

func InsertModeratorApplication(ctx context.Context, app *infra.Deps, application ModeratorApplication) error {

}

func ListModeratorApplicationsDB(ctx context.Context, app *infra.Deps, query string, args []any, result *[]ModeratorApplication) error {

}

func UpdateModeratorApplicationStatus(ctx context.Context, app *infra.Deps, id, status string) (any, error) {

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
