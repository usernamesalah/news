CREATE TABLE IF NOT EXISTS news (
    id SERIAL PRIMARY KEY,
    author TEXT,
    body TEXT,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);