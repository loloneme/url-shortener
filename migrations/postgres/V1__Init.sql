CREATE TABLE shortened_urls (
    id BIGSERIAL PRIMARY KEY,
    short VARCHAR(10) UNIQUE CHECK(LENGTH(short) = 10),
    original TEXT UNIQUE NOT NULL
);