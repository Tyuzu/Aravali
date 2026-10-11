// File: internal/cart/cartSQLDB.go

package cart

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/internal/auth"
	"scav/internal/pay"

	"github.com/jackc/pgx/v5"
)

var (
	cartTable       = config.Tables.CartTable
	couponTable     = config.Tables.CouponTable
	farmOrdersTable = config.Tables.FarmOrdersTable
	ordersTable     = config.Tables.OrderTable
)

/* ───────────────────────── Cart Operations ───────────────────────── */

func findUserByID(ctx context.Context, app *infra.Deps, userID string) (auth.User, bool) {
	userID = strings.TrimSpace(userID)
	if userID == "" || app == nil || app.SQLDB == nil {
		return auth.User{}, false
	}

	var user auth.User
	err := app.SQLDB.QueryRow(
		ctx,
		`SELECT userid, username, email, name, phone_number, address, avatar, banner, role, online, is_verified, email_verified, followerscount, followscount, wallet_balance, created_at, updated_at, last_login FROM `+config.Tables.UserTable+` WHERE userid = $1 LIMIT 1`,
		userID,
	).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.PhoneNumber,
		&user.Address,
		&user.Avatar,
		&user.Banner,
		&user.Role,
		&user.Online,
		&user.IsVerified,
		&user.EmailVerified,
		&user.FollowersCount,
		&user.FollowingCount,
		&user.WalletBalance,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.User{}, false
		}
		return auth.User{}, false
	}
	return user, true
}

func findCouponByFilter(ctx context.Context, app *infra.Deps, filter map[string]any) (Coupon, error) {
	if app == nil || app.SQLDB == nil {
		return Coupon{}, errors.New("database unavailable")
	}
	query, args := buildSQLFilter(filter)
	rows, err := app.SQLDB.Query(ctx, `SELECT code, discount, expiresat, active, entityid, entitytype FROM `+couponTable+` WHERE `+query+` LIMIT 1`, args...)
	if err != nil {
		return Coupon{}, err
	}
	defer rows.Close()
	if !rows.Next() {
		return Coupon{}, errors.New("coupon not found")
	}
	var coupon Coupon
	var expiresAt time.Time
	if err := rows.Scan(&coupon.Code, &coupon.Discount, &expiresAt, &coupon.Active, &coupon.EntityID, &coupon.EntityType); err != nil {
		return Coupon{}, err
	}
	coupon.ExpiresAt = expiresAt
	return coupon, nil
}

func findCouponByCode(ctx context.Context, app *infra.Deps, code string) (Coupon, error) {
	return findCouponByFilter(ctx, app, map[string]any{"code": strings.TrimSpace(strings.ToLower(code))})
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
	if app == nil || app.SQLDB == nil {
		return nil
	}
	payload, err := json.Marshal(map[string]any{
		"orderid":  order.OrderID,
		"farmid":   order.FarmID,
		"cropid":   order.CropID,
		"userid":   order.UserID,
		"items":    order.Items,
		"quantity": order.Quantity,
		"subtotal": order.Subtotal,
		"discount": order.Discount,
		"tax":      order.Tax,
		"delivery": order.Delivery,
		"total":    order.Total,
		"address":  order.Address,
		"name":     order.Name,
		"phone":    order.Phone,
		"status":   string(order.Status),
		"approved": order.ApprovedBy,
	})
	if err != nil {
		return err
	}
	_, err = app.SQLDB.Exec(
		ctx,
		`INSERT INTO `+farmOrdersTable+` (orderid, farmid, userid, status, total, created_at, updated_at, metadata) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (orderid) DO UPDATE SET farmid = EXCLUDED.farmid, userid = EXCLUDED.userid, status = EXCLUDED.status, total = EXCLUDED.total, updated_at = NOW(), metadata = EXCLUDED.metadata`,
		order.OrderID,
		order.FarmID,
		order.UserID,
		string(order.Status),
		float64(order.Total)/100,
		order.CreatedAt,
		time.Now(),
		payload,
	)
	return err
}

func insertGeneralOrderRecord(ctx context.Context, app *infra.Deps, order Order) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	payload, err := json.Marshal(map[string]any{
		"orderid":       order.OrderID,
		"ordertype":     order.OrderType,
		"userid":        order.UserID,
		"items":         order.Items,
		"address":       order.Address,
		"paymentmethod": order.PaymentMethod,
		"status":        order.Status,
		"approvedby":    order.ApprovedBy,
		"subtotal":      order.Subtotal,
		"discount":      order.Discount,
		"tax":           order.Tax,
		"delivery":      order.Delivery,
		"total":         order.Total,
		"name":          order.Name,
		"phone":         order.Phone,
	})
	if err != nil {
		return err
	}
	_, err = app.SQLDB.Exec(
		ctx,
		`INSERT INTO `+ordersTable+` (orderid, userid, total, status, created_at, updated_at, metadata) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (orderid) DO UPDATE SET userid = EXCLUDED.userid, total = EXCLUDED.total, status = EXCLUDED.status, updated_at = NOW(), metadata = EXCLUDED.metadata`,
		order.OrderID,
		order.UserID,
		float64(order.Total)/100,
		order.Status,
		order.CreatedAt,
		time.Now(),
		payload,
	)
	return err
}

func fetchTransactionsByOrderIDs(ctx context.Context, app *infra.Deps, orderIDs []string) map[string]pay.Transaction {
	if app == nil || app.SQLDB == nil || len(orderIDs) == 0 {
		return map[string]pay.Transaction{}
	}
	rows, err := app.SQLDB.Query(
		ctx,
		`SELECT entity_id, userid, type, method, status, amount FROM transactions WHERE entity_type = 'order' AND entity_id = ANY($1)`,
		orderIDs,
	)
	if err != nil {
		return map[string]pay.Transaction{}
	}
	defer rows.Close()

	result := make(map[string]pay.Transaction)
	for rows.Next() {
		var txn pay.Transaction
		if err := rows.Scan(&txn.EntityID, &txn.UserID, &txn.Type, &txn.Method, &txn.Status, &txn.Amount); err != nil {
			continue
		}
		result[txn.EntityID] = txn
	}
	return result
}

func fetchUserNamesByIDs(ctx context.Context, app *infra.Deps, userIDs map[string]struct{}) map[string]string {
	if app == nil || app.SQLDB == nil || len(userIDs) == 0 {
		return map[string]string{}
	}
	ids := make([]string, 0, len(userIDs))
	for id := range userIDs {
		if strings.TrimSpace(id) != "" {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return map[string]string{}
	}

	rows, err := app.SQLDB.Query(ctx, `SELECT userid, name FROM `+config.Tables.UserTable+` WHERE userid = ANY($1)`, ids)
	if err != nil {
		return map[string]string{}
	}
	defer rows.Close()

	result := make(map[string]string, len(ids))
	for rows.Next() {
		var userID, name string
		if err := rows.Scan(&userID, &name); err != nil {
			continue
		}
		result[userID] = name
	}
	return result
}

func coerceCartDocs(docs []any) ([]CartItem, error) {
	items := make([]CartItem, 0, len(docs))
	for _, doc := range docs {
		switch v := doc.(type) {
		case CartItem:
			items = append(items, v)
		case map[string]any:
			blob, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			var item CartItem
			if err := json.Unmarshal(blob, &item); err != nil {
				return nil, err
			}
			items = append(items, item)
		default:
			blob, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			var item CartItem
			if err := json.Unmarshal(blob, &item); err != nil {
				return nil, err
			}
			items = append(items, item)
		}
	}
	return items, nil
}

func normalizeCartItems(userID string, items []CartItem) []CartItem {
	normalized := make([]CartItem, 0, len(items))
	for _, item := range items {
		item.UserID = userID
		if strings.TrimSpace(item.ItemID) == "" {
			continue
		}
		if item.AddedAt.IsZero() && !item.UpdatedAt.IsZero() {
			item.AddedAt = item.UpdatedAt
		}
		if item.UpdatedAt.IsZero() {
			item.UpdatedAt = item.AddedAt
		}
		normalized = append(normalized, item)
	}
	return normalized
}

func getCartItemsFromDB(
	ctx context.Context,
	userID string,
	app *infra.Deps,
) ([]CartItem, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("invalid user id")
	}
	if app == nil || app.SQLDB == nil {
		return []CartItem{}, nil
	}

	var raw []byte
	err := app.SQLDB.QueryRow(ctx, `SELECT items FROM `+cartTable+` WHERE userid = $1 LIMIT 1`, userID).Scan(&raw)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []CartItem{}, nil
		}
		return nil, err
	}
	if len(raw) == 0 || string(raw) == "null" || string(raw) == "[]" || strings.TrimSpace(string(raw)) == "" {
		return []CartItem{}, nil
	}

	var items []CartItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	for i := range items {
		items[i].UserID = userID
	}
	return items, nil
}

func replaceCartItemsInDB(
	ctx context.Context,
	userID string,
	docs []any,
	app *infra.Deps,
) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return errors.New("invalid user id")
	}
	if app == nil || app.SQLDB == nil {
		return nil
	}

	items, err := coerceCartDocs(docs)
	if err != nil {
		return err
	}
	items = normalizeCartItems(userID, items)
	payload, err := json.Marshal(items)
	if err != nil {
		return err
	}
	cartID := "cart:" + userID
	_, err = app.SQLDB.Exec(
		ctx,
		`INSERT INTO `+cartTable+` (cartid, userid, items, subtotal, discount, total, created_at, updated_at) VALUES ($1, $2, $3::jsonb, 0, 0, 0, NOW(), NOW()) ON CONFLICT (cartid) DO UPDATE SET userid = EXCLUDED.userid, items = EXCLUDED.items, updated_at = NOW()`,
		cartID,
		userID,
		string(payload),
	)
	return err
}

func upsertCartItemInDB(
	ctx context.Context,
	userID string,
	item CartItem,
	app *infra.Deps,
) error {
	items, err := getCartItemsFromDB(ctx, userID, app)
	if err != nil {
		return err
	}

	item.UserID = userID
	if item.AddedAt.IsZero() {
		item.AddedAt = time.Now()
	}
	if item.UpdatedAt.IsZero() {
		item.UpdatedAt = item.AddedAt
	}

	updated := false
	for i, current := range items {
		if current.ItemID == item.ItemID && current.Category == item.Category && current.EntityID == item.EntityID && current.EntityType == item.EntityType {
			items[i] = item
			updated = true
			break
		}
	}
	if !updated {
		items = append(items, item)
	}

	docs := make([]any, 0, len(items))
	for i := range items {
		docs = append(docs, items[i])
	}
	return replaceCartItemsInDB(ctx, userID, docs, app)
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
	items, err := getCartItemsFromDB(ctx, userID, app)
	if err != nil {
		return 0, err
	}
	updated := int64(0)
	filtered := make([]CartItem, 0, len(items))
	for _, item := range items {
		match := item.ItemID == itemID && strings.EqualFold(strings.TrimSpace(item.Category), strings.TrimSpace(category)) && strings.EqualFold(strings.TrimSpace(item.EntityID), strings.TrimSpace(entityID)) && strings.EqualFold(strings.TrimSpace(item.EntityType), strings.TrimSpace(entityType))
		if match {
			if quantity <= 0 {
				updated = 1
				continue
			}
			item.Quantity = quantity
			item.UpdatedAt = time.Now()
			filtered = append(filtered, item)
			updated = 1
			continue
		}
		filtered = append(filtered, item)
	}
	if updated == 0 {
		return 0, errors.New("item not found")
	}
	docs := make([]any, 0, len(filtered))
	for i := range filtered {
		docs = append(docs, filtered[i])
	}
	if err := replaceCartItemsInDB(ctx, userID, docs, app); err != nil {
		return 0, err
	}
	return updated, nil
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
	items, err := getCartItemsFromDB(ctx, userID, app)
	if err != nil {
		return err
	}
	filtered := make([]CartItem, 0, len(items))
	for _, item := range items {
		if item.ItemID == itemID && strings.EqualFold(strings.TrimSpace(item.Category), strings.TrimSpace(category)) && strings.EqualFold(strings.TrimSpace(item.EntityID), strings.TrimSpace(entityID)) && strings.EqualFold(strings.TrimSpace(item.EntityType), strings.TrimSpace(entityType)) {
			continue
		}
		filtered = append(filtered, item)
	}
	docs := make([]any, 0, len(filtered))
	for i := range filtered {
		docs = append(docs, filtered[i])
	}
	return replaceCartItemsInDB(ctx, userID, docs, app)
}

func clearCartForUser(
	ctx context.Context,
	userID string,
	app *infra.Deps,
) error {
	if strings.TrimSpace(userID) == "" {
		return errors.New("invalid user id")
	}
	if app == nil || app.SQLDB == nil {
		return nil
	}
	cartID := "cart:" + userID
	_, err := app.SQLDB.Exec(
		ctx,
		`INSERT INTO `+cartTable+` (cartid, userid, items, subtotal, discount, total, created_at, updated_at) VALUES ($1, $2, '[]'::jsonb, 0, 0, 0, NOW(), NOW()) ON CONFLICT (cartid) DO UPDATE SET userid = EXCLUDED.userid, items = EXCLUDED.items, subtotal = 0, discount = 0, total = 0, updated_at = NOW()`,
		cartID,
		userID,
	)
	return err
}

func getGroupedCart(
	ctx context.Context,
	userID string,
	category string,
	app *infra.Deps,
) (map[string][]CartItem, error) {
	items, err := getCartItemsFromDB(ctx, userID, app)
	if err != nil {
		return nil, err
	}
	grouped := make(map[string][]CartItem)
	for _, item := range items {
		if strings.TrimSpace(category) != "" && !strings.EqualFold(strings.TrimSpace(item.Category), strings.TrimSpace(category)) {
			continue
		}
		grouped[item.Category] = append(grouped[item.Category], item)
	}
	return grouped, nil
}

/* ───────────────────────── Orders Operations ───────────────────────── */

func fetchUserOrdersFromDB(
	ctx context.Context,
	userID string,
	app *infra.Deps,
) ([]Order, []FarmOrder, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, nil, errors.New("invalid user id")
	}
	if app == nil || app.SQLDB == nil {
		return []Order{}, []FarmOrder{}, nil
	}

	regularOrders := make([]Order, 0)
	rows, err := app.SQLDB.Query(ctx, `SELECT orderid, userid, total, status, created_at, metadata FROM `+ordersTable+` WHERE userid = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order Order
		var metadata []byte
		if err := rows.Scan(&order.OrderID, &order.UserID, &order.Total, &order.Status, &order.CreatedAt, &metadata); err != nil {
			continue
		}
		if len(metadata) > 0 {
			_ = json.Unmarshal(metadata, &order)
		}
		regularOrders = append(regularOrders, order)
	}

	farmOrders := make([]FarmOrder, 0)
	farmRows, err := app.SQLDB.Query(ctx, `SELECT orderid, farmid, userid, status, total, created_at, metadata FROM `+farmOrdersTable+` WHERE userid = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return regularOrders, nil, err
	}
	defer farmRows.Close()
	for farmRows.Next() {
		var order FarmOrder
		var metadata []byte
		if err := farmRows.Scan(&order.OrderID, &order.FarmID, &order.UserID, &order.Status, &order.Total, &order.CreatedAt, &metadata); err != nil {
			continue
		}
		if len(metadata) > 0 {
			_ = json.Unmarshal(metadata, &order)
		}
		farmOrders = append(farmOrders, order)
	}
	return regularOrders, farmOrders, nil
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
	_ context.Context,
	productID string,
	app *infra.Deps,
) (*ItemDetails, error) {
	if strings.TrimSpace(productID) == "" {
		return nil, errors.New("product id is required")
	}
	if app == nil || app.SQLDB == nil {
		return nil, errors.New("product lookup unavailable")
	}

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
	_ context.Context,
	cropID string,
	app *infra.Deps,
) (*ItemDetails, error) {
	if strings.TrimSpace(cropID) == "" {
		return nil, errors.New("crop id is required")
	}
	if app == nil || app.SQLDB == nil {
		return nil, errors.New("crop lookup unavailable")
	}

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
	if strings.TrimSpace(menuID) == "" {
		return nil, errors.New("menu id is required")
	}
	if app == nil || app.SQLDB == nil {
		return nil, errors.New("menu lookup unavailable")
	}

	var menu struct {
		MenuID      string  `db:"menuid"`
		PlaceID     string  `db:"placeid"`
		Name        string  `db:"name"`
		Price       float64 `db:"price"`
		Discount    float64 `db:"discount"`
		Stock       int     `db:"stock"`
		MenuPhoto   string  `db:"menu_pic"`
		Description string  `db:"description"`
		UserID      string  `db:"userid"`
	}

	row := app.SQLDB.QueryRow(
		ctx,
		`SELECT menuid, placeid, name, price, discount, stock, menu_pic, description, userid FROM `+config.Tables.MenuTable+` WHERE menuid = $1 LIMIT 1`,
		menuID,
	)
	if err := row.Scan(&menu.MenuID, &menu.PlaceID, &menu.Name, &menu.Price, &menu.Discount, &menu.Stock, &menu.MenuPhoto, &menu.Description, &menu.UserID); err != nil {
		return nil, err
	}
	if menu.Stock <= 0 {
		return nil, errors.New("menu out of stock")
	}
	if menu.Price < 0 {
		return nil, errors.New("invalid menu price")
	}

	return &ItemDetails{
		Name:       menu.Name,
		Type:       "menu",
		Category:   "menu",
		Price:      menu.Price,
		Discount:   SQLclampDiscount(menu.Discount),
		Unit:       "unit",
		EntityID:   menu.PlaceID,
		EntityType: "place",
		Available:  menu.Stock,
	}, nil
}

/* ───────────────────────── Merchandise ───────────────────────── */

func lookupMerchandise(
	_ context.Context,
	merchID string,
	app *infra.Deps,
) (*ItemDetails, error) {
	if strings.TrimSpace(merchID) == "" {
		return nil, errors.New("merchandise id is required")
	}
	if app == nil || app.SQLDB == nil {
		return nil, errors.New("merchandise lookup unavailable")
	}

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
