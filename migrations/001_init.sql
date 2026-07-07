CREATE TABLE IF NOT EXISTS spaces (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    owner      TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id         TEXT PRIMARY KEY,
    space_id   TEXT NOT NULL REFERENCES spaces(id),
    author     TEXT NOT NULL,
    content    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_messages_space_id ON messages(space_id);
