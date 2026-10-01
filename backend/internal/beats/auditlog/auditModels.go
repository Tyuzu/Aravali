// File: internal/beats/auditlog/auditModels.go

package auditlog

import "time"

// AuditLog tracks sensitive operations for compliance and security
type AuditLog struct {
	ID         string                 `db:"_id,omitempty" json:"id"`
	UserID     string                 `db:"userid" json:"userid"`
	Action     string                 `db:"action" json:"action"` // e.g., "TICKET_PURCHASE", "MERCH_DELETE", "ORDER_MARK_PAID"
	EntityType string                 `db:"entityType" json:"entityType"`
	EntityID   string                 `db:"entityId" json:"entityId"`
	Changes    map[string]interface{} `db:"changes,omitempty" json:"changes,omitempty"` // What changed
	IPAddress  string                 `db:"ipAddress" json:"ipAddress"`
	UserAgent  string                 `db:"userAgent" json:"userAgent"`
	Status     string                 `db:"status" json:"status"`                     // "success", "failed", "attempted"
	Reason     string                 `db:"reason,omitempty" json:"reason,omitempty"` // Why it failed
	CreatedAt  time.Time              `db:"createdAt" json:"createdAt"`
}

// AuditAction constants for common operations
const (
	AuditActionTicketPurchase   = "TICKET_PURCHASE"
	AuditActionTicketCancel     = "TICKET_CANCEL"
	AuditActionMerchCreate      = "MERCH_CREATE"
	AuditActionMerchUpdate      = "MERCH_UPDATE"
	AuditActionMerchDelete      = "MERCH_DELETE"
	AuditActionMerchPurchase    = "MERCH_PURCHASE"
	AuditActionOrderAccept      = "ORDER_ACCEPT"
	AuditActionOrderReject      = "ORDER_REJECT"
	AuditActionOrderMarkPaid    = "ORDER_MARK_PAID"
	AuditActionOrderMarkDeliver = "ORDER_MARK_DELIVERED"
	AuditActionPaymentProcess   = "PAYMENT_PROCESS"
	AuditActionPayment          = "PAYMENT"
	AuditActionTopUp            = "TOPUP"
	AuditActionFarmCreate       = "FARM_CREATE"
	AuditActionFarmDelete       = "FARM_DELETE"
)
