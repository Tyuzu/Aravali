package pay

import (
	"encoding/json"
	"net/http"
	"time"

	"scav/config/mqevent"
	"scav/infra/mq"
	"scav/internal/beats/auditlog"
	"scav/utils"
	log "scav/utils/logger"
)

func (p *PaymentService) Pay(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromRequest(r)

	var req PayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		auditlog.LogAction(
			ctx, p.app, r, userID,
			auditlog.AuditActionPayment,
			"payment_error", "json_decode", "failed",
			map[string]interface{}{
				"error":        err.Error(),
				"content_type": r.Header.Get("Content-Type"),
			},
		)
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	// Set default payment method if not provided
	if req.Method == "" {
		req.Method = "wallet"
	}
	req.Method = NormalizePaymentMethod(req.Method)

	// Validate allowed methods globally or forward to CashOnDelivery handler if req.Method == "cod"
	if req.Method == "cod" || req.Method == "cash_on_delivery" {
		p.handleCashOnDelivery(w, r, req, userID)
		return
	}

	// Validate required fields
	if req.PaymentType == "" || req.EntityType == "" || req.EntityID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing required fields: paymentType, entityType, or entityId")
		return
	}

	// ────────── PAYMENT RULES ──────────
	rule, ok := PaymentRules[req.PaymentType]
	if !ok {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid payment type")
		return
	}

	if !rule.AllowedEntities[req.EntityType] {
		utils.RespondWithError(w, http.StatusBadRequest, "entity not allowed for payment type")
		return
	}

	if !rule.AllowedMethods[req.Method] {
		utils.RespondWithError(w, http.StatusBadRequest, "payment method not allowed")
		return
	}

	// ────────── PRICE RESOLUTION ──────────
	resolver, err := p.resolver(req.EntityType)
	if err != nil {
		auditlog.LogAction(
			ctx, p.app, r, userID,
			auditlog.AuditActionPayment,
			"payment_error", "resolver_failed", "failed",
			map[string]interface{}{
				"entity_type": req.EntityType,
				"error":       err.Error(),
			},
		)
		utils.RespondWithError(w, http.StatusBadRequest, "unsupported entity")
		return
	}

	price, err := resolver(ctx, req.EntityID)
	if err != nil {
		auditlog.LogAction(
			ctx, p.app, r, userID,
			auditlog.AuditActionPayment,
			"payment_error", "entity_not_found", "failed",
			map[string]interface{}{
				"entity_type": req.EntityType,
				"entity_id":   req.EntityID,
				"error":       err.Error(),
			},
		)
		utils.RespondWithError(w, http.StatusNotFound, "entity not found: "+req.EntityType+" ("+req.EntityID+")")
		return
	}

	// SECURITY: Handle custom amounts carefully
	if req.Amount > 0 {
		if !rule.AllowCustomAmt {
			utils.RespondWithError(w, http.StatusBadRequest, "custom amount not allowed")
			return
		}

		// Only allow custom amounts for specific payment types (funding/donations)
		// not for purchases, orders, etc
		if req.PaymentType != "funding" && req.PaymentType != "donation" {
			utils.RespondWithError(w, http.StatusBadRequest, "custom amounts only allowed for donations")
			return
		}

		// SECURITY: Set reasonable limits on custom amounts
		const maxCustomAmount = 1000000 // 10 lakh rupees max
		if req.Amount > maxCustomAmount {
			utils.RespondWithError(w, http.StatusBadRequest, "custom amount exceeds maximum limit")
			return
		}

		price = req.Amount
	}

	if price <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid amount")
		return
	}

	// ────────── REDIS LOCK ──────────
	lockKey := "payment_lock:" + userID
	lockToken := utils.GetUUID()

	locked, err := p.app.Cache.SetNX(ctx, lockKey, []byte(lockToken), 30*time.Second)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "lock error")
		return
	}
	if !locked {
		utils.RespondWithError(w, http.StatusTooManyRequests, "retry")
		return
	}

	defer func() {
		_ = p.app.Cache.Del(ctx, lockKey)
	}()

	// ────────── ACCOUNT RESOLUTION ──────────
	userAcc, err := p.getOrCreateAccount(ctx, userID)
	if err != nil {
		if err.Error() == "user_not_found" {
			utils.RespondWithError(w, http.StatusInternalServerError, "user account not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "account error")
		return
	}

	userAccount, err := p.getAccountByID(ctx, userAcc)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "account error")
		return
	}
	if err := p.ensureAccountActive(userAccount); err != nil {
		utils.RespondWithError(w, http.StatusForbidden, "user account is not active")
		return
	}

	// ────────── PREVENT SELF-FUNDING ──────────
	if req.PaymentType == "funding" && userID == req.EntityID {
		utils.RespondWithError(w, http.StatusForbidden, "self funding not allowed")
		return
	}

	var destinationAcc string
	if req.PaymentType == "funding" {
		destinationAcc, err = p.getOrCreateAccount(ctx, req.EntityID)
		if err != nil {
			if err.Error() == "user_not_found" {
				utils.RespondWithError(w, http.StatusBadRequest, "destination user not found")
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, "destination account error")
			return
		}
	} else {
		destinationAcc, err = p.getOrCreateAccount(ctx, "merchant")
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "destination account error")
			return
		}
	}

	destinationAccount, err := p.getAccountByID(ctx, destinationAcc)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "destination account error")
		return
	}
	if err := p.ensureAccountActive(destinationAccount); err != nil {
		utils.RespondWithError(w, http.StatusForbidden, "destination account is not active")
		return
	}

	// ────────── BALANCE CHECK (WALLET ONLY) ──────────
	if req.Method == "wallet" {
		var acc Account
		if err := p.app.DB.FindOne(ctx, accountsCollection, map[string]any{"_id": userAcc}, &acc); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "account error")
			return
		}

		if acc.CachedBalance < price {
			utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
				"success": false,
				"message": "insufficient balance",
			})
			return
		}
	}

	// ────────── TRANSACTION + LEDGER ──────────
	txnID := utils.GetUUID()
	now := time.Now()

	txn := Transaction{
		ID:          txnID,
		UserID:      userID,
		Type:        "payment",
		Method:      req.Method,
		EntityType:  req.EntityType,
		EntityID:    req.EntityID,
		FromAccount: userAcc,
		ToAccount:   destinationAcc,
		Amount:      price,
		Currency:    "INR",
		Status:      "initiated",
		CreatedAt:   now,
		UpdatedAt:   now,
		Meta:        Meta{"payment_type": req.PaymentType},
	}

	if err := p.app.DB.InsertOne(ctx, transactionsCollection, txn); err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}

	j := JournalEntry{
		ID:            utils.GetUUID(),
		TxnID:         txnID,
		DebitAccount:  userAcc,
		CreditAccount: destinationAcc,
		Amount:        price,
		Currency:      "INR",
		CreatedAt:     now,
	}

	if err := p.app.DB.InsertOne(ctx, journalCollection, j); err != nil {
		p.failTxn(ctx, txnID)
		utils.RespondWithError(w, http.StatusInternalServerError, "failed")
		return
	}

	// ────────── BALANCE UPDATES ──────────
	if req.Method == "wallet" {
		if err := p.app.DB.Inc(ctx, accountsCollection, map[string]any{"_id": userAcc}, "cached_balance", -price); err != nil {
			p.failTxn(ctx, txnID)
			utils.RespondWithError(w, http.StatusInternalServerError, "failed")
			return
		}

		if err := p.app.DB.Inc(ctx, accountsCollection, map[string]any{"_id": destinationAcc}, "cached_balance", price); err != nil {
			p.failTxn(ctx, txnID)
			utils.RespondWithError(w, http.StatusInternalServerError, "failed")
			return
		}
	}

	p.successTxn(ctx, txnID)

	// Log audit trail for payment transaction
	auditlog.LogAction(
		ctx, p.app, r, userID,
		auditlog.AuditActionPayment,
		"transaction", txnID, "success",
		map[string]interface{}{
			"amount":       price,
			"method":       req.Method,
			"entity_type":  req.EntityType,
			"entity_id":    req.EntityID,
			"payment_type": req.PaymentType,
		},
	)

	if err := mq.PublishWithMeta(ctx, p.app.MQ, mqevent.PaymentDoneEvent, mqevent.PaymentDonePayload{}); err != nil {
		log.Printf("failed to publish payment event: %v", err)
	}

	// If this payment is for an order, mark the order as paid and decrement inventory.
	if req.EntityType == "order" {
		// Try regular orders first (use generic map to avoid import cycles)
		var ord map[string]any
		if err := p.app.DB.FindOne(ctx, ordersCollection, map[string]any{"orderId": req.EntityID}, &ord); err == nil {
			_, _ = p.app.DB.UpdateOne(ctx, ordersCollection, map[string]any{"orderId": req.EntityID}, map[string]any{"$set": map[string]any{"status": "paid"}})

			if itemsRaw, ok := ord["items"].(map[string]any); ok {
				for category, raw := range itemsRaw {
					itemsSlice, ok := raw.([]any)
					if !ok {
						continue
					}
					for _, it := range itemsSlice {
						itMap, ok := it.(map[string]any)
						if !ok {
							continue
						}
						itemID, _ := itMap["itemId"].(string)
						// quantity may decode as float64
						qty := 0
						if qf, ok := itMap["quantity"].(float64); ok {
							qty = int(qf)
						} else if qi, ok := itMap["quantity"].(int); ok {
							qty = qi
						}

						switch category {
						case "crops":
							if itemID != "" && qty > 0 {
								_, _ = p.app.DB.UpdateOne(ctx, cropsCollection, map[string]any{"cropid": itemID}, map[string]any{"$inc": map[string]any{"quantity": -qty}})
							}
						case "menu":
							if itemID != "" && qty > 0 {
								_, _ = p.app.DB.UpdateOne(ctx, menuCollection, map[string]any{"menuid": itemID}, map[string]any{"$inc": map[string]any{"stock": -qty}})
							}
						case "merch":
							if itemID != "" && qty > 0 {
								_, _ = p.app.DB.UpdateOne(ctx, merchCollection, map[string]any{"merchid": itemID}, map[string]any{"$inc": map[string]any{"stock": -qty}})
							}
						default:
							if itemID != "" && qty > 0 {
								_, _ = p.app.DB.UpdateOne(ctx, productCollection, map[string]any{"productid": itemID}, map[string]any{"$inc": map[string]any{"quantity": -qty}})
							}
						}
					}
				}
			}

		} else {
			// Try farm orders with generic map
			var ford map[string]any
			if err := p.app.DB.FindOne(ctx, farmOrdersCollection, map[string]any{"orderid": req.EntityID}, &ford); err == nil {
				_, _ = p.app.DB.UpdateOne(ctx, farmOrdersCollection, map[string]any{"orderid": req.EntityID}, map[string]any{"$set": map[string]any{"status": "paid"}})
				if itemsRaw, ok := ford["items"].(map[string]any); ok {
					if cropsRaw, ok := itemsRaw["crops"].([]any); ok {
						for _, it := range cropsRaw {
							itMap, ok := it.(map[string]any)
							if !ok {
								continue
							}
							itemID, _ := itMap["itemId"].(string)
							qty := 0
							if qf, ok := itMap["quantity"].(float64); ok {
								qty = int(qf)
							} else if qi, ok := itMap["quantity"].(int); ok {
								qty = qi
							}
							if itemID != "" && qty > 0 {
								_, _ = p.app.DB.UpdateOne(ctx, cropsCollection, map[string]any{"cropid": itemID}, map[string]any{"$inc": map[string]any{"quantity": -qty}})
							}
						}
					}
				}
			}
		}
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success":        true,
		"transaction_id": txnID,
	})
}
