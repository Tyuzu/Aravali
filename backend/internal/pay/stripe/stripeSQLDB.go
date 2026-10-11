// File: internal/pay/stripe/stripeSQLDB.go

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
	if app == nil || app.SQLDB == nil {
		return 0, errors.New("database not initialized")
	}

	var tableName string
	var idField string

	switch entityType {
	case "funding":
		tableName = fundingTable
		idField = "fundingid"
	case "order":
		tableName = stripeOrdersTable
		idField = "orderid"
	default:
		return 0, errors.New("invalid entityType")
	}

	if tableName == "" || idField == "" {
		return 0, errors.New("missing table metadata")
	}

	result, err := app.SQLDB.Exec(ctx,
		`UPDATE `+tableName+` SET status = $1, paymentintentid = $2, amount = $3, updated_at = NOW() WHERE `+idField+` = $4`,
		"paid",
		paymentIntentId,
		amount,
		entityId,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
