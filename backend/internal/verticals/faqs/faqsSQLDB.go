// File: internal/verticals/faqs/faqsSQLDB.go

package faqs

import (
	"context"
	"strconv"
	"strings"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

var faqsTable = config.Tables.FAQsTable

func insertFAQ(ctx context.Context, app *infra.Deps, faq FAQ) error {
	return app.SQLDB.Insert(ctx, faqsTable, faq)
}

func findFAQByID(ctx context.Context, app *infra.Deps, faqID string, faq *FAQ) error {
	query := "faqid = $1"
	args := []any{faqID}

	return app.SQLDB.FindOne(ctx, faqsTable, query, args, faq)
}

func updateFAQContent(ctx context.Context, app *infra.Deps, faqID string, update map[string]any) (int64, error) {
	query := "faqid = $1"
	args := []any{faqID}

	return app.SQLDB.UpdateOne(ctx, faqsTable, query, args, update)
}

func deleteFAQ(ctx context.Context, app *infra.Deps, faqID, userID string) (int64, error) {
	query := "faqid = $1 AND createdby = $2"
	args := []any{faqID, userID}

	return app.SQLDB.Delete(ctx, faqsTable, query, args)
}

func findFAQsByEntity(
	ctx context.Context,
	app *infra.Deps,
	entityType string,
	entityID string,
	opts sqldb.FindManyOptions,
	faqs *[]FAQ,
) error {
	query := "entity_type = $1 AND entity_id = $2"
	args := []any{entityType, entityID}

	return app.SQLDB.FindManyWithOptions(ctx, faqsTable, query, args, opts, faqs)
}

func buildFAQListOptions(page, limit int, sortBy string) sqldb.FindManyOptions {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	skip := (page - 1) * limit
	orderBy := "created_at DESC, faqid DESC"
	switch sortBy {
	case "old":
		orderBy = "created_at ASC, faqid ASC"
	case "likes":
		orderBy = "likes DESC, created_at DESC"
	}

	return sqldb.FindManyOptions{
		Limit:   int64(limit),
		Offset:  int64(skip),
		OrderBy: orderBy,
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
