package farms

import (
	"context"

	"scav/config"
	"scav/infra/sqldb"
	"scav/internal/cart"
)

var (
	cropsTable      = config.Tables.CropsTable
	farmsTable      = config.Tables.FarmsTable
	farmOrdersTable = config.Tables.FarmOrdersTable
)

func SQLinsertFarm(ctx context.Context, database sqldb.Database, farm Farm) error {
	return database.InsertOne(ctx, farmsTable, farm)
}

func SQLgetFarmByID(ctx context.Context, database sqldb.Database, farmID string) (Farm, error) {
	var farm Farm
	query := "farmid = $1"
	args := []any{farmID}

	err := database.FindOne(ctx, farmsTable, query, args, &farm)
	return farm, err
}

func SQLgetFarmByCreatedBy(ctx context.Context, database sqldb.Database, userID string) (Farm, error) {
	var farm Farm
	query := "created_by = $1"
	args := []any{userID}

	err := database.FindOne(ctx, farmsTable, query, args, &farm)
	return farm, err
}

func SQLgetCropsByFarmID(ctx context.Context, database sqldb.Database, farmID string) ([]Crop, error) {
	var crops []Crop
	query := "farmid = $1"
	args := []any{farmID}

	if err := database.FindMany(ctx, cropsTable, query, args, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []Crop{}, nil
	}
	return crops, nil
}

func SQLgetFarmOrdersByFarmID(ctx context.Context, database sqldb.Database, farmID string) ([]cart.FarmOrder, error) {
	var orders []cart.FarmOrder
	query := "farmid = $1"
	args := []any{farmID}

	if err := database.FindMany(ctx, farmOrdersTable, query, args, &orders); err != nil {
		return nil, err
	}
	if orders == nil {
		return []cart.FarmOrder{}, nil
	}
	return orders, nil
}

func SQLgetCropsByCropID(ctx context.Context, database sqldb.Database, cropID string) ([]Crop, error) {
	var crops []Crop
	query := "cropid = $1"
	args := []any{cropID}

	if err := database.FindMany(ctx, cropsTable, query, args, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []Crop{}, nil
	}
	return crops, nil
}

func SQLgetFarmsByIDs(ctx context.Context, database sqldb.Database, farmIDs []string) ([]Farm, error) {
	var farms []Farm
	query := "farmid = ANY($1)"
	args := []any{farmIDs}

	if err := database.FindMany(ctx, farmsTable, query, args, &farms); err != nil {
		return nil, err
	}
	if farms == nil {
		return []Farm{}, nil
	}
	return farms, nil
}

func SQLgetCropsByNameFilter(ctx context.Context, database sqldb.Database, cropName string) ([]Crop, error) {
	query := "name ILIKE $1"
	args := []any{cropName}

	var crops []Crop
	if err := database.FindMany(ctx, cropsTable, query, args, &crops); err != nil {
		return nil, err
	}
	if crops == nil {
		return []Crop{}, nil
	}
	return crops, nil
}

func SQLgetPaginatedFarms(ctx context.Context, database sqldb.Database, search string, offset, limit int) ([]Farm, int64, error) {
	whereClause := ""
	args := []any{}

	if search != "" {
		whereClause = "WHERE name ILIKE $1 OR location ILIKE $1 OR owner ILIKE $1"
		args = append(args, "%"+search+"%")
	}

	argOffsetIndex := len(args) + 1
	argLimitIndex := len(args) + 2

	rawQuery := `
		SELECT f.*, COALESCE(json_agg(c) FILTER (WHERE c.cropid IS NOT NULL), '[]') AS crops
		FROM ` + farmsTable + ` f
		LEFT JOIN ` + cropsTable + ` c ON f.farmid = c.farmid
		` + whereClause + `
		GROUP BY f.farmid
		ORDER BY f.created_at DESC
		OFFSET $` + string(rune('0'+argOffsetIndex)) + ` LIMIT $` + string(rune('0'+argLimitIndex))

	queryArgs := append(args, offset, limit)

	var farms []Farm
	if err := database.QueryRaw(ctx, rawQuery, queryArgs, &farms); err != nil {
		return nil, 0, err
	}

	total, err := database.Count(ctx, farmsTable, whereClause, args)
	if err != nil {
		if search == "" {
			return farms, 0, nil
		}
		return farms, 0, err
	}
	return farms, total, nil
}

func SQLgetMyFarmsPage(ctx context.Context, database sqldb.Database, userID string, offset, limit int) ([]Farm, int64, error) {
	rawQuery := `
		SELECT f.*, COALESCE(json_agg(c) FILTER (WHERE c.cropid IS NOT NULL), '[]') AS crops
		FROM ` + farmsTable + ` f
		LEFT JOIN ` + cropsTable + ` c ON f.farmid = c.farmid
		WHERE f.created_by = $1
		GROUP BY f.farmid
		ORDER BY f.created_at DESC
		OFFSET $2 LIMIT $3
	`

	var farms []Farm
	if err := database.QueryRaw(ctx, rawQuery, []any{userID, offset, limit}, &farms); err != nil {
		return nil, 0, err
	}

	total, err := database.Count(ctx, farmsTable, "WHERE created_by = $1", []any{userID})
	if err != nil {
		return farms, 0, nil
	}
	return farms, total, nil
}

func SQLupdateOwnedFarm(ctx context.Context, database sqldb.Database, farmID, userID string, update map[string]any) (int64, error) {
	query := "farmid = $1 AND created_by = $2"
	args := []any{farmID, userID}

	return database.UpdateOne(ctx, farmsTable, query, args, update)
}

func SQLdeleteFarmByID(ctx context.Context, database sqldb.Database, farmID string) (int64, error) {
	query := "farmid = $1"
	args := []any{farmID}

	return database.DeleteOne(ctx, farmsTable, query, args)
}
