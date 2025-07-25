-- +goose Up
CREATE INDEX idx_user_id ON likes (user_id);
CREATE INDEX idx_post_id ON likes (post_id);

-- +goose Down
DROP INDEX idx_user_id;
DROP INDEX idx_post_id;