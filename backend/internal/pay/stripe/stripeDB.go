package stripe

import (
	"context"
	"errors"
	"time"

	"scav/config"
	"scav/infra"
)

var fundingCollection = config.Collections.FundingCollection
var stripeOrdersCollection = config.Collections.StripeOrdersCollection

func updatePaymentStatus(
	ctx context.Context,
	entityType string,
	entityId string,
	amount int64,
	paymentIntentId string,
	app *infra.Deps,
) (any, error) {
	var collection string
	var idField string

	switch entityType {
	case "funding":
		collection = fundingCollection
		idField = "fundingid"
	case "order":
		collection = stripeOrdersCollection
		idField = "orderid"
	default:
		return nil, errors.New("invalid entityType")
	}

	update := map[string]any{
		"paid":            true,
		"amount":          amount,
		"paymentIntentId": paymentIntentId,
		"paidAt":          time.Now().UTC(),
	}

	return app.DB.Update(
		ctx,
		collection,
		map[string]any{idField: entityId},
		update,
	)
}
