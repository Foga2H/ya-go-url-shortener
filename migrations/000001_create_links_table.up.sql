CREATE TABLE IF NOT EXISTS links (
    uuid TEXT PRIMARY KEY,
    short_url TEXT NOT NULL UNIQUE,
    original_url TEXT NOT NULL
);

CREATE INDEX idx_links_short_url ON links(short_url);