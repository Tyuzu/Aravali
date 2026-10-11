// File: internal/faqs/faqsSQLDB.go

package faqs

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
)

var faqsTable = config.Tables.FAQsTable

func insertFAQ(ctx context.Context, app *infra.Deps, faq FAQ) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if strings.TrimSpace(faq.FAQID) == "" {
		faq.FAQID = fmt.Sprintf("faq_%d", time.Now().UnixNano())
	}
	if faq.CreatedAt.IsZero() {
		faq.CreatedAt = time.Now()
	}
	if faq.UpdatedAt.IsZero() {
		faq.UpdatedAt = faq.CreatedAt
	}

	_, err := app.SQLDB.Exec(
		ctx,
		`INSERT INTO `+faqsTable+` (faqid, entity_type, entity_id, content, created_by, likes, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		faq.FAQID,
		faq.EntityType,
		faq.EntityID,
		faq.Content,
		faq.CreatedBy,
		faq.Likes,
		faq.CreatedAt,
		faq.UpdatedAt,
	)
	return err
}

func findFAQByID(ctx context.Context, app *infra.Deps, faqID string, faq *FAQ) error {
	if app == nil || app.SQLDB == nil {
		return fmt.Errorf("faq %s not found", faqID)
	}
	if faq == nil {
		return fmt.Errorf("faq pointer is nil")
	}
	if strings.TrimSpace(faqID) == "" {
		return fmt.Errorf("faq id is required")
	}

	var createdAt, updatedAt time.Time
	row := app.SQLDB.QueryRow(
		ctx,
		`SELECT faqid, entity_type, entity_id, content, created_by, likes, created_at, updated_at FROM `+faqsTable+` WHERE faqid = $1 LIMIT 1`,
		faqID,
	)
	if err := row.Scan(&faq.FAQID, &faq.EntityType, &faq.EntityID, &faq.Content, &faq.CreatedBy, &faq.Likes, &createdAt, &updatedAt); err != nil {
		return err
	}
	faq.CreatedAt = createdAt
	faq.UpdatedAt = updatedAt
	return nil
}

func updateFAQContent(ctx context.Context, app *infra.Deps, faqID string, update map[string]any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if strings.TrimSpace(faqID) == "" {
		return 0, fmt.Errorf("faq id is required")
	}
	if len(update) == 0 {
		return 0, nil
	}

	parts := make([]string, 0, len(update))
	args := make([]any, 0, len(update)+1)
	for key, value := range update {
		if key == "" || value == nil {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value)
	}
	if len(parts) == 0 {
		return 0, nil
	}
	args = append(args, faqID)
	res, err := app.SQLDB.Exec(
		ctx,
		`UPDATE `+faqsTable+` SET `+strings.Join(parts, ", ")+` WHERE faqid = $`+fmt.Sprintf("%d", len(args)),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func deleteFAQ(ctx context.Context, app *infra.Deps, faqID, userID string) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if strings.TrimSpace(faqID) == "" {
		return 0, fmt.Errorf("faq id is required")
	}

	if strings.TrimSpace(userID) == "" {
		return 0, nil
	}
	res, err := app.SQLDB.Exec(ctx, `DELETE FROM `+faqsTable+` WHERE faqid = $1 AND created_by = $2`, faqID, userID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func findFAQsByEntity(
	ctx context.Context,
	app *infra.Deps,
	entityType string,
	entityID string,
	opts map[string]any,
	faqs *[]FAQ,
) error {
	if app == nil || app.SQLDB == nil || faqs == nil {
		return nil
	}
	if strings.TrimSpace(entityType) == "" || strings.TrimSpace(entityID) == "" {
		return fmt.Errorf("entity type and entity id are required")
	}

	sortBy := "created_at DESC"
	if v, ok := opts["sortBy"].(string); ok && strings.TrimSpace(v) != "" {
		sortBy = v
	}
	limit := 10
	if v, ok := opts["limit"].(int); ok && v > 0 {
		limit = v
	}
	page := 1
	if v, ok := opts["page"].(int); ok && v > 0 {
		page = v
	}
	offset := (page - 1) * limit

	rows, err := app.SQLDB.Query(
		ctx,
		`SELECT faqid, entity_type, entity_id, content, created_by, likes, created_at, updated_at FROM `+faqsTable+` WHERE entity_type = $1 AND entity_id = $2 ORDER BY `+sortBy+` LIMIT $3 OFFSET $4`,
		entityType, entityID, limit, offset,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	items := make([]FAQ, 0)
	for rows.Next() {
		var item FAQ
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&item.FAQID, &item.EntityType, &item.EntityID, &item.Content, &item.CreatedBy, &item.Likes, &createdAt, &updatedAt); err != nil {
			return err
		}
		item.CreatedAt = createdAt
		item.UpdatedAt = updatedAt
		items = append(items, item)
	}
	*faqs = items
	return rows.Err()
}

func buildFAQListOptions(page, limit int, sortBy string) map[string]any {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	if sortBy == "" {
		sortBy = "created_at DESC"
	}
	return map[string]any{
		"page":   page,
		"limit":  limit,
		"sortBy": sortBy,
	}
}

func getFAQPage(ctx context.Context, app *infra.Deps, entityType, entityID string, page, limit int, sortBy string) ([]FAQ, error) {
	var faqs []FAQ
	opts := buildFAQListOptions(page, limit, sortBy)
	if err := findFAQsByEntity(ctx, app, entityType, entityID, opts, &faqs); err != nil {
		return nil, err
	}
	if faqs == nil {
		return []FAQ{}, nil
	}
	return faqs, nil
}

func parseFAQPageLimit(rawPage, rawLimit string) (int, int) {
	page := 1
	if v := strings.TrimSpace(rawPage); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if v := strings.TrimSpace(rawLimit); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	return page, limit
}
