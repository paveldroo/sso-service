CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    pass_hash BLOB NOT NULL,
    is_admin BOOL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS apps (
    id text,
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);
