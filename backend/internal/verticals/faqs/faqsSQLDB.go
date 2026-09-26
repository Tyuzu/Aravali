package faqs

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

var faqsTable = config.Tables.FAQsTable

func SQLinsertFAQ(ctx context.Context, app *infra.Deps, faq FAQ) error {
	return app.SQLDB.Insert(ctx, faqsTable, faq)
}

func SQLfindFAQByID(ctx context.Context, app *infra.Deps, faqID string, faq *FAQ) error {
	query := "faqid = $1"
	args := []any{faqID}

	return app.SQLDB.FindOne(ctx, faqsTable, query, args, faq)
}

func SQLupdateFAQContent(ctx context.Context, app *infra.Deps, faqID string, update map[string]any) (int64, error) {
	query := "faqid = $1"
	args := []any{faqID}

	return app.SQLDB.UpdateOne(ctx, faqsTable, query, args, update)
}

func SQLdeleteFAQ(ctx context.Context, app *infra.Deps, faqID, userID string) (int64, error) {
	query := "faqid = $1 AND createdby = $2"
	args := []any{faqID, userID}

	return app.SQLDB.Delete(ctx, faqsTable, query, args)
}

func SQLfindFAQsByEntity(
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
