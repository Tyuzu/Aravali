-- Domain: messaging, moderation, admin flows, and operational tables

CREATE TABLE IF NOT EXISTS chats (
    chatid TEXT PRIMARY KEY,
    users TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    participants TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_message JSONB,
    read_status JSONB NOT NULL DEFAULT '{}'::JSONB,
    entitytype TEXT,
    entityid TEXT,
    last_seq BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS merechats (
    chatid TEXT PRIMARY KEY,
    users TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    participants TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_message JSONB,
    read_status JSONB NOT NULL DEFAULT '{}'::JSONB,
    entitytype TEXT,
    entityid TEXT,
    last_seq BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS messages (
    messageid TEXT PRIMARY KEY,
    chatid TEXT NOT NULL,
    roomid TEXT,
    userid TEXT NOT NULL,
    text TEXT,
    fileurl TEXT,
    filetype TEXT,
    content TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    read_by TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    reactions TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    chat_ids TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    sendername TEXT,
    avatarurl TEXT,
    media JSONB,
    editedat TIMESTAMPTZ,
    status TEXT,
    nonce TEXT,
    seq BIGINT NOT NULL DEFAULT 0,
    sender TEXT,
    senderid TEXT,
    room TEXT,
    timestamp BIGINT,
    files JSONB NOT NULL DEFAULT '[]'::JSONB
);

CREATE TABLE IF NOT EXISTS reports (
    reportid TEXT PRIMARY KEY,
    userid TEXT,
    entity_type TEXT,
    entity_id TEXT,
    reason TEXT,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS appeals (
    appealid TEXT PRIMARY KEY,
    reportid TEXT,
    userid TEXT,
    status TEXT,
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS role_applications (
    roleapplicationid TEXT PRIMARY KEY,
    userid TEXT,
    role TEXT,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS modapps (
    appid TEXT PRIMARY KEY,
    userid TEXT,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS subscribers (
    subscriberid TEXT PRIMARY KEY,
    userid TEXT,
    target_userid TEXT,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS baitoapply (
    baitoappid TEXT PRIMARY KEY,
    userid TEXT,
    baitoid TEXT,
    status TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS baitos (
    baitoid TEXT PRIMARY KEY,
    title TEXT,
    description TEXT,
    userid TEXT,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS baitoworkers (
    baitoworkerid TEXT PRIMARY KEY,
    userid TEXT,
    baitoid TEXT,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE INDEX IF NOT EXISTS idx_chats_users ON chats USING GIN (users);
CREATE INDEX IF NOT EXISTS idx_chats_participants ON chats USING GIN (participants);
CREATE INDEX IF NOT EXISTS idx_chats_updated_at ON chats (updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_merechats_users ON merechats USING GIN (users);
CREATE INDEX IF NOT EXISTS idx_merechats_participants ON merechats USING GIN (participants);
CREATE INDEX IF NOT EXISTS idx_merechats_updated_at ON merechats (updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_messages_chatid ON messages (chatid);
CREATE INDEX IF NOT EXISTS idx_messages_userid ON messages (userid);
CREATE INDEX IF NOT EXISTS idx_messages_deleted ON messages (deleted);
CREATE INDEX IF NOT EXISTS idx_messages_read_by ON messages USING GIN (read_by);
CREATE INDEX IF NOT EXISTS idx_messages_reactions ON messages USING GIN (reactions);
CREATE INDEX IF NOT EXISTS idx_messages_chat_ids ON messages USING GIN (chat_ids);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_reports_entity ON reports (entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_appeals_reportid ON appeals (reportid);
