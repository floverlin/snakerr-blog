package sqlite

import (
	"blog/internal/model"
	"context"
	"database/sql"
	"fmt"
	"slices"
)

type SmokeRepository struct {
	db *sql.DB
}

func (r *SmokeRepository) Save(ctx context.Context, id uint64, msg string) (uint64, error) {
	const op = "save"
	var wrap = noWrap
	query := `INSERT INTO smoking (user_id, body) VALUES (?, ?)`
	res, err := r.db.ExecContext(ctx, query, id, msg)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, wrap(err))
	}
	resID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, noWrap(err))
	}
	return uint64(resID), nil
}
func (r *SmokeRepository) Get(ctx context.Context, limit int) ([]*model.SmokeMessage, error) {
	const op = "get"
	var wrap = noWrap
	res := []*model.SmokeMessage{}
	query := `SELECT smoking.id, user_id, users.username, body, smoking.created_at
		FROM smoking JOIN users ON smoking.user_id = users.id
		ORDER BY smoking.created_at DESC LIMIT ?`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, wrap(err))
	}
	for rows.Next() {
		m := model.SmokeMessage{}
		rows.Scan(&m.ID, &m.UserID, &m.Username, &m.Body, &m.CreatedAt)
		res = append(res, &m)
	}
	slices.Reverse(res)
	return res, nil
}
