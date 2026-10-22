// File: internal/crops/cropsSQLDB.go

package crops

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/internal/farms"
)

var (
	cropsTable      = config.Tables.CropsTable
	cropsAboutTable = config.Tables.CropsAboutTable
	catalogueTable  = config.Tables.CatalogueTable
)

func _cropMetadataMap(crop farms.Crop) map[string]any {
	meta := map[string]any{
		"name":         crop.Name,
		"farmid":       crop.FarmID,
		"farmName":     crop.FarmName,
		"category":     crop.Category,
		"price":        crop.Price,
		"discount":     crop.Discount,
		"quantity":     crop.Quantity,
		"unit":         crop.Unit,
		"notes":        crop.Notes,
		"featured":     crop.Featured,
		"out_of_stock": crop.OutOfStock,
		"createdby":    crop.CreatedBy,
		"owner":        crop.CreatedBy,
	}
	if crop.HarvestDate != nil {
		meta["harvestDate"] = crop.HarvestDate.Format(time.RFC3339)
	}
	if crop.ExpiryDate != nil {
		meta["expiryDate"] = crop.ExpiryDate.Format(time.RFC3339)
	}
	if crop.Banner != "" {
		meta["banner"] = crop.Banner
	}
	if crop.CatalogueId != "" {
		meta["catalogueid"] = crop.CatalogueId
	}
	return meta
}

func _decodeCropMetadata(raw []byte, dst *farms.Crop) {
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
	if v, ok := meta["farmName"]; ok {
		if s, ok2 := v.(string); ok2 {
			dst.FarmName = s
		}
	}
	if v, ok := meta["createdby"]; ok {
		if s, ok2 := v.(string); ok2 {
			dst.CreatedBy = s
		}
	}
	if v, ok := meta["owner"]; ok {
		if s, ok2 := v.(string); ok2 {
			dst.CreatedBy = s
		}
	}
}

func _insertCrop(ctx context.Context, app *infra.Deps, crop farms.Crop) error {
	if app == nil || app.SQLDB == nil {
		return errors.New("database unavailable")
	}
	if strings.TrimSpace(crop.CropId) == "" {
		return errors.New("crop id is required")
	}
	if strings.TrimSpace(crop.FarmID) == "" {
		return errors.New("farm id is required")
	}

	meta, err := json.Marshal(_cropMetadataMap(crop))
	if err != nil {
		return err
	}

	_, err = app.SQLDB.Exec(
		ctx,
		`INSERT INTO `+cropsTable+` (cropid, name, description, farmid, status, created_at, updated_at, metadata) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (cropid) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, farmid = EXCLUDED.farmid, status = EXCLUDED.status, updated_at = NOW(), metadata = EXCLUDED.metadata`,
		crop.CropId,
		crop.Name,
		"",
		crop.FarmID,
		"active",
		crop.CreatedAt,
		time.Now(),
		meta,
	)
	return err
}

func _updateCrop(ctx context.Context, app *infra.Deps, cropID string, update map[string]any) error {
	if app == nil || app.SQLDB == nil {
		return errors.New("database unavailable")
	}
	if strings.TrimSpace(cropID) == "" {
		return errors.New("crop id is required")
	}
	if len(update) == 0 {
		return nil
	}

	setParts := make([]string, 0, len(update))
	args := make([]any, 0, len(update))
	for key, value := range update {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value)
	}

	if len(setParts) == 0 {
		return nil
	}
	_, err := app.SQLDB.Exec(
		ctx,
		`UPDATE `+cropsTable+` SET `+strings.Join(setParts, ", ")+` WHERE cropid = $`+fmt.Sprintf("%d", len(args)+1),
		append(args, cropID)...,
	)
	return err
}

func _findFilteredCrops(ctx context.Context, app *infra.Deps, query string, args []any) ([]farms.Crop, error) {
	if app == nil || app.SQLDB == nil {
		return nil, errors.New("database unavailable")
	}
	if strings.TrimSpace(query) == "" {
		query = "1 = 1"
	}

	rows, err := app.SQLDB.Query(
		ctx,
		`SELECT cropid, name, farmid, metadata, created_at, updated_at FROM `+cropsTable+` WHERE `+query+` ORDER BY created_at DESC`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]farms.Crop, 0)
	for rows.Next() {
		var crop farms.Crop
		var metadata []byte
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&crop.CropId, &crop.Name, &crop.FarmID, &metadata, &createdAt, &updatedAt); err != nil {
			continue
		}
		crop.CreatedAt = createdAt
		crop.UpdatedAt = updatedAt
		_decodeCropMetadata(metadata, &crop)
		items = append(items, crop)
	}
	return items, rows.Err()
}

func _findCatalogueItems(ctx context.Context, app *infra.Deps, query string, args []any) ([]farms.CropCatalogueItem, error) {
	if app == nil || app.SQLDB == nil {
		return nil, errors.New("database unavailable")
	}
	if strings.TrimSpace(query) == "" {
		query = "1 = 1"
	}

	rows, err := app.SQLDB.Query(
		ctx,
		`SELECT catalogueid, name, description, metadata, created_at, updated_at FROM `+catalogueTable+` WHERE `+query+` ORDER BY created_at DESC`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]farms.CropCatalogueItem, 0)
	for rows.Next() {
		var id, name, description string
		var metadata []byte
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&id, &name, &description, &metadata, &createdAt, &updatedAt); err != nil {
			continue
		}
		item := farms.CropCatalogueItem{
			Name:       name,
			Category:   description,
			Banner:     "",
			Featured:   false,
			PriceRange: []int{},
		}
		if len(metadata) > 0 {
			var raw map[string]any
			if err := json.Unmarshal(metadata, &raw); err == nil {
				if v, ok := raw["banner"].(string); ok {
					item.Banner = v
				}
				if v, ok := raw["category"].(string); ok {
					item.Category = v
				}
				if v, ok := raw["featured"].(bool); ok {
					item.Featured = v
				}
				if v, ok := raw["priceRange"].([]any); ok {
					item.PriceRange = make([]int, 0, len(v))
					for _, entry := range v {
						if n, ok := entry.(float64); ok {
							item.PriceRange = append(item.PriceRange, int(n))
						}
					}
				}
			}
		}
		_ = id
		_ = createdAt
		_ = updatedAt
		items = append(items, item)
	}
	return items, rows.Err()
}

func _buildCropFilterQuery(filter map[string]any) (string, []any) {
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
			if gt, ok := v["$gt"]; ok {
				clauses = append(clauses, fmt.Sprintf("%s > $%d", key, len(args)+1))
				args = append(args, gt)
				continue
			}
			if gte, ok := v["$gte"]; ok {
				clauses = append(clauses, fmt.Sprintf("%s >= $%d", key, len(args)+1))
				args = append(args, gte)
				continue
			}
			if lte, ok := v["$lte"]; ok {
				clauses = append(clauses, fmt.Sprintf("%s <= $%d", key, len(args)+1))
				args = append(args, lte)
				continue
			}
		}
		clauses = append(clauses, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value)
	}
	if len(clauses) == 0 {
		return "1 = 1", nil
	}
	return strings.Join(clauses, " AND "), args
}

func _getAllCrops(ctx context.Context, app *infra.Deps) ([]farms.Crop, error) {
	return _findFilteredCrops(ctx, app, "1 = 1", nil)
}

func _createCropAbout(ctx context.Context, app *infra.Deps, crop *CropAbout) error {
	if app == nil || app.SQLDB == nil {
		return errors.New("database unavailable")
	}
	if crop == nil {
		return errors.New("crop about is required")
	}
	if strings.TrimSpace(crop.ID) == "" {
		crop.ID = "cropabout_" + time.Now().Format("20060102150405")
	}

	metadata, err := json.Marshal(map[string]any{
		"id":                 crop.ID,
		"commonName":         crop.CommonName,
		"scientificName":     crop.ScientificName,
		"image":              crop.Image,
		"imageAlt":           crop.ImageAlt,
		"description":        crop.Description,
		"nutritionalValues":  crop.NutritionalValues,
		"growingConditions":  crop.GrowingConditions,
		"plantingHarvesting": crop.PlantingHarvesting,
		"careTips":           crop.CareTips,
		"varieties":          crop.Varieties,
		"usage":              crop.Usage,
		"funFacts":           crop.FunFacts,
	})
	if err != nil {
		return err
	}

	_, err = app.SQLDB.Exec(
		ctx,
		`INSERT INTO `+cropsAboutTable+` (cropsaboutid, cropid, title, content, created_at, updated_at, metadata) VALUES ($1, $2, $3, $4, NOW(), NOW(), $5) ON CONFLICT (cropsaboutid) DO UPDATE SET cropid = EXCLUDED.cropid, title = EXCLUDED.title, content = EXCLUDED.content, updated_at = NOW(), metadata = EXCLUDED.metadata`,
		crop.ID,
		crop.ID,
		crop.CommonName,
		crop.Description,
		metadata,
	)
	return err
}

func _getCropAboutByID(ctx context.Context, app *infra.Deps, cropID string) (*CropAbout, error) {
	if app == nil || app.SQLDB == nil {
		return nil, errors.New("database unavailable")
	}
	cropID = strings.TrimSpace(cropID)
	if cropID == "" {
		return nil, errors.New("crop id is required")
	}

	var about CropAbout
	var cropAboutID, dbCropID string
	var title, content string
	var metadata []byte
	var createdAt, updatedAt time.Time

	row := app.SQLDB.QueryRow(
		ctx,
		`SELECT cropsaboutid, cropid, title, content, created_at, updated_at, metadata FROM `+cropsAboutTable+` WHERE cropid = $1 OR cropsaboutid = $1 LIMIT 1`,
		cropID,
	)
	if err := row.Scan(&cropAboutID, &dbCropID, &title, &content, &createdAt, &updatedAt, &metadata); err != nil {
		return nil, err
	}
	about.ID = cropAboutID
	if about.ID == "" {
		about.ID = dbCropID
	}
	about.CommonName = title
	about.Description = content
	about.ID = strings.TrimSpace(about.ID)
	if about.ID == "" {
		about.ID = cropID
	}
	if len(metadata) > 0 {
		var raw map[string]any
		if err := json.Unmarshal(metadata, &raw); err == nil {
			if v, ok := raw["commonName"].(string); ok && about.CommonName == "" {
				about.CommonName = v
			}
			if v, ok := raw["scientificName"].(string); ok {
				about.ScientificName = v
			}
			if v, ok := raw["image"].(string); ok {
				about.Image = v
			}
			if v, ok := raw["imageAlt"].(string); ok {
				about.ImageAlt = v
			}
			if v, ok := raw["description"].(string); ok && about.Description == "" {
				about.Description = v
			}
			if v, ok := raw["plantingHarvesting"].(string); ok {
				about.PlantingHarvesting = v
			}
			if v, ok := raw["usage"].(string); ok {
				about.Usage = v
			}
			if v, ok := raw["careTips"].([]any); ok {
				about.CareTips = make([]string, 0, len(v))
				for _, item := range v {
					if s, ok := item.(string); ok {
						about.CareTips = append(about.CareTips, s)
					}
				}
			}
			if v, ok := raw["varieties"].([]any); ok {
				about.Varieties = make([]string, 0, len(v))
				for _, item := range v {
					if s, ok := item.(string); ok {
						about.Varieties = append(about.Varieties, s)
					}
				}
			}
			if v, ok := raw["funFacts"].([]any); ok {
				about.FunFacts = make([]string, 0, len(v))
				for _, item := range v {
					if s, ok := item.(string); ok {
						about.FunFacts = append(about.FunFacts, s)
					}
				}
			}
			if rawNutritional, ok := raw["nutritionalValues"].([]any); ok {
				about.NutritionalValues = make([]NutritionalValue, 0, len(rawNutritional))
				for _, item := range rawNutritional {
					entryMap, ok := item.(map[string]any)
					if !ok {
						continue
					}
					about.NutritionalValues = append(about.NutritionalValues, NutritionalValue{
						Label: fmt.Sprint(entryMap["label"]),
						Value: fmt.Sprint(entryMap["value"]),
					})
				}
			}
		}
	}
	_ = createdAt
	_ = updatedAt
	return &about, nil
}

func _getAllCropAbouts(ctx context.Context, app *infra.Deps) ([]CropAbout, error) {
	if app == nil || app.SQLDB == nil {
		return nil, errors.New("database unavailable")
	}

	rows, err := app.SQLDB.Query(
		ctx,
		`SELECT cropsaboutid, cropid, title, content, created_at, updated_at, metadata FROM `+cropsAboutTable+` ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]CropAbout, 0)
	for rows.Next() {
		var about CropAbout
		var cropAboutID, cropID string
		var title, content string
		var metadata []byte
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&cropAboutID, &cropID, &title, &content, &createdAt, &updatedAt, &metadata); err != nil {
			continue
		}
		about.ID = cropAboutID
		if about.ID == "" {
			about.ID = cropID
		}
		about.CommonName = title
		about.Description = content
		if len(metadata) > 0 {
			var raw map[string]any
			if err := json.Unmarshal(metadata, &raw); err == nil {
				if v, ok := raw["commonName"].(string); ok {
					about.CommonName = v
				}
				if v, ok := raw["description"].(string); ok {
					about.Description = v
				}
				if v, ok := raw["scientificName"].(string); ok {
					about.ScientificName = v
				}
				if v, ok := raw["image"].(string); ok {
					about.Image = v
				}
				if v, ok := raw["imageAlt"].(string); ok {
					about.ImageAlt = v
				}
			}
		}
		_ = createdAt
		_ = updatedAt
		items = append(items, about)
	}
	return items, rows.Err()
}

func _updateCropAbout(ctx context.Context, app *infra.Deps, cropID string, crop *CropAbout) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, errors.New("database unavailable")
	}
	if crop == nil {
		return 0, errors.New("crop about is required")
	}
	cropID = strings.TrimSpace(cropID)
	if cropID == "" {
		return 0, errors.New("crop id is required")
	}

	metadata, err := json.Marshal(map[string]any{
		"id":                 crop.ID,
		"commonName":         crop.CommonName,
		"scientificName":     crop.ScientificName,
		"image":              crop.Image,
		"imageAlt":           crop.ImageAlt,
		"description":        crop.Description,
		"nutritionalValues":  crop.NutritionalValues,
		"growingConditions":  crop.GrowingConditions,
		"plantingHarvesting": crop.PlantingHarvesting,
		"careTips":           crop.CareTips,
		"varieties":          crop.Varieties,
		"usage":              crop.Usage,
		"funFacts":           crop.FunFacts,
	})
	if err != nil {
		return 0, err
	}

	result, err := app.SQLDB.Exec(
		ctx,
		`UPDATE `+cropsAboutTable+` SET cropsaboutid = COALESCE(NULLIF($1, ''), cropsaboutid), cropid = COALESCE(NULLIF($2, ''), cropid), title = COALESCE(NULLIF($3, ''), title), content = COALESCE(NULLIF($4, ''), content), updated_at = NOW(), metadata = $5 WHERE cropid = $6 OR cropsaboutid = $6`,
		crop.ID,
		cropID,
		crop.CommonName,
		crop.Description,
		metadata,
		cropID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func _deleteCropAbout(ctx context.Context, app *infra.Deps, cropID string) error {
	if app == nil || app.SQLDB == nil {
		return errors.New("database unavailable")
	}
	cropID = strings.TrimSpace(cropID)
	if cropID == "" {
		return errors.New("crop id is required")
	}
	_, err := app.SQLDB.Exec(
		ctx,
		`DELETE FROM `+cropsAboutTable+` WHERE cropid = $1 OR cropsaboutid = $1`,
		cropID,
	)
	return err
}
