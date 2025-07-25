package sqlite

import (
	"context"
	"database/sql"
	"fmt"
)

type LikeRepository struct {
	db *sql.DB
}

func (r *LikeRepository) IsLiked(ctx context.Context, userID, postID uint64) (bool, error) {
	const op = "is liked"
	var wrap = noWrap
	var res bool
	query := `SELECT EXISTS (SELECT 1 FROM likes WHERE user_id = ? AND post_id = ?)`
	if err := r.db.QueryRowContext(ctx, query, userID, postID).Scan(&res); err != nil {
		return false, fmt.Errorf("%s: %w", op, wrap(err))
	}
	return res, nil
}

func (r *LikeRepository) IsDisliked(ctx context.Context, userID, postID uint64) (bool, error) {
	const op = "is disliked"
	var wrap = noWrap
	var res bool
	query := `SELECT EXISTS (SELECT 1 FROM dislikes WHERE user_id = ? AND post_id = ?)`
	if err := r.db.QueryRowContext(ctx, query, userID, postID).Scan(&res); err != nil {
		return false, fmt.Errorf("%s: %w", op, wrap(err))
	}
	return res, nil
}

func (r *LikeRepository) Like(ctx context.Context, userID, postID uint64) error {
	const op = "like"
	var wrap = noWrap
	query := `INSERT INTO likes (user_id, post_id) VALUES (?, ?)`
	if _, err := r.db.ExecContext(ctx, query, userID, postID); err != nil {
		return fmt.Errorf("%s: %w", op, wrap(err))
	}
	return nil
}

func (r *LikeRepository) Dislike(ctx context.Context, userID, postID uint64) error {
	const op = "dislike"
	var wrap = noWrap
	query := `INSERT INTO dislikes (user_id, post_id) VALUES (?, ?)`
	if _, err := r.db.ExecContext(ctx, query, userID, postID); err != nil {
		return fmt.Errorf("%s: %w", op, wrap(err))
	}
	return nil
}

func (r *LikeRepository) Unlike(ctx context.Context, userID, postID uint64) error {
	const op = "unlike"
	var wrap = noWrap
	query := `DELETE FROM likes WHERE user_id = ? AND post_id = ?`
	if _, err := r.db.ExecContext(ctx, query, userID, postID); err != nil {
		return fmt.Errorf("%s: %w", op, wrap(err))
	}
	return nil
}

func (r *LikeRepository) Undislike(ctx context.Context, userID, postID uint64) error {
	const op = "undislike"
	var wrap = noWrap
	query := `DELETE FROM dislikes WHERE user_id = ? AND post_id = ?`
	if _, err := r.db.ExecContext(ctx, query, userID, postID); err != nil {
		return fmt.Errorf("%s: %w", op, wrap(err))
	}
	return nil
}

func (r *LikeRepository) Count(ctx context.Context, postID uint64) (int, int, error) {
	const op = "count"
	var wrap = noWrap
	var likes, dislikes int

	query := `SELECT COUNT(*) FROM likes WHERE post_id = ?`
	if err := r.db.QueryRowContext(ctx, query, postID).Scan(&likes); err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, wrap(err))
	}

	query = `SELECT COUNT(*) FROM dislikes WHERE post_id = ?`
	if err := r.db.QueryRowContext(ctx, query, postID).Scan(&dislikes); err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, wrap(err))
	}
	return likes, dislikes, nil
}
