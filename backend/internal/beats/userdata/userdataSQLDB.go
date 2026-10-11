// File: internal/beats/userdata/userdataSQLDB.go

package userdata

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/utils"
)

var userdataTable = config.Tables.UserDataTable

func userDataValue(content UserData) map[string]any {
	return map[string]any{
		"entity_id":   content.EntityID,
		"entity_type": content.EntityType,
		"item_id":     content.ItemID,
		"item_type":   content.ItemType,
		"userid":      content.UserID,
		"created_at":  content.CreatedAt,
	}
}

func decodeUserDataValue(raw []byte) (UserData, error) {
	var payload map[string]any
	if len(raw) == 0 || string(raw) == "null" {
		return UserData{}, nil
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return UserData{}, err
	}
	content := UserData{}
	if v, ok := payload["userid"].(string); ok {
		content.UserID = v
	}
	if v, ok := payload["entity_id"].(string); ok {
		content.EntityID = v
	}
	if v, ok := payload["entity_type"].(string); ok {
		content.EntityType = v
	}
	if v, ok := payload["item_id"].(string); ok {
		content.ItemID = v
	}
	if v, ok := payload["item_type"].(string); ok {
		content.ItemType = v
	}
	if v, ok := payload["created_at"].(string); ok {
		content.CreatedAt = v
	}
	return content, nil
}

// SQLInsertUserData inserts a single user data document.
func InsertUserData(ctx context.Context, app *infra.Deps, content UserData) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if content.UserID == "" {
		return errors.New("userid required")
	}
	if content.EntityID == "" {
		return errors.New("entity_id required")
	}
	if content.EntityType == "" {
		return errors.New("entity_type required")
	}
	if content.CreatedAt == "" {
		content.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if content.ItemID == "" && content.ItemType == "" {
		content.ItemID = content.EntityID
	}

	payload, err := json.Marshal(userDataValue(content))
	if err != nil {
		return err
	}
	userdataID := utils.GetUUID()
	_, err = app.SQLDB.Exec(ctx,
		`INSERT INTO userdata (userdataid, userid, key_name, value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (userdataid) DO UPDATE SET
			userid = EXCLUDED.userid,
			key_name = EXCLUDED.key_name,
			value = EXCLUDED.value,
			updated_at = NOW()`,
		userdataID,
		content.UserID,
		content.EntityType+":"+content.EntityID,
		payload,
	)
	return err
}

// SQLDeleteUserData removes user data matching a SQL WHERE condition.
func DeleteUserData(ctx context.Context, app *infra.Deps, where string, args []any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if strings.TrimSpace(where) == "" {
		return 0, errors.New("where clause required")
	}
	result, err := app.SQLDB.Exec(ctx, `DELETE FROM userdata WHERE `+where, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// SQLInsertUserDataMany inserts many user data documents.
func InsertUserDataMany(ctx context.Context, app *infra.Deps, docs []any) error {
	if app == nil || app.SQLDB == nil || len(docs) == 0 {
		return nil
	}
	for _, raw := range docs {
		if doc, ok := raw.(UserData); ok {
			if err := InsertUserData(ctx, app, doc); err != nil {
				return err
			}
		}
	}
	return nil
}

// SQLFindUserData finds user data for a given SQL WHERE query and arguments.
func FindUserData(ctx context.Context, app *infra.Deps, where string, args []any, out *[]UserData) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if out == nil {
		return errors.New("nil result")
	}
	if strings.TrimSpace(where) == "" {
		where = "1 = 1"
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT userdataid, userid, key_name, value, created_at, updated_at FROM userdata WHERE `+where, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	results := make([]UserData, 0)
	for rows.Next() {
		var userdataID, userID, keyName string
		var value []byte
		var createdAt time.Time
		var updatedAt time.Time
		if err := rows.Scan(&userdataID, &userID, &keyName, &value, &createdAt, &updatedAt); err != nil {
			return err
		}
		decoded, err := decodeUserDataValue(value)
		if err != nil {
			continue
		}
		if decoded.UserID == "" {
			decoded.UserID = userID
		}
		if decoded.CreatedAt == "" {
			decoded.CreatedAt = createdAt.Format(time.RFC3339)
		}
		results = append(results, decoded)
	}
	*out = results
	return rows.Err()
}

// Database Helpers

// SQLFetchUserDataByEntity queries user data based on entity type and user ID.
func FetchUserDataByEntity(ctx context.Context, app *infra.Deps, entityType, username string) ([]UserData, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	results := make([]UserData, 0)
	rows, err := app.SQLDB.Query(ctx, `SELECT userdataid, userid, key_name, value, created_at, updated_at FROM userdata WHERE userid = $1 AND value->>'entity_type' = $2 ORDER BY created_at DESC`, username, entityType)
	if err != nil {
		return results, err
	}
	defer rows.Close()
	for rows.Next() {
		var userdataID, userID, keyName string
		var value []byte
		var createdAt time.Time
		var updatedAt time.Time
		if err := rows.Scan(&userdataID, &userID, &keyName, &value, &createdAt, &updatedAt); err != nil {
			return results, err
		}
		decoded, err := decodeUserDataValue(value)
		if err != nil {
			continue
		}
		decoded.UserID = userID
		if decoded.CreatedAt == "" {
			decoded.CreatedAt = createdAt.Format(time.RFC3339)
		}
		results = append(results, decoded)
	}
	return results, rows.Err()
}

// SQLFetchOtherUserFeedPosts retrieves posts for a given user from the database.
func FetchOtherUserFeedPosts(ctx context.Context, app *infra.Deps, username string) ([]postDoc, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	return []postDoc{}, nil
}
