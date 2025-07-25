package storage

import (
	"blog/internal/model"
	"context"
)

type Storage interface {
	User() UserRepository
	Post() PostRepository
	Follow() FollowRepository
	Like() LikeRepository
	Smoke() SmokeRepository
	Snake() SnakeRepository
	Chat() ChatRepository
	Comment() CommentRepository
}

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *model.Comment) (uint64, error)
	GetComments(ctx context.Context, postID uint64) ([]*model.Comment, error)
}

type ChatRepository interface {
	AddMessage(ctx context.Context, message *model.Message) error
	GetAllChats(ctx context.Context, id uint64) ([]*model.Chat, error)
	UpdateChat(ctx context.Context, message *model.Message, readed bool) error
	GetMessages(ctx context.Context, id uint64, dialogistID uint64, offset int, limit int) ([]*model.Message, error)
}

type UserRepository interface {
	Create(context.Context, *model.User) (uint64, error)
	Delete(ctx context.Context, id uint64) (uint64, error)
	Update(context.Context, *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id uint64) (*model.User, error)
}

type PostRepository interface {
	Create(context.Context, *model.Post) (uint64, error)
	Delete(ctx context.Context, id uint64) (uint64, error)
	GetAllPaginate(ctx context.Context, id uint64, page int, limit int) ([]*model.Post, int, error)
	GetMyPaginate(ctx context.Context, id uint64, page int, limit int) ([]*model.Post, int, error)
	GetByIDPaginate(ctx context.Context, id uint64, pageID uint64, page int, limit int) ([]*model.Post, int, error)
}

type FollowRepository interface {
	Follow(ctx context.Context, from *model.User, to *model.User) error
	Unfollow(ctx context.Context, from *model.User, to *model.User) error
	IsFollower(ctx context.Context, follower uint64, followed uint64) (bool, error)
	Count(ctx context.Context, id uint64) (int, int, error)
	GetFollowers(ctx context.Context, id uint64) ([]*model.User, error)
	GetFollows(ctx context.Context, id uint64) ([]*model.User, error)
}

type LikeRepository interface {
	IsLiked(ctx context.Context, userID, postID uint64) (bool, error)
	Like(ctx context.Context, userID, postID uint64) error
	Unlike(ctx context.Context, userID, postID uint64) error
	IsDisliked(ctx context.Context, userID, postID uint64) (bool, error)
	Dislike(ctx context.Context, userID, postID uint64) error
	Undislike(ctx context.Context, userID, postID uint64) error
	Count(ctx context.Context, postID uint64) (int, int, error)
}

type SmokeRepository interface {
	Save(ctx context.Context, id uint64, msg string) (uint64, error)
	Get(context.Context, int) ([]*model.SmokeMessage, error)
}

type SnakeRepository interface {
	Save(ctx context.Context, id uint64, record int) error
	GetPersonalBest(ctx context.Context, id uint64) (int, error)
	GetGlobalBest(context.Context) (int, error)
	Leaders(ctx context.Context, limit int) ([]*model.MetaUser, error)
}
