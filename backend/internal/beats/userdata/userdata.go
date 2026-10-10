// File: internal/beats/userdata/userdata.go

package userdata

import (
	"context"
	"scav/infra"
	log "scav/infra/logger"
	"time"
)

var ValidEntityTypes = map[string]bool{
	"userhome":  true,
	"place":     true,
	"event":     true,
	"feedpost":  true,
	"media":     true,
	"ticket":    true,
	"merch":     true,
	"review":    true,
	"comment":   true,
	"like":      true,
	"favourite": true,
	"booking":   true,
	"blogpost":  true,
}

func IsValidEntityType(entityType string) bool {
	return ValidEntityTypes[entityType]
}

func SetUserData(dataType, dataId, userId, itemType, itemId string, app *infra.Deps) {
	AddUserData(dataType, dataId, userId, itemType, itemId, app)
}

func DelUserData(dataType, dataId, userId string, app *infra.Deps) {
	RemUserData(dataType, dataId, userId, app)
}

func AddUserData(entityType, entityId, userId, itemType, itemId string, app *infra.Deps) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	content := UserData{
		EntityID:   entityId,
		EntityType: entityType,
		ItemID:     itemId,
		ItemType:   itemType,
		UserID:     userId,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	if err := InsertUserData(ctx, app, content); err != nil {
		log.Printf("Error inserting user data: %v", err)
	}
}

func RemUserData(entityType, entityId, userId string, app *infra.Deps) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	where := "entity_type = $1 AND entity_id = $2 AND userid = $3"
	args := []any{entityType, entityId, userId}

	if _, err := DeleteUserData(ctx, app, where, args); err != nil {
		log.Printf("Error deleting user data: %v", err)
	}
}

func AddUserDataBatch(docs []UserData, app *infra.Deps) {
	if len(docs) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var toInsert []any
	for _, doc := range docs {
		toInsert = append(toInsert, doc)
	}

	if err := InsertUserDataMany(ctx, app, toInsert); err != nil {
		log.Printf("Error inserting batch user data: %v", err)
	}
}
