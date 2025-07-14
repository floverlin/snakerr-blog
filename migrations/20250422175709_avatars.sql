-- +goose Up
ALTER TABLE users ADD COLUMN avatar TEXT DEFAULT "default";

-- +goose Down
ALTER TABLE users DROP COLUMN avatar;
