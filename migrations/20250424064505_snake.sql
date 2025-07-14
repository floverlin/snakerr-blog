-- +goose Up
CREATE TABLE snake_game (
    user_id INTEGER PRIMARY KEY,
    record INTEGER NOT NULL,
    created_at INTEGER DEFAULT (UNIXEPOCH()),
    CONSTRAINT user_fk FOREIGN KEY (user_id) REFERENCES users (id)
);
INSERT INTO snake_game (user_id, record) SELECT id, 0 FROM users;

-- +goose Down
DROP TABLE snake_game;