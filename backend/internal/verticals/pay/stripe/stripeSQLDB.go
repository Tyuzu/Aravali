// File: internal/verticals/pay/stripe/stripeSQLDB.go

package stripe

import (
	"context"
	"errors"

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
) (int64, error) {

	switch entityType {
	case "funding":
		Table = fundingTable
		idField = "fundingid"
	case "order":
		Table = stripeOrdersTable
		idField = "orderid"
	default:
		return 0, errors.New("invalid entityType")
	}

}
