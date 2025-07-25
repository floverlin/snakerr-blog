-- +goose Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL,
    description TEXT,
    password TEXT NOT NULL,
    created_at INTEGER DEFAULT (UNIXEPOCH())
);
CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    created_at INTEGER DEFAULT (UNIXEPOCH()),
    CONSTRAINT user_fk FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE posts;
DROP TABLE users;