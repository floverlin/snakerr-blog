-- +goose Up
CREATE TABLE comments(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    body TEXT NOT NULL,
    created_at INTEGER DEFAULT (UNIXEPOCH()),
    CONSTRAINT post_fk FOREIGN KEY (post_id) REFERENCES posts (id),
    CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE comments;