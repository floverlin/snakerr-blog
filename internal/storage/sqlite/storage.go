package sqlite

import (
	"blog/internal/storage"
	"database/sql"
)

type Storage struct {
	db      *sql.DB
	user    storage.UserRepository
	post    storage.PostRepository
	follow  storage.FollowRepository
	like    storage.LikeRepository
	smoke   storage.SmokeRepository
	snake   storage.SnakeRepository
	chat    storage.ChatRepository
	comment storage.CommentRepository
}

func New(db *sql.DB) *Storage {

	return &Storage{
		db:      db,
		user:    &UserRepository{db: db},
		post:    &PostRepository{db: db},
		follow:  &FollowRepository{db: db},
		like:    &LikeRepository{db: db},
		smoke:   &SmokeRepository{db: db},
		snake:   &SnakeRepository{db: db},
		chat:    &ChatRepository{db: db},
		comment: &CommentRepository{db: db},
	}
}

func (s *Storage) User() storage.UserRepository {
	return s.user
}

func (s *Storage) Post() storage.PostRepository {
	return s.post
}

func (s *Storage) Follow() storage.FollowRepository {
	return s.follow
}

func (s *Storage) Like() storage.LikeRepository {
	return s.like
}

func (s *Storage) Smoke() storage.SmokeRepository {
	return s.smoke
}

func (s *Storage) Snake() storage.SnakeRepository {
	return s.snake
}

func (s *Storage) Chat() storage.ChatRepository {
	return s.chat
}

func (s *Storage) Comment() storage.CommentRepository {
	return s.comment
}
