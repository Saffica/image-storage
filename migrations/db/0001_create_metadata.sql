CREATE TABLE IF NOT EXISTS metadata (
    metadata_id SERIAL PRIMARY KEY,
    download_link VARCHAR (2048) UNIQUE NOT NULL,
    downloaded BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
)
---- create above / drop below ----
DROP TABLE IF EXISTS metadata