package sqlite

import (
	"blog/internal/model"
	"context"
	"database/sql"
	"fmt"
)

type SnakeRepository struct {
	db *sql.DB
}

func (r *SnakeRepository) Save(ctx context.Context, id uint64, record int) error {
	const op = "save"
	var wrap = noWrap
	query := `INSERT INTO snake_game (user_id, record, created_at) VALUES (?, ?, UNIXEPOCH())
		ON CONFLICT (user_id) DO
		UPDATE SET record = ?, created_at = UNIXEPOCH() WHERE user_id = ?`
	if _, err := r.db.ExecContext(ctx, query, id, record, record, id); err != nil {
		return fmt.Errorf("%s: %w", op, wrap(err))
	}
	return nil
}
func (r *SnakeRepository) GetPersonalBest(ctx context.Context, id uint64) (int, error) {
	const op = "get personal best"
	var wrap = noWrap
	var record int
	query := `SELECT COALESCE (
		(SELECT record FROM snake_game WHERE user_id = ?), 0
		)`
	if err := r.db.QueryRowContext(ctx, query, id, record).Scan(&record); err != nil {
		return 0, fmt.Errorf("%s: %w", op, wrap(err))
	}
	return record, nil
}
func (r *SnakeRepository) GetGlobalBest(ctx context.Context) (int, error) {
	const op = "get global best"
	var wrap = noWrap
	var record int
	query := `SELECT COALESCE (
		(SELECT record FROM snake_game ORDER BY record DESC), 0
		)`
	if err := r.db.QueryRowContext(ctx, query).Scan(&record); err != nil {
		return 0, fmt.Errorf("%s: %w", op, wrap(err))
	}
	return record, nil
}
func (r *SnakeRepository) Leaders(ctx context.Context, limit int) ([]*model.MetaUser, error) {
	const op = "leaderboard"
	var wrap = noWrap
	res := []*model.MetaUser{}
	query := `SELECT users.username, users.id, record, snake_game.created_at
		FROM snake_game JOIN users ON snake_game.user_id = users.id
		ORDER BY record DESC, snake_game.created_at ASC
		LIMIT ?`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, wrap(err))
	}
	for rows.Next() {
		mu := model.MetaUser{}
		rows.Scan(&mu.Username, &mu.ID, &mu.Record, &mu.RecordCreatedAt)
		res = append(res, &mu)
	}
	return res, nil
}
