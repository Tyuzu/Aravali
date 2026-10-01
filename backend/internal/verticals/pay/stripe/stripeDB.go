// File: internal/verticals/pay/stripe/stripeDB.go

package stripe

import (
	"context"
	"errors"
	"fmt"
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
) (int64, error) {
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
		return 0, errors.New("invalid entityType")
	}

	where := fmt.Sprintf("%s = $1", idField)
	args := []any{entityId}
	update := map[string]any{
		"paid":            true,
		"amount":          amount,
		"paymentIntentId": paymentIntentId,
		"paidAt":          time.Now().UTC(),
	}

	return app.SQLDB.Update(ctx, table, where, args, update)
}
