package baito

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

var UsersTable = config.Tables.UserTable
var BaitoTable = config.Tables.BaitoTable
var BaitoAppTable = config.Tables.BaitoApplicationsTable

func SQLdeleteBaitoRecord(ctx context.Context, app *infra.Deps, baitoID, userID string) (int64, error) {
	query := "baitoid = $1 AND ownerid = $2"
	args := []any{baitoID, userID}

	return app.SQLDB.DeleteOne(ctx, BaitoTable, query, args)
}

func SQLsaveBaitoApplication(ctx context.Context, app *infra.Deps, application BaitoApplication) error {
	return app.SQLDB.Insert(ctx, BaitoAppTable, application)
}

func SQLincrementBaitoApplicationCount(ctx context.Context, app *infra.Deps, baitoID string) error {
	query := "baitoid = $1"
	args := []any{baitoID}

	return app.SQLDB.Inc(ctx, BaitoTable, query, args, "application_count", 1)
}

func SQLbuildMyApplicationsResult(applications []map[string]any, jobs []Baito) []map[string]any {
	jobByID := make(map[string]Baito, len(jobs))
	for _, job := range jobs {
		jobByID[job.BaitoId] = job
	}

	results := make([]map[string]any, 0, len(applications))
	for _, application := range applications {
		result := map[string]any{}
		for k, v := range application {
			result[k] = v
		}

		jobID, _ := application["baitoid"].(string)
		if job, ok := jobByID[jobID]; ok {
			result["jobId"] = job.BaitoId
			result["title"] = job.Title
			result["location"] = job.Location
			result["wage"] = job.Wage
		}

		if _, ok := result["id"]; !ok {
			if id, ok := application["_id"]; ok {
				result["id"] = id
			}
		}
		if _, ok := result["jobId"]; !ok {
			result["jobId"] = jobID
		}
		results = append(results, result)
	}
	return results
}

func SQLcreateBaitoRecord(ctx context.Context, app *infra.Deps, baito Baito) error {
	return app.SQLDB.Insert(ctx, BaitoTable, baito)
}

func SQLupdateBaitoRecord(ctx context.Context, app *infra.Deps, baitoID, userID string, update map[string]any) (int64, error) {
	query := "baitoid = $1 AND ownerid = $2"
	args := []any{baitoID, userID}

	return app.SQLDB.UpdateOne(ctx, BaitoTable, query, args, update)
}

func SQLfindLatestBaitosFromDB(ctx context.Context, app *infra.Deps, query string, args []any, limit int) ([]BaitosResponse, error) {
	var baitos []BaitosResponse
	err := app.SQLDB.FindManyWithOptions(ctx, BaitoTable, query, args, sqldb.FindManyOptions{
		Limit: limit,
	}, &baitos)
	return baitos, err
}

func SQLfindRelatedBaitosFromDB(ctx context.Context, app *infra.Deps, query string, args []any, limit int) ([]BaitosResponse, error) {
	var baitos []BaitosResponse
	err := app.SQLDB.FindManyWithOptions(ctx, BaitoTable, query, args, sqldb.FindManyOptions{
		Limit: limit,
	}, &baitos)
	return baitos, err
}

func SQLfindBaitoByIDFromDB(ctx context.Context, app *infra.Deps, baitoID string) (Baito, error) {
	var baito Baito
	query := "baitoid = $1"
	args := []any{baitoID}

	err := app.SQLDB.FindOne(ctx, BaitoTable, query, args, &baito)
	return baito, err
}

func SQLfindMyBaitosFromDB(ctx context.Context, app *infra.Deps, userID string) ([]BaitosResponse, error) {
	var baitos []BaitosResponse
	query := "ownerid = $1"
	args := []any{userID}

	err := app.SQLDB.FindManyWithOptions(ctx, BaitoTable, query, args, sqldb.FindManyOptions{}, &baitos)
	return baitos, err
}

func SQLcountApplicationsForBaito(ctx context.Context, app *infra.Deps, baitoID string) (int64, error) {
	query := "baitoid = $1"
	args := []any{baitoID}

	return app.SQLDB.CountDocuments(ctx, BaitoAppTable, query, args)
}

func SQLfindBaitoApplicantsFromDB(ctx context.Context, app *infra.Deps, baitoID string) ([]map[string]any, error) {
	var results []map[string]any
	query := "baitoid = $1"
	args := []any{baitoID}

	err := app.SQLDB.FindMany(ctx, BaitoAppTable, query, args, &results)
	return results, err
}

func SQLfindMyApplicationsFromDB(ctx context.Context, app *infra.Deps, userID string) ([]map[string]any, error) {
	var applications []map[string]any
	query := "userid = $1"
	args := []any{userID}

	if err := app.SQLDB.FindMany(ctx, BaitoAppTable, query, args, &applications); err != nil {
		return nil, err
	}
	if len(applications) == 0 {
		return []map[string]any{}, nil
	}

	jobIDs := make([]string, 0, len(applications))
	seen := make(map[string]struct{}, len(applications))
	for _, application := range applications {
		jobID, _ := application["baitoid"].(string)
		if jobID == "" {
			continue
		}
		if _, ok := seen[jobID]; ok {
			continue
		}
		seen[jobID] = struct{}{}
		jobIDs = append(jobIDs, jobID)
	}

	jobs := make([]Baito, 0, len(jobIDs))
	for _, jobID := range jobIDs {
		var job Baito
		jobQuery := "baitoid = $1"
		jobArgs := []any{jobID}

		if err := app.SQLDB.FindOne(ctx, BaitoTable, jobQuery, jobArgs, &job); err != nil {
			continue
		}
		jobs = append(jobs, job)
	}

	return buildMyApplicationsResult(applications, jobs), nil
}
