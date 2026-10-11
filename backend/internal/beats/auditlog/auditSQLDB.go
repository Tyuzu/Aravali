// File: internal/beats/auditlog/auditSQLDB.go

package auditlog

import (
	"context"
	"encoding/json"
	"time"

	"scav/config"
	"scav/infra"
	"scav/utils"
)

var auditTable = config.Tables.AuditlogsTable

func InsertAuditLog(ctx context.Context, app *infra.Deps, logEntry AuditLog) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if logEntry.ID == "" {
		logEntry.ID = utils.GetUUID()
	}
	if logEntry.CreatedAt.IsZero() {
		logEntry.CreatedAt = time.Now().UTC()
	}

	payload, err := json.Marshal(map[string]any{
		"userID":     logEntry.UserID,
		"action":     logEntry.Action,
		"entityType": logEntry.EntityType,
		"entityID":   logEntry.EntityID,
		"status":     logEntry.Status,
		"reason":     logEntry.Reason,
		"changes":    logEntry.Changes,
		"ipAddress":  logEntry.IPAddress,
		"userAgent":  logEntry.UserAgent,
	})
	if err != nil {
		return err
	}

	_, err = app.SQLDB.Exec(ctx,
		`INSERT INTO audit_logs (id, actorid, action, entity_type, entity_id, payload, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			actorid = EXCLUDED.actorid,
			action = EXCLUDED.action,
			entity_type = EXCLUDED.entity_type,
			entity_id = EXCLUDED.entity_id,
			payload = EXCLUDED.payload`,
		logEntry.ID,
		logEntry.UserID,
		logEntry.Action,
		logEntry.EntityType,
		logEntry.EntityID,
		payload,
		logEntry.CreatedAt,
	)
	return err
}
