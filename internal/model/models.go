package model

type User struct {
	ID          uint64 `json:"id"`
	Email       string `json:"-"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	Description string `json:"-"`
	Password    string `json:"-"`
}

type MetaUser struct {
	User
	Record          int
	RecordCreatedAt int
}

type Post struct {
	ID           uint64 `sql:"id"`
	UserID       uint64 `sql:"user_id"`
	Author       string
	Title        string `sql:"title"`
	Body         string `sql:"body"`
	CreatedAt    int    `sql:"created_at"`
	Liked        bool
	LikeCount    int
	Disliked     bool
	DislikeCount int
	CommentCount int
}

type Comment struct {
	ID        uint64 `json:"id"`
	User      User   `json:"user"`
	PostID    uint64 `json:"post_id"`
	Body      string `json:"body"`
	CreatedAt int    `json:"created_at"`
}

type SmokeMessage struct {
	ID        uint64 `json:"-"`
	UserID    uint64 `json:"user_id"`
	Username  string `json:"username"`
	Body      string `json:"body"`
	CreatedAt int    `json:"created_at"`
}

type messageType string

const (
	MessageToUser messageType = "message_to_user"
	NewMessage    messageType = "new_message"
	Info          messageType = "info"
)

type Message struct {
	Type      messageType `json:"type"`
	From      User        `json:"from"`
	To        User        `json:"to"`
	Body      string      `json:"body"`
	CreatedAt int         `json:"created_at"`
}

type Chat struct {
	Me        User `json:"me"`
	Dialogist User `json:"dialogist"`
	Readed    bool `json:"readed"`
	UpdatedAt int  `json:"updated_at"`
}
