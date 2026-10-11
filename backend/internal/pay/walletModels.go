// File: internal/pay/walletModels.go

package pay

import "time"

// Meta is a generic key-value map for transaction metadata
// Keep this flexible, but never put money here.
type Meta map[string]interface{}

// =====================
// MONEY & TRANSACTIONS
// =====================

// Transaction represents a wallet or payment transaction.
// Amount is ALWAYS stored in the smallest currency unit (e.g. paise).
type Transaction struct {
	ID        string `db:"_id,omitempty" json:"id"`
	UserID    string `db:"userid,omitempty" json:"userid,omitempty"` // owner (viewer) of this txn
	ParentTxn string `db:"parent_txn,omitempty" json:"parent_txn,omitempty"`

	Type string `db:"type" json:"type"`
	// allowed:
	// topup, payment, transfer, refund
	// debit, credit (derived / per-user views)

	Method string `db:"method" json:"method"`
	// wallet, card, upi, cod, transfer, refund

	EntityType string `db:"entity_type,omitempty" json:"entity_type,omitempty"`
	EntityID   string `db:"entity_id,omitempty" json:"entity_id,omitempty"`

	FromAccount string `db:"from_account,omitempty" json:"from_account,omitempty"`
	ToAccount   string `db:"to_account,omitempty" json:"to_account,omitempty"`

	Amount   int64  `db:"amount" json:"amount"` // SMALLEST UNIT (paise)
	Currency string `db:"currency" json:"currency"`

	Status string `db:"status" json:"status"`
	// initiated, success, failed, reversed

	IdempotencyKey string `db:"external_ref,omitempty" json:"external_ref,omitempty"`

	Meta Meta `db:"meta,omitempty" json:"meta,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// =====================
// LEDGER (SOURCE OF TRUTH)
// =====================

// JournalEntry represents a double-entry ledger record.
// This is the real source of truth for money movement.
type JournalEntry struct {
	ID            string `db:"_id,omitempty" json:"id"`
	TxnID         string `db:"txn_id" json:"txn_id"`
	DebitAccount  string `db:"debit_account" json:"debit_account"`
	CreditAccount string `db:"credit_account" json:"credit_account"`

	Amount   int64  `db:"amount" json:"amount"` // SMALLEST UNIT
	Currency string `db:"currency" json:"currency"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Meta      Meta      `db:"meta,omitempty" json:"meta,omitempty"`
}

// GlobalLedger tracks total money additions and deletions across the system.
// This is used for auditing and reporting total money in circulation.
type GlobalLedger struct {
	ID       string `db:"_id,omitempty" json:"id"`
	TxnID    string `db:"txn_id" json:"txn_id"`
	Type     string `db:"type" json:"type"`     // addition | deletion | transfer
	Reason   string `db:"reason" json:"reason"` // topup | refund | payment | transfer | correction
	Amount   int64  `db:"amount" json:"amount"` // SMALLEST UNIT, always positive
	Currency string `db:"currency" json:"currency"`

	// For additions: which account received money (from external/system)
	// For deletions: which account lost money (to external/system)
	AccountID string `db:"account_id" json:"account_id"`
	UserID    string `db:"userid,omitempty" json:"userid,omitempty"`

	// Reference to the transaction or journal entry
	JournalEntryID string `db:"journal_entry_id,omitempty" json:"journal_entry_id,omitempty"`

	// Running totals for quick reporting
	TotalAdditionsUpto int64 `db:"total_additions_upto" json:"total_additions_upto"` // cumulative additions
	TotalDeletionsUpto int64 `db:"total_deletions_upto" json:"total_deletions_upto"` // cumulative deletions
	NetBalanceUpto     int64 `db:"net_balance_upto" json:"net_balance_upto"`         // additions - deletions

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Meta      Meta      `db:"meta,omitempty" json:"meta,omitempty"`
}

// =====================
// ACCOUNTS
// =====================

// Account represents a wallet account.
// CachedBalance is a PERFORMANCE CACHE, not a source of truth.
type Account struct {
	ID       string `db:"_id,omitempty" json:"id"`
	UserID   string `db:"userid" json:"userid"`
	Currency string `db:"currency" json:"currency"`

	Status string `db:"status" json:"status"`
	// active, frozen, closed

	CachedBalance int64 `db:"cached_balance" json:"cached_balance"` // SMALLEST UNIT

	Version int `db:"version" json:"version"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// =====================
// PAYMENTS
// =====================

// PayRequest is the request payload for /wallet/pay
// Amount is optional and ONLY used when entity allows custom pricing.
type PayRequest struct {
	PaymentType string `json:"paymentType"` // funding | purchase

	EntityType string `json:"entityType"`
	EntityID   string `json:"entityId"`

	Method string `json:"method"` // wallet, card, upi, cod

	Amount int64 `json:"amount,omitempty"` // SMALLEST UNIT
}

// =====================
// IDEMPOTENCY
// =====================

// IdempotencyRecord stores cached responses for safe retries.
type IdempotencyRecord struct {
	Key         string `db:"key" json:"key"`
	Method      string `db:"method" json:"method"`
	Path        string `db:"path" json:"path"`
	UserID      string `db:"userid" json:"userid"`
	RequestHash string `db:"request_hash" json:"request_hash"`

	Response map[string]interface{} `db:"response,omitempty" json:"response,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
}
