CREATE TABLE IF NOT EXISTS urls (
    url_id SERIAL PRIMARY KEY,
    shortened VARCHAR(255) UNIQUE NOT NULL,
    original_url VARCHAR(255) NOT NULL
    );

CREATE INDEX ix_urls_shortened ON urls (shortened);