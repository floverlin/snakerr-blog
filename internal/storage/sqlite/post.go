package sqlite

import (
	"blog/internal/model"
	"context"
	"database/sql"
	"fmt"
)

type PostRepository struct {
	db *sql.DB
}

func (r *PostRepository) Create(ctx context.Context, p *model.Post) (uint64, error) {
	const op = "create"
	var wrap = wrapNotNull
	query := `INSERT INTO posts(user_id, title, body) VALUES (?, ?, ?) RETURNING id`
	var id uint64
	err := r.db.QueryRowContext(ctx, query, p.UserID, p.Title, p.Body).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, wrap(err))
	}
	return id, nil
}

func (r *PostRepository) Delete(ctx context.Context, id uint64) (uint64, error) {
	return 0, nil
}

func (r *PostRepository) GetAllPaginate(ctx context.Context, id uint64, page int, limit int) ([]*model.Post, int, error) {
	const op = "get all paginate"
	var wrap = noWrap
	posts := make([]*model.Post, 0, limit)
	query := `SELECT posts.id, posts.user_id, posts.title, posts.body, posts.created_at, users.username,
		(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id) as like_count,
		(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id AND likes.user_id = ?) as liked,
		(SELECT COUNT(*) FROM dislikes WHERE dislikes.post_id = posts.id) as dislike_count,
		(SELECT COUNT(*) FROM dislikes WHERE dislikes.post_id = posts.id AND dislikes.user_id = ?) as disliked
		FROM posts JOIN users ON posts.user_id = users.id
		GROUP BY posts.id
		ORDER BY posts.created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, id, id, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, wrap(err))
	}
	for rows.Next() {
		p := &model.Post{}
		rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Body, &p.CreatedAt, &p.Author, &p.LikeCount, &p.Liked, &p.DislikeCount, &p.Disliked)
		posts = append(posts, p)
	}
	var count int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM posts").Scan(&count); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, noWrap(err))
	}
	r.countComments(ctx, posts)
	return posts, count, nil
}

func (r *PostRepository) GetByIDPaginate(ctx context.Context, id uint64, pageID uint64, page int, limit int) ([]*model.Post, int, error) {
	const op = "get by id paginate"
	var wrap = noWrap
	posts := make([]*model.Post, 0, limit)
	query := `SELECT posts.id, posts.user_id, posts.title, posts.body, posts.created_at, users.username,
		(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id) as like_count,
		(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id AND likes.user_id = ?) as liked,
		(SELECT COUNT(*) FROM dislikes WHERE dislikes.post_id = posts.id) as dislike_count,
		(SELECT COUNT(*) FROM dislikes WHERE dislikes.post_id = posts.id AND dislikes.user_id = ?) as disliked
		FROM posts JOIN users ON posts.user_id = users.id
		WHERE posts.user_id = ?
		GROUP BY posts.id
		ORDER BY posts.created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, id, id, pageID, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, wrap(err))
	}
	for rows.Next() {
		p := &model.Post{}
		rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Body, &p.CreatedAt, &p.Author, &p.LikeCount, &p.Liked, &p.DislikeCount, &p.Disliked)

		posts = append(posts, p)
	}
	var count int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM posts WHERE user_id = ?", pageID).Scan(&count); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, noWrap(err))
	}
	r.countComments(ctx, posts)
	return posts, count, nil
}

func (r *PostRepository) GetMyPaginate(ctx context.Context, id uint64, page int, limit int) ([]*model.Post, int, error) {
	const op = "get my paginate"
	var wrap = wrapNoRows
	posts := make([]*model.Post, 0, limit)
	query := `SELECT posts.id, posts.user_id, posts.title, posts.body, posts.created_at, users.username,
		(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id) as like_count,
		(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id AND likes.user_id = ?) as liked,
		(SELECT COUNT(*) FROM dislikes WHERE dislikes.post_id = posts.id) as dislike_count,
		(SELECT COUNT(*) FROM dislikes WHERE dislikes.post_id = posts.id AND dislikes.user_id = ?) as disliked
		FROM posts JOIN followers ON followers.followed = posts.user_id
		JOIN users ON posts.user_id = users.id
		WHERE followers.follower = ?
		GROUP BY posts.id
		ORDER BY posts.created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, id, id, id, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, wrap(err))
	}
	for rows.Next() {
		p := &model.Post{}
		rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Body, &p.CreatedAt, &p.Author, &p.LikeCount, &p.Liked, &p.DislikeCount, &p.Disliked)
		posts = append(posts, p)
	}
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM posts JOIN followers ON followers.followed = posts.user_id
		JOIN users ON posts.user_id = users.id
		WHERE followers.follower = ?`, id).Scan(&count); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, noWrap(err))
	}
	r.countComments(ctx, posts)
	return posts, count, nil
}

func (r *PostRepository) countComments(ctx context.Context, posts []*model.Post) error {
	const op = "count comments"
	var wrap = noWrap
	for _, post := range posts {
		query := `SELECT COUNT(*) FROM comments WHERE post_id = ?`
		err := r.db.QueryRowContext(ctx, query, post.ID).Scan(&post.CommentCount)
		if err != nil {
			return fmt.Errorf("%s: %w", op, wrap(err))
		}
	}
	return nil
}
