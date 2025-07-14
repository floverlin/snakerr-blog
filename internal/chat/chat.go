package chat

import (
	"blog/internal/model"
	"blog/internal/storage"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Chat struct {
	storage     storage.ChatRepository
	onliners    map[uint64][]*onliner
	mutex       *sync.RWMutex
	upgrader    *websocket.Upgrader
	logger      *slog.Logger
	pingRate    time.Duration
	pongTimeout time.Duration
}

func New(storage storage.ChatRepository, logger *slog.Logger, pingRate time.Duration, pongTimeout time.Duration) *Chat {
	return &Chat{
		storage:  storage,
		onliners: map[uint64][]*onliner{},
		mutex:    &sync.RWMutex{},
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		logger:      logger,
		pingRate:    pingRate,
		pongTimeout: pongTimeout,
	}
}

func (c *Chat) HandleWS(w http.ResponseWriter, r *http.Request, user *model.User) error {
	const op = "handle chat websocket"
	ws, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer ws.Close()

	c.newOnliner(ws, user).handleWS()

	return nil
}

func (c *Chat) newOnliner(ws *websocket.Conn, user *model.User) *onliner {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	o := &onliner{
		ID:      c.newOnlinerID(user),
		WS:      ws,
		User:    user,
		Chat:    c,
		Timeout: make(chan struct{}),
		Done:    make(chan struct{}),
		Mutex:   &sync.Mutex{},
	}
	c.addOnliner(o)
	return o
}

func (c *Chat) newOnlinerID(user *model.User) int {
	userOnliners, ok := c.onliners[user.ID]
	if !ok {
		return 1
	}
	lastID := userOnliners[len(userOnliners)-1].ID
	return lastID + 1
}

func (c *Chat) sendMessage(id uint64, message *model.Message) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	user_onliners, ok := c.onliners[id]

	if ok {
		message.Type = model.NewMessage
		for _, onliner := range user_onliners {
			onliner.Mutex.Lock()
			_ = onliner.WS.NetConn().SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := onliner.WS.WriteJSON(message); err != nil {
				c.logger.Error("chat", "write json", err)
			}
			onliner.Mutex.Unlock()
		}
	}
}

func (c *Chat) addOnliner(o *onliner) {
	userOnliners, ok := c.onliners[o.User.ID]
	if ok {
		userOnliners = append(userOnliners, o)
		c.onliners[o.User.ID] = userOnliners
	} else {
		c.onliners[o.User.ID] = []*onliner{o}
	}
}

func (c *Chat) removeOnliner(o *onliner) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	userOnliners, ok := c.onliners[o.User.ID]
	if !ok {
		c.logger.Error("no onliners to remove")
		return
	}
	userOnliners = slices.DeleteFunc(userOnliners, func(e *onliner) bool {
		return e.ID == o.ID
	})
	if len(userOnliners) == 0 {
		delete(c.onliners, o.User.ID)
	} else {
		c.onliners[o.User.ID] = userOnliners
	}
}

func (c *Chat) LogOnliners() {
	for {
		c.mutex.RLock()
		c.logger.Debug("log onliners",
			"all", len(c.onliners),
		)
		c.mutex.RUnlock()
		time.Sleep(1 * time.Hour)
	}
}
