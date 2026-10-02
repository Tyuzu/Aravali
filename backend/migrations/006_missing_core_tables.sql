-- Domain: activity streams, artist events, itineraries, recipes, maps, memberships,
-- bookings tiers, and related operational tables that are referenced in config.Tables

CREATE TABLE IF NOT EXISTS activities (
    activityid TEXT PRIMARY KEY,
    placeId TEXT,
    action TEXT,
    performedBy TEXT,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    details TEXT,
    ipAddress TEXT,
    deviceInfo TEXT,
    activity_type TEXT,
    entity_id TEXT,
    entity_type TEXT,
    userid TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS artistevents (
    eventid TEXT PRIMARY KEY,
    artistid TEXT,
    title TEXT,
    date TEXT,
    venue TEXT,
    city TEXT,
    country TEXT,
    creatorid TEXT,
    ticket_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS bts (
    btsid TEXT PRIMARY KEY,
    title TEXT,
    description TEXT,
    content TEXT,
    userid TEXT,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS itinerary (
    itineraryid TEXT PRIMARY KEY,
    userid TEXT,
    name TEXT,
    description TEXT,
    start_date TEXT,
    end_date TEXT,
    status TEXT,
    published BOOLEAN NOT NULL DEFAULT FALSE,
    forked_from TEXT,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    days JSONB NOT NULL DEFAULT '[]'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS journals (
    _id TEXT PRIMARY KEY,
    txn_id TEXT,
    debit_account TEXT,
    credit_account TEXT,
    amount BIGINT NOT NULL DEFAULT 0,
    currency TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    meta JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS maps (
    mapid TEXT PRIMARY KEY,
    entity TEXT,
    title TEXT,
    config JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS memberships (
    _id TEXT PRIMARY KEY,
    placeId TEXT,
    name TEXT,
    price NUMERIC(18,2) NOT NULL DEFAULT 0,
    description TEXT,
    createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS recipes (
    recipeid TEXT PRIMARY KEY,
    userid TEXT,
    title TEXT,
    description TEXT,
    cookTime TEXT,
    cuisine TEXT,
    dietary TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    portionSize TEXT,
    season TEXT,
    tags TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    images TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    ingredients JSONB NOT NULL DEFAULT '[]'::JSONB,
    steps TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    difficulty TEXT,
    banner TEXT,
    servings INTEGER NOT NULL DEFAULT 0,
    videoUrl TEXT,
    notes TEXT,
    createdAt BIGINT NOT NULL DEFAULT 0,
    views INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS tiers (
    id TEXT PRIMARY KEY,
    entityType TEXT,
    entityId TEXT,
    name TEXT,
    price NUMERIC(18,2) NOT NULL DEFAULT 0,
    capacity INTEGER NOT NULL DEFAULT 0,
    timeRange TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    daysOfWeek INTEGER[] NOT NULL DEFAULT ARRAY[]::INTEGER[],
    features TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    createdAt BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE INDEX IF NOT EXISTS idx_activities_userid ON activities (userid);
CREATE INDEX IF NOT EXISTS idx_activities_timestamp ON activities (timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_activities_entity ON activities (entity_type, entity_id);

CREATE INDEX IF NOT EXISTS idx_artistevents_artistid ON artistevents (artistid);
CREATE INDEX IF NOT EXISTS idx_artistevents_creatorid ON artistevents (creatorid);

CREATE INDEX IF NOT EXISTS idx_bts_userid ON bts (userid);
CREATE INDEX IF NOT EXISTS idx_bts_updated_at ON bts (updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_itinerary_userid ON itinerary (userid);
CREATE INDEX IF NOT EXISTS idx_itinerary_status ON itinerary (status);
CREATE INDEX IF NOT EXISTS idx_itinerary_deleted ON itinerary (deleted);

CREATE INDEX IF NOT EXISTS idx_journals_txn_id ON journals (txn_id);
CREATE INDEX IF NOT EXISTS idx_journals_created_at ON journals (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_maps_entity ON maps (entity);
CREATE INDEX IF NOT EXISTS idx_maps_updated_at ON maps (updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_memberships_placeid ON memberships (placeId);
CREATE INDEX IF NOT EXISTS idx_memberships_created_at ON memberships (createdAt DESC);

CREATE INDEX IF NOT EXISTS idx_recipes_userid ON recipes (userid);
CREATE INDEX IF NOT EXISTS idx_recipes_title ON recipes (title);
CREATE INDEX IF NOT EXISTS idx_recipes_created_at ON recipes (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_tiers_entity ON tiers (entityType, entityId);
CREATE INDEX IF NOT EXISTS idx_tiers_created_at ON tiers (created_at DESC);
