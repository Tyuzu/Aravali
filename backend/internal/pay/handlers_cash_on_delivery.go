package pay

import (
	"encoding/json"
	"net/http"
	"scav/config/mqevent"
	"scav/infra/mq"
	"scav/internal/beats/auditlog"
	"scav/utils"
	log "scav/utils/logger"
	"time"
)

func (p *PaymentService) handleCashOnDelivery(w http.ResponseWriter, r *http.Request, req PayRequest, userID string) {
	ctx := r.Context()
	if req.PaymentType == "" || req.EntityType == "" || req.EntityID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing required fields: paymentType, entityType, or entityId")
		return
	}

	req.Method = NormalizePaymentMethod(req.Method)
	if req.Method == "" {
		req.Method = "cash_on_delivery"
	}

	rule, ok := PaymentRules[req.PaymentType]
	if !ok {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid payment type")
		return
	}
	if !rule.AllowedEntities[req.EntityType] {
		utils.RespondWithError(w, http.StatusBadRequest, "entity not allowed for payment type")
		return
	}
	if req.EntityType != "order" && req.EntityType != "cart" {
		utils.RespondWithError(w, http.StatusBadRequest, "cash on delivery is only supported for orders and carts")
		return
	}
	if !rule.AllowedMethods[req.Method] && !rule.AllowedMethods["cod"] && !rule.AllowedMethods["cash_on_delivery"] {
		utils.RespondWithError(w, http.StatusBadRequest, "cash on delivery not allowed for this payment type")
		return
	}

	resolver, err := p.resolver(req.EntityType)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "unsupported entity")
		return
	}

	price, err := resolver(ctx, req.EntityID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "entity not found: "+req.EntityType+" ("+req.EntityID+")")
		return
	}

	txnID := utils.GetUUID()
	now := time.Now()
	txn := Transaction{
		ID:         txnID,
		UserID:     userID,
		Type:       "payment",
		Method:     "cash_on_delivery",
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
		Amount:     price,
		Currency:   "INR",
		Status:     "pending",
		CreatedAt:  now,
		UpdatedAt:  now,
		Meta:       Meta{"payment_type": req.PaymentType},
	}
	if err := p.app.DB.InsertOne(ctx, transactionsCollection, txn); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create transaction")
		return
	}

	if req.EntityType == "order" {
		filter := map[string]any{"orderId": req.EntityID}
		update := map[string]any{
			"$set": map[string]any{
				"paymentMethod": "cash_on_delivery",
				"status":        "cod_pending",
				"updatedAt":     now,
			},
		}
		if _, err := p.app.DB.UpdateOne(ctx, "orders", filter, update); err != nil {
			auditlog.LogAction(
				ctx, p.app, r, userID,
				auditlog.AuditActionPayment,
				"payment_error", "order_update_failed", "warning",
				map[string]interface{}{
					"entity_type": req.EntityType,
					"entity_id":   req.EntityID,
					"error":       err.Error(),
				},
			)
		}
	}

	auditlog.LogAction(
		ctx, p.app, r, userID,
		auditlog.AuditActionPayment,
		"cash_on_delivery", req.EntityID, "success",
		map[string]interface{}{
			"amount":       price,
			"entity_type":  req.EntityType,
			"payment_type": req.PaymentType,
			"transaction":  txnID,
		},
	)

	if err := mq.PublishWithMeta(ctx, p.app.MQ, mqevent.CashOnDeliveryProcessedEvent, mqevent.CashOnDeliveryProcessedPayload{}); err != nil {
		log.Printf("failed to publish cash on delivery processed event: %v", err)
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// CashOnDelivery handles cash-on-delivery payment requests
func (p *PaymentService) CashOnDelivery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromRequest(r)

	var req PayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if req.Method != "" {
		req.Method = NormalizePaymentMethod(req.Method)
	}
	p.handleCashOnDelivery(w, r, req, userID)
	_ = ctx
}
