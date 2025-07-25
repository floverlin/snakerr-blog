package sqlite

import (
	"blog/internal/model"
	"context"
	"database/sql"
	"fmt"
)

type FollowRepository struct {
	db *sql.DB
}

func (r *FollowRepository) Follow(ctx context.Context, from *model.User, to *model.User) error {
	const op = "follow"
	var wrap = noWrap
	query := `INSERT INTO followers (follower, followed) VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, query, from.ID, to.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, wrap(err))
	}
	return nil
}

func (r *FollowRepository) Unfollow(ctx context.Context, from *model.User, to *model.User) error {
	const op = "unfollow"
	var wrap = noWrap
	query := `DELETE FROM followers WHERE follower = ? AND followed = ?`
	_, err := r.db.ExecContext(ctx, query, from.ID, to.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, wrap(err))
	}
	return nil
}

func (r *FollowRepository) IsFollower(ctx context.Context, follower uint64, followed uint64) (bool, error) {
	const op = "is follower"
	var wrap = noWrap
	var res bool
	query := `SELECT EXISTS (SELECT 1 FROM followers WHERE follower = ? AND followed = ?)`
	err := r.db.QueryRowContext(ctx, query, follower, followed).Scan(&res)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, wrap(err))
	}
	return res, nil
}

func (r *FollowRepository) Count(ctx context.Context, id uint64) (int, int, error) {
	const op = "count"
	var wrap = noWrap
	var subscribers, subscribes int
	query := `SELECT
		COALESCE(SUM(CASE WHEN followed = ? THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN follower = ? THEN 1 ELSE 0 END), 0)
		FROM followers`
	err := r.db.QueryRowContext(ctx, query, id, id).Scan(&subscribers, &subscribes)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, wrap(err))
	}
	return subscribers, subscribes, nil
}

func (r *FollowRepository) GetFollowers(ctx context.Context, id uint64) ([]*model.User, error) {
	const op = "get followers"
	var wrap = wrapNoRows
	users := []*model.User{}
	query := `SELECT users.id, users.username, users.avatar FROM users
		JOIN followers ON followers.follower = users.id
		WHERE followers.followed = ?`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, wrap(err))
	}
	for rows.Next() {
		user := model.User{}
		rows.Scan(&user.ID, &user.Username, &user.Avatar)
		users = append(users, &user)
	}
	return users, nil
}

func (r *FollowRepository) GetFollows(ctx context.Context, id uint64) ([]*model.User, error) {
	const op = "get follows"
	var wrap = wrapNoRows
	users := []*model.User{}
	query := `SELECT users.id, users.username, users.avatar FROM users
		JOIN followers ON followers.followed = users.id
		WHERE followers.follower = ?`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, wrap(err))
	}
	for rows.Next() {
		user := model.User{}
		rows.Scan(&user.ID, &user.Username, &user.Avatar)
		users = append(users, &user)
	}
	return users, nil
}
