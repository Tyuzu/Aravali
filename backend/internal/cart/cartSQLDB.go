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

func SQLfindUserByID(ctx context.Context, app *infra.Deps, userID string) (auth.User, bool) {
	var user auth.User
	query := "userid = $1"
	args := []any{userID}

	if err := app.SQLDB.FindOne(ctx, "users", query, args, &user); err != nil {
		return auth.User{}, false
	}
	return user, true
}

func SQLfindCouponByFilter(ctx context.Context, app *infra.Deps, query string, args []any) (Coupon, error) {
	var coupon Coupon
	if err := app.SQLDB.FindOne(ctx, couponTable, query, args, &coupon); err != nil {
		return Coupon{}, err
	}
	return coupon, nil
}

func SQLfindCouponByCode(ctx context.Context, app *infra.Deps, code string) (Coupon, error) {
	var coupon Coupon
	query := "code = $1"
	args := []any{code}

	if err := app.SQLDB.FindOne(ctx, couponTable, query, args, &coupon); err != nil {
		return Coupon{}, err
	}
	return coupon, nil
}

func SQLvalidateCouponServer(ctx context.Context, code string, subtotal int64, app *infra.Deps) (*CouponResult, error) {
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

func SQLinsertFarmOrderRecord(ctx context.Context, app *infra.Deps, order FarmOrder) error {
	return app.SQLDB.InsertOne(ctx, farmOrdersTable, order)
}

func SQLinsertGeneralOrderRecord(ctx context.Context, app *infra.Deps, order Order) error {
	return app.SQLDB.InsertOne(ctx, ordersTable, order)
}

func SQLfetchTransactionsByOrderIDs(ctx context.Context, app *infra.Deps, orderIDs []string) map[string]pay.Transaction {
	txnMap := make(map[string]pay.Transaction)
	if len(orderIDs) == 0 {
		return txnMap
	}

	query := "entity_type = $1 AND entity_id = ANY($2)"
	args := []any{"order", orderIDs}

	var txns []pay.Transaction
	if err := app.SQLDB.FindMany(ctx, "transactions", query, args, &txns); err != nil {
		return txnMap
	}

	for _, t := range txns {
		if t.EntityID != "" {
			txnMap[t.EntityID] = t
		}
	}
	return txnMap
}

func SQLfetchUserNamesByIDs(ctx context.Context, app *infra.Deps, userIDs map[string]struct{}) map[string]string {
	nameMap := make(map[string]string)
	if len(userIDs) == 0 {
		return nameMap
	}

	ids := make([]string, 0, len(userIDs))
	for id := range userIDs {
		ids = append(ids, id)
	}

	query := "userid = ANY($1)"
	args := []any{ids}

	var users []auth.User
	if err := app.SQLDB.FindMany(ctx, "users", query, args, &users); err != nil {
		return nameMap
	}

	for _, u := range users {
		if u.Name != "" {
			nameMap[u.UserID] = u.Name
		}
	}
	return nameMap
}

func SQLgetCartItemsFromDB(
	ctx context.Context,
	userID string,
	app *infra.Deps,
) ([]CartItem, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("invalid user id")
	}

	items := make([]CartItem, 0)
	query := "userid = $1"
	args := []any{userID}

	err := app.SQLDB.FindMany(
		ctx,
		cartTable,
		query,
		args,
		&items,
	)

	if err != nil {
		return nil, err
	}

	return items, nil
}

func SQLreplaceCartItemsInDB(
	ctx context.Context,
	userID string,
	docs []any,
	app *infra.Deps,
) error {
	if strings.TrimSpace(userID) == "" {
		return errors.New("invalid user id")
	}

	query := "userid = $1"
	args := []any{userID}

	if _, err := app.SQLDB.DeleteMany(
		ctx,
		cartTable,
		query,
		args,
	); err != nil {
		return err
	}

	if len(docs) == 0 {
		return nil
	}

	return app.SQLDB.InsertMany(
		ctx,
		cartTable,
		docs,
	)
}

func SQLupsertCartItemInDB(
	ctx context.Context,
	userID string,
	item CartItem,
	app *infra.Deps,
) error {
	if strings.TrimSpace(userID) == "" {
		return errors.New("invalid user id")
	}

	if strings.TrimSpace(item.ItemID) == "" {
		return errors.New("invalid item id")
	}

	if item.Quantity <= 0 || item.Quantity > maxCartQuantity {
		return errors.New("invalid quantity")
	}

	whereClause, args := SQLbuildCartFilter(
		userID,
		item.ItemID,
		item.Category,
		item.EntityID,
		item.EntityType,
	)

	now := time.Now()

	// Check if item exists to handle increment ($inc) vs set on insert ($setOnInsert)
	var existing CartItem
	err := app.SQLDB.FindOne(ctx, cartTable, whereClause, args, &existing)

	if err == nil {
		// Existing record: Increment quantity and update metadata
		updateValues := map[string]any{
			"quantity":   existing.Quantity + item.Quantity,
			"userid":     userID,
			"itemId":     item.ItemID,
			"itemName":   item.ItemName,
			"itemType":   item.ItemType,
			"unit":       item.Unit,
			"category":   item.Category,
			"entityId":   item.EntityID,
			"entityType": item.EntityType,
			"price":      item.Price,
			"discount":   item.Discount,
			"updatedAt":  now,
		}
		_, err = app.SQLDB.UpdateOne(ctx, cartTable, whereClause, args, updateValues)
		return err
	}

	// New record: Insert with addedAt timestamp
	newCartItem := map[string]any{
		"userid":     userID,
		"itemId":     item.ItemID,
		"itemName":   item.ItemName,
		"itemType":   item.ItemType,
		"quantity":   item.Quantity,
		"unit":       item.Unit,
		"category":   item.Category,
		"entityId":   item.EntityID,
		"entityType": item.EntityType,
		"price":      item.Price,
		"discount":   item.Discount,
		"addedAt":    now,
		"updatedAt":  now,
	}

	return app.SQLDB.InsertOne(ctx, cartTable, newCartItem)
}

func SQLupdateCartItemQuantityInDB(
	ctx context.Context,
	userID string,
	itemID string,
	category string,
	quantity int,
	entityID string,
	entityType string,
	app *infra.Deps,
) (int64, error) {
	if quantity <= 0 || quantity > maxCartQuantity {
		return 0, errors.New("invalid quantity")
	}

	whereClause, args := SQLbuildCartFilter(
		userID,
		itemID,
		category,
		entityID,
		entityType,
	)

	updateValues := map[string]any{
		"quantity":  quantity,
		"updatedAt": time.Now(),
	}

	return app.SQLDB.UpdateOne(
		ctx,
		cartTable,
		whereClause,
		args,
		updateValues,
	)
}

func SQLdeleteCartItemFromDB(
	ctx context.Context,
	userID string,
	itemID string,
	category string,
	entityID string,
	entityType string,
	app *infra.Deps,
) error {
	whereClause, args := SQLbuildCartFilter(
		userID,
		itemID,
		category,
		entityID,
		entityType,
	)

	_, err := app.SQLDB.DeleteMany(
		ctx,
		cartTable,
		whereClause,
		args,
	)

	return err
}

func SQLclearCartForUser(
	ctx context.Context,
	userID string,
	app *infra.Deps,
) error {
	if strings.TrimSpace(userID) == "" {
		return errors.New("invalid user id")
	}

	query := "userid = $1"
	args := []any{userID}

	_, err := app.SQLDB.DeleteMany(
		ctx,
		cartTable,
		query,
		args,
	)

	return err
}

func SQLgetGroupedCart(
	ctx context.Context,
	userID string,
	category string,
	app *infra.Deps,
) (map[string][]CartItem, error) {
	items, err := SQLgetCartItemsFromDB(ctx, userID, app)
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]CartItem)

	for _, item := range items {
		if category != "" && item.Category != category {
			continue
		}

		grouped[item.Category] = append(
			grouped[item.Category],
			item,
		)
	}

	return grouped, nil
}

/* ───────────────────────── Orders Operations ───────────────────────── */

func SQLfetchUserOrdersFromDB(
	ctx context.Context,
	userID string,
	app *infra.Deps,
) ([]Order, []FarmOrder, error) {
	regularOrders := make([]Order, 0)
	query := "userid = $1"
	args := []any{userID}

	if err := app.SQLDB.FindMany(
		ctx,
		ordersTable,
		query,
		args,
		&regularOrders,
	); err != nil {
		return nil, nil, err
	}

	farmOrders := make([]FarmOrder, 0)

	if err := app.SQLDB.FindMany(
		ctx,
		farmOrdersTable,
		query,
		args,
		&farmOrders,
	); err != nil {
		return regularOrders, farmOrders, nil
	}

	return regularOrders, farmOrders, nil
}

/* ───────────────────────── Item Resolution ───────────────────────── */

func SQLresolveLookupTypeAlias(itemType string, category string) string {
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

func SQLlookupItemDetailsByType(
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

func SQLlookupItemDetails(
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

func SQLlookupProduct(
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

	query := "productid = $1"
	args := []any{productID}

	if err := app.SQLDB.FindOne(
		ctx,
		"products",
		query,
		args,
		&product,
	); err != nil {
		return nil, err
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

func SQLlookupCrop(
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

	query := "cropid = $1"
	args := []any{cropID}

	if err := app.SQLDB.FindOne(
		ctx,
		"crops",
		query,
		args,
		&crop,
	); err != nil {
		return nil, err
	}

	if crop.AvailableQty <= 0 {
		return nil, errors.New("crop out of stock")
	}

	if crop.Price < 0 {
		return nil, errors.New("invalid crop price")
	}

	farmName := crop.FarmName

	if farmName == "" && crop.FarmID != "" {
		var farm struct {
			Name string `db:"name"`
		}

		farmQuery := "farmid = $1"
		farmArgs := []any{crop.FarmID}

		if err := app.SQLDB.FindOne(
			ctx,
			"farms",
			farmQuery,
			farmArgs,
			&farm,
		); err == nil {
			farmName = farm.Name
		}
	}

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

func SQLlookupMenu(
	ctx context.Context,
	menuID string,
	app *infra.Deps,
) (*ItemDetails, error) {
	var menu struct {
		MenuID   string  `db:"menuid"`
		Name     string  `db:"name"`
		Price    float64 `db:"price"`
		Discount float64 `db:"discount"`
		Stock    int     `db:"stock"`
		PlaceID  string  `db:"placeid"`
		Place    string  `db:"place"`
	}

	query := "menuid = $1"
	args := []any{menuID}

	if err := app.SQLDB.FindOne(
		ctx,
		"menu",
		query,
		args,
		&menu,
	); err != nil {
		return nil, err
	}

	if menu.Stock <= 0 {
		return nil, errors.New("menu item out of stock")
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
		EntityName: menu.Place,
		EntityType: "place",
		Available:  menu.Stock,
	}, nil
}

/* ───────────────────────── Merchandise ───────────────────────── */

func SQLlookupMerchandise(
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

	query := "merchid = $1"
	args := []any{merchID}

	if err := app.SQLDB.FindOne(
		ctx,
		"merchandise",
		query,
		args,
		&merch,
	); err != nil {
		return nil, err
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

func SQLclampDiscount(discount float64) float64 {
	if discount < 0 {
		return 0
	}

	if discount > 100 {
		return 100
	}

	return discount
}

func SQLbuildCartFilter(
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
