-- Domain: artists, media, posts, and content discovery

CREATE TABLE IF NOT EXISTS artists (
    artistid TEXT PRIMARY KEY,
    name TEXT,
    bio TEXT,
    userids TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS albums (
    albumid TEXT PRIMARY KEY,
    title TEXT,
    artistid TEXT,
    cover_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS songs (
    songid TEXT PRIMARY KEY,
    title TEXT,
    artistid TEXT,
    albumid TEXT,
    url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS playlists (
    playlistid TEXT PRIMARY KEY,
    title TEXT,
    userid TEXT,
    songs TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS blogposts (
    postid TEXT PRIMARY KEY,
    title TEXT,
    content TEXT,
    userid TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    tags TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS feedposts (
    postid TEXT PRIMARY KEY,
    title TEXT,
    content TEXT,
    userid TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    tags TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS media (
    mediaid TEXT PRIMARY KEY,
    title TEXT,
    url TEXT,
    type TEXT,
    ownerid TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS files (
    fileid TEXT PRIMARY KEY,
    name TEXT,
    url TEXT,
    mime_type TEXT,
    size BIGINT NOT NULL DEFAULT 0,
    ownerid TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS hashtags (
    hashtagid TEXT PRIMARY KEY,
    name TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS comments (
    commentid TEXT PRIMARY KEY,
    entity_type TEXT,
    entity_id TEXT,
    userid TEXT,
    content TEXT,
    parent_id TEXT,
    likes BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS reviews (
    reviewid TEXT PRIMARY KEY,
    userid TEXT,
    entity_type TEXT,
    entity_id TEXT,
    rating INTEGER,
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS faqs (
    faqid TEXT PRIMARY KEY,
    question TEXT,
    answer TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS notices (
    noticeid TEXT PRIMARY KEY,
    title TEXT,
    content TEXT,
    userid TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS adverts (
    id TEXT PRIMARY KEY,
    title TEXT,
    description TEXT,
    image_url TEXT,
    target_url TEXT,
    status TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS analytics (
    id TEXT PRIMARY KEY,
    entity_type TEXT,
    entity_id TEXT,
    event_name TEXT,
    payload JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_artists_name ON artists (name);
CREATE INDEX IF NOT EXISTS idx_songs_artistid ON songs (artistid);
CREATE INDEX IF NOT EXISTS idx_albums_artistid ON albums (artistid);
CREATE INDEX IF NOT EXISTS idx_blogposts_userid ON blogposts (userid);
CREATE INDEX IF NOT EXISTS idx_feedposts_userid ON feedposts (userid);
CREATE INDEX IF NOT EXISTS idx_media_ownerid ON media (ownerid);
CREATE INDEX IF NOT EXISTS idx_files_ownerid ON files (ownerid);
CREATE INDEX IF NOT EXISTS idx_comments_entity ON comments (entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_reviews_entity ON reviews (entity_type, entity_id);
