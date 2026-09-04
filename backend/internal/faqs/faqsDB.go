package faqs

import (
	"context"

	"scav/config"
	db "scav/infra/db"
)

var faqsCollection = config.Collections.FAQsCollection

func insertFAQ(ctx context.Context, database db.Database, faq FAQ) error {
	return database.Insert(ctx, faqsCollection, faq)
}

func findFAQByID(ctx context.Context, database db.Database, faqID string, faq *FAQ) error {
	return database.FindOne(ctx, faqsCollection, map[string]any{"faqid": faqID}, faq)
}

func updateFAQContent(ctx context.Context, database db.Database, faqID string, update map[string]any) (any, error) {
	return database.UpdateOne(ctx, faqsCollection, map[string]any{"faqid": faqID}, update)
}

func deleteFAQ(ctx context.Context, database db.Database, faqID, userID string) (int64, error) {
	return database.Delete(ctx, faqsCollection, map[string]any{"faqid": faqID, "createdby": userID})
}

func findFAQsByEntity(
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
	return database.FindManyWithOptions(ctx, faqsCollection, filter, opts, faqs)
}
