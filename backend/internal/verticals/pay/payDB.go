package pay

import (
	"context"
	"errors"
	"time"

	"scav/config"
	"scav/internal/auth"
	"scav/internal/tickets"
	"scav/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var journalCollection = config.Collections.JournalCollection
var accountsCollection = config.Collections.AccountsCollection
var transactionsCollection = config.Collections.TransactionCollection
var globalLedgerCollection = config.Collections.GlobalLedgerCollection
var farmOrdersCollection = config.Collections.FarmOrdersCollection
var ticketsCollection = config.Collections.TicketsCollection
var menuCollection = config.Collections.MenuCollection
var serviceCollection = config.Collections.ServiceCollection
var productCollection = config.Collections.ProductCollection
var bookingsCollection = config.Collections.BookingsCollection
var merchCollection = config.Collections.MerchCollection
var cropsCollection = config.Collections.CropsCollection
var ordersCollection = config.Collections.OrderCollection
var usersCollection = config.Collections.UserCollection

func (p *PaymentService) getOrCreateAccount(ctx context.Context, userID string) (string, error) {
	var acc Account
	err := p.app.DB.FindOne(ctx, accountsCollection, map[string]any{"userid": userID}, &acc)
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

	if err := p.app.DB.InsertOne(ctx, accountsCollection, newAcc); err != nil {
		// race: retry read
		err = p.app.DB.FindOne(ctx, accountsCollection, map[string]any{"userid": userID}, &acc)
		return acc.ID, err
	}

	return newAcc.ID, nil
}

func (p *PaymentService) userExists(ctx context.Context, userID string) bool {
	if userID == "" {
		return false
	}

	var user auth.User
	return p.app.DB.FindOne(ctx, usersCollection, map[string]any{"userid": userID}, &user) == nil
}

func (p *PaymentService) getAccountByID(ctx context.Context, accountID string) (Account, error) {
	var acc Account
	err := p.app.DB.FindOne(ctx, accountsCollection, map[string]any{"_id": accountID}, &acc)
	return acc, err
}

func (p *PaymentService) getAccountByUserID(ctx context.Context, userID string) (Account, error) {
	var acc Account
	err := p.app.DB.FindOne(ctx, accountsCollection, map[string]any{"userid": userID}, &acc)
	return acc, err
}

func (p *PaymentService) listUserTransactions(ctx context.Context, userID string, skip int64, limit int64) ([]Transaction, error) {
	var txns []Transaction
	filter := map[string]any{
		"$or": []map[string]any{
			{"userid": userID},
			{"meta.recipient": userID},
		},
	}

	err := p.app.DB.FindMany(
		ctx,
		transactionsCollection,
		filter,
		&txns,
		options.Find().
			SetSort(map[string]int{"created_at": -1}).
			SetSkip(skip).
			SetLimit(limit),
	)
	return txns, err
}

func (p *PaymentService) fetchPriceByField(ctx context.Context, collectionName string, fieldName string, entityID string) (int64, error) {
	var v struct{ Price int64 }
	err := p.app.DB.FindOne(ctx, collectionName, map[string]any{fieldName: entityID}, &v)
	return v.Price, err
}

func (p *PaymentService) findOrderTotalByUser(ctx context.Context, orderID string, userID string) (string, int64, error) {
	var regularOrder struct {
		Total int64 `bson:"total"`
	}
	err := p.app.DB.FindOne(ctx, ordersCollection, map[string]any{"orderId": orderID, "userid": userID}, &regularOrder)
	if err == nil {
		return "regular", regularOrder.Total, nil
	}
	if err != mongo.ErrNoDocuments {
		return "", 0, err
	}

	var farmOrder struct {
		Total int64 `bson:"total"`
	}
	err = p.app.DB.FindOne(ctx, farmOrdersCollection, map[string]any{"orderid": orderID, "userid": userID}, &farmOrder)
	if err == nil {
		return "farm", farmOrder.Total, nil
	}
	if err == mongo.ErrNoDocuments {
		return "", 0, mongo.ErrNoDocuments
	}
	return "", 0, err
}

func (p *PaymentService) findOrderTotalByID(ctx context.Context, id string) (int64, error) {
	var o struct {
		Total int64 `bson:"total"`
	}
	err := p.app.DB.FindOne(ctx, ordersCollection, map[string]any{"orderId": id}, &o)
	if err == nil {
		return o.Total, nil
	}
	if err != mongo.ErrNoDocuments {
		return 0, err
	}

	var fo struct {
		PriceAtPurchase float64 `bson:"priceAtPurchase"`
	}
	err = p.app.DB.FindOne(ctx, farmOrdersCollection, map[string]any{"orderid": id}, &fo)
	if err != nil {
		return 0, err
	}
	return int64(fo.PriceAtPurchase * 100), nil
}

func (p *PaymentService) findRefundRequestByOrderID(ctx context.Context, orderID string) (OrderRefundRequest, error) {
	var refund OrderRefundRequest
	err := p.app.DB.FindOne(ctx, RefundsCollection, map[string]any{
		"order_id": orderID,
		"status":   map[string]any{"$in": []string{"pending", "approved"}},
	}, &refund)
	return refund, err
}

func (p *PaymentService) findActiveRefundRequestByOrderID(ctx context.Context, orderID string) (OrderRefundRequest, error) {
	var refund OrderRefundRequest
	err := p.app.DB.FindOne(ctx, RefundsCollection, map[string]any{
		"order_id": orderID,
		"status":   map[string]any{"$in": []string{"pending", "approved"}},
	}, &refund)
	return refund, err
}

func (p *PaymentService) createRefundRequestRecord(ctx context.Context, refundReq OrderRefundRequest) error {
	return p.app.DB.InsertOne(ctx, RefundsCollection, refundReq)
}

func (p *PaymentService) countRefundRequestsByUser(ctx context.Context, userID string) (int64, error) {
	return p.app.DB.Count(ctx, RefundsCollection, map[string]any{"userid": userID})
}

func (p *PaymentService) listRefundRequestsByUser(ctx context.Context, userID string, skip int, limit int) ([]tickets.RefundRequest, error) {
	var refunds []tickets.RefundRequest
	err := p.app.DB.FindMany(
		ctx,
		RefundsCollection,
		map[string]any{"userid": userID},
		&refunds,
		options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(map[string]any{"created_at": -1}),
	)
	return refunds, err
}

func (p *PaymentService) countRefundRequests(ctx context.Context, filter map[string]any) (int64, error) {
	return p.app.DB.Count(ctx, RefundsCollection, filter)
}

func (p *PaymentService) listRefundRequests(ctx context.Context, filter map[string]any, skip int, limit int) ([]tickets.RefundRequest, error) {
	var refunds []tickets.RefundRequest
	err := p.app.DB.FindMany(
		ctx,
		RefundsCollection,
		filter,
		&refunds,
		options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(map[string]any{"created_at": -1}),
	)
	return refunds, err
}

func (p *PaymentService) findRefundRequestByID(ctx context.Context, refundID string) (OrderRefundRequest, error) {
	var refund OrderRefundRequest
	err := p.app.DB.FindOne(ctx, RefundsCollection, map[string]any{"_id": refundID}, &refund)
	return refund, err
}

func (p *PaymentService) createRefundTransactionRecord(ctx context.Context, refundTxn Transaction) error {
	return p.app.DB.InsertOne(ctx, transactionsCollection, refundTxn)
}

func (p *PaymentService) updateRefundRequestStatus(ctx context.Context, refundID string, update map[string]any) error {
	_, err := p.app.DB.UpdateOne(ctx, RefundsCollection, map[string]any{"_id": refundID}, map[string]any{"$set": update})
	return err
}

func (p *PaymentService) createWalletAccount(ctx context.Context, userID string) (Account, error) {
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

	if err := p.app.DB.InsertOne(ctx, accountsCollection, newAcc); err != nil {
		return Account{}, err
	}
	return newAcc, nil
}

func (p *PaymentService) createTransactionRecord(ctx context.Context, txn Transaction) error {
	return p.app.DB.InsertOne(ctx, transactionsCollection, txn)
}

func (p *PaymentService) createJournalEntryRecord(ctx context.Context, entry JournalEntry) error {
	return p.app.DB.InsertOne(ctx, journalCollection, entry)
}

func (p *PaymentService) applyBalanceDelta(ctx context.Context, accountID string, delta int64) error {
	return p.app.DB.Inc(ctx, accountsCollection, map[string]any{"_id": accountID}, "cached_balance", delta)
}

func (p *PaymentService) setTransactionStatus(ctx context.Context, txnID string, status string) error {
	_, err := p.app.DB.UpdateOne(ctx, transactionsCollection,
		map[string]any{"_id": txnID},
		map[string]any{"$set": map[string]any{"status": status, "updated_at": time.Now()}},
	)
	return err
}

func (p *PaymentService) updateTransactionStatus(ctx context.Context, txnID string, status string, updatedAt time.Time) error {
	_, err := p.app.DB.UpdateOne(ctx, transactionsCollection,
		map[string]any{"_id": txnID},
		map[string]any{"$set": map[string]any{"status": status, "updated_at": updatedAt}},
	)
	return err
}

func (p *PaymentService) findTransactionByID(ctx context.Context, txnID string) (Transaction, error) {
	var txn Transaction
	err := p.app.DB.FindOne(ctx, transactionsCollection, map[string]any{"_id": txnID}, &txn)
	return txn, err
}

func (p *PaymentService) markTransactionReversed(ctx context.Context, txnID string, updatedAt time.Time) error {
	_, err := p.app.DB.UpdateOne(
		ctx,
		transactionsCollection,
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

func (p *PaymentService) createTransferViews(ctx context.Context, txnID string, senderID string, recipientID string, amount int64, now time.Time) error {
	return p.app.DB.InsertMany(ctx, transactionsCollection, []interface{}{
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

func (p *PaymentService) recordWebhookProcessing(ctx context.Context, payload *PaymentWebhookPayload) error {
	return p.app.DB.InsertOne(ctx, webhookCollection, map[string]any{
		"transactionId": payload.TransactionID,
		"orderId":       payload.OrderID,
		"userid":        payload.UserID,
		"status":        payload.Status,
		"amount":        payload.Amount,
		"processedAt":   time.Now(),
	})
}

func (p *PaymentService) hasWebhookBeenProcessed(ctx context.Context, transactionID string) (bool, error) {
	var existingWebhook map[string]any
	err := p.app.DB.FindOne(ctx, webhookCollection, map[string]any{
		"transactionId": transactionID,
	}, &existingWebhook)
	return err == nil, err
}

func (p *PaymentService) incrementTopupBalanceByUser(ctx context.Context, userID string, amount float64, updatedAt time.Time) error {
	_, err := p.app.DB.UpdateOne(ctx, accountsCollection, map[string]any{
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

func (p *PaymentService) setTransactionStatusByID(ctx context.Context, txnID string, status string, updatedAt time.Time) error {
	_, err := p.app.DB.UpdateOne(ctx, transactionsCollection, map[string]any{
		"_id": txnID,
	}, map[string]any{
		"$set": map[string]any{
			"status":     status,
			"updated_at": updatedAt,
		},
	})
	return err
}

func (p *PaymentService) updateOrderStatus(ctx context.Context, collectionName string, lookupField string, orderID string, status string) error {
	_, err := p.app.DB.UpdateOne(ctx, collectionName, map[string]any{lookupField: orderID}, map[string]any{"$set": map[string]any{"status": status}})
	return err
}

func (p *PaymentService) updateOrderSet(ctx context.Context, collectionName string, lookupField string, orderID string, update map[string]any) error {
	_, err := p.app.DB.UpdateOne(ctx, collectionName, map[string]any{lookupField: orderID}, map[string]any{"$set": update})
	return err
}

func (p *PaymentService) decrementInventory(ctx context.Context, collectionName string, lookupField string, itemID string, incField string, qty int) error {
	if itemID == "" || qty <= 0 {
		return nil
	}
	_, err := p.app.DB.UpdateOne(ctx, collectionName, map[string]any{lookupField: itemID}, map[string]any{"$inc": map[string]any{incField: -qty}})
	return err
}

func (p *PaymentService) findOrderByID(ctx context.Context, collectionName string, lookupField string, orderID string, out any) error {
	return p.app.DB.FindOne(ctx, collectionName, map[string]any{lookupField: orderID}, out)
}

func (p *PaymentService) failTxn(ctx context.Context, txnID string) {
	_, _ = p.app.DB.UpdateOne(ctx, transactionsCollection,
		map[string]any{"_id": txnID},
		map[string]any{"$set": map[string]any{"status": "failed", "updated_at": time.Now()}},
	)
}

func (p *PaymentService) successTxn(ctx context.Context, txnID string) {
	_, _ = p.app.DB.UpdateOne(ctx, transactionsCollection,
		map[string]any{"_id": txnID},
		map[string]any{"$set": map[string]any{"status": "success", "updated_at": time.Now()}},
	)
}

func (p *PaymentService) recordGlobalLedger(ctx context.Context, txnID string, journalEntryID string, ledgerType string, reason string, amount int64, accountID string, userID string) error {
	var entries []GlobalLedger

	totalAdditions := int64(0)
	totalDeletions := int64(0)

	err := p.app.DB.FindMany(ctx, globalLedgerCollection,
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

	return p.app.DB.InsertOne(ctx, globalLedgerCollection, entry)
}
