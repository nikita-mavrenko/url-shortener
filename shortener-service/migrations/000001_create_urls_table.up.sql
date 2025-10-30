CREATE TABLE IF NOT EXISTS urls (
    id VARCHAR(50) PRIMARY KEY,
    original_url TEXT NOT NULL
    );

CREATE INDEX ix_urls_id_hash ON urls USING HASH (id);