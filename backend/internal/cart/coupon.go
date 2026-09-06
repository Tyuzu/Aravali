package cart

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"scav/config/mqevent"
	"scav/infra"
	"scav/infra/mq"
	"scav/utils"
	log "scav/utils/logger"
)

const (
	couponTimeout = 5 * time.Second
)

/* ───────────────────────── Coupon Models ───────────────────────── */

func (r CouponRequest) validate() error {
	if strings.TrimSpace(r.Code) == "" {
		return errors.New("coupon code missing")
	}
	if strings.TrimSpace(r.EntityID) == "" || strings.TrimSpace(r.EntityType) == "" {
		return errors.New("entity details required")
	}
	return nil
}

type CouponResponse struct {
	Valid    bool    `json:"valid"`
	Discount float64 `json:"discount"`
	Message  string  `json:"message"`
}

/* ───────────────────────── Coupon Validation (SERVER) ───────────────────────── */

type CouponResult struct {
	DiscountAmount int64
}

type dbCoupon struct {
	Code        string  `bson:"code"`
	Active      bool    `bson:"active"`
	ExpiresAt   int64   `bson:"expiresat"`
	Type        string  `bson:"type"`  // "flat" or "percent"
	Value       float64 `bson:"value"` // ₹ or %
	MaxDiscount float64 `bson:"maxdiscount"`
}

/* ───────────────────────── Validate Coupon Handler ───────────────────────── */

// ValidateCouponHandler checks if a coupon code is applicable for an entity and returns calculated discounts.
func ValidateCouponHandler(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CouponRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, CouponResponse{
				Valid:   false,
				Message: "Invalid request body",
			})
			return
		}

		if err := req.validate(); err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, CouponResponse{
				Valid:   false,
				Message: err.Error(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), couponTimeout)
		defer cancel()

		code := strings.TrimSpace(strings.ToLower(req.Code))
		filter := map[string]any{
			"code":       code,
			"entityId":   req.EntityID,
			"entityType": strings.ToLower(req.EntityType),
			"active":     true,
		}

		coupon, err := findCouponByFilter(ctx, app, filter)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusNotFound, CouponResponse{
				Valid:   false,
				Message: "Coupon not valid for this entity",
			})
			return
		}

		if !coupon.ExpiresAt.IsZero() && time.Now().After(coupon.ExpiresAt) {
			utils.RespondWithJSON(w, http.StatusGone, CouponResponse{
				Valid:   false,
				Message: "Coupon expired",
			})
			return
		}

		discount := 0.0
		if req.Cart > 0 {
			discount = (req.Cart * coupon.Discount) / 100
		}

		if err := mq.PublishWithMeta(ctx, app.MQ, mqevent.CouponValidatedEvent, mqevent.CouponValidatedPayload{}); err != nil {
			log.Printf("ValidateCouponHandler: failed to publish coupon event: %v", err)
		}

		utils.RespondWithJSON(w, http.StatusOK, CouponResponse{
			Valid:    true,
			Discount: discount,
			Message:  "Coupon applied",
		})
	}
}
