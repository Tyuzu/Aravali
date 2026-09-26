package admin

import (
	"context"
	"scav/infra"
	"time"
)

// Role applications
func SQLFindPendingRoleApplication(ctx context.Context, app *infra.Deps, userID, role string, result *RoleApplication) error {
	query := "userid = $1 AND role = $2 AND status = $3"
	args := []any{userID, role, "pending"}
	return app.SQLDB.FindOne(ctx, roleApplicationsCollection, query, args, result)
}

func SQLInsertRoleApplication(ctx context.Context, app *infra.Deps, application RoleApplication) error {
	return app.SQLDB.Insert(ctx, roleApplicationsCollection, application)
}

func SQLFindRoleApplicationsByUser(ctx context.Context, app *infra.Deps, userID string, result *[]RoleApplication) error {
	query := "userid = $1"
	args := []any{userID}
	return app.SQLDB.FindMany(ctx, roleApplicationsCollection, query, args, result)
}

func SQLListRoleApplicationsDB(ctx context.Context, app *infra.Deps, query string, args []any, result *[]RoleApplication) error {
	return app.SQLDB.FindMany(ctx, roleApplicationsCollection, query, args, result)
}

func SQLGetRoleApplicationByID(ctx context.Context, app *infra.Deps, id string, result *RoleApplication) error {
	query := "id = $1"
	args := []any{id}
	return app.SQLDB.FindOne(ctx, roleApplicationsCollection, query, args, result)
}

func SQLGetUserRoles(ctx context.Context, app *infra.Deps, userID string, result any) error {
	query := "userid = $1"
	args := []any{userID}
	return app.SQLDB.FindOne(ctx, usersCollection, query, args, result)
}

func SQLUpdateUserRoles(ctx context.Context, app *infra.Deps, userID string, roles []string) (int64, error) {
	query := "userid = $1"
	args := []any{userID}

	updateValues := map[string]any{
		"role":       roles,
		"updated_at": time.Now().UTC(),
	}
	return app.SQLDB.UpdateOne(ctx, usersCollection, query, args, updateValues)
}

func SQLUpdateRoleApplicationStatus(ctx context.Context, app *infra.Deps, appID, status string) (int64, error) {
	query := "id = $1"
	args := []any{appID}

	updateValues := map[string]any{
		"status":     status,
		"updated_at": time.Now().UTC(),
	}
	return app.SQLDB.UpdateOne(ctx, roleApplicationsCollection, query, args, updateValues)
}

// Moderator applications
func SQLFindModeratorApplicationByUser(ctx context.Context, app *infra.Deps, userID string, result *ModeratorApplication) error {
	query := "userid = $1"
	args := []any{userID}
	return app.SQLDB.FindOne(ctx, moderatorApplicationsCollection, query, args, result)
}

func SQLInsertModeratorApplication(ctx context.Context, app *infra.Deps, application ModeratorApplication) error {
	return app.SQLDB.Insert(ctx, moderatorApplicationsCollection, application)
}

func SQLListModeratorApplicationsDB(ctx context.Context, app *infra.Deps, query string, args []any, result *[]ModeratorApplication) error {
	return app.SQLDB.FindMany(ctx, moderatorApplicationsCollection, query, args, result)
}

func SQLUpdateModeratorApplicationStatus(ctx context.Context, app *infra.Deps, id, status string) (any, error) {
	query := "id = $1"
	args := []any{id}

	updateValues := map[string]any{
		"status":     status,
		"updated_at": time.Now().UTC(),
	}
	return app.SQLDB.UpdateOne(ctx, moderatorApplicationsCollection, query, args, updateValues)
}
