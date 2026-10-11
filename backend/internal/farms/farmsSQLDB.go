// File: internal/farms/farmsSQLDB.go

package farms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/internal/cart"

	"github.com/jackc/pgx/v5"
)

var (
	cropsTable      = config.Tables.CropsTable
	farmsTable      = config.Tables.FarmsTable
	farmOrdersTable = config.Tables.FarmOrdersTable
)

func ensureFarmDB(app *infra.Deps) error {
	if app == nil || app.SQLDB == nil {
		return errors.New("database unavailable")
	}
	return nil
}

func parseFarmMetadata(raw []byte, dst *Farm) {
	if dst == nil || len(raw) == 0 || string(raw) == "null" {
		return
	}
	var meta map[string]any
	if err := json.Unmarshal(raw, &meta); err != nil {
		return
	}
	if v, ok := meta["owner"]; ok {
		if s, ok2 := v.(string); ok2 && s != "" {
			dst.Owner = s
		}
	}
	if v, ok := meta["contact"]; ok {
		if s, ok2 := v.(string); ok2 && s != "" {
			dst.Contact = s
		}
	}
	if v, ok := meta["social"]; ok {
		if s, ok2 := v.(string); ok2 && s != "" {
			dst.Social = s
		}
	}
	if v, ok := meta["practice"]; ok {
		if s, ok2 := v.(string); ok2 && s != "" {
			dst.Practice = s
		}
	}
	if v, ok := meta["availability"]; ok {
		if b, err := json.Marshal(v); err == nil {
			var availability WeeklyAvailability
			if err := json.Unmarshal(b, &availability); err == nil {
				dst.Availability = availability
			}
		}
	}
	if v, ok := meta["tags"]; ok {
		if b, err := json.Marshal(v); err == nil {
			var tags []string
			if err := json.Unmarshal(b, &tags); err == nil {
				dst.Tags = tags
			}
		}
	}
}

func parseCropMetadata(raw []byte, dst *Crop) {
	if dst == nil || len(raw) == 0 || string(raw) == "null" {
		return
	}
	var meta map[string]any
	if err := json.Unmarshal(raw, &meta); err != nil {
		return
	}
	if v, ok := meta["price"]; ok {
		if n, ok2 := v.(float64); ok2 {
			dst.Price = n
		}
	}
	if v, ok := meta["discount"]; ok {
		if n, ok2 := v.(float64); ok2 {
			dst.Discount = n
		}
	}
	if v, ok := meta["quantity"]; ok {
		if n, ok2 := v.(float64); ok2 {
			dst.Quantity = int(n)
		}
	}
	if v, ok := meta["unit"]; ok {
		if s, ok2 := v.(string); ok2 {
			dst.Unit = s
		}
	}
	if v, ok := meta["banner"]; ok {
		if s, ok2 := v.(string); ok2 {
			dst.Banner = s
		}
	}
	if v, ok := meta["category"]; ok {
		if s, ok2 := v.(string); ok2 {
			dst.Category = s
		}
	}
	if v, ok := meta["featured"]; ok {
		if b, ok2 := v.(bool); ok2 {
			dst.Featured = b
		}
	}
	if v, ok := meta["out_of_stock"]; ok {
		if b, ok2 := v.(bool); ok2 {
			dst.OutOfStock = b
		}
	}
	if v, ok := meta["farmid"]; ok {
		if s, ok2 := v.(string); ok2 {
			dst.FarmID = s
		}
	}
}

func makeUpdateSet(update map[string]any) (string, []any) {
	parts := make([]string, 0, len(update))
	args := make([]any, 0, len(update))
	for key, value := range update {
		parts = append(parts, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value)
	}
	return strings.Join(parts, ", "), args
}

func insertFarm(ctx context.Context, app *infra.Deps, farm Farm) error {
	if strings.TrimSpace(farm.FarmID) == "" {
		return errors.New("farm id is required")
	}
	if err := ensureFarmDB(app); err != nil {
		return err
	}

	meta, err := json.Marshal(map[string]any{
		"owner":        farm.Owner,
		"contact":      farm.Contact,
		"social":       farm.Social,
		"practice":     farm.Practice,
		"tags":         farm.Tags,
		"availability": farm.Availability,
	})
	if err != nil {
		return err
	}

	_, err = app.SQLDB.Exec(
		ctx,
		`INSERT INTO `+farmsTable+` (farmid, name, description, userid, status, created_at, updated_at, metadata) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (farmid) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, userid = EXCLUDED.userid, status = EXCLUDED.status, updated_at = NOW(), metadata = EXCLUDED.metadata`,
		farm.FarmID,
		farm.Name,
		farm.Description,
		farm.CreatedBy,
		"active",
		farm.CreatedAt,
		time.Now(),
		meta,
	)
	return err
}

func getFarmByID(ctx context.Context, app *infra.Deps, farmID string) (Farm, error) {
	farmID = strings.TrimSpace(farmID)
	if farmID == "" {
		return Farm{}, errors.New("farm id is required")
	}
	if err := ensureFarmDB(app); err != nil {
		return Farm{}, err
	}

	var farm Farm
	var metadata []byte
	var createdAt, updatedAt time.Time
	if err := app.SQLDB.QueryRow(
		ctx,
		`SELECT farmid, name, description, userid, status, created_at, updated_at, metadata FROM `+farmsTable+` WHERE farmid = $1 LIMIT 1`,
		farmID,
	).Scan(&farm.FarmID, &farm.Name, &farm.Description, &farm.CreatedBy, &farm.Owner, &createdAt, &updatedAt, &metadata); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Farm{}, errors.New("farm not found")
		}
		return Farm{}, err
	}
	farm.CreatedAt = createdAt
	farm.UpdatedAt = updatedAt
	farm.Owner = farm.CreatedBy
	parseFarmMetadata(metadata, &farm)
	return farm, nil
}

func getFarmByCreatedBy(ctx context.Context, app *infra.Deps, userID string) (Farm, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return Farm{}, errors.New("user id is required")
	}
	if err := ensureFarmDB(app); err != nil {
		return Farm{}, err
	}

	var farm Farm
	var metadata []byte
	var createdAt, updatedAt time.Time
	if err := app.SQLDB.QueryRow(
		ctx,
		`SELECT farmid, name, description, userid, status, created_at, updated_at, metadata FROM `+farmsTable+` WHERE userid = $1 LIMIT 1`,
		userID,
	).Scan(&farm.FarmID, &farm.Name, &farm.Description, &farm.CreatedBy, &farm.Owner, &createdAt, &updatedAt, &metadata); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Farm{}, errors.New("farm not found")
		}
		return Farm{}, err
	}
	farm.CreatedAt = createdAt
	farm.UpdatedAt = updatedAt
	farm.Owner = farm.CreatedBy
	parseFarmMetadata(metadata, &farm)
	return farm, nil
}

func getCropsByFarmID(ctx context.Context, app *infra.Deps, farmID string) ([]Crop, error) {
	farmID = strings.TrimSpace(farmID)
	if farmID == "" {
		return nil, errors.New("farm id is required")
	}
	if err := ensureFarmDB(app); err != nil {
		return nil, err
	}

	rows, err := app.SQLDB.Query(
		ctx,
		`SELECT cropid, name, farmid, metadata, created_at, updated_at FROM `+cropsTable+` WHERE farmid = $1 ORDER BY created_at DESC`,
		farmID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Crop, 0)
	for rows.Next() {
		var crop Crop
		var metadata []byte
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&crop.CropId, &crop.Name, &crop.FarmID, &metadata, &createdAt, &updatedAt); err != nil {
			continue
		}
		crop.CreatedAt = createdAt
		crop.UpdatedAt = updatedAt
		parseCropMetadata(metadata, &crop)
		items = append(items, crop)
	}
	return items, nil
}

func getFarmOrdersByFarmID(ctx context.Context, app *infra.Deps, farmID string) ([]cart.FarmOrder, error) {
	farmID = strings.TrimSpace(farmID)
	if farmID == "" {
		return nil, errors.New("farm id is required")
	}
	if err := ensureFarmDB(app); err != nil {
		return nil, err
	}

	rows, err := app.SQLDB.Query(
		ctx,
		`SELECT orderid, farmid, userid, status, total, created_at, metadata FROM `+farmOrdersTable+` WHERE farmid = $1 ORDER BY created_at DESC`,
		farmID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]cart.FarmOrder, 0)
	for rows.Next() {
		var order cart.FarmOrder
		var metadata []byte
		var createdAt time.Time
		if err := rows.Scan(&order.OrderID, &order.FarmID, &order.UserID, &order.Status, &order.Total, &createdAt, &metadata); err != nil {
			continue
		}
		order.CreatedAt = createdAt
		if len(metadata) > 0 {
			_ = json.Unmarshal(metadata, &order)
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func getCropsByCropID(ctx context.Context, app *infra.Deps, cropID string) ([]Crop, error) {
	cropID = strings.TrimSpace(cropID)
	if cropID == "" {
		return nil, errors.New("crop id is required")
	}
	if err := ensureFarmDB(app); err != nil {
		return nil, err
	}

	rows, err := app.SQLDB.Query(
		ctx,
		`SELECT cropid, name, farmid, metadata, created_at, updated_at FROM `+cropsTable+` WHERE cropid = $1 ORDER BY created_at DESC`,
		cropID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Crop, 0)
	for rows.Next() {
		var crop Crop
		var metadata []byte
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&crop.CropId, &crop.Name, &crop.FarmID, &metadata, &createdAt, &updatedAt); err != nil {
			continue
		}
		crop.CreatedAt = createdAt
		crop.UpdatedAt = updatedAt
		parseCropMetadata(metadata, &crop)
		items = append(items, crop)
	}
	return items, nil
}

func getFarmsByIDs(ctx context.Context, app *infra.Deps, farmIDs []string) ([]Farm, error) {
	if len(farmIDs) == 0 {
		return []Farm{}, nil
	}
	if err := ensureFarmDB(app); err != nil {
		return nil, err
	}

	filtered := make([]string, 0, len(farmIDs))
	for _, id := range farmIDs {
		if s := strings.TrimSpace(id); s != "" {
			filtered = append(filtered, s)
		}
	}
	if len(filtered) == 0 {
		return []Farm{}, nil
	}

	rows, err := app.SQLDB.Query(
		ctx,
		`SELECT farmid, name, description, userid, status, created_at, updated_at, metadata FROM `+farmsTable+` WHERE farmid = ANY($1) ORDER BY created_at DESC`,
		filtered,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	farms := make([]Farm, 0)
	for rows.Next() {
		var farm Farm
		var metadata []byte
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&farm.FarmID, &farm.Name, &farm.Description, &farm.CreatedBy, &farm.Owner, &createdAt, &updatedAt, &metadata); err != nil {
			continue
		}
		farm.CreatedAt = createdAt
		farm.UpdatedAt = updatedAt
		farm.Owner = farm.CreatedBy
		parseFarmMetadata(metadata, &farm)
		farms = append(farms, farm)
	}
	return farms, nil
}

func getCropsByNameFilter(ctx context.Context, app *infra.Deps, cropName string) ([]Crop, error) {
	cropName = strings.TrimSpace(cropName)
	if cropName == "" {
		return []Crop{}, nil
	}
	if err := ensureFarmDB(app); err != nil {
		return nil, err
	}

	rows, err := app.SQLDB.Query(
		ctx,
		`SELECT cropid, name, farmid, metadata, created_at, updated_at FROM `+cropsTable+` WHERE LOWER(name) LIKE $1 ORDER BY created_at DESC`,
		"%"+strings.ToLower(cropName)+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Crop, 0)
	for rows.Next() {
		var crop Crop
		var metadata []byte
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&crop.CropId, &crop.Name, &crop.FarmID, &metadata, &createdAt, &updatedAt); err != nil {
			continue
		}
		crop.CreatedAt = createdAt
		crop.UpdatedAt = updatedAt
		parseCropMetadata(metadata, &crop)
		items = append(items, crop)
	}
	return items, nil
}

func getPaginatedFarms(ctx context.Context, app *infra.Deps, search string, offset, limit int) ([]Farm, int64, error) {
	if err := ensureFarmDB(app); err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var total int64
	countQuery := `SELECT COUNT(*) FROM ` + farmsTable
	if search != "" {
		countQuery += ` WHERE LOWER(name) LIKE $1`
		if err := app.SQLDB.QueryRow(ctx, countQuery, "%"+strings.ToLower(search)+"%").Scan(&total); err != nil {
			return nil, 0, err
		}
	} else {
		if err := app.SQLDB.QueryRow(ctx, countQuery).Scan(&total); err != nil {
			return nil, 0, err
		}
	}

	query := `SELECT farmid, name, description, userid, status, created_at, updated_at, metadata FROM ` + farmsTable
	args := []any{}
	if search != "" {
		query += ` WHERE LOWER(name) LIKE $1`
		args = append(args, "%"+strings.ToLower(search)+"%")
	}
	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := app.SQLDB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]Farm, 0)
	for rows.Next() {
		var farm Farm
		var metadata []byte
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&farm.FarmID, &farm.Name, &farm.Description, &farm.CreatedBy, &farm.Owner, &createdAt, &updatedAt, &metadata); err != nil {
			continue
		}
		farm.CreatedAt = createdAt
		farm.UpdatedAt = updatedAt
		farm.Owner = farm.CreatedBy
		parseFarmMetadata(metadata, &farm)
		items = append(items, farm)
	}
	return items, total, nil
}

func getMyFarmsPage(ctx context.Context, app *infra.Deps, userID string, offset, limit int) ([]Farm, int64, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, 0, errors.New("user id is required")
	}
	if err := ensureFarmDB(app); err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	countQuery := `SELECT COUNT(*) FROM ` + farmsTable + ` WHERE userid = $1`
	var total int64
	if err := app.SQLDB.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := app.SQLDB.Query(
		ctx,
		`SELECT farmid, name, description, userid, status, created_at, updated_at, metadata FROM `+farmsTable+` WHERE userid = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]Farm, 0)
	for rows.Next() {
		var farm Farm
		var metadata []byte
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&farm.FarmID, &farm.Name, &farm.Description, &farm.CreatedBy, &farm.Owner, &createdAt, &updatedAt, &metadata); err != nil {
			continue
		}
		farm.CreatedAt = createdAt
		farm.UpdatedAt = updatedAt
		farm.Owner = farm.CreatedBy
		parseFarmMetadata(metadata, &farm)
		items = append(items, farm)
	}
	return items, total, nil
}

func updateOwnedFarm(ctx context.Context, app *infra.Deps, farmID, userID string, update map[string]any) (int64, error) {
	farmID = strings.TrimSpace(farmID)
	userID = strings.TrimSpace(userID)
	if farmID == "" || userID == "" {
		return 0, errors.New("farm id and user id are required")
	}
	if err := ensureFarmDB(app); err != nil {
		return 0, err
	}
	if len(update) == 0 {
		return 0, nil
	}

	setClause, setArgs := makeUpdateSet(update)
	if setClause == "" {
		return 0, nil
	}
	result, err := app.SQLDB.Exec(
		ctx,
		`UPDATE `+farmsTable+` SET `+setClause+`, updated_at = NOW() WHERE farmid = $`+fmt.Sprintf("%d", len(setArgs)+1)+` AND userid = $`+fmt.Sprintf("%d", len(setArgs)+2),
		append(append(setArgs, farmID), userID)...,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func deleteFarmByID(ctx context.Context, app *infra.Deps, farmID string) (int64, error) {
	farmID = strings.TrimSpace(farmID)
	if farmID == "" {
		return 0, errors.New("farm id is required")
	}
	if err := ensureFarmDB(app); err != nil {
		return 0, err
	}
	result, err := app.SQLDB.Exec(ctx, `DELETE FROM `+farmsTable+` WHERE farmid = $1`, farmID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
