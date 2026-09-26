package pay

import (
	"context"
	"errors"
	"fmt"
	"time"

	"scav/config"
	"scav/internal/auth"
	"scav/internal/verticals/tickets"
	"scav/utils"
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
	err := p.app.SQLDB.FindOne(ctx, accountsTable, "userid = $1", []any{userID}, &acc)
	if err == nil {
		return acc.ID, nil
	}

	if userID != "merchant" && userID != "external" {
		if !p.SQLuserExists(ctx, userID) {
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

	if err := p.app.SQLDB.InsertOne(ctx, accountsTable, newAcc); err != nil {
		// race: retry read
		err = p.app.SQLDB.FindOne(ctx, accountsTable, "userid = $1", []any{userID}, &acc)
		return acc.ID, err
	}

	return newAcc.ID, nil
}

func (p *PaymentService) SQLuserExists(ctx context.Context, userID string) bool {
	if userID == "" {
		return false
	}

	var user auth.User
	return p.app.SQLDB.FindOne(ctx, usersTable, "userid = $1", []any{userID}, &user) == nil
}

func (p *PaymentService) SQLgetAccountByID(ctx context.Context, accountID string) (Account, error) {
	var acc Account
	err := p.app.SQLDB.FindOne(ctx, accountsTable, "id = $1", []any{accountID}, &acc)
	return acc, err
}

func (p *PaymentService) SQLgetAccountByUserID(ctx context.Context, userID string) (Account, error) {
	var acc Account
	err := p.app.SQLDB.FindOne(ctx, accountsTable, "userid = $1", []any{userID}, &acc)
	return acc, err
}

func (p *PaymentService) SQLlistUserTransactions(ctx context.Context, userID string, skip int64, limit int64) ([]Transaction, error) {
	var txns []Transaction
	query := "(userid = $1 OR meta->>'recipient' = $1) ORDER BY created_at DESC OFFSET $2 LIMIT $3"
	args := []any{userID, skip, limit}

	err := p.app.SQLDB.FindMany(ctx, transactionsTable, query, args, &txns)
	return txns, err
}

func (p *PaymentService) SQLfetchPriceByField(ctx context.Context, tableName string, fieldName string, entityID string) (int64, error) {
	var v struct{ Price int64 }
	query := fmt.Sprintf("%s = $1", fieldName)
	err := p.app.SQLDB.FindOne(ctx, tableName, query, []any{entityID}, &v)
	return v.Price, err
}

func (p *PaymentService) SQLfindOrderTotalByUser(ctx context.Context, orderID string, userID string) (string, int64, error) {
	var regularOrder struct {
		Total int64 `db:"total"`
	}
	err := p.app.SQLDB.FindOne(ctx, ordersTable, "orderid = $1 AND userid = $2", []any{orderID, userID}, &regularOrder)
	if err == nil {
		return "regular", regularOrder.Total, nil
	}

	var farmOrder struct {
		Total int64 `db:"total"`
	}
	err = p.app.SQLDB.FindOne(ctx, farmOrdersTable, "orderid = $1 AND userid = $2", []any{orderID, userID}, &farmOrder)
	if err == nil {
		return "farm", farmOrder.Total, nil
	}

	return "", 0, err
}

func (p *PaymentService) SQLfindOrderTotalByID(ctx context.Context, id string) (int64, error) {
	var o struct {
		Total int64 `db:"total"`
	}
	err := p.app.SQLDB.FindOne(ctx, ordersTable, "orderid = $1", []any{id}, &o)
	if err == nil {
		return o.Total, nil
	}

	var fo struct {
		PriceAtPurchase float64 `db:"priceatpurchase"`
	}
	err = p.app.SQLDB.FindOne(ctx, farmOrdersTable, "orderid = $1", []any{id}, &fo)
	if err != nil {
		return 0, err
	}
	return int64(fo.PriceAtPurchase * 100), nil
}

func (p *PaymentService) SQLfindRefundRequestByOrderID(ctx context.Context, orderID string) (OrderRefundRequest, error) {
	var refund OrderRefundRequest
	query := "order_id = $1 AND status IN ('pending', 'approved')"
	err := p.app.SQLDB.FindOne(ctx, RefundsTable, query, []any{orderID}, &refund)
	return refund, err
}

func (p *PaymentService) SQLfindActiveRefundRequestByOrderID(ctx context.Context, orderID string) (OrderRefundRequest, error) {
	var refund OrderRefundRequest
	query := "order_id = $1 AND status IN ('pending', 'approved')"
	err := p.app.SQLDB.FindOne(ctx, RefundsTable, query, []any{orderID}, &refund)
	return refund, err
}

func (p *PaymentService) SQLcreateRefundRequestRecord(ctx context.Context, refundReq OrderRefundRequest) error {
	return p.app.SQLDB.InsertOne(ctx, RefundsTable, refundReq)
}

func (p *PaymentService) SQLcountRefundRequestsByUser(ctx context.Context, userID string) (int64, error) {
	return p.app.SQLDB.Count(ctx, RefundsTable, "userid = $1", []any{userID})
}

func (p *PaymentService) SQLlistRefundRequestsByUser(ctx context.Context, userID string, skip int, limit int) ([]tickets.RefundRequest, error) {
	var refunds []tickets.RefundRequest
	query := "userid = $1 ORDER BY created_at DESC OFFSET $2 LIMIT $3"
	args := []any{userID, skip, limit}

	err := p.app.SQLDB.FindMany(ctx, RefundsTable, query, args, &refunds)
	return refunds, err
}

func (p *PaymentService) SQLcountRefundRequests(ctx context.Context, query string, args []any) (int64, error) {
	return p.app.SQLDB.Count(ctx, RefundsTable, query, args)
}

func (p *PaymentService) SQLlistRefundRequests(ctx context.Context, query string, args []any, skip int, limit int) ([]tickets.RefundRequest, error) {
	var refunds []tickets.RefundRequest
	fullQuery := fmt.Sprintf("%s ORDER BY created_at DESC OFFSET %d LIMIT %d", query, skip, limit)

	err := p.app.SQLDB.FindMany(ctx, RefundsTable, fullQuery, args, &refunds)
	return refunds, err
}

func (p *PaymentService) SQLfindRefundRequestByID(ctx context.Context, refundID string) (OrderRefundRequest, error) {
	var refund OrderRefundRequest
	err := p.app.SQLDB.FindOne(ctx, RefundsTable, "id = $1", []any{refundID}, &refund)
	return refund, err
}

func (p *PaymentService) SQLcreateRefundTransactionRecord(ctx context.Context, refundTxn Transaction) error {
	return p.app.SQLDB.InsertOne(ctx, transactionsTable, refundTxn)
}

func (p *PaymentService) SQLupdateRefundRequestStatus(ctx context.Context, refundID string, update map[string]any) error {
	_, err := p.app.SQLDB.UpdateOne(ctx, RefundsTable, "id = $1", []any{refundID}, update)
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

	if err := p.app.SQLDB.InsertOne(ctx, accountsTable, newAcc); err != nil {
		return Account{}, err
	}
	return newAcc, nil
}

func (p *PaymentService) SQLcreateTransactionRecord(ctx context.Context, txn Transaction) error {
	return p.app.SQLDB.InsertOne(ctx, transactionsTable, txn)
}

func (p *PaymentService) SQLcreateJournalEntryRecord(ctx context.Context, entry JournalEntry) error {
	return p.app.SQLDB.InsertOne(ctx, journalTable, entry)
}

func (p *PaymentService) SQLapplyBalanceDelta(ctx context.Context, accountID string, delta int64) error {
	return p.app.SQLDB.Inc(ctx, accountsTable, "id = $1", []any{accountID}, "cached_balance", delta)
}

func (p *PaymentService) SQLsetTransactionStatus(ctx context.Context, txnID string, status string) error {
	update := map[string]any{
		"status":     status,
		"updated_at": time.Now(),
	}
	_, err := p.app.SQLDB.UpdateOne(ctx, transactionsTable, "id = $1", []any{txnID}, update)
	return err
}

func (p *PaymentService) SQLupdateTransactionStatus(ctx context.Context, txnID string, status string, updatedAt time.Time) error {
	update := map[string]any{
		"status":     status,
		"updated_at": updatedAt,
	}
	_, err := p.app.SQLDB.UpdateOne(ctx, transactionsTable, "id = $1", []any{txnID}, update)
	return err
}

func (p *PaymentService) SQLfindTransactionByID(ctx context.Context, txnID string) (Transaction, error) {
	var txn Transaction
	err := p.app.SQLDB.FindOne(ctx, transactionsTable, "id = $1", []any{txnID}, &txn)
	return txn, err
}

func (p *PaymentService) SQLmarkTransactionReversed(ctx context.Context, txnID string, updatedAt time.Time) error {
	update := map[string]any{
		"status":     "reversed",
		"updated_at": updatedAt,
	}
	_, err := p.app.SQLDB.UpdateOne(ctx, transactionsTable, "id = $1", []any{txnID}, update)
	return err
}

func (p *PaymentService) SQLcreateTransferViews(ctx context.Context, txnID string, senderID string, recipientID string, amount int64, now time.Time) error {
	return p.app.SQLDB.InsertMany(ctx, transactionsTable, []any{
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
	return p.app.SQLDB.InsertOne(ctx, webhookTable, map[string]any{
		"transactionid": payload.TransactionID,
		"orderid":       payload.OrderID,
		"userid":        payload.UserID,
		"status":        payload.Status,
		"amount":        payload.Amount,
		"processedat":   time.Now(),
	})
}

func (p *PaymentService) SQLhasWebhookBeenProcessed(ctx context.Context, transactionID string) (bool, error) {
	var existingWebhook map[string]any
	err := p.app.SQLDB.FindOne(ctx, webhookTable, "transactionid = $1", []any{transactionID}, &existingWebhook)
	return err == nil, err
}

func (p *PaymentService) SQLincrementTopupBalanceByUser(ctx context.Context, userID string, amount float64, updatedAt time.Time) error {
	err := p.app.SQLDB.Inc(ctx, accountsTable, "userid = $1", []any{userID}, "cached_balance", int64(amount))
	if err != nil {
		return err
	}

	_, err = p.app.SQLDB.UpdateOne(ctx, accountsTable, "userid = $1", []any{userID}, map[string]any{
		"updated_at": updatedAt,
	})
	return err
}

func (p *PaymentService) SQLsetTransactionStatusByID(ctx context.Context, txnID string, status string, updatedAt time.Time) error {
	update := map[string]any{
		"status":     status,
		"updated_at": updatedAt,
	}
	_, err := p.app.SQLDB.UpdateOne(ctx, transactionsTable, "id = $1", []any{txnID}, update)
	return err
}

func (p *PaymentService) SQLupdateOrderStatus(ctx context.Context, tableName string, lookupField string, orderID string, status string) error {
	query := fmt.Sprintf("%s = $1", lookupField)
	_, err := p.app.SQLDB.UpdateOne(ctx, tableName, query, []any{orderID}, map[string]any{"status": status})
	return err
}

func (p *PaymentService) SQLupdateOrderSet(ctx context.Context, tableName string, lookupField string, orderID string, update map[string]any) error {
	query := fmt.Sprintf("%s = $1", lookupField)
	_, err := p.app.SQLDB.UpdateOne(ctx, tableName, query, []any{orderID}, update)
	return err
}

func (p *PaymentService) SQLdecrementInventory(ctx context.Context, tableName string, lookupField string, itemID string, incField string, qty int) error {
	if itemID == "" || qty <= 0 {
		return nil
	}
	query := fmt.Sprintf("%s = $1", lookupField)
	return p.app.SQLDB.Inc(ctx, tableName, query, []any{itemID}, incField, int64(-qty))
}

func (p *PaymentService) SQLfindOrderByID(ctx context.Context, tableName string, lookupField string, orderID string, out any) error {
	query := fmt.Sprintf("%s = $1", lookupField)
	return p.app.SQLDB.FindOne(ctx, tableName, query, []any{orderID}, out)
}

func (p *PaymentService) SQLfailTxn(ctx context.Context, txnID string) {
	_, _ = p.app.SQLDB.UpdateOne(ctx, transactionsTable, "id = $1", []any{txnID}, map[string]any{
		"status":     "failed",
		"updated_at": time.Now(),
	})
}

func (p *PaymentService) SQLsuccessTxn(ctx context.Context, txnID string) {
	_, _ = p.app.SQLDB.UpdateOne(ctx, transactionsTable, "id = $1", []any{txnID}, map[string]any{
		"status":     "success",
		"updated_at": time.Now(),
	})
}

func (p *PaymentService) SQLrecordGlobalLedger(ctx context.Context, txnID string, journalEntryID string, ledgerType string, reason string, amount int64, accountID string, userID string) error {
	var entries []GlobalLedger

	totalAdditions := int64(0)
	totalDeletions := int64(0)

	err := p.app.SQLDB.FindMany(ctx, globalLedgerTable, "1=1 ORDER BY created_at ASC", []any{}, &entries)

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

	return p.app.SQLDB.InsertOne(ctx, globalLedgerTable, entry)
}
