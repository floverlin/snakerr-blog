-- +goose Up
CREATE TABLE smoking (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    body TEXT,
    created_at INTEGER DEFAULT (UNIXEPOCH())

);
CREATE INDEX idx_created_at ON smoking (created_at);

-- +goose Down
DROP INDEX idx_created_at;
DROP TABLE smoking;
