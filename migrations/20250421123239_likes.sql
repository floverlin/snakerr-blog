-- +goose Up
CREATE TABLE likes (
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, post_id),
    CONSTRAINT user_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT post_fk FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE likes;
