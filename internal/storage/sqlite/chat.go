package sqlite

import (
	"blog/internal/model"
	"context"
	"database/sql"
	"fmt"
	"slices"
)

type ChatRepository struct {
	db *sql.DB
}

func (r *ChatRepository) GetAllChats(ctx context.Context, id uint64) ([]*model.Chat, error) {
	const op = "get all chats"
	var wrap = wrapNoRows
	chats := []*model.Chat{}
	query := `SELECT readed, updated_at,
		m.username, m.avatar, m.id, d.username, d.avatar, d.id
		FROM chats JOIN users m ON me = m.id
		JOIN users d ON dialogist = d.id
		WHERE me = ?
		ORDER BY updated_at DESC`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, wrap(err))
	}
	for rows.Next() {
		c := &model.Chat{}
		m := model.User{}
		d := model.User{}
		rows.Scan(
			&c.Readed, &c.UpdatedAt,
			&m.Username, &m.Avatar, &m.ID,
			&d.Username, &d.Avatar, &d.ID,
		)
		c.Me = m
		c.Dialogist = d
		chats = append(chats, c)
	}
	return chats, nil
}

func (r *ChatRepository) UpdateChat(ctx context.Context, message *model.Message, read bool) error {
	const op = "update chat"
	var wrap = noWrap
	if read {
		query := `UPDATE chats SET readed = ?
		WHERE me = ? AND dialogist = ?`
		_, err := r.db.ExecContext(ctx, query, true, message.To.ID, message.From.ID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, wrap(err))
		}
		return nil
	}
	query := `INSERT INTO chats (me, dialogist, readed, updated_at) VALUES (?, ?, ?, UNIXEPOCH())
		ON CONFLICT (me, dialogist) DO
		UPDATE SET me = ?, dialogist = ?, readed = ?, updated_at = UNIXEPOCH()
		WHERE me = ? AND dialogist = ?`
	_, err := r.db.ExecContext(
		ctx, query,
		message.To.ID, message.From.ID, false,
		message.To.ID, message.From.ID, false,
		message.To.ID, message.From.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, wrap(err))
	}
	query = `INSERT INTO chats (me, dialogist, readed, updated_at) VALUES (?, ?, ?, UNIXEPOCH())
		ON CONFLICT (me, dialogist) DO
		UPDATE SET me = ?, dialogist = ?, readed = ?, updated_at = UNIXEPOCH()
		WHERE me = ? AND dialogist = ?`
	_, err = r.db.ExecContext(
		ctx, query,
		message.From.ID, message.To.ID, true,
		message.From.ID, message.To.ID, true,
		message.From.ID, message.To.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, wrap(err))
	}

	return nil
}

func (r *ChatRepository) AddMessage(ctx context.Context, message *model.Message) error {
	const op = "add message"
	var wrap = noWrap
	query := `INSERT INTO messages (from_, to_, body) VALUES (?, ?, ?)`
	if _, err := r.db.ExecContext(ctx, query, message.From.ID, message.To.ID, message.Body); err != nil {
		return fmt.Errorf("%s: %w", op, wrap(err))
	}
	return nil
}

func (r *ChatRepository) GetMessages(ctx context.Context, id uint64, dialogistID uint64, offset int, limit int) ([]*model.Message, error) {
	const op = "get messages"
	var wrap = wrapNoRows
	messages := []*model.Message{}
	query := `SELECT from_, to_, body, messages.created_at, f.username, f.avatar, t.username, t.avatar
		FROM messages JOIN users f ON f.id = from_
		JOIN users t ON t.id = to_
		WHERE (from_ = ? AND to_ = ?) OR (to_ = ? AND from_ = ?)
		ORDER BY messages.created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, id, dialogistID, id, dialogistID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, wrap(err))
	}
	for rows.Next() {
		m := &model.Message{}
		from := model.User{}
		to := model.User{}
		rows.Scan(
			&from.ID, &to.ID, &m.Body, &m.CreatedAt,
			&from.Username, &from.Avatar,
			&to.Username, &to.Avatar,
		)
		m.From = from
		m.To = to
		messages = append(messages, m)
	}
	slices.Reverse(messages)
	return messages, nil
}
