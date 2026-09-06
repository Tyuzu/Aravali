package admin

import (
	"context"
	"time"

	"scav/infra/db"
)

// Role applications
func FindPendingRoleApplication(ctx context.Context, d db.Database, userID, role string, result *RoleApplication) error {
	return d.FindOne(ctx, roleApplicationsCollection, map[string]any{"userid": userID, "role": role, "status": "pending"}, result)
}

func InsertRoleApplication(ctx context.Context, d db.Database, application RoleApplication) error {
	return d.Insert(ctx, roleApplicationsCollection, application)
}

func FindRoleApplicationsByUser(ctx context.Context, d db.Database, userID string, result *[]RoleApplication) error {
	return d.FindMany(ctx, roleApplicationsCollection, map[string]any{"userid": userID}, result)
}

func ListRoleApplicationsDB(ctx context.Context, d db.Database, filter map[string]any, result *[]RoleApplication) error {
	return d.FindMany(ctx, roleApplicationsCollection, filter, result)
}

func GetRoleApplicationByID(ctx context.Context, d db.Database, id string, result *RoleApplication) error {
	return d.FindOne(ctx, roleApplicationsCollection, map[string]any{"id": id}, result)
}

func GetUserRoles(ctx context.Context, d db.Database, userID string, result any) error {
	return d.FindOne(ctx, usersCollection, map[string]any{"userid": userID}, result)
}

func UpdateUserRoles(ctx context.Context, d db.Database, userID string, roles []string) (any, error) {
	return d.UpdateOne(ctx, usersCollection, map[string]any{"userid": userID}, map[string]any{"$set": map[string]any{"role": roles, "updated_at": time.Now().UTC()}})
}

func UpdateRoleApplicationStatus(ctx context.Context, d db.Database, appID, status string) (any, error) {
	return d.UpdateOne(ctx, roleApplicationsCollection, map[string]any{"id": appID}, map[string]any{"$set": map[string]any{"status": status, "updated_at": time.Now().UTC()}})
}

// Moderator applications
func FindModeratorApplicationByUser(ctx context.Context, d db.Database, userID string, result *ModeratorApplication) error {
	return d.FindOne(ctx, moderatorApplicationsCollection, map[string]any{"userid": userID}, result)
}

func InsertModeratorApplication(ctx context.Context, d db.Database, application ModeratorApplication) error {
	return d.Insert(ctx, moderatorApplicationsCollection, application)
}

func ListModeratorApplicationsDB(ctx context.Context, d db.Database, filter map[string]any, result *[]ModeratorApplication) error {
	return d.FindMany(ctx, moderatorApplicationsCollection, filter, result)
}

func UpdateModeratorApplicationStatus(ctx context.Context, d db.Database, id, status string) (any, error) {
	return d.UpdateOne(ctx, moderatorApplicationsCollection, map[string]any{"id": id}, map[string]any{"$set": map[string]any{"status": status, "updatedAt": time.Now().UTC(), "updated_at": time.Now().UTC()}})
}
