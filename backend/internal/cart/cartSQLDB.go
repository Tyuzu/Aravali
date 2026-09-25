package cart

import (
	"context"
	"errors"
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
	if err := app.DB.FindOne(ctx, "users", map[string]any{"userid": userID}, &user); err != nil {
		return auth.User{}, false
	}
	return user, true
}

func SQLfindCouponByFilter(ctx context.Context, app *infra.Deps, filter map[string]any) (Coupon, error) {
	var coupon Coupon
	if err := app.DB.FindOne(ctx, couponTable, filter, &coupon); err != nil {
		return Coupon{}, err
	}
	return coupon, nil
}

func SQLfindCouponByCode(ctx context.Context, app *infra.Deps, code string) (Coupon, error) {
	var coupon Coupon
	if err := app.DB.FindOne(ctx, couponTable, map[string]any{"code": code}, &coupon); err != nil {
		return Coupon{}, err
	}
	return coupon, nil
}

func SQLvalidateCouponServer(ctx context.Context, code string, subtotal int64, app *infra.Deps) (*CouponResult, error) {
	code = strings.TrimSpace(strings.ToLower(code))
	if code == "" {
		return &CouponResult{DiscountAmount: 0}, nil
	}

	coupon, err := findCouponByCode(ctx, app, code)
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
	return app.DB.Insert(ctx, farmOrdersTable, order)
}

func SQLinsertGeneralOrderRecord(ctx context.Context, app *infra.Deps, order Order) error {
	return app.DB.Insert(ctx, ordersTable, order)
}

func SQLfetchTransactionsByOrderIDs(ctx context.Context, app *infra.Deps, orderIDs []string) map[string]pay.Transaction {
	txnMap := make(map[string]pay.Transaction)
	if len(orderIDs) == 0 {
		return txnMap
	}

	var txns []pay.Transaction
	if err := app.DB.FindMany(ctx, "transactions", map[string]any{
		"entity_type": "order",
		"entity_id":   map[string]any{"$in": orderIDs},
	}, &txns); err != nil {
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

	var users []auth.User
	if err := app.DB.FindMany(ctx, "users", map[string]any{"userid": map[string]any{"$in": ids}}, &users); err != nil {
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

	err := app.DB.FindMany(
		ctx,
		cartTable,
		map[string]any{"userid": userID},
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

	/*
		IMPORTANT:

		This operation is still delete + insert because the current DB
		abstraction shown in your code does not expose Mongo transactions.

		For a production system, the preferred implementation is a MongoDB
		transaction.

		We at least validate the complete document set before reaching here,
		so malformed data does not partially enter the cart.
	*/

	if _, err := app.DB.Delete(
		ctx,
		cartTable,
		map[string]any{"userid": userID},
	); err != nil {
		return err
	}

	if len(docs) == 0 {
		return nil
	}

	return app.DB.InsertMany(
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

	/*
		Cart identity is:

		    user + item + category + entity

		This prevents two unrelated entities using the same item ID from
		colliding in the cart.
	*/
	filter := buildCartFilter(
		userID,
		item.ItemID,
		item.Category,
		item.EntityID,
		item.EntityType,
	)

	now := time.Now()

	update := map[string]any{
		"$inc": map[string]any{
			"quantity": item.Quantity,
		},
		"$set": map[string]any{
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
		},
		"$setOnInsert": map[string]any{
			"addedAt": now,
		},
	}

	return app.DB.Upsert(
		ctx,
		cartTable,
		filter,
		update,
	)
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
) (any, error) {
	if quantity <= 0 || quantity > maxCartQuantity {
		return nil, errors.New("invalid quantity")
	}

	filter := buildCartFilter(
		userID,
		itemID,
		category,
		entityID,
		entityType,
	)

	update := map[string]any{
		"$set": map[string]any{
			"quantity":  quantity,
			"updatedAt": time.Now(),
		},
	}

	return app.DB.Update(
		ctx,
		cartTable,
		filter,
		update,
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
	filter := buildCartFilter(
		userID,
		itemID,
		category,
		entityID,
		entityType,
	)

	_, err := app.DB.Delete(
		ctx,
		cartTable,
		filter,
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

	_, err := app.DB.Delete(
		ctx,
		cartTable,
		map[string]any{"userid": userID},
	)

	return err
}

func SQLgetGroupedCart(
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

	if err := app.DB.FindMany(
		ctx,
		ordersTable,
		map[string]any{"userid": userID},
		&regularOrders,
	); err != nil {
		return nil, nil, err
	}

	farmOrders := make([]FarmOrder, 0)

	if err := app.DB.FindMany(
		ctx,
		farmOrdersTable,
		map[string]any{"userid": userID},
		&farmOrders,
	); err != nil {
		/*
				Preserve your existing compatibility behaviour here.
			 Ideally this should eventually return the error instead of
			 silently pretending there are no farm orders.
		*/
		return regularOrders, farmOrders, nil
	}

	return regularOrders, farmOrders, nil
}

/* ───────────────────────── Item Resolution ───────────────────────── */

/*
	NEVER resolve an item by searching every table.

	The caller knows the item type. Use that information.

	This prevents:

	    productid == cropid

	from accidentally resolving a crop as a product.
*/

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

	lookupType := resolveLookupTypeAlias(itemType, category)

	switch lookupType {
	case "crop":
		return lookupCrop(ctx, itemID, app)
	case "product":
		return lookupProduct(ctx, itemID, app)
	case "menu":
		return lookupMenu(ctx, itemID, app)
	case "merch":
		return lookupMerchandise(ctx, itemID, app)
	default:
		return nil, errors.New("unsupported item type")
	}
}

/*
Compatibility helper.

Do not use this for security-sensitive operations when itemType is
available. It remains useful for old internal callers.
*/
func SQLlookupItemDetails(
	ctx context.Context,
	itemID string,
	app *infra.Deps,
) (*ItemDetails, error) {
	lookups := []func(context.Context, string, *infra.Deps) (*ItemDetails, error){
		lookupCrop,
		lookupProduct,
		lookupMenu,
		lookupMerchandise,
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
		ProductID string  `bson:"productid"`
		Name      string  `bson:"name"`
		Type      string  `bson:"type"`
		Category  string  `bson:"category"`
		Price     float64 `bson:"price"`
		Discount  float64 `bson:"discount"`
		Unit      string  `bson:"unit"`
		Quantity  int     `bson:"quantity"`
		UserID    string  `bson:"userid"`
	}

	if err := app.DB.FindOne(
		ctx,
		"products",
		map[string]any{"productid": productID},
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
		Discount:   clampDiscount(product.Discount),
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
		CropID       string  `bson:"cropid"`
		Name         string  `bson:"name"`
		Category     string  `bson:"category"`
		Price        float64 `bson:"price"`
		Discount     float64 `bson:"discount"`
		AvailableQty int     `bson:"quantity"`
		Unit         string  `bson:"unit"`
		FarmID       string  `bson:"farmid"`
		FarmName     string  `bson:"farmname"`
	}

	if err := app.DB.FindOne(
		ctx,
		"crops",
		map[string]any{"cropid": cropID},
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
			Name string `bson:"name"`
		}

		if err := app.DB.FindOne(
			ctx,
			"farms",
			map[string]any{"farmid": crop.FarmID},
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
		Discount:   clampDiscount(crop.Discount),
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
		MenuID   string  `bson:"menuid"`
		Name     string  `bson:"name"`
		Price    float64 `bson:"price"`
		Discount float64 `bson:"discount"`
		Stock    int     `bson:"stock"`
		PlaceID  string  `bson:"placeid"`
		Place    string  `bson:"place"`
	}

	if err := app.DB.FindOne(
		ctx,
		"menu",
		map[string]any{"menuid": menuID},
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
		Discount:   clampDiscount(menu.Discount),
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
		MerchID    string  `bson:"merchid"`
		Name       string  `bson:"name"`
		Price      float64 `bson:"price"`
		Discount   float64 `bson:"discount"`
		Stock      int     `bson:"stock"`
		EntityID   string  `bson:"entity_id"`
		EntityType string  `bson:"entity_type"`
	}

	if err := app.DB.FindOne(
		ctx,
		"merchandise",
		map[string]any{"merchid": merchID},
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
		Discount:   clampDiscount(merch.Discount),
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
) map[string]any {
	filter := map[string]any{
		"userid": userID,
		"itemId": itemID,
	}

	if category != "" {
		filter["category"] = category
	}

	if entityID != "" {
		filter["entityId"] = entityID
	}

	if entityType != "" {
		filter["entityType"] = entityType
	}

	return filter
}
