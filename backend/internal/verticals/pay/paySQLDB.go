package pay

import (
	"context"
	"errors"
	"time"

	"scav/config"
	"scav/internal/auth"
	"scav/internal/verticals/tickets"
	"scav/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var journalTable = config.Tables.JournalTable
var accountsTable = config.Tables.AccountsTable
var transactionsTable = config.Tables.TransactionTable
var globalLedgerTable = config.Tables.GlobalLedgerTable
var farmOrdersTable = config.Tables.FarmOrdersTable
var ticketsTable = config.Tables.TicketsTable
var menuTable = config.Tables.MenuTable
var serviceTable = config.Tables.ServiceTable
var productTable = config.Tables.ProductTable
var bookingsTable = config.Tables.BookingsTable
var merchTable = config.Tables.MerchTable
var cropsTable = config.Tables.CropsTable
var ordersTable = config.Tables.OrderTable
var usersTable = config.Tables.UserTable
var RefundsTable = config.Tables.RefundsTable
var webhookTable = config.Tables.DeliveryWebhooksTable

func (p *PaymentService) SQLgetOrCreateAccount(ctx context.Context, userID string) (string, error) {
	var acc Account
	err := p.app.DB.FindOne(ctx, accountsTable, map[string]any{"userid": userID}, &acc)
	if err == nil {
		return acc.ID, nil
	}

	if userID != "merchant" && userID != "external" {
		if !p.userExists(ctx, userID) {
			return "", errors.New("user_not_found")
		}
	}

	newAcc := Account{
		ID:            utils.GetUUID(),
		UserID:        userID,
		Currency:      "INR",
		Status:        "active",
		CachedBalance: 0,
		Version:       1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := p.app.DB.InsertOne(ctx, accountsTable, newAcc); err != nil {
		// race: retry read
		err = p.app.DB.FindOne(ctx, accountsTable, map[string]any{"userid": userID}, &acc)
		return acc.ID, err
	}

	return newAcc.ID, nil
}

func (p *PaymentService) SQLuserExists(ctx context.Context, userID string) bool {
	if userID == "" {
		return false
	}

	var user auth.User
	return p.app.DB.FindOne(ctx, usersTable, map[string]any{"userid": userID}, &user) == nil
}

func (p *PaymentService) SQLgetAccountByID(ctx context.Context, accountID string) (Account, error) {
	var acc Account
	err := p.app.DB.FindOne(ctx, accountsTable, map[string]any{"_id": accountID}, &acc)
	return acc, err
}

func (p *PaymentService) SQLgetAccountByUserID(ctx context.Context, userID string) (Account, error) {
	var acc Account
	err := p.app.DB.FindOne(ctx, accountsTable, map[string]any{"userid": userID}, &acc)
	return acc, err
}

func (p *PaymentService) SQLlistUserTransactions(ctx context.Context, userID string, skip int64, limit int64) ([]Transaction, error) {
	var txns []Transaction
	filter := map[string]any{
		"$or": []map[string]any{
			{"userid": userID},
			{"meta.recipient": userID},
		},
	}

	err := p.app.DB.FindMany(
		ctx,
		transactionsTable,
		filter,
		&txns,
		options.Find().
			SetSort(map[string]int{"created_at": -1}).
			SetSkip(skip).
			SetLimit(limit),
	)
	return txns, err
}

func (p *PaymentService) SQLfetchPriceByField(ctx context.Context, tableName string, fieldName string, entityID string) (int64, error) {
	var v struct{ Price int64 }
	err := p.app.DB.FindOne(ctx, tableName, map[string]any{fieldName: entityID}, &v)
	return v.Price, err
}

func (p *PaymentService) SQLfindOrderTotalByUser(ctx context.Context, orderID string, userID string) (string, int64, error) {
	var regularOrder struct {
		Total int64 `bson:"total"`
	}
	err := p.app.DB.FindOne(ctx, ordersTable, map[string]any{"orderId": orderID, "userid": userID}, &regularOrder)
	if err == nil {
		return "regular", regularOrder.Total, nil
	}
	if err != mongo.ErrNoDocuments {
		return "", 0, err
	}

	var farmOrder struct {
		Total int64 `bson:"total"`
	}
	err = p.app.DB.FindOne(ctx, farmOrdersTable, map[string]any{"orderid": orderID, "userid": userID}, &farmOrder)
	if err == nil {
		return "farm", farmOrder.Total, nil
	}
	if err == mongo.ErrNoDocuments {
		return "", 0, mongo.ErrNoDocuments
	}
	return "", 0, err
}

func (p *PaymentService) SQLfindOrderTotalByID(ctx context.Context, id string) (int64, error) {
	var o struct {
		Total int64 `bson:"total"`
	}
	err := p.app.DB.FindOne(ctx, ordersTable, map[string]any{"orderId": id}, &o)
	if err == nil {
		return o.Total, nil
	}
	if err != mongo.ErrNoDocuments {
		return 0, err
	}

	var fo struct {
		PriceAtPurchase float64 `bson:"priceAtPurchase"`
	}
	err = p.app.DB.FindOne(ctx, farmOrdersTable, map[string]any{"orderid": id}, &fo)
	if err != nil {
		return 0, err
	}
	return int64(fo.PriceAtPurchase * 100), nil
}

func (p *PaymentService) SQLfindRefundRequestByOrderID(ctx context.Context, orderID string) (OrderRefundRequest, error) {
	var refund OrderRefundRequest
	err := p.app.DB.FindOne(ctx, RefundsTable, map[string]any{
		"order_id": orderID,
		"status":   map[string]any{"$in": []string{"pending", "approved"}},
	}, &refund)
	return refund, err
}

func (p *PaymentService) SQLfindActiveRefundRequestByOrderID(ctx context.Context, orderID string) (OrderRefundRequest, error) {
	var refund OrderRefundRequest
	err := p.app.DB.FindOne(ctx, RefundsTable, map[string]any{
		"order_id": orderID,
		"status":   map[string]any{"$in": []string{"pending", "approved"}},
	}, &refund)
	return refund, err
}

func (p *PaymentService) SQLcreateRefundRequestRecord(ctx context.Context, refundReq OrderRefundRequest) error {
	return p.app.DB.InsertOne(ctx, RefundsTable, refundReq)
}

func (p *PaymentService) SQLcountRefundRequestsByUser(ctx context.Context, userID string) (int64, error) {
	return p.app.DB.Count(ctx, RefundsTable, map[string]any{"userid": userID})
}

func (p *PaymentService) SQLlistRefundRequestsByUser(ctx context.Context, userID string, skip int, limit int) ([]tickets.RefundRequest, error) {
	var refunds []tickets.RefundRequest
	err := p.app.DB.FindMany(
		ctx,
		RefundsTable,
		map[string]any{"userid": userID},
		&refunds,
		options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(map[string]any{"created_at": -1}),
	)
	return refunds, err
}

func (p *PaymentService) SQLcountRefundRequests(ctx context.Context, filter map[string]any) (int64, error) {
	return p.app.DB.Count(ctx, RefundsTable, filter)
}

func (p *PaymentService) SQLlistRefundRequests(ctx context.Context, filter map[string]any, skip int, limit int) ([]tickets.RefundRequest, error) {
	var refunds []tickets.RefundRequest
	err := p.app.DB.FindMany(
		ctx,
		RefundsTable,
		filter,
		&refunds,
		options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(map[string]any{"created_at": -1}),
	)
	return refunds, err
}

func (p *PaymentService) SQLfindRefundRequestByID(ctx context.Context, refundID string) (OrderRefundRequest, error) {
	var refund OrderRefundRequest
	err := p.app.DB.FindOne(ctx, RefundsTable, map[string]any{"_id": refundID}, &refund)
	return refund, err
}

func (p *PaymentService) SQLcreateRefundTransactionRecord(ctx context.Context, refundTxn Transaction) error {
	return p.app.DB.InsertOne(ctx, transactionsTable, refundTxn)
}

func (p *PaymentService) SQLupdateRefundRequestStatus(ctx context.Context, refundID string, update map[string]any) error {
	_, err := p.app.DB.UpdateOne(ctx, RefundsTable, map[string]any{"_id": refundID}, map[string]any{"$set": update})
	return err
}

func (p *PaymentService) SQLcreateWalletAccount(ctx context.Context, userID string) (Account, error) {
	newAcc := Account{
		ID:            utils.GetUUID(),
		UserID:        userID,
		Currency:      "INR",
		Status:        "active",
		CachedBalance: 0,
		Version:       1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := p.app.DB.InsertOne(ctx, accountsTable, newAcc); err != nil {
		return Account{}, err
	}
	return newAcc, nil
}

func (p *PaymentService) SQLcreateTransactionRecord(ctx context.Context, txn Transaction) error {
	return p.app.DB.InsertOne(ctx, transactionsTable, txn)
}

func (p *PaymentService) SQLcreateJournalEntryRecord(ctx context.Context, entry JournalEntry) error {
	return p.app.DB.InsertOne(ctx, journalTable, entry)
}

func (p *PaymentService) SQLapplyBalanceDelta(ctx context.Context, accountID string, delta int64) error {
	return p.app.DB.Inc(ctx, accountsTable, map[string]any{"_id": accountID}, "cached_balance", delta)
}

func (p *PaymentService) SQLsetTransactionStatus(ctx context.Context, txnID string, status string) error {
	_, err := p.app.DB.UpdateOne(ctx, transactionsTable,
		map[string]any{"_id": txnID},
		map[string]any{"$set": map[string]any{"status": status, "updated_at": time.Now()}},
	)
	return err
}

func (p *PaymentService) SQLupdateTransactionStatus(ctx context.Context, txnID string, status string, updatedAt time.Time) error {
	_, err := p.app.DB.UpdateOne(ctx, transactionsTable,
		map[string]any{"_id": txnID},
		map[string]any{"$set": map[string]any{"status": status, "updated_at": updatedAt}},
	)
	return err
}

func (p *PaymentService) SQLfindTransactionByID(ctx context.Context, txnID string) (Transaction, error) {
	var txn Transaction
	err := p.app.DB.FindOne(ctx, transactionsTable, map[string]any{"_id": txnID}, &txn)
	return txn, err
}

func (p *PaymentService) SQLmarkTransactionReversed(ctx context.Context, txnID string, updatedAt time.Time) error {
	_, err := p.app.DB.UpdateOne(
		ctx,
		transactionsTable,
		map[string]any{"_id": txnID},
		map[string]any{
			"$set": map[string]any{
				"status":     "reversed",
				"updated_at": updatedAt,
			},
		},
	)
	return err
}

func (p *PaymentService) SQLcreateTransferViews(ctx context.Context, txnID string, senderID string, recipientID string, amount int64, now time.Time) error {
	return p.app.DB.InsertMany(ctx, transactionsTable, []interface{}{
		Transaction{
			ID:        utils.GetUUID(),
			ParentTxn: txnID,
			UserID:    senderID,
			Type:      "debit",
			Amount:    amount,
			Status:    "success",
			CreatedAt: now,
		},
		Transaction{
			ID:        utils.GetUUID(),
			ParentTxn: txnID,
			UserID:    recipientID,
			Type:      "credit",
			Amount:    amount,
			Status:    "success",
			CreatedAt: now,
		},
	})
}

func (p *PaymentService) SQLrecordWebhookProcessing(ctx context.Context, payload *PaymentWebhookPayload) error {
	return p.app.DB.InsertOne(ctx, webhookTable, map[string]any{
		"transactionId": payload.TransactionID,
		"orderId":       payload.OrderID,
		"userid":        payload.UserID,
		"status":        payload.Status,
		"amount":        payload.Amount,
		"processedAt":   time.Now(),
	})
}

func (p *PaymentService) SQLhasWebhookBeenProcessed(ctx context.Context, transactionID string) (bool, error) {
	var existingWebhook map[string]any
	err := p.app.DB.FindOne(ctx, webhookTable, map[string]any{
		"transactionId": transactionID,
	}, &existingWebhook)
	return err == nil, err
}

func (p *PaymentService) SQLincrementTopupBalanceByUser(ctx context.Context, userID string, amount float64, updatedAt time.Time) error {
	_, err := p.app.DB.UpdateOne(ctx, accountsTable, map[string]any{
		"userid": userID,
	}, map[string]any{
		"$inc": map[string]any{
			"cached_balance": int64(amount),
		},
		"$set": map[string]any{
			"updated_at": updatedAt,
		},
	})
	return err
}

func (p *PaymentService) SQLsetTransactionStatusByID(ctx context.Context, txnID string, status string, updatedAt time.Time) error {
	_, err := p.app.DB.UpdateOne(ctx, transactionsTable, map[string]any{
		"_id": txnID,
	}, map[string]any{
		"$set": map[string]any{
			"status":     status,
			"updated_at": updatedAt,
		},
	})
	return err
}

func (p *PaymentService) SQLupdateOrderStatus(ctx context.Context, tableName string, lookupField string, orderID string, status string) error {
	_, err := p.app.DB.UpdateOne(ctx, tableName, map[string]any{lookupField: orderID}, map[string]any{"$set": map[string]any{"status": status}})
	return err
}

func (p *PaymentService) SQLupdateOrderSet(ctx context.Context, tableName string, lookupField string, orderID string, update map[string]any) error {
	_, err := p.app.DB.UpdateOne(ctx, tableName, map[string]any{lookupField: orderID}, map[string]any{"$set": update})
	return err
}

func (p *PaymentService) SQLdecrementInventory(ctx context.Context, tableName string, lookupField string, itemID string, incField string, qty int) error {
	if itemID == "" || qty <= 0 {
		return nil
	}
	_, err := p.app.DB.UpdateOne(ctx, tableName, map[string]any{lookupField: itemID}, map[string]any{"$inc": map[string]any{incField: -qty}})
	return err
}

func (p *PaymentService) SQLfindOrderByID(ctx context.Context, tableName string, lookupField string, orderID string, out any) error {
	return p.app.DB.FindOne(ctx, tableName, map[string]any{lookupField: orderID}, out)
}

func (p *PaymentService) SQLfailTxn(ctx context.Context, txnID string) {
	_, _ = p.app.DB.UpdateOne(ctx, transactionsTable,
		map[string]any{"_id": txnID},
		map[string]any{"$set": map[string]any{"status": "failed", "updated_at": time.Now()}},
	)
}

func (p *PaymentService) SQLsuccessTxn(ctx context.Context, txnID string) {
	_, _ = p.app.DB.UpdateOne(ctx, transactionsTable,
		map[string]any{"_id": txnID},
		map[string]any{"$set": map[string]any{"status": "success", "updated_at": time.Now()}},
	)
}

func (p *PaymentService) SQLrecordGlobalLedger(ctx context.Context, txnID string, journalEntryID string, ledgerType string, reason string, amount int64, accountID string, userID string) error {
	var entries []GlobalLedger

	totalAdditions := int64(0)
	totalDeletions := int64(0)

	err := p.app.DB.FindMany(ctx, globalLedgerTable,
		map[string]any{},
		&entries)

	if err == nil && len(entries) > 0 {
		lastEntry := entries[len(entries)-1]
		totalAdditions = lastEntry.TotalAdditionsUpto
		totalDeletions = lastEntry.TotalDeletionsUpto
	}

	switch ledgerType {
	case "addition":
		totalAdditions += amount
	case "deletion":
		totalDeletions += amount
	}

	entry := GlobalLedger{
		ID:                 utils.GetUUID(),
		TxnID:              txnID,
		Type:               ledgerType,
		Reason:             reason,
		Amount:             amount,
		Currency:           "INR",
		AccountID:          accountID,
		UserID:             userID,
		JournalEntryID:     journalEntryID,
		TotalAdditionsUpto: totalAdditions,
		TotalDeletionsUpto: totalDeletions,
		NetBalanceUpto:     totalAdditions - totalDeletions,
		CreatedAt:          time.Now(),
	}

	return p.app.DB.InsertOne(ctx, globalLedgerTable, entry)
}
