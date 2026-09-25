package auditlog

import (
	"context"
	"scav/config"
	"scav/infra"
)

var auditCollection = config.Collections.AuditlogsCollection

func InsertAuditLog(ctx context.Context, app *infra.Deps, logEntry AuditLog) error {
	if app == nil || app.DB == nil {
		return nil
	}
	return app.DB.InsertOne(ctx, auditCollection, logEntry)
}
