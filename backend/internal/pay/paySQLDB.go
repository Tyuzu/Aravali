// File: internal/pay/paySQLDB.go

package pay

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/internal/tickets"
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

func (p *PaymentService) getOrCreateAccount(ctx context.Context, userID string) (string, error) {
	return p.SQLgetOrCreateAccount(ctx, userID)
}

func (p *PaymentService) getAccountByID(ctx context.Context, accountID string) (Account, error) {
	return p.SQLgetAccountByID(ctx, accountID)
}

func (p *PaymentService) getAccountByUserID(ctx context.Context, userID string) (Account, error) {
	return p.SQLgetAccountByUserID(ctx, userID)
}

func (p *PaymentService) listUserTransactions(ctx context.Context, userID string, skip int64, limit int64) ([]Transaction, error) {
	return p.SQLlistUserTransactions(ctx, userID, skip, limit)
}

func (p *PaymentService) fetchPriceByField(ctx context.Context, tableName string, fieldName string, entityID string) (int64, error) {
	return p.SQLfetchPriceByField(ctx, tableName, fieldName, entityID)
}

func (p *PaymentService) findOrderTotalByID(ctx context.Context, id string) (int64, error) {
	return p.SQLfindOrderTotalByID(ctx, id)
}

func (p *PaymentService) createTransactionRecord(ctx context.Context, txn Transaction) error {
	return p.SQLcreateTransactionRecord(ctx, txn)
}

func (p *PaymentService) applyBalanceDelta(ctx context.Context, accountID string, delta int64) error {
	return p.SQLapplyBalanceDelta(ctx, accountID, delta)
}

func (p *PaymentService) setTransactionStatus(ctx context.Context, txnID string, status string) error {
	return p.SQLsetTransactionStatus(ctx, txnID, status)
}

func (p *PaymentService) updateTransactionStatus(ctx context.Context, txnID string, status string, updatedAt time.Time) error {
	return p.SQLupdateTransactionStatus(ctx, txnID, status, updatedAt)
}

func (p *PaymentService) findTransactionByID(ctx context.Context, txnID string) (Transaction, error) {
	return p.SQLfindTransactionByID(ctx, txnID)
}

func (p *PaymentService) markTransactionReversed(ctx context.Context, txnID string, updatedAt time.Time) error {
	return p.SQLmarkTransactionReversed(ctx, txnID, updatedAt)
}

func (p *PaymentService) recordWebhookProcessing(ctx context.Context, payload *PaymentWebhookPayload) error {
	return p.SQLrecordWebhookProcessing(ctx, payload)
}

func (p *PaymentService) hasWebhookBeenProcessed(ctx context.Context, transactionID string) (bool, error) {
	return p.SQLhasWebhookBeenProcessed(ctx, transactionID)
}

func (p *PaymentService) incrementTopupBalanceByUser(ctx context.Context, userID string, amount float64, updatedAt time.Time) error {
	return p.SQLincrementTopupBalanceByUser(ctx, userID, amount, updatedAt)
}

func (p *PaymentService) setTransactionStatusByID(ctx context.Context, txnID string, status string, updatedAt time.Time) error {
	return p.SQLsetTransactionStatusByID(ctx, txnID, status, updatedAt)
}

func (p *PaymentService) updateOrderStatus(ctx context.Context, tableName string, lookupField string, orderID string, status string) error {
	return p.SQLupdateOrderStatus(ctx, tableName, lookupField, orderID, status)
}

func (p *PaymentService) updateOrderSet(ctx context.Context, tableName string, lookupField string, orderID string, update map[string]any) error {
	return p.SQLupdateOrderSet(ctx, tableName, lookupField, orderID, update)
}

func (p *PaymentService) decrementInventory(ctx context.Context, tableName string, lookupField string, itemID string, incField string, qty int) error {
	return p.SQLdecrementInventory(ctx, tableName, lookupField, itemID, incField, qty)
}

func (p *PaymentService) findOrderByID(ctx context.Context, tableName string, lookupField string, orderID string, out any) error {
	return p.SQLfindOrderByID(ctx, tableName, lookupField, orderID, out)
}

func (p *PaymentService) findOrderTotalByUser(ctx context.Context, orderID string, userID string) (string, int64, error) {
	return p.SQLfindOrderTotalByUser(ctx, orderID, userID)
}

func (p *PaymentService) findRefundRequestByOrderID(ctx context.Context, orderID string) (OrderRefundRequest, error) {
	return p.SQLfindRefundRequestByOrderID(ctx, orderID)
}

func (p *PaymentService) findActiveRefundRequestByOrderID(ctx context.Context, orderID string) (OrderRefundRequest, error) {
	return p.SQLfindActiveRefundRequestByOrderID(ctx, orderID)
}

func (p *PaymentService) createRefundRequestRecord(ctx context.Context, refundReq OrderRefundRequest) error {
	return p.SQLcreateRefundRequestRecord(ctx, refundReq)
}

func (p *PaymentService) countRefundRequestsByUser(ctx context.Context, userID string) (int64, error) {
	return p.SQLcountRefundRequestsByUser(ctx, userID)
}

func (p *PaymentService) listRefundRequestsByUser(ctx context.Context, userID string, skip int, limit int) ([]tickets.RefundRequest, error) {
	return p.SQLlistRefundRequestsByUser(ctx, userID, skip, limit)
}

func (p *PaymentService) countRefundRequests(ctx context.Context, filter map[string]any) (int64, error) {
	where, args := buildSQLFilter(filter)
	return p.SQLcountRefundRequests(ctx, where, args)
}

func (p *PaymentService) listRefundRequests(ctx context.Context, filter map[string]any, skip int, limit int) ([]tickets.RefundRequest, error) {
	where, args := buildSQLFilter(filter)
	return p.SQLlistRefundRequests(ctx, where, args, skip, limit)
}

func buildSQLFilter(filter map[string]any) (string, []any) {
	if len(filter) == 0 {
		return "1=1", nil
	}
	clauses := make([]string, 0, len(filter))
	args := make([]any, 0, len(filter))
	idx := 1
	for key, val := range filter {
		if val == nil {
			clauses = append(clauses, fmt.Sprintf("%s IS NULL", key))
			continue
		}
		clauses = append(clauses, fmt.Sprintf("%s = $%d", key, idx))
		args = append(args, val)
		idx++
	}
	return strings.Join(clauses, " AND "), args
}

func (p *PaymentService) findRefundRequestByID(ctx context.Context, refundID string) (OrderRefundRequest, error) {
	return p.SQLfindRefundRequestByID(ctx, refundID)
}

func (p *PaymentService) createRefundTransactionRecord(ctx context.Context, refundTxn Transaction) error {
	return p.SQLcreateRefundTransactionRecord(ctx, refundTxn)
}

func (p *PaymentService) updateRefundRequestStatus(ctx context.Context, refundID string, update map[string]any) error {
	return p.SQLupdateRefundRequestStatus(ctx, refundID, update)
}

func (p *PaymentService) recordGlobalLedger(ctx context.Context, txnID string, journalEntryID string, ledgerType string, reason string, amount int64, accountID string, userID string) error {
	return p.SQLrecordGlobalLedger(ctx, txnID, journalEntryID, ledgerType, reason, amount, accountID, userID)
}

func (p *PaymentService) createJournalEntryRecord(ctx context.Context, entry JournalEntry) error {
	return p.SQLcreateJournalEntryRecord(ctx, entry)
}

func (p *PaymentService) failTxn(ctx context.Context, txnID string) {
	p.SQLfailTxn(ctx, txnID)
}

func (p *PaymentService) successTxn(ctx context.Context, txnID string) {
	p.SQLsuccessTxn(ctx, txnID)
}

func (p *PaymentService) createTransferViews(ctx context.Context, txnID string, senderID string, recipientID string, amount int64, now time.Time) error {
	return p.SQLcreateTransferViews(ctx, txnID, senderID, recipientID, amount, now)
}

func (p *PaymentService) userExists(ctx context.Context, userID string) bool {
	return p.SQLuserExists(ctx, userID)
}

func (p *PaymentService) createWalletAccount(ctx context.Context, userID string) (Account, error) {
	return p.SQLcreateWalletAccount(ctx, userID)
}

func (p *PaymentService) SQLgetOrCreateAccount(ctx context.Context, userID string) (string, error) {
	if p == nil || p.app == nil || p.app.SQLDB == nil {
		return "", nil
	}
	if userID != "merchant" && userID != "external" && !p.SQLuserExists(ctx, userID) {
		return "", errors.New("user_not_found")
	}
	if userID == "merchant" || userID == "external" {
		acc := Account{ID: utils.GetUUID(), UserID: userID, Currency: "INR", Status: "active", CachedBalance: 0, Version: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		_, err := p.app.SQLDB.Exec(ctx, `INSERT INTO `+accountsTable+` (_id, userid, currency, status, cached_balance, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) ON CONFLICT (_id) DO NOTHING`, acc.ID, acc.UserID, acc.Currency, acc.Status, acc.CachedBalance, acc.Version)
		if err != nil {
			return "", err
		}
		return acc.ID, nil
	}
	acc, err := p.SQLgetAccountByUserID(ctx, userID)
	if err == nil && acc.ID != "" {
		return acc.ID, nil
	}
	newAcc := Account{ID: utils.GetUUID(), UserID: userID, Currency: "INR", Status: "active", CachedBalance: 0, Version: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	_, err = p.app.SQLDB.Exec(ctx, `INSERT INTO `+accountsTable+` (_id, userid, currency, status, cached_balance, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) ON CONFLICT (_id) DO NOTHING`, newAcc.ID, newAcc.UserID, newAcc.Currency, newAcc.Status, newAcc.CachedBalance, newAcc.Version)
	if err != nil {
		return "", err
	}
	return newAcc.ID, nil
}

func (p *PaymentService) SQLuserExists(ctx context.Context, userID string) bool {
	if p == nil || p.app == nil || p.app.SQLDB == nil || userID == "" {
		return false
	}
	var exists int
	err := p.app.SQLDB.QueryRow(ctx, `SELECT 1 FROM `+usersTable+` WHERE userid = $1 LIMIT 1`, userID).Scan(&exists)
	return err == nil && exists == 1
}

func (p *PaymentService) SQLgetAccountByID(ctx context.Context, accountID string) (Account, error) {
	var acc Account
	if p == nil || p.app == nil || p.app.SQLDB == nil || accountID == "" {
		return acc, nil
	}
	err := p.app.SQLDB.QueryRow(ctx, `SELECT _id, userid, currency, status, cached_balance, version, created_at, updated_at FROM `+accountsTable+` WHERE _id = $1 LIMIT 1`, accountID).Scan(&acc.ID, &acc.UserID, &acc.Currency, &acc.Status, &acc.CachedBalance, &acc.Version, &acc.CreatedAt, &acc.UpdatedAt)
	return acc, err
}

func (p *PaymentService) SQLgetAccountByUserID(ctx context.Context, userID string) (Account, error) {
	var acc Account
	if p == nil || p.app == nil || p.app.SQLDB == nil || userID == "" {
		return acc, nil
	}
	err := p.app.SQLDB.QueryRow(ctx, `SELECT _id, userid, currency, status, cached_balance, version, created_at, updated_at FROM `+accountsTable+` WHERE userid = $1 LIMIT 1`, userID).Scan(&acc.ID, &acc.UserID, &acc.Currency, &acc.Status, &acc.CachedBalance, &acc.Version, &acc.CreatedAt, &acc.UpdatedAt)
	return acc, err
}

func (p *PaymentService) SQLlistUserTransactions(ctx context.Context, userID string, skip int64, limit int64) ([]Transaction, error) {
	if p == nil || p.app == nil || p.app.SQLDB == nil || userID == "" {
		return nil, nil
	}
	if limit <= 0 {
		limit = 25
	}
	rows, err := p.app.SQLDB.Query(ctx, `SELECT * FROM `+transactionsTable+` WHERE userid = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Transaction, 0)
	for rows.Next() {
		var txn Transaction
		if err := rows.Scan(&txn.ID, &txn.UserID, &txn.ParentTxn, &txn.Type, &txn.Method, &txn.EntityType, &txn.EntityID, &txn.FromAccount, &txn.ToAccount, &txn.Amount, &txn.Currency, &txn.Status, &txn.IdempotencyKey, &txn.CreatedAt, &txn.UpdatedAt); err == nil {
			items = append(items, txn)
		}
	}
	return items, rows.Err()
}

func (p *PaymentService) SQLfetchPriceByField(ctx context.Context, tableName string, fieldName string, entityID string) (int64, error) {
	if p == nil || p.app == nil || p.app.SQLDB == nil || tableName == "" || fieldName == "" || entityID == "" {
		return 0, nil
	}
	var price int64
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = $1 LIMIT 1`, fieldName, tableName, fieldName)
	err := p.app.SQLDB.QueryRow(ctx, query, entityID).Scan(&price)
	return price, err
}

func (p *PaymentService) SQLfindOrderTotalByUser(ctx context.Context, orderID string, userID string) (string, int64, error) {
	if p == nil || p.app == nil || p.app.SQLDB == nil || orderID == "" || userID == "" {
		return "", 0, nil
	}
	var orderType string
	var total int64
	if err := p.app.SQLDB.QueryRow(ctx, `SELECT type, total FROM `+ordersTable+` WHERE orderid = $1 AND userid = $2 LIMIT 1`, orderID, userID).Scan(&orderType, &total); err == nil {
		return orderType, total, nil
	}
	if err := p.app.SQLDB.QueryRow(ctx, `SELECT type, total FROM `+farmOrdersTable+` WHERE orderid = $1 AND userid = $2 LIMIT 1`, orderID, userID).Scan(&orderType, &total); err == nil {
		return orderType, total, nil
	}
	return "", 0, nil
}

func (p *PaymentService) SQLfindOrderTotalByID(ctx context.Context, id string) (int64, error) {
	if p == nil || p.app == nil || p.app.SQLDB == nil || id == "" {
		return 0, nil
	}
	var total int64
	if err := p.app.SQLDB.QueryRow(ctx, `SELECT total FROM `+ordersTable+` WHERE orderid = $1 LIMIT 1`, id).Scan(&total); err == nil {
		return total, nil
	}
	if err := p.app.SQLDB.QueryRow(ctx, `SELECT total FROM `+farmOrdersTable+` WHERE orderid = $1 LIMIT 1`, id).Scan(&total); err == nil {
		return total, nil
	}
	return 0, nil
}

func (p *PaymentService) SQLfindRefundRequestByOrderID(ctx context.Context, orderID string) (OrderRefundRequest, error) {
	var refund OrderRefundRequest
	if p == nil || p.app == nil || p.app.SQLDB == nil || orderID == "" {
		return refund, nil
	}
	err := p.app.SQLDB.QueryRow(ctx, `SELECT * FROM `+RefundsTable+` WHERE order_id = $1 LIMIT 1`, orderID).Scan(&refund.ID, &refund.OrderID, &refund.UserID, &refund.OrderType, &refund.Amount, &refund.Reason, &refund.Status, &refund.TransactionID, &refund.ReviewedBy, &refund.ReviewedAt, &refund.ReviewNotes, &refund.CreatedAt, &refund.UpdatedAt)
	return refund, err
}

func (p *PaymentService) SQLfindActiveRefundRequestByOrderID(ctx context.Context, orderID string) (OrderRefundRequest, error) {
	return p.SQLfindRefundRequestByOrderID(ctx, orderID)
}

func (p *PaymentService) SQLcreateRefundRequestRecord(ctx context.Context, refundReq OrderRefundRequest) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil {
		return nil
	}
	if refundReq.ID == "" {
		refundReq.ID = utils.GetUUID()
	}
	_, err := p.app.SQLDB.Exec(ctx,
		`INSERT INTO `+RefundsTable+` (_id, order_id, userid, order_type, amount, reason, status, transaction_id, reviewed_by, reviewed_at, review_notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())`,
		refundReq.ID, refundReq.OrderID, refundReq.UserID, refundReq.OrderType, refundReq.Amount, refundReq.Reason, refundReq.Status, refundReq.TransactionID, refundReq.ReviewedBy, refundReq.ReviewedAt, refundReq.ReviewNotes,
	)
	return err
}

func (p *PaymentService) SQLcountRefundRequestsByUser(ctx context.Context, userID string) (int64, error) {
	if p == nil || p.app == nil || p.app.SQLDB == nil || userID == "" {
		return 0, nil
	}
	var count int64
	err := p.app.SQLDB.QueryRow(ctx, `SELECT COUNT(*) FROM `+RefundsTable+` WHERE userid = $1`, userID).Scan(&count)
	return count, err
}

func (p *PaymentService) SQLlistRefundRequestsByUser(ctx context.Context, userID string, skip int, limit int) ([]tickets.RefundRequest, error) {
	if p == nil || p.app == nil || p.app.SQLDB == nil || userID == "" {
		return nil, nil
	}
	if limit <= 0 {
		limit = 25
	}
	rows, err := p.app.SQLDB.Query(ctx, `SELECT * FROM `+RefundsTable+` WHERE userid = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	refunds := make([]tickets.RefundRequest, 0)
	for rows.Next() {
		var refund tickets.RefundRequest
		if err := rows.Scan(&refund.EventID, &refund.TicketID, &refund.UserID, &refund.UniqueCode, &refund.RequestDate, &refund.Status, &refund.Amount, &refund.ProcessedAt, &refund.RefundedAt); err == nil {
			refunds = append(refunds, refund)
		}
	}
	return refunds, rows.Err()
}

func (p *PaymentService) SQLcountRefundRequests(ctx context.Context, query string, args []any) (int64, error) {
	if p == nil || p.app == nil || p.app.SQLDB == nil {
		return 0, nil
	}
	if query == "" {
		query = "1 = 1"
	}
	var count int64
	err := p.app.SQLDB.QueryRow(ctx, `SELECT COUNT(*) FROM `+RefundsTable+` WHERE `+query, args...).Scan(&count)
	return count, err
}

func (p *PaymentService) SQLlistRefundRequests(ctx context.Context, query string, args []any, skip int, limit int) ([]tickets.RefundRequest, error) {
	if p == nil || p.app == nil || p.app.SQLDB == nil {
		return nil, nil
	}
	if limit <= 0 {
		limit = 25
	}
	if query == "" {
		query = "1 = 1"
	}
	rows, err := p.app.SQLDB.Query(ctx, `SELECT * FROM `+RefundsTable+` WHERE `+query+` ORDER BY created_at DESC LIMIT $1 OFFSET $2`, append(args, limit, skip)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	refunds := make([]tickets.RefundRequest, 0)
	for rows.Next() {
		var refund tickets.RefundRequest
		if err := rows.Scan(&refund.EventID, &refund.TicketID, &refund.UserID, &refund.UniqueCode, &refund.RequestDate, &refund.Status, &refund.Amount, &refund.ProcessedAt, &refund.RefundedAt); err == nil {
			refunds = append(refunds, refund)
		}
	}
	return refunds, rows.Err()
}

func (p *PaymentService) SQLfindRefundRequestByID(ctx context.Context, refundID string) (OrderRefundRequest, error) {
	var refund OrderRefundRequest
	if p == nil || p.app == nil || p.app.SQLDB == nil || refundID == "" {
		return refund, nil
	}
	err := p.app.SQLDB.QueryRow(ctx, `SELECT * FROM `+RefundsTable+` WHERE _id = $1 OR id = $1 LIMIT 1`, refundID).Scan(&refund.ID, &refund.OrderID, &refund.UserID, &refund.OrderType, &refund.Amount, &refund.Reason, &refund.Status, &refund.TransactionID, &refund.ReviewedBy, &refund.ReviewedAt, &refund.ReviewNotes, &refund.CreatedAt, &refund.UpdatedAt)
	return refund, err
}

func (p *PaymentService) SQLcreateRefundTransactionRecord(ctx context.Context, refundTxn Transaction) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil {
		return nil
	}
	if refundTxn.ID == "" {
		refundTxn.ID = utils.GetUUID()
	}
	_, err := p.app.SQLDB.Exec(ctx,
		`INSERT INTO `+transactionsTable+` (_id, userid, parent_txn, type, method, entity_type, entity_id, from_account, to_account, amount, currency, status, external_ref, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW())`,
		refundTxn.ID, refundTxn.UserID, refundTxn.ParentTxn, refundTxn.Type, refundTxn.Method, refundTxn.EntityType, refundTxn.EntityID, refundTxn.FromAccount, refundTxn.ToAccount, refundTxn.Amount, refundTxn.Currency, refundTxn.Status, refundTxn.IdempotencyKey,
	)
	return err
}

func (p *PaymentService) SQLupdateRefundRequestStatus(ctx context.Context, refundID string, update map[string]any) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil || refundID == "" || len(update) == 0 {
		return nil
	}
	parts := make([]string, 0, len(update))
	args := make([]any, 0, len(update)+1)
	for key, value := range update {
		parts = append(parts, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value)
	}
	args = append(args, refundID)
	_, err := p.app.SQLDB.Exec(ctx, `UPDATE `+RefundsTable+` SET `+strings.Join(parts, ", ")+` WHERE _id = $`+fmt.Sprintf("%d", len(args)), args...)
	return err
}

func (p *PaymentService) SQLcreateWalletAccount(ctx context.Context, userID string) (Account, error) {
	if p == nil || p.app == nil || p.app.SQLDB == nil {
		return Account{}, nil
	}
	if _, err := p.SQLgetOrCreateAccount(ctx, userID); err != nil {
		return Account{}, err
	}
	return p.SQLgetAccountByUserID(ctx, userID)
}

func (p *PaymentService) SQLcreateTransactionRecord(ctx context.Context, txn Transaction) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil {
		return nil
	}
	if txn.ID == "" {
		txn.ID = utils.GetUUID()
	}
	_, err := p.app.SQLDB.Exec(ctx,
		`INSERT INTO `+transactionsTable+` (_id, userid, parent_txn, type, method, entity_type, entity_id, from_account, to_account, amount, currency, status, external_ref, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW())`,
		txn.ID, txn.UserID, txn.ParentTxn, txn.Type, txn.Method, txn.EntityType, txn.EntityID, txn.FromAccount, txn.ToAccount, txn.Amount, txn.Currency, txn.Status, txn.IdempotencyKey,
	)
	return err
}

func (p *PaymentService) SQLcreateJournalEntryRecord(ctx context.Context, entry JournalEntry) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil {
		return nil
	}
	if entry.ID == "" {
		entry.ID = utils.GetUUID()
	}
	_, err := p.app.SQLDB.Exec(ctx,
		`INSERT INTO `+journalTable+` (_id, txn_id, debit_account, credit_account, amount, currency, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())`,
		entry.ID, entry.TxnID, entry.DebitAccount, entry.CreditAccount, entry.Amount, entry.Currency,
	)
	return err
}

func (p *PaymentService) SQLapplyBalanceDelta(ctx context.Context, accountID string, delta int64) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil || accountID == "" {
		return nil
	}
	_, err := p.app.SQLDB.Exec(ctx, `UPDATE `+accountsTable+` SET cached_balance = cached_balance + $1, updated_at = NOW() WHERE _id = $2`, delta, accountID)
	return err
}

func (p *PaymentService) SQLsetTransactionStatus(ctx context.Context, txnID string, status string) error {
	return p.SQLsetTransactionStatusByID(ctx, txnID, status, time.Now())
}

func (p *PaymentService) SQLupdateTransactionStatus(ctx context.Context, txnID string, status string, updatedAt time.Time) error {
	return p.SQLsetTransactionStatusByID(ctx, txnID, status, updatedAt)
}

func (p *PaymentService) SQLfindTransactionByID(ctx context.Context, txnID string) (Transaction, error) {
	var txn Transaction
	if p == nil || p.app == nil || p.app.SQLDB == nil || txnID == "" {
		return txn, nil
	}
	err := p.app.SQLDB.QueryRow(ctx, `SELECT * FROM `+transactionsTable+` WHERE _id = $1 OR txn_id = $1 LIMIT 1`, txnID).Scan(&txn.ID, &txn.UserID, &txn.ParentTxn, &txn.Type, &txn.Method, &txn.EntityType, &txn.EntityID, &txn.FromAccount, &txn.ToAccount, &txn.Amount, &txn.Currency, &txn.Status, &txn.IdempotencyKey, &txn.CreatedAt, &txn.UpdatedAt)
	return txn, err
}

func (p *PaymentService) SQLmarkTransactionReversed(ctx context.Context, txnID string, updatedAt time.Time) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil || txnID == "" {
		return nil
	}
	_, err := p.app.SQLDB.Exec(ctx, `UPDATE `+transactionsTable+` SET status = 'reversed', updated_at = $1 WHERE _id = $2 OR txn_id = $2`, updatedAt, txnID)
	return err
}

func (p *PaymentService) SQLcreateTransferViews(ctx context.Context, txnID string, senderID string, recipientID string, amount int64, now time.Time) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil {
		return nil
	}
	_, err := p.app.SQLDB.Exec(ctx, `INSERT INTO transfer_views (txn_id, senderid, recipientid, amount, created_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`, txnID, senderID, recipientID, amount, now)
	return err
}

func (p *PaymentService) SQLrecordWebhookProcessing(ctx context.Context, payload *PaymentWebhookPayload) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil || payload == nil {
		return nil
	}
	_, err := p.app.SQLDB.Exec(ctx, `INSERT INTO `+webhookTable+` (transaction_id, payload, processed_at) VALUES ($1, $2, NOW()) ON CONFLICT DO NOTHING`, payload.TransactionID, payload)
	return err
}

func (p *PaymentService) SQLhasWebhookBeenProcessed(ctx context.Context, transactionID string) (bool, error) {
	if p == nil || p.app == nil || p.app.SQLDB == nil || transactionID == "" {
		return false, nil
	}
	var count int64
	err := p.app.SQLDB.QueryRow(ctx, `SELECT COUNT(*) FROM `+webhookTable+` WHERE transaction_id = $1`, transactionID).Scan(&count)
	return count > 0, err
}

func (p *PaymentService) SQLincrementTopupBalanceByUser(ctx context.Context, userID string, amount float64, updatedAt time.Time) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil || userID == "" {
		return nil
	}
	_, err := p.app.SQLDB.Exec(ctx, `UPDATE `+accountsTable+` SET cached_balance = cached_balance + $1, updated_at = $2 WHERE userid = $3`, int64(amount*100), updatedAt, userID)
	return err
}

func (p *PaymentService) SQLsetTransactionStatusByID(ctx context.Context, txnID string, status string, updatedAt time.Time) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil || txnID == "" {
		return nil
	}
	_, err := p.app.SQLDB.Exec(ctx, `UPDATE `+transactionsTable+` SET status = $1, updated_at = $2 WHERE _id = $3 OR txn_id = $3`, status, updatedAt, txnID)
	return err
}

func (p *PaymentService) SQLupdateOrderStatus(ctx context.Context, tableName string, lookupField string, orderID string, status string) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil || tableName == "" || lookupField == "" || orderID == "" {
		return nil
	}
	_, err := p.app.SQLDB.Exec(ctx, fmt.Sprintf(`UPDATE %s SET status = $1 WHERE %s = $2`, tableName, lookupField), status, orderID)
	return err
}

func (p *PaymentService) SQLupdateOrderSet(ctx context.Context, tableName string, lookupField string, orderID string, update map[string]any) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil || tableName == "" || lookupField == "" || orderID == "" || len(update) == 0 {
		return nil
	}
	parts := make([]string, 0, len(update))
	args := make([]any, 0, len(update)+2)
	for key, value := range update {
		parts = append(parts, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value)
	}
	args = append(args, orderID)
	_, err := p.app.SQLDB.Exec(ctx, fmt.Sprintf(`UPDATE %s SET %s WHERE %s = $%d`, tableName, strings.Join(parts, ", "), lookupField, len(args)), args...)
	return err
}

func (p *PaymentService) SQLdecrementInventory(ctx context.Context, tableName string, lookupField string, itemID string, incField string, qty int) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil || tableName == "" || lookupField == "" || itemID == "" || incField == "" {
		return nil
	}
	_, err := p.app.SQLDB.Exec(ctx, fmt.Sprintf(`UPDATE %s SET %s = %s - $1 WHERE %s = $2`, tableName, incField, incField, lookupField), qty, itemID)
	return err
}

func (p *PaymentService) SQLfindOrderByID(ctx context.Context, tableName string, lookupField string, orderID string, out any) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil || tableName == "" || lookupField == "" || orderID == "" || out == nil {
		return nil
	}
	_, err := p.app.SQLDB.Query(ctx, fmt.Sprintf(`SELECT * FROM %s WHERE %s = $1 LIMIT 1`, tableName, lookupField), orderID)
	return err
}

func (p *PaymentService) SQLfailTxn(ctx context.Context, txnID string) {
	if p == nil || p.app == nil || p.app.SQLDB == nil || txnID == "" {
		return
	}
	_, _ = p.app.SQLDB.Exec(ctx, `UPDATE `+transactionsTable+` SET status = 'failed', updated_at = NOW() WHERE _id = $1 OR txn_id = $1`, txnID)
}

func (p *PaymentService) SQLsuccessTxn(ctx context.Context, txnID string) {
	if p == nil || p.app == nil || p.app.SQLDB == nil || txnID == "" {
		return
	}
	_, _ = p.app.SQLDB.Exec(ctx, `UPDATE `+transactionsTable+` SET status = 'success', updated_at = NOW() WHERE _id = $1 OR txn_id = $1`, txnID)
}

func (p *PaymentService) SQLrecordGlobalLedger(ctx context.Context, txnID string, journalEntryID string, ledgerType string, reason string, amount int64, accountID string, userID string) error {
	if p == nil || p.app == nil || p.app.SQLDB == nil {
		return nil
	}
	_, err := p.app.SQLDB.Exec(ctx,
		`INSERT INTO `+globalLedgerTable+` (txn_id, journal_entry_id, type, reason, amount, account_id, userid, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`,
		txnID, journalEntryID, ledgerType, reason, amount, accountID, userID,
	)
	return err
}
