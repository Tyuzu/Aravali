package stripe

import (
	"context"
	"errors"
	"fmt"
	"time"

	"scav/config"
	"scav/infra"
)

var fundingCollection = config.Collections.FundingCollection
var stripeOrdersCollection = config.Collections.StripeOrdersCollection

func SQLupdatePaymentStatus(
	ctx context.Context,
	entityType string,
	entityId string,
	amount int64,
	paymentIntentId string,
	app *infra.Deps,
) (int64, error) {
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
		return 0, errors.New("invalid entityType")
	}

	query := fmt.Sprintf("%s = $1", idField)
	args := []any{entityId}

	update := map[string]any{
		"paid":            true,
		"amount":          amount,
		"paymentIntentId": paymentIntentId,
		"paidAt":          time.Now().UTC(),
	}

	return app.SQLDB.Update(
		ctx,
		collection,
		query,
		args,
		update,
	)
}
