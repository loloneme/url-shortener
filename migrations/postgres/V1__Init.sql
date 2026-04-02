CREATE TABLE shortened_urls (
    id BIGSERIAL PRIMARY KEY,
    short VARCHAR(10) UNIQUE NOT NULL CHECK (LENGTH(short) = 10),
    original TEXT UNIQUE NOT NULL
);

CREATE INDEX idx_short_url ON shortened_urls(short);