-- +goose Up
CREATE TABLE followers(
    follower INTEGER NOT NULL,
    followed INTEGER NOT NULL,
    PRIMARY KEY (follower, followed),
    CONSTRAINT follower_fk FOREIGN KEY (follower) REFERENCES users (id),
    CONSTRAINT followed_fk FOREIGN KEY (followed) REFERENCES users (id)
);

-- +goose Down
DROP TABLE followers;
