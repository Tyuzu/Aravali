package auditlog

import (
	"context"
	"scav/infra"
)

func InsertAuditLog(ctx context.Context, app *infra.Deps, logEntry AuditLog) error {
	if app == nil || app.DB == nil {
		return nil
	}
	return app.DB.InsertOne(ctx, auditCollection, logEntry)
}
