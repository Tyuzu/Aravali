package baito

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/db"
)

var UsersCollection = config.Collections.UserCollection
var BaitoCollection = config.Collections.BaitoCollection
var BaitoAppCollection = config.Collections.BaitoApplicationsCollection

func deleteBaitoRecord(ctx context.Context, app *infra.Deps, baitoID, userID string) (int64, error) {
	rows, err := app.DB.DeleteOne(ctx, BaitoCollection, map[string]any{
		"baitoid": baitoID,
		"ownerid": userID,
	})
	return rows, err
}

func saveBaitoApplication(ctx context.Context, app *infra.Deps, application BaitoApplication) error {
	return app.DB.Insert(ctx, BaitoAppCollection, application)
}

func incrementBaitoApplicationCount(ctx context.Context, app *infra.Deps, baitoID string) error {
	var baito Baito
	if err := app.DB.FindOne(ctx, BaitoCollection, map[string]any{"baitoid": baitoID}, &baito); err != nil {
		return err
	}

	baito.ApplicationCount++
	_, err := app.DB.UpdateOne(ctx, BaitoCollection, map[string]any{"baitoid": baitoID}, baito)
	return err
}

func buildMyApplicationsResult(applications []map[string]any, jobs []Baito) []map[string]any {
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

func createBaitoRecord(ctx context.Context, app *infra.Deps, baito Baito) error {
	return app.DB.Insert(ctx, BaitoCollection, baito)
}

func updateBaitoRecord(ctx context.Context, app *infra.Deps, baitoID, userID string, update map[string]any) (int64, error) {
	_, err := app.DB.UpdateOne(ctx, BaitoCollection, map[string]any{
		"baitoid": baitoID,
		"ownerid": userID,
	}, update)
	return 0, err
}

func findLatestBaitosFromDB(ctx context.Context, app *infra.Deps, filter any, limit int) ([]BaitosResponse, error) {
	var baitos []BaitosResponse
	err := app.DB.FindManyWithOptions(ctx, BaitoCollection, filter, db.FindManyOptions{
		Limit: limit,
	}, &baitos)
	return baitos, err
}

func findRelatedBaitosFromDB(ctx context.Context, app *infra.Deps, filter any, limit int) ([]BaitosResponse, error) {
	var baitos []BaitosResponse
	err := app.DB.FindManyWithOptions(ctx, BaitoCollection, filter, db.FindManyOptions{
		Limit: limit,
	}, &baitos)
	return baitos, err
}

func findBaitoByIDFromDB(ctx context.Context, app *infra.Deps, baitoID string) (Baito, error) {
	var baito Baito
	err := app.DB.FindOne(ctx, BaitoCollection, map[string]any{"baitoid": baitoID}, &baito)
	return baito, err
}

func findMyBaitosFromDB(ctx context.Context, app *infra.Deps, userID string) ([]BaitosResponse, error) {
	var baitos []BaitosResponse
	err := app.DB.FindManyWithOptions(ctx, BaitoCollection, map[string]any{"ownerId": userID}, db.FindManyOptions{}, &baitos)
	return baitos, err
}

func findBaitoApplicantsFromDB(ctx context.Context, app *infra.Deps, baitoID string) ([]map[string]any, error) {
	var results []map[string]any
	err := app.DB.FindMany(ctx, BaitoAppCollection, map[string]any{"baitoid": baitoID}, &results)
	return results, err
}

func findMyApplicationsFromDB(ctx context.Context, app *infra.Deps, userID string) ([]map[string]any, error) {
	var applications []map[string]any
	if err := app.DB.FindMany(ctx, BaitoAppCollection, map[string]any{"userid": userID}, &applications); err != nil {
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
		if err := app.DB.FindOne(ctx, BaitoCollection, map[string]any{"baitoid": jobID}, &job); err != nil {
			continue
		}
		jobs = append(jobs, job)
	}

	return buildMyApplicationsResult(applications, jobs), nil
}
