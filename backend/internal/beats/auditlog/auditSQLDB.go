package auditlog

import (
	"context"
	"scav/config"
	"scav/infra"
)

var auditTable = config.Tables.AuditlogsTable

func SQLInsertAuditLog(ctx context.Context, app *infra.Deps, logEntry AuditLog) error {
	return app.SQLDB.InsertOne(ctx, auditTable, logEntry)
}
