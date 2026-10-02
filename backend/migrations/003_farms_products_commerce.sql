-- Domain: farms, commerce, products, and payments

CREATE TABLE IF NOT EXISTS farms (
    farmid TEXT PRIMARY KEY,
    name TEXT,
    description TEXT,
    userid TEXT,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS crops (
    cropid TEXT PRIMARY KEY,
    name TEXT,
    description TEXT,
    farmid TEXT,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS cropsabout (
    cropsaboutid TEXT PRIMARY KEY,
    cropid TEXT,
    title TEXT,
    content TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS catalogue (
    catalogueid TEXT PRIMARY KEY,
    name TEXT,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS products (
    productid TEXT PRIMARY KEY,
    name TEXT,
    description TEXT,
    farmid TEXT,
    price NUMERIC(18,2) NOT NULL DEFAULT 0,
    stock BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS cart (
    cartid TEXT PRIMARY KEY,
    userid TEXT,
    items JSONB NOT NULL DEFAULT '[]'::JSONB,
    subtotal NUMERIC(18,2) NOT NULL DEFAULT 0,
    discount NUMERIC(18,2) NOT NULL DEFAULT 0,
    total NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS coupons (
    couponid TEXT PRIMARY KEY,
    code TEXT,
    discount_type TEXT,
    discount_value NUMERIC(18,2) NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    valid_from TIMESTAMPTZ,
    valid_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS orders (
    orderid TEXT PRIMARY KEY,
    userid TEXT,
    total NUMERIC(18,2) NOT NULL DEFAULT 0,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS forders (
    orderid TEXT PRIMARY KEY,
    farmid TEXT,
    userid TEXT,
    status TEXT,
    total NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS funding (
    fundingid TEXT PRIMARY KEY,
    userid TEXT,
    amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS transactions (
    txn_id TEXT PRIMARY KEY,
    userid TEXT,
    amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    type TEXT,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS globalledger (
    globalledgerid TEXT PRIMARY KEY,
    account_id TEXT,
    entry_type TEXT,
    amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS refunds (
    refundid TEXT PRIMARY KEY,
    userid TEXT,
    orderid TEXT,
    amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS stripeorders (
    stripeorderid TEXT PRIMARY KEY,
    orderid TEXT,
    payment_intent_id TEXT,
    status TEXT,
    amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS merch (
    merchid TEXT PRIMARY KEY,
    name TEXT,
    description TEXT,
    ownerid TEXT,
    price NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS menu (
    menuid TEXT PRIMARY KEY,
    name TEXT,
    description TEXT,
    vendorid TEXT,
    price NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS service (
    serviceid TEXT PRIMARY KEY,
    name TEXT,
    category TEXT,
    price NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS tickets (
    ticketid TEXT PRIMARY KEY,
    eventid TEXT,
    title TEXT,
    price NUMERIC(18,2) NOT NULL DEFAULT 0,
    quantity BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS purticks (
    purchasedticketid TEXT PRIMARY KEY,
    userid TEXT,
    eventid TEXT,
    ticketid TEXT,
    status TEXT,
    amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE INDEX IF NOT EXISTS idx_farms_userid ON farms (userid);
CREATE INDEX IF NOT EXISTS idx_crops_farmid ON crops (farmid);
CREATE INDEX IF NOT EXISTS idx_products_farmid ON products (farmid);
CREATE INDEX IF NOT EXISTS idx_cart_userid ON cart (userid);
CREATE INDEX IF NOT EXISTS idx_orders_userid ON orders (userid);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status);
CREATE INDEX IF NOT EXISTS idx_transactions_userid ON transactions (userid);
CREATE INDEX IF NOT EXISTS idx_refunds_orderid ON refunds (orderid);
CREATE INDEX IF NOT EXISTS idx_merch_ownerid ON merch (ownerid);
