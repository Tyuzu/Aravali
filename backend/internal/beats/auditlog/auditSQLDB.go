package auditlog

import (
	"context"
	"scav/config"
	"scav/infra"
)

var auditTable = config.Tables.AuditlogsTable

func SQLInsertAuditLog(ctx context.Context, app *infra.Deps, logEntry AuditLog) error {
	if app == nil || app.DB == nil {
		return nil
	}
	return app.DB.InsertOne(ctx, auditTable, logEntry)
}
