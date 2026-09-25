package workers

import (
	"context"
	"errors"
	"scav/config"
	"scav/infra"
	"scav/infra/db"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var BaitoWorkersTable = config.Tables.BaitoWorkerTable
var UsersTable = config.Tables.UserTable

func SQLfindWorkerByIDFromDB(ctx context.Context, app *infra.Deps, workerID string) (BaitoWorker, error) {
	var worker BaitoWorker
	err := app.DB.FindOne(ctx, BaitoWorkersTable, map[string]any{"baitoWorkerId": workerID}, &worker)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return BaitoWorker{}, nil
	}
	return worker, err
}

func SQLgetUniqueWorkerSkillsFromDB(ctx context.Context, app *infra.Deps) ([]string, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$preferredRoles"}},
		{{Key: "$group", Value: map[string]any{"_id": "$preferredRoles"}}},
		{{Key: "$project", Value: map[string]any{"_id": 0, "skill": "$_id"}}},
	}

	var results []map[string]any
	err := app.DB.Aggregate(ctx, BaitoWorkersTable, pipeline, &results)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []string{}, nil
		}
		return nil, err
	}

	skills := make([]string, 0, len(results))
	for _, r := range results {
		if s, ok := r["skill"].(string); ok && s != "" {
			skills = append(skills, s)
		}
	}

	return skills, nil
}

func SQLfindWorkersFromDB(ctx context.Context, app *infra.Deps, filter any, opts db.FindManyOptions) ([]BaitoWorkersResponse, error) {
	var workers []BaitoWorkersResponse
	err := app.DB.FindManyWithOptions(ctx, BaitoWorkersTable, filter, opts, &workers)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return []BaitoWorkersResponse{}, nil
	}
	return workers, err
}

func SQLcountWorkersFromDB(ctx context.Context, app *infra.Deps, filter any) (int64, error) {
	return app.DB.CountDocuments(ctx, BaitoWorkersTable, filter)
}

func SQLfindExistingWorkerProfile(ctx context.Context, app *infra.Deps, userID string, result any) error {
	err := app.DB.FindOne(ctx, BaitoWorkersTable, map[string]any{"userid": userID}, result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil
	}
	return err
}

func SQLcreateWorkerProfileRecord(ctx context.Context, app *infra.Deps, worker BaitoWorker) error {
	return app.DB.Insert(ctx, BaitoWorkersTable, worker)
}

func SQLupdateWorkerProfileRecord(ctx context.Context, app *infra.Deps, workerID, userID string, update map[string]any) error {
	_, err := app.DB.UpdateOne(ctx, BaitoWorkersTable, map[string]any{
		"baitoWorkerId": workerID,
		"userid":        userID,
	}, update)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil
	}
	return err
}

func SQLaddWorkerRoleToUser(ctx context.Context, app *infra.Deps, userID string) error {
	err := app.DB.AddToSet(ctx, UsersTable, map[string]any{"userid": userID}, "role", "worker")
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil
	}
	return err
}

func SQLtouchUserUpdatedAt(ctx context.Context, app *infra.Deps, userID string) error {
	_, err := app.DB.UpdateOne(ctx, UsersTable, map[string]any{"userid": userID}, map[string]any{"updated_at": time.Now()})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil
	}
	return err
}
