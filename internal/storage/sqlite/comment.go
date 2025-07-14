package sqlite

import (
	"blog/internal/model"
	"context"
	"database/sql"
	"fmt"
	"slices"
)

type CommentRepository struct {
	db *sql.DB
}

func (r *CommentRepository) CreateComment(ctx context.Context, c *model.Comment) (uint64, error) {
	const op = "create"
	var wrap = wrapNotNull
	query := `INSERT INTO comments(user_id, post_id, body) 
		VALUES (?, ?, ?) RETURNING id`
	row := r.db.QueryRowContext(ctx, query, c.User.ID, c.PostID, c.Body)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, wrap(err))
	}
	return id, nil
}
func (r *CommentRepository) GetComments(ctx context.Context, postID uint64) ([]*model.Comment, error) {
	const op = "get"
	var wrap = wrapNoRows
	res := []*model.Comment{}
	query := `SELECT comments.id, comments.post_id, comments.body, comments.created_at,
		comments.user_id, users.username, users.avatar
		FROM comments JOIN users ON users.id = comments.user_id
		WHERE (post_id == ?) ORDER BY comments.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, wrap(err))
	}
	for rows.Next() {
		c := model.Comment{User: model.User{}}
		err := rows.Scan(
			&c.ID, &c.PostID, &c.Body, &c.CreatedAt,
			&c.User.ID, &c.User.Username, &c.User.Avatar,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		res = append(res, &c)
	}
	slices.Reverse(res)
	return res, nil
}
