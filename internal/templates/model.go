package templates

import (
	"blog/internal/locales"
	"blog/internal/model"
	"context"
	"log"
)

type BaseTempl struct {
	Header  HeaderDef
	Locale  locales.Locale
	Version string
}

type HeaderDef struct {
	Locale       locales.Locale
	Title        string
	Flash        string
	IsLogged     bool
	HaveUnreaded bool
	Version      string
}

type ErrorTempl struct {
	Code   int
	Status string
	Error  string
}

type PaginateDef struct {
	Addr string
	Prev bool
	Next bool
	Page int
}

type ctxKey int8

const (
	CtxLocaleKey   ctxKey = iota // map[string]string
	CtxIsLoggedKey               // bool
	CtxFlashKey                  // string
	CtxUnreadedKey               // int
	CtxVersionKey                // string
)

const (
	DEFAULT_AVATAR      = "default"
	DEFAULT_AVATAR_PATH = "/static/img/default_avatar.jpg"
)

func baseTempData(ctx context.Context, title string) BaseTempl {
	var (
		locale       locales.Locale
		isLogged     bool
		flash        string
		haveUnreaded bool
		version      string
	)
	if v, ok := ctx.Value(CtxLocaleKey).(locales.Locale); !ok {
		locale = map[string]string{}
		log.Println("no locale in context") // todo
	} else {
		locale = v
	}
	if v, ok := ctx.Value(CtxIsLoggedKey).(bool); !ok {
		log.Println("no isLogged in context") // todo
	} else {
		isLogged = v
	}
	if v, ok := ctx.Value(CtxFlashKey).(string); ok {
		flash = v
	}
	if v, ok := ctx.Value(CtxUnreadedKey).(int); ok {
		if v != 0 {
			haveUnreaded = true
		}
	}
	if v, ok := ctx.Value(CtxVersionKey).(string); ok {
		version = v
	}
	return BaseTempl{
		Header: HeaderDef{
			Title:        title,
			Locale:       locale,
			Flash:        flash,
			IsLogged:     isLogged,
			HaveUnreaded: haveUnreaded,
			Version:      version,
		},
		Locale:  locale,
		Version: version,
	}
}

func NewsPageData(
	ctx context.Context,
	title string,
	addr string,
	page int,
	totalPages int,
	posts []*model.Post,
) any {
	base := baseTempData(ctx, title)
	paginate := PaginateDef{Page: page, Addr: addr}
	if totalPages > page {
		paginate.Next = true
	}
	if page > 1 {
		paginate.Prev = true
	}
	return struct {
		BaseTempl
		PaginateDef
		Posts []*model.Post
	}{
		BaseTempl:   base,
		PaginateDef: paginate,
		Posts:       posts,
	}
}

func UserPageData(
	ctx context.Context,
	title string,
	addr string,
	page, totalPages int,
	posts []*model.Post,
	user *model.User,
	subscribers, subscribes int,
	isFollowed bool,
	csrfToken string,
) any {
	base := baseTempData(ctx, title)
	user.Avatar = avatarPath(user.Avatar)
	paginate := PaginateDef{Page: page, Addr: addr}
	if totalPages > page {
		paginate.Next = true
	}
	if page > 1 {
		paginate.Prev = true
	}
	return struct {
		BaseTempl
		PaginateDef
		*model.User
		Posts            []*model.Post
		SubscribersCount int
		SubscribesCount  int
		IsFollowed       bool
		CSRFToken        string
	}{
		BaseTempl:        base,
		PaginateDef:      paginate,
		User:             user,
		Posts:            posts,
		SubscribersCount: subscribers,
		SubscribesCount:  subscribes,
		IsFollowed:       isFollowed,
		CSRFToken:        csrfToken,
	}
}

func SmokingPageData(ctx context.Context, title string) any {
	base := baseTempData(ctx, title)
	return struct {
		BaseTempl
	}{
		BaseTempl: base,
	}
}

func LoginPageData(ctx context.Context, title string, csrfToken string) any {
	base := baseTempData(ctx, title)
	return struct {
		BaseTempl
		CSRFToken string
	}{
		BaseTempl: base,
		CSRFToken: csrfToken,
	}
}
func LogoutPageData(ctx context.Context, title string) any {
	base := baseTempData(ctx, title)
	return struct {
		BaseTempl
	}{
		BaseTempl: base,
	}
}
func RegisterPageData(ctx context.Context, title string, csrfToken string) any {
	base := baseTempData(ctx, title)
	return struct {
		BaseTempl
		CSRFToken string
	}{
		BaseTempl: base,
		CSRFToken: csrfToken,
	}
}
func HelloPageData(ctx context.Context, title string, time int64) any {
	base := baseTempData(ctx, title)
	return struct {
		BaseTempl
		Time int64
	}{
		BaseTempl: base,
		Time:      time,
	}
}
func EditUserPageData(ctx context.Context, title string, user *model.User, csrfToken string) any {
	base := baseTempData(ctx, title)
	return struct {
		BaseTempl
		*model.User
		CSRFToken string
	}{
		BaseTempl: base,
		User:      user,
		CSRFToken: csrfToken,
	}
}

func ErrorPageData(ctx context.Context, code int, status string) any {
	base := baseTempData(ctx, "Error")
	return struct {
		BaseTempl
		Code   int
		Status string
	}{
		BaseTempl: base,
		Code:      code,
		Status:    status,
	}
}

func SnakeGamePageData(
	ctx context.Context,
	title string,
	personal int,
	global int,
	csrfToken string,
) any {
	base := baseTempData(ctx, title)
	return struct {
		BaseTempl
		CSRFToken    string
		PersonalBest int
		GlobalBest   int
	}{
		BaseTempl:    base,
		CSRFToken:    csrfToken,
		PersonalBest: personal,
		GlobalBest:   global,
	}
}

func LeaderboardPageData(
	ctx context.Context,
	title string,
	users []*model.MetaUser,
) any {
	base := baseTempData(ctx, title)
	return struct {
		BaseTempl
		Users []*model.MetaUser
	}{
		BaseTempl: base,
		Users:     users,
	}
}

func NotFoundPageData(ctx context.Context) any {
	base := baseTempData(ctx, "Oups...")
	return struct {
		BaseTempl
	}{
		BaseTempl: base,
	}
}

func AllChatsPageData(
	ctx context.Context,
	title string,
	chats []*model.Chat,
) any {
	for _, chat := range chats {
		chat.Dialogist.Avatar = avatarPath(chat.Dialogist.Avatar)
	}
	base := baseTempData(ctx, title)
	return struct {
		BaseTempl
		Chats []*model.Chat
	}{
		BaseTempl: base,
		Chats:     chats,
	}
}

func ChatPageData(
	ctx context.Context,
	title string,
	messages []*model.Message,
	dialogist *model.User,
) any {
	base := baseTempData(ctx, title)
	return struct {
		BaseTempl
		Messages  []*model.Message
		Dialogist *model.User
	}{
		BaseTempl: base,
		Messages:  messages,
		Dialogist: dialogist,
	}
}

func UsersListPageData(
	ctx context.Context,
	title string,
	users []*model.User,
) any {
	for _, user := range users {
		user.Avatar = avatarPath(user.Avatar)
	}
	base := baseTempData(ctx, title)
	return struct {
		BaseTempl
		Users []*model.User
	}{
		BaseTempl: base,
		Users:     users,
	}
}

func avatarPath(avatar string) string {
	if avatar == DEFAULT_AVATAR {
		return DEFAULT_AVATAR_PATH
	} else {
		return "/uploads/" + avatar + ".jpg"
	}
}
