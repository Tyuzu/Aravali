-- Domain: events, places, vendors, bookings, and deliveries

CREATE TABLE IF NOT EXISTS places (
    placeid TEXT PRIMARY KEY,
    name TEXT,
    description TEXT,
    address TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS events (
    eventid TEXT PRIMARY KEY,
    title TEXT,
    description TEXT,
    userid TEXT,
    venue TEXT,
    starts_at TIMESTAMPTZ,
    ends_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB,
    tags TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[]
);

CREATE TABLE IF NOT EXISTS vendors (
    vendorid TEXT PRIMARY KEY,
    name TEXT,
    description TEXT,
    userid TEXT,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS vendor_availability (
    vendoravailabilityid TEXT PRIMARY KEY,
    vendorid TEXT,
    day_of_week INTEGER,
    starts_at TIME,
    ends_at TIME,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS slots (
    slotid TEXT PRIMARY KEY,
    vendorid TEXT,
    starts_at TIMESTAMPTZ,
    ends_at TIMESTAMPTZ,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS date_caps (
    datecapid TEXT PRIMARY KEY,
    title TEXT,
    start_at TIMESTAMPTZ,
    end_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS bookings (
    bookingid TEXT PRIMARY KEY,
    userid TEXT,
    slotid TEXT,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS deliveries (
    deliveryid TEXT PRIMARY KEY,
    orderid TEXT,
    userid TEXT,
    status TEXT,
    driverid TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS deliveryEvents (
    deliveryeventid TEXT PRIMARY KEY,
    deliveryid TEXT,
    event_type TEXT,
    message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS delwebhooks (
    id TEXT PRIMARY KEY,
    deliveryid TEXT,
    endpoint TEXT,
    payload JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS drivers (
    driverid TEXT PRIMARY KEY,
    userid TEXT,
    status TEXT,
    location JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS driver_rejection (
    id TEXT PRIMARY KEY,
    driverid TEXT,
    reason TEXT,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS hirings (
    hiringid TEXT PRIMARY KEY,
    userid TEXT,
    vendorid TEXT,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE INDEX IF NOT EXISTS idx_events_userid ON events (userid);
CREATE INDEX IF NOT EXISTS idx_places_name ON places (name);
CREATE INDEX IF NOT EXISTS idx_vendors_status ON vendors (status);
CREATE INDEX IF NOT EXISTS idx_slots_vendorid ON slots (vendorid);
CREATE INDEX IF NOT EXISTS idx_bookings_userid ON bookings (userid);
CREATE INDEX IF NOT EXISTS idx_deliveries_orderid ON deliveries (orderid);
CREATE INDEX IF NOT EXISTS idx_deliveries_driverid ON deliveries (driverid);
