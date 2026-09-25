package faqs

import (
	"context"

	"scav/config"
	db "scav/infra/db"
)

var faqsTable = config.Tables.FAQsTable

func SQLinsertFAQ(ctx context.Context, database db.Database, faq FAQ) error {
	return database.Insert(ctx, faqsTable, faq)
}

func SQLfindFAQByID(ctx context.Context, database db.Database, faqID string, faq *FAQ) error {
	return database.FindOne(ctx, faqsTable, map[string]any{"faqid": faqID}, faq)
}

func SQLupdateFAQContent(ctx context.Context, database db.Database, faqID string, update map[string]any) (any, error) {
	return database.UpdateOne(ctx, faqsTable, map[string]any{"faqid": faqID}, update)
}

func SQLdeleteFAQ(ctx context.Context, database db.Database, faqID, userID string) (int64, error) {
	return database.Delete(ctx, faqsTable, map[string]any{"faqid": faqID, "createdby": userID})
}

func SQLfindFAQsByEntity(
	ctx context.Context,
	database db.Database,
	entityType string,
	entityID string,
	opts db.FindManyOptions,
	faqs *[]FAQ,
) error {
	filter := map[string]any{
		"entity_type": entityType,
		"entity_id":   entityID,
	}
	return database.FindManyWithOptions(ctx, faqsTable, filter, opts, faqs)
}
