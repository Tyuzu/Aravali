// File: internal/verticals/faqs/faqsSQLDB.go

package faqs

import (
	"context"
	"strconv"
	"strings"

	"scav/config"
	"scav/infra"
)

var faqsTable = config.Tables.FAQsTable

func insertFAQ(ctx context.Context, app *infra.Deps, faq FAQ) error {
}

func findFAQByID(ctx context.Context, app *infra.Deps, faqID string, faq *FAQ) error {
}

func updateFAQContent(ctx context.Context, app *infra.Deps, faqID string, update map[string]any) (int64, error) {
}

func deleteFAQ(ctx context.Context, app *infra.Deps, faqID, userID string) (int64, error) {
}

func findFAQsByEntity(
	ctx context.Context,
	app *infra.Deps,
	entityType string,
	entityID string,
	opts map[string]any,
	faqs *[]FAQ,
) error {
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
