// File: internal/verticals/pay/stripe/stripeDB.go

package stripe

import (
	"context"
	"errors"
	"time"

	"scav/config"
	"scav/infra"
)

var fundingTable = config.Tables.FundingTable
var stripeOrdersTable = config.Tables.StripeOrdersTable

func updatePaymentStatus(
	ctx context.Context,
	entityType string,
	entityId string,
	amount int64,
	paymentIntentId string,
	app *infra.Deps,
) (any, error) {
	var table string
	var idField string

	switch entityType {
	case "funding":
		table = fundingTable
		idField = "fundingid"
	case "order":
		table = stripeOrdersTable
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
		table,
		map[string]any{idField: entityId},
		update,
	)
}
