-- +goose Up
CREATE TABLE messages_new(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_ INTEGER NOT NULL,
    to_ INTEGER NOT NULL,
    body TEXT NOT NULL,
    created_at INTEGER DEFAULT (UNIXEPOCH()),
    CONSTRAINT from_fk FOREIGN KEY (from_) REFERENCES users (id),
    CONSTRAINT to_fk FOREIGN KEY (to_) REFERENCES users (id)
);

INSERT INTO messages_new(from_, to_, body, created_at)
    SELECT from_, to_, body, created_at FROM messages;

DROP TABLE messages;

ALTER TABLE messages_new RENAME TO messages;

CREATE TABLE chats_new (
    me INTEGER NOT NULL,
    dialogist INTEGER NOT NULL,
    readed INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    PRIMARY KEY (me, dialogist),
    CONSTRAINT me_fk FOREIGN KEY (me) REFERENCES users (id),
    CONSTRAINT dialogist_fk FOREIGN KEY (dialogist) REFERENCES users (id)
);

INSERT INTO chats_new(me, dialogist, readed, updated_at)
    SELECT me, dialogist, readed, updated_at FROM chats;

DROP TABLE chats;

ALTER TABLE chats_new RENAME TO chats;

-- +goose Down
