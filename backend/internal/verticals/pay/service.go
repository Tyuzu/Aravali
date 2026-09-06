package pay

import (
	"context"
	"errors"

	"scav/infra"
)

// ===== Price Resolver =====

type PriceResolver func(ctx context.Context, entityID string) (int64, error)

// ===== Payment Service =====

type PaymentService struct {
	app       *infra.Deps
	resolvers map[string]PriceResolver
}

// Constructor
func NewPaymentService(app *infra.Deps) *PaymentService {
	return &PaymentService{
		app:       app,
		resolvers: make(map[string]PriceResolver),
	}
}

// Register resolver
func (p *PaymentService) RegisterResolver(entityType string, r PriceResolver) {
	p.resolvers[entityType] = r
}

func (p *PaymentService) resolver(entityType string) (PriceResolver, error) {
	r, ok := p.resolvers[entityType]
	if !ok {
		return nil, errors.New("unsupported entity type")
	}
	return r, nil
}

// ===== Default Resolvers =====

func (p *PaymentService) RegisterDefaultResolvers() {
	p.RegisterResolver("ticket", func(ctx context.Context, id string) (int64, error) {
		return p.fetchPriceByField(ctx, ticketsCollection, "ticketid", id)
	})

	p.RegisterResolver("menu", func(ctx context.Context, id string) (int64, error) {
		return p.fetchPriceByField(ctx, menuCollection, "menuid", id)
	})

	p.RegisterResolver("service", func(ctx context.Context, id string) (int64, error) {
		return p.fetchPriceByField(ctx, serviceCollection, "serviceid", id)
	})

	// donations / tips
	p.RegisterResolver("post", func(ctx context.Context, id string) (int64, error) {
		return 0, nil
	})

	// orders - fetch total from order or farmOrders
	p.RegisterResolver("order", func(ctx context.Context, id string) (int64, error) {
		return p.findOrderTotalByID(ctx, id)
	})

	// cart - custom entity, no fixed price
	p.RegisterResolver("cart", func(ctx context.Context, id string) (int64, error) {
		return 0, nil
	})

	// product - treat like menu item
	p.RegisterResolver("product", func(ctx context.Context, id string) (int64, error) {
		return p.fetchPriceByField(ctx, productCollection, "productid", id)
	})

	// booking - has a price
	p.RegisterResolver("booking", func(ctx context.Context, id string) (int64, error) {
		return p.fetchPriceByField(ctx, bookingsCollection, "bookingid", id)
	})

	// merch - has a price
	p.RegisterResolver("merch", func(ctx context.Context, id string) (int64, error) {
		return p.fetchPriceByField(ctx, merchCollection, "merchid", id)
	})

	// crop - has a price
	p.RegisterResolver("crop", func(ctx context.Context, id string) (int64, error) {
		return p.fetchPriceByField(ctx, cropsCollection, "cropid", id)
	})

	// farm - custom entity
	p.RegisterResolver("farm", func(ctx context.Context, id string) (int64, error) {
		return 0, nil
	})

	// beat - has a price
	p.RegisterResolver("beat", func(ctx context.Context, id string) (int64, error) {
		return p.fetchPriceByField(ctx, "beats", "beatid", id)
	})

	// donation - custom amount
	p.RegisterResolver("donation", func(ctx context.Context, id string) (int64, error) {
		return 0, nil
	})

	// funding - custom amount
	p.RegisterResolver("funding", func(ctx context.Context, id string) (int64, error) {
		return 0, nil
	})
}

// ===== Account Helpers =====

func (p *PaymentService) ensureAccountActive(acc Account) error {
	if acc.Status != "active" {
		return errors.New("account_not_active")
	}
	return nil
}

// ===== Payment Rules =====

type PaymentRule struct {
	AllowedEntities map[string]bool
	AllowedMethods  map[string]bool
	AllowCustomAmt  bool
}

func NormalizePaymentMethod(method string) string {
	switch method {
	case "", "wallet":
		return "wallet"
	case "card", "card_payment", "stripe":
		return "card"
	case "cod", "cash_on_delivery":
		return "cash_on_delivery"
	case "transfer", "wallet_transfer":
		return "transfer"
	default:
		return method
	}
}

var PaymentRules = map[string]PaymentRule{
	"funding": {
		AllowedEntities: map[string]bool{"artist": true},
		AllowedMethods:  map[string]bool{"card": true, "wallet": true},
		AllowCustomAmt:  true,
	},
	"donation": {
		AllowedEntities: map[string]bool{"post": true, "artist": true},
		AllowedMethods:  map[string]bool{"wallet": true, "card": true},
		AllowCustomAmt:  true,
	},
	"purchase": {
		AllowedEntities: map[string]bool{
			"order":   true,
			"cart":    true,
			"ticket":  true,
			"menu":    true,
			"service": true,
			"product": true,
			"booking": true,
			"merch":   true,
			"crop":    true,
			"farm":    true,
			"beat":    true,
		},
		AllowedMethods: map[string]bool{
			"wallet":           true,
			"card":             true,
			"transfer":         true,
			"cod":              true,
			"cash_on_delivery": true,
		},
		AllowCustomAmt: false,
	},
}
