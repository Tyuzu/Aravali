package workers

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

var BaitoWorkersTable = config.Tables.BaitoWorkerTable
var UsersTable = config.Tables.UserTable

func SQLfindWorkerByIDFromDB(ctx context.Context, app *infra.Deps, workerID string) (BaitoWorker, error) {
	var worker BaitoWorker
	query := "baitoWorkerId = $1"
	args := []any{workerID}

	err := app.SQLDB.FindOne(ctx, BaitoWorkersTable, query, args, &worker)
	if errors.Is(err, sql.ErrNoRows) {
		return BaitoWorker{}, nil
	}
	return worker, err
}

func SQLgetUniqueWorkerSkillsFromDB(ctx context.Context, app *infra.Deps) ([]string, error) {
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

func SQLfindWorkersFromDB(ctx context.Context, app *infra.Deps, query string, args []any, opts sqldb.FindManyOptions) ([]BaitoWorkersResponse, error) {
	var workers []BaitoWorkersResponse
	err := app.SQLDB.FindManyWithOptions(ctx, BaitoWorkersTable, query, args, opts, &workers)
	if errors.Is(err, sql.ErrNoRows) {
		return []BaitoWorkersResponse{}, nil
	}
	return workers, err
}

func SQLcountWorkersFromDB(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.CountDocuments(ctx, BaitoWorkersTable, query, args)
}

func SQLfindExistingWorkerProfile(ctx context.Context, app *infra.Deps, userID string, result any) error {
	query := "userid = $1"
	args := []any{userID}

	err := app.SQLDB.FindOne(ctx, BaitoWorkersTable, query, args, result)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

func SQLcreateWorkerProfileRecord(ctx context.Context, app *infra.Deps, worker BaitoWorker) error {
	return app.SQLDB.Insert(ctx, BaitoWorkersTable, worker)
}

func SQLupdateWorkerProfileRecord(ctx context.Context, app *infra.Deps, workerID, userID string, update map[string]any) error {
	query := "baitoWorkerId = $1 AND userid = $2"
	args := []any{workerID, userID}

	_, err := app.SQLDB.UpdateOne(ctx, BaitoWorkersTable, query, args, update)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

func SQLaddWorkerRoleToUser(ctx context.Context, app *infra.Deps, userID string) error {
	query := "userid = $1"
	args := []any{userID}

	err := app.SQLDB.AddToSet(ctx, UsersTable, query, args, "role", "worker")
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

func SQLtouchUserUpdatedAt(ctx context.Context, app *infra.Deps, userID string) error {
	query := "userid = $1"
	args := []any{userID}

	update := map[string]any{"updated_at": time.Now()}
	_, err := app.SQLDB.UpdateOne(ctx, UsersTable, query, args, update)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}
