package verticals

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"scav/internal/verticals/booking"
	"scav/internal/verticals/comments"
	"scav/internal/verticals/faqs"
	"scav/internal/verticals/filemgr"
	"scav/internal/verticals/media"
	"scav/internal/verticals/menu"
	"scav/internal/verticals/merch"
	"scav/internal/verticals/notices"
	"scav/internal/verticals/pay"
	"scav/internal/verticals/reviews"
	"scav/internal/verticals/tickets"
)

// EnsureSchema runs all nested Verticals schema bootstrap hooks.
func EnsureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return fmt.Errorf("nil postgres pool")
	}

	initializers := []func(context.Context, *pgxpool.Pool) error{
		booking.EnsureSchema,
		comments.EnsureSchema,
		faqs.EnsureSchema,
		filemgr.EnsureSchema,
		media.EnsureSchema,
		menu.EnsureSchema,
		merch.EnsureSchema,
		notices.EnsureSchema,
		pay.EnsureSchema,
		reviews.EnsureSchema,
		tickets.EnsureSchema,
	}

	for _, initFn := range initializers {
		if err := initFn(ctx, pool); err != nil {
			return err
		}
	}

	return nil
}
