// File: internal/cart/cartSQLDB.go

package cart

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/internal/auth"
	"scav/internal/verticals/pay"
)

var (
	cartTable       = config.Tables.CartTable
	couponTable     = config.Tables.CouponTable
	farmOrdersTable = config.Tables.FarmOrdersTable
	ordersTable     = config.Tables.OrderTable
)

/* ───────────────────────── Cart Operations ───────────────────────── */

func findUserByID(ctx context.Context, app *infra.Deps, userID string) (auth.User, bool) {
	var user auth.User
	return user, true
}

func findCouponByFilter(ctx context.Context, app *infra.Deps, filter map[string]any) (Coupon, error) {

	var coupon Coupon
	return coupon, nil
}

func findCouponByCode(ctx context.Context, app *infra.Deps, code string) (Coupon, error) {
	var coupon Coupon
	return coupon, nil
}

func validateCouponServer(ctx context.Context, code string, subtotal int64, app *infra.Deps) (*CouponResult, error) {
	code = strings.TrimSpace(strings.ToLower(code))
	if code == "" {
		return &CouponResult{DiscountAmount: 0}, nil
	}

	coupon, err := SQLfindCouponByCode(ctx, app, code)
	if err != nil || !coupon.Active {
		return nil, errors.New("invalid coupon")
	}

	if !coupon.ExpiresAt.IsZero() && time.Now().After(coupon.ExpiresAt) {
		return nil, errors.New("coupon expired")
	}

	discount := int64(0)
	if coupon.Discount > 0 {
		raw := float64(subtotal) * (coupon.Discount / 100)
		discount = int64(raw)
	}

	if discount > subtotal {
		discount = subtotal
	}

	return &CouponResult{DiscountAmount: discount}, nil
}

func insertFarmOrderRecord(ctx context.Context, app *infra.Deps, order FarmOrder) error {

}

func insertGeneralOrderRecord(ctx context.Context, app *infra.Deps, order Order) error {

}

func fetchTransactionsByOrderIDs(ctx context.Context, app *infra.Deps, orderIDs []string) map[string]pay.Transaction {
	txnMap := make(map[string]pay.Transaction)

	return txnMap
}

func fetchUserNamesByIDs(ctx context.Context, app *infra.Deps, userIDs map[string]struct{}) map[string]string {
	nameMap := make(map[string]string)

	return nameMap
}

func getCartItemsFromDB(
	ctx context.Context,
	userID string,
	app *infra.Deps,
) ([]CartItem, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("invalid user id")
	}

	items := make([]CartItem, 0)

	return items, nil
}

func replaceCartItemsInDB(
	ctx context.Context,
	userID string,
	docs []any,
	app *infra.Deps,
) error {

}

func upsertCartItemInDB(
	ctx context.Context,
	userID string,
	item CartItem,
	app *infra.Deps,
) error {

}

func updateCartItemQuantityInDB(
	ctx context.Context,
	userID string,
	itemID string,
	category string,
	quantity int,
	entityID string,
	entityType string,
	app *infra.Deps,
) (int64, error) {
}

func deleteCartItemFromDB(
	ctx context.Context,
	userID string,
	itemID string,
	category string,
	entityID string,
	entityType string,
	app *infra.Deps,
) error {
}

func clearCartForUser(
	ctx context.Context,
	userID string,
	app *infra.Deps,
) error {
}

func getGroupedCart(
	ctx context.Context,
	userID string,
	category string,
	app *infra.Deps,
) (map[string][]CartItem, error) {
}

/* ───────────────────────── Orders Operations ───────────────────────── */

func fetchUserOrdersFromDB(
	ctx context.Context,
	userID string,
	app *infra.Deps,
) ([]Order, []FarmOrder, error) {
}

/* ───────────────────────── Item Resolution ───────────────────────── */

func resolveLookupTypeAlias(itemType string, category string) string {
	itemType = strings.ToLower(strings.TrimSpace(itemType))
	category = strings.ToLower(strings.TrimSpace(category))

	switch itemType {
	case "crop", "farm":
		return "crop"
	case "product", "book", "tool", "tools":
		return "product"
	case "menu", "food":
		return "menu"
	case "merch", "merchandise", "merchandises":
		return "merch"
	}

	switch {
	case strings.Contains(category, "crop"):
		return "crop"
	case strings.Contains(category, "product") || strings.Contains(category, "tool"):
		return "product"
	case strings.Contains(category, "menu") || strings.Contains(category, "food"):
		return "menu"
	case strings.Contains(category, "merch"):
		return "merch"
	default:
		return ""
	}
}

func lookupItemDetailsByType(
	ctx context.Context,
	itemID string,
	itemType string,
	category string,
	app *infra.Deps,
) (*ItemDetails, error) {
	itemID = strings.TrimSpace(itemID)
	if itemID == "" {
		return nil, errors.New("item id is required")
	}

	lookupType := SQLresolveLookupTypeAlias(itemType, category)

	switch lookupType {
	case "crop":
		return SQLlookupCrop(ctx, itemID, app)
	case "product":
		return SQLlookupProduct(ctx, itemID, app)
	case "menu":
		return SQLlookupMenu(ctx, itemID, app)
	case "merch":
		return SQLlookupMerchandise(ctx, itemID, app)
	default:
		return nil, errors.New("unsupported item type")
	}
}

func lookupItemDetails(
	ctx context.Context,
	itemID string,
	app *infra.Deps,
) (*ItemDetails, error) {
	lookups := []func(context.Context, string, *infra.Deps) (*ItemDetails, error){
		SQLlookupCrop,
		SQLlookupProduct,
		SQLlookupMenu,
		SQLlookupMerchandise,
	}

	for _, lookup := range lookups {
		if details, err := lookup(ctx, itemID, app); err == nil && details != nil {
			return details, nil
		}
	}

	return nil, errors.New("item not found")
}

/* ───────────────────────── Product ───────────────────────── */

func lookupProduct(
	ctx context.Context,
	productID string,
	app *infra.Deps,
) (*ItemDetails, error) {
	var product struct {
		ProductID string  `db:"productid"`
		Name      string  `db:"name"`
		Type      string  `db:"type"`
		Category  string  `db:"category"`
		Price     float64 `db:"price"`
		Discount  float64 `db:"discount"`
		Unit      string  `db:"unit"`
		Quantity  int     `db:"quantity"`
		UserID    string  `db:"userid"`
	}

	if product.Quantity <= 0 {
		return nil, errors.New("product out of stock")
	}

	if product.Price < 0 {
		return nil, errors.New("invalid product price")
	}

	category := strings.TrimSpace(product.Category)
	if category == "" {
		category = "products"
	}

	itemType := strings.TrimSpace(product.Type)
	if itemType == "" {
		itemType = "product"
	}

	return &ItemDetails{
		Name:       product.Name,
		Type:       itemType,
		Category:   category,
		Price:      product.Price,
		Discount:   SQLclampDiscount(product.Discount),
		Unit:       product.Unit,
		EntityID:   product.UserID,
		EntityType: "vendor",
		Available:  product.Quantity,
	}, nil
}

/* ───────────────────────── Crop ───────────────────────── */

func lookupCrop(
	ctx context.Context,
	cropID string,
	app *infra.Deps,
) (*ItemDetails, error) {
	var crop struct {
		CropID       string  `db:"cropid"`
		Name         string  `db:"name"`
		Category     string  `db:"category"`
		Price        float64 `db:"price"`
		Discount     float64 `db:"discount"`
		AvailableQty int     `db:"quantity"`
		Unit         string  `db:"unit"`
		FarmID       string  `db:"farmid"`
		FarmName     string  `db:"farmname"`
	}

	if crop.AvailableQty <= 0 {
		return nil, errors.New("crop out of stock")
	}

	if crop.Price < 0 {
		return nil, errors.New("invalid crop price")
	}

	farmName := crop.FarmName

	unit := strings.TrimSpace(crop.Unit)
	if unit == "" {
		unit = "kg"
	}

	itemType := strings.TrimSpace(crop.Category)
	if itemType == "" {
		itemType = "crop"
	}

	return &ItemDetails{
		Name:       crop.Name,
		Type:       itemType,
		Category:   "crops",
		Price:      crop.Price,
		Discount:   SQLclampDiscount(crop.Discount),
		Unit:       unit,
		EntityID:   crop.FarmID,
		EntityName: farmName,
		EntityType: "farm",
		Available:  crop.AvailableQty,
	}, nil
}

/* ───────────────────────── Menu ───────────────────────── */

func lookupMenu(
	ctx context.Context,
	menuID string,
	app *infra.Deps,
) (*ItemDetails, error) {

	return &ItemDetails{
		Name:       menu.Name,
		Type:       "menu",
		Category:   "menu",
		Price:      menu.Price,
		Discount:   SQLclampDiscount(menu.Discount),
		Unit:       "unit",
		EntityID:   menu.PlaceID,
		EntityName: menu.Place,
		EntityType: "place",
		Available:  menu.Stock,
	}, nil
}

/* ───────────────────────── Merchandise ───────────────────────── */

func lookupMerchandise(
	ctx context.Context,
	merchID string,
	app *infra.Deps,
) (*ItemDetails, error) {
	var merch struct {
		MerchID    string  `db:"merchid"`
		Name       string  `db:"name"`
		Price      float64 `db:"price"`
		Discount   float64 `db:"discount"`
		Stock      int     `db:"stock"`
		EntityID   string  `db:"entity_id"`
		EntityType string  `db:"entity_type"`
	}

	if merch.Stock <= 0 {
		return nil, errors.New("merchandise out of stock")
	}

	if merch.Price < 0 {
		return nil, errors.New("invalid merchandise price")
	}

	return &ItemDetails{
		Name:       merch.Name,
		Type:       "merchandise",
		Category:   "merchandise",
		Price:      merch.Price,
		Discount:   SQLclampDiscount(merch.Discount),
		Unit:       "unit",
		EntityID:   merch.EntityID,
		EntityType: merch.EntityType,
		Available:  merch.Stock,
	}, nil
}

/* ───────────────────────── Helpers ───────────────────────── */

func clampDiscount(discount float64) float64 {
	if discount < 0 {
		return 0
	}

	if discount > 100 {
		return 100
	}

	return discount
}

func buildSQLFilter(filter map[string]any) (string, []any) {
	if len(filter) == 0 {
		return "1 = 1", nil
	}
	clauses := make([]string, 0, len(filter))
	args := make([]any, 0, len(filter))
	for key, value := range filter {
		if value == nil {
			continue
		}
		clauses = append(clauses, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value)
	}
	if len(clauses) == 0 {
		return "1 = 1", nil
	}
	return strings.Join(clauses, " AND "), args
}

func buildCartFilter(
	userID,
	itemID,
	category,
	entityID,
	entityType string,
) (string, []any) {
	conditions := []string{"userid = $1", "itemId = $2"}
	args := []any{userID, itemID}
	argIdx := 3

	if category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, category)
		argIdx++
	}

	if entityID != "" {
		conditions = append(conditions, fmt.Sprintf("entityId = $%d", argIdx))
		args = append(args, entityID)
		argIdx++
	}

	if entityType != "" {
		conditions = append(conditions, fmt.Sprintf("entityType = $%d", argIdx))
		args = append(args, entityType)
		argIdx++
	}

	return strings.Join(conditions, " AND "), args
}

const maxCartQuantity = 100

func SQLfindCouponByCode(ctx context.Context, app *infra.Deps, code string) (Coupon, error) {
	return findCouponByCode(ctx, app, code)
}

func SQLbuildCartFilter(userID, itemID, category, entityID, entityType string) (string, []any) {
	return buildCartFilter(userID, itemID, category, entityID, entityType)
}

func SQLgetCartItemsFromDB(ctx context.Context, userID string, app *infra.Deps) ([]CartItem, error) {
	return getCartItemsFromDB(ctx, userID, app)
}

func SQLresolveLookupTypeAlias(itemType, category string) string {
	return resolveLookupTypeAlias(itemType, category)
}

func SQLlookupCrop(ctx context.Context, cropID string, app *infra.Deps) (*ItemDetails, error) {
	return lookupCrop(ctx, cropID, app)
}

func SQLlookupProduct(ctx context.Context, productID string, app *infra.Deps) (*ItemDetails, error) {
	return lookupProduct(ctx, productID, app)
}

func SQLlookupMenu(ctx context.Context, menuID string, app *infra.Deps) (*ItemDetails, error) {
	return lookupMenu(ctx, menuID, app)
}

func SQLlookupMerchandise(ctx context.Context, merchID string, app *infra.Deps) (*ItemDetails, error) {
	return lookupMerchandise(ctx, merchID, app)
}

func SQLclampDiscount(discount float64) float64 {
	return clampDiscount(discount)
}
