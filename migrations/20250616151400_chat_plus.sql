-- +goose Up
CREATE TABLE messages (
    from_ INTEGER NOT NULL,
    to_ INTEGER NOT NULL,
    body TEXT NOT NULL,
    created_at INTEGER DEFAULT (UNIXEPOCH())
);
CREATE TABLE chats (
    me INTEGER NOT NULL,
    dialogist INTEGER NOT NULL,
    readed INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    PRIMARY KEY (me, dialogist)
);

-- +goose Down
DROP TABLE messages;
DROP TABLE chats;
