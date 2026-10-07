package bootstrap

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"scav/internal/admin"
	"scav/internal/artists"
	"scav/internal/auth"
	"scav/internal/baito"
	"scav/internal/beats"
	"scav/internal/cart"
	"scav/internal/crops"
	"scav/internal/deliveries"
	"scav/internal/events"
	"scav/internal/farms"
	"scav/internal/home"
	"scav/internal/itinerary"
	"scav/internal/maps"
	"scav/internal/mechat"
	"scav/internal/newchat"
	"scav/internal/places"
	"scav/internal/posts"
	"scav/internal/products"
	"scav/internal/profile"
	"scav/internal/recipes"
	"scav/internal/reports"
	"scav/internal/search"
	"scav/internal/settings"
	"scav/internal/vendors"
	"scav/internal/verticals"
	"scav/internal/workers"
)

func EnsureModuleSchemas(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return fmt.Errorf("nil postgres pool")
	}

	initializers := []func(context.Context, *pgxpool.Pool) error{
		admin.EnsureSchema,
		artists.EnsureSchema,
		auth.EnsureUserColumns,
		baito.EnsureSchema,
		beats.EnsureSchema,
		cart.EnsureSchema,
		crops.EnsureSchema,
		deliveries.EnsureSchema,
		events.EnsureSchema,
		farms.EnsureFarmColumns,
		home.EnsureSchema,
		itinerary.EnsureSchema,
		maps.EnsureSchema,
		mechat.EnsureSchema,
		newchat.EnsureSchema,
		places.EnsureSchema,
		posts.EnsureSchema,
		products.EnsureSchema,
		profile.EnsureSchema,
		recipes.EnsureSchema,
		reports.EnsureSchema,
		search.EnsureSchema,
		settings.EnsureSchema,
		vendors.EnsureSchema,
		verticals.EnsureSchema,
		workers.EnsureSchema,
	}

	for _, initFn := range initializers {
		if err := initFn(ctx, pool); err != nil {
			return err
		}
	}

	return nil
}
