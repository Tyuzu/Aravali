package beats

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"scav/internal/beats/activity"
	"scav/internal/beats/ads"
	"scav/internal/beats/analytics"
	"scav/internal/beats/auditlog"
	"scav/internal/beats/autocomplete"
	"scav/internal/beats/follows"
	"scav/internal/beats/hashtags"
	"scav/internal/beats/notifications"
	"scav/internal/beats/suggestions"
	"scav/internal/beats/userdata"
)

// EnsureSchema runs all nested Beats schema bootstrap hooks.
func EnsureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return fmt.Errorf("nil postgres pool")
	}

	initializers := []func(context.Context, *pgxpool.Pool) error{
		activity.EnsureSchema,
		ads.EnsureSchema,
		analytics.EnsureSchema,
		auditlog.EnsureSchema,
		autocomplete.EnsureSchema,
		follows.EnsureSchema,
		hashtags.EnsureSchema,
		notifications.EnsureSchema,
		suggestions.EnsureSchema,
		userdata.EnsureSchema,
	}

	for _, initFn := range initializers {
		if err := initFn(ctx, pool); err != nil {
			return err
		}
	}

	return nil
}
