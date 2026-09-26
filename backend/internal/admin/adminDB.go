package admin

import (
	"context"
	"scav/infra"
	"time"
)

// Role applications
func FindPendingRoleApplication(ctx context.Context, app *infra.Deps, userID, role string, result *RoleApplication) error {
	return app.DB.FindOne(ctx, roleApplicationsCollection, map[string]any{"userid": userID, "role": role, "status": "pending"}, result)
}

func InsertRoleApplication(ctx context.Context, app *infra.Deps, application RoleApplication) error {
	return app.DB.Insert(ctx, roleApplicationsCollection, application)
}

func FindRoleApplicationsByUser(ctx context.Context, app *infra.Deps, userID string, result *[]RoleApplication) error {
	return app.DB.FindMany(ctx, roleApplicationsCollection, map[string]any{"userid": userID}, result)
}

func ListRoleApplicationsDB(ctx context.Context, app *infra.Deps, filter map[string]any, result *[]RoleApplication) error {
	return app.DB.FindMany(ctx, roleApplicationsCollection, filter, result)
}

func GetRoleApplicationByID(ctx context.Context, app *infra.Deps, id string, result *RoleApplication) error {
	return app.DB.FindOne(ctx, roleApplicationsCollection, map[string]any{"id": id}, result)
}

func GetUserRoles(ctx context.Context, app *infra.Deps, userID string, result any) error {
	return app.DB.FindOne(ctx, usersCollection, map[string]any{"userid": userID}, result)
}

func UpdateUserRoles(ctx context.Context, app *infra.Deps, userID string, roles []string) (any, error) {
	return app.DB.UpdateOne(ctx, usersCollection, map[string]any{"userid": userID}, map[string]any{"$set": map[string]any{"role": roles, "updated_at": time.Now().UTC()}})
}

func UpdateRoleApplicationStatus(ctx context.Context, app *infra.Deps, appID, status string) (any, error) {
	return app.DB.UpdateOne(ctx, roleApplicationsCollection, map[string]any{"id": appID}, map[string]any{"$set": map[string]any{"status": status, "updated_at": time.Now().UTC()}})
}

// Moderator applications
func FindModeratorApplicationByUser(ctx context.Context, app *infra.Deps, userID string, result *ModeratorApplication) error {
	return app.DB.FindOne(ctx, moderatorApplicationsCollection, map[string]any{"userid": userID}, result)
}

func InsertModeratorApplication(ctx context.Context, app *infra.Deps, application ModeratorApplication) error {
	return app.DB.Insert(ctx, moderatorApplicationsCollection, application)
}

func ListModeratorApplicationsDB(ctx context.Context, app *infra.Deps, filter map[string]any, result *[]ModeratorApplication) error {
	return app.DB.FindMany(ctx, moderatorApplicationsCollection, filter, result)
}

func UpdateModeratorApplicationStatus(ctx context.Context, app *infra.Deps, id, status string) (any, error) {
	return app.DB.UpdateOne(ctx, moderatorApplicationsCollection, map[string]any{"id": id}, map[string]any{"$set": map[string]any{"status": status, "updatedAt": time.Now().UTC(), "updated_at": time.Now().UTC()}})
}
