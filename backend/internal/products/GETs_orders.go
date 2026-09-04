package products

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"scav/infra"
	"scav/internal/auth"
	"scav/internal/cart"
	"scav/internal/farms"
	"scav/internal/pay"
	"scav/utils"
)

/* ---------------------------------------------------- */
/* Orders placed BY the current user (buyer)            */
/* ---------------------------------------------------- */

func GetMyFarmOrders(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		userID := utils.GetUserIDFromRequest(r)

		// Parse pagination parameters
		skip := 0
		limit := 10
		if s := r.URL.Query().Get("skip"); s != "" {
			if parsed, err := strconv.Atoi(s); err == nil && parsed >= 0 {
				skip = parsed
			}
		}
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
				limit = parsed
			}
		}

		var orders []cart.FarmOrder
		if err := app.DB.FindMany(
			ctx,
			farmOrdersCollection,
			map[string]any{"userid": userID},
			&orders,
		); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, utils.M{
				"success": false,
				"message": "Failed to fetch orders",
			})
			return
		}

		// Apply pagination
		total := len(orders)
		start := skip
		end := skip + limit
		if start > total {
			start = total
		}
		if end > total {
			end = total
		}

		paginatedOrders := orders[start:end]

		// Enrich orders: populate entityName for items when possible
		enriched := make([]map[string]any, 0, len(paginatedOrders))

		// cache farm names by id to avoid repeated DB lookups
		farmNameCache := map[string]string{}

		for _, o := range paginatedOrders {
			ordMap := map[string]any{
				"orderId":         o.OrderID,
				"userid":          o.UserID,
				"farmId":          o.FarmID,
				"cropId":          o.CropID,
				"quantity":        o.Quantity,
				"priceAtPurchase": o.PriceAtPurchase,
				"createdAt":       o.CreatedAt,
				"status":          o.Status,
				"approvedBy":      o.ApprovedBy,
				"subtotal":        o.Subtotal,
				"discount":        o.Discount,
				"tax":             o.Tax,
				"delivery":        o.Delivery,
				"total":           o.Total,
				"address":         o.Address,
				"name":            o.Name,
				"phone":           o.Phone,
			}

			// copy and enrich items
			itemsOut := map[string]any{}
			for cat, items := range o.Items {
				arr := make([]map[string]any, 0, len(items))
				for _, it := range items {
					itMap := map[string]any{
						"itemId":     it.ItemID,
						"itemType":   it.ItemType,
						"entityId":   it.EntityID,
						"entityType": it.EntityType,
						"itemName":   it.ItemName,
						"quantity":   it.Quantity,
						"price":      it.Price,
						"discount":   it.Discount,
						"category":   it.Category,
						"addedAt":    it.AddedAt,
						"updatedAt":  it.UpdatedAt,
					}

					// resolve entityName if entityId present and entityType is farm (or empty)
					if it.EntityID != "" {
						name := ""
						if v, ok := farmNameCache[it.EntityID]; ok {
							name = v
						} else if it.EntityType == "farm" || it.EntityType == "" {
							f := fetchFarmByID(r.Context(), it.EntityID, app)
							if f.FarmID != "" {
								name = f.Name
							}
							farmNameCache[it.EntityID] = name
						}
						if name != "" {
							itMap["entityName"] = name
						}
					}

					arr = append(arr, itMap)
				}
				itemsOut[cat] = arr
			}

			ordMap["items"] = itemsOut

			enriched = append(enriched, ordMap)
		}

		utils.RespondWithJSON(w, http.StatusOK, utils.M{
			"success": true,
			"orders":  enriched,
			"total":   total,
			"skip":    skip,
			"limit":   limit,
		})
	}
}

/* ---------------------------------------------------- */
/* Orders coming INTO farms owned by the farmer         */
/* ---------------------------------------------------- */

func GetIncomingFarmOrders(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		userID := utils.GetUserIDFromRequest(r)

		// 1. Fetch farms owned by this user
		var myfarms []farms.Farm
		if err := app.DB.FindMany(
			ctx,
			farmsCollection,
			map[string]any{"createdBy": userID},
			&myfarms,
		); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, utils.M{
				"success": false,
				"message": "Failed to fetch farms",
			})
			return
		}

		farmIDs := make([]string, 0, len(myfarms))
		for _, f := range myfarms {
			farmIDs = append(farmIDs, f.FarmID)
		}

		if len(farmIDs) == 0 {
			utils.RespondWithJSON(w, http.StatusOK, utils.M{
				"success": true,
				"orders":  []OrderDisplay{},
			})
			return
		}

		// 2. Build filter query from URL params
		filter := map[string]any{"farmid": map[string]any{"$in": farmIDs}}

		// Filter by status
		if status := r.URL.Query().Get("status"); status != "" {
			filter["status"] = status
		}

		// Filter by date range
		if dateFrom := r.URL.Query().Get("dateFrom"); dateFrom != "" {
			if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
				filter["createdat"] = map[string]any{"$gte": t}
			}
		}

		if dateTo := r.URL.Query().Get("dateTo"); dateTo != "" {
			if t, err := time.Parse("2006-01-02", dateTo); err == nil {
				// Add one day to include all orders on that date
				t = t.Add(24 * time.Hour)
				if dateFrom := r.URL.Query().Get("dateFrom"); dateFrom != "" {
					// If there's already a $gte, we need to use $lte
					if existingDateFilter, ok := filter["createdat"].(map[string]any); ok {
						existingDateFilter["$lte"] = t
						filter["createdat"] = existingDateFilter
					}
				} else {
					filter["createdat"] = map[string]any{"$lte": t}
				}
			}
		}

		// 2. Fetch orders for those farms
		var orders []cart.FarmOrder
		if err := app.DB.FindMany(
			ctx,
			farmOrdersCollection,
			filter,
			&orders,
		); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, utils.M{
				"success": false,
				"message": "Failed to fetch orders",
			})
			return
		}

		// 3. Build frontend-friendly response and apply client-side filters
		displayOrders := make([]OrderDisplay, 0, len(orders))
		cropFilter := r.URL.Query().Get("crop")
		paymentFilter := r.URL.Query().Get("payment")

		// Pre-fetch transactions for these orders to derive accurate payment status
		orderIDs := make([]string, 0, len(orders))
		for _, o := range orders {
			orderIDs = append(orderIDs, o.OrderID)
		}

		txnByOrder := map[string]pay.Transaction{}
		if len(orderIDs) > 0 {
			var txns []pay.Transaction
			_ = app.DB.FindMany(ctx, "transactions", map[string]any{
				"entity_type": "order",
				"entity_id":   map[string]any{"$in": orderIDs},
			}, &txns)

			for _, t := range txns {
				if t.EntityID != "" {
					txnByOrder[t.EntityID] = t
				}
			}
		}
		for _, o := range orders {
			user := fetchUserByID(ctx, o.UserID, app)
			crop := fetchCropByID(ctx, o.CropID, app)
			farm := fetchFarmByID(ctx, o.FarmID, app)

			// Client-side filtering for crop (since we filter by crop name)
			if cropFilter != "" && crop.Name != cropFilter {
				continue
			}

			// Client-side filtering for payment status (prefer transaction-derived status)
			var paymentStatus string
			if txn, ok := txnByOrder[o.OrderID]; ok {
				paymentStatus = derivePaymentStatusFromTxn(&txn, o.Status)
			} else {
				paymentStatus = derivePaymentStatus(o.Status)
			}
			if paymentFilter != "" && paymentStatus != paymentFilter {
				continue
			}

			displayOrders = append(displayOrders, OrderDisplay{
				ID:           o.OrderID,
				Buyer:        user.UserID,
				Farm:         firstNonEmpty(farm.FarmID, farm.Name),
				Contact:      user.Email,
				Crop:         firstNonEmpty(crop.Name, crop.CropId),
				CropID:       crop.CropId,
				Qty:          o.Quantity,
				Unit:         crop.Unit,
				OrderDate:    o.CreatedAt.Format("2006-01-02"),
				DeliveryDate: estimateDeliveryDate(o.CreatedAt),
				Address:      firstNonEmpty(o.Address, user.Address),
				Payment:      paymentStatus,
				Status:       string(o.Status),
			})
		}
		utils.RespondWithJSON(w, http.StatusOK, utils.M{
			"success": true,
			"orders":  displayOrders,
		})
	}
}

// helper: return first non-empty string from args
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func fetchFarmByID(ctx context.Context, id string, app *infra.Deps) farms.Farm {
	var farm farms.Farm

	if id == "" {
		return farm
	}

	err := app.DB.FindOne(
		ctx,
		farmsCollection,
		map[string]any{"farmid": id},
		&farm,
	)
	if err != nil {
		return farms.Farm{}
	}

	return farm
}

func fetchUserByID(ctx context.Context, id string, app *infra.Deps) auth.User {
	var user auth.User

	if id == "" {
		return user
	}

	err := app.DB.FindOne(
		ctx,
		usersCollection,
		map[string]any{"userid": id},
		&user,
	)
	if err != nil {
		return auth.User{}
	}

	return user
}

func derivePaymentStatus(status cart.OrderStatus) string {
	normalized := string(status)

	switch normalized {
	case "paid", "delivered":
		return "paid"
	case "rejected":
		return "unpaid"
	default:
		return "pending"
	}
}

// derivePaymentStatusFromTxn returns a payment label using transaction information
func derivePaymentStatusFromTxn(txn *pay.Transaction, status cart.OrderStatus) string {
	if txn == nil {
		return derivePaymentStatus(status)
	}

	// If transaction succeeded, consider it paid
	if strings.ToLower(txn.Status) == "success" {
		return "paid"
	}

	// If method is cod and txn exists but not success, treat as pending
	if strings.ToLower(txn.Method) == "cod" {
		return "pending"
	}

	// Fallback to order-based derivation
	return derivePaymentStatus(status)
}

func fetchCropByID(ctx context.Context, id string, app *infra.Deps) farms.Crop {
	var crop farms.Crop

	if id == "" {
		return crop
	}

	err := app.DB.FindOne(
		ctx,
		cropsCollection,
		map[string]any{"cropid": id},
		&crop,
	)
	if err != nil {
		return farms.Crop{}
	}

	return crop
}
