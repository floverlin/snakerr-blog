package apiserver

import (
	"blog/internal/model"
	"blog/internal/templates"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/gorilla/websocket"
)

type smokingRoom struct {
	smokers []*smoker
}

func (r *smokingRoom) letIn(moker *smoker) error {
	if slices.ContainsFunc(r.smokers, func(elem *smoker) bool {
		return moker.id == elem.id
	}) {
		return errors.New("smoker with this id already in room")
	}
	joinMsg := map[string]any{
		"user_id":  -1,
		"username": "!admin",
		"body":     fmt.Sprintf("%s join room", moker.username),
	}
	for _, sm := range r.smokers {
		if err := sm.conn.WriteJSON(joinMsg); err != nil {
			log.Println("letin: ", err)
		}
	}
	r.smokers = append(r.smokers, moker)
	return nil
}

func (r *smokingRoom) kickOut(moker *smoker) {
	idx := slices.IndexFunc(r.smokers, func(s *smoker) bool { return s.id == moker.id })
	if idx == -1 {
		log.Println("smoker not found")
		return
	}
	kickMsg := map[string]any{
		"user_id":  -1,
		"username": "admin",
		"body":     fmt.Sprintf("%s leave out room", moker.username),
	}
	r.smokers = append(r.smokers[:idx], r.smokers[idx+1:]...)
	for _, sm := range r.smokers {
		if err := sm.conn.WriteJSON(kickMsg); err != nil {
			log.Println("kickout: ", err)
		}
	}
}

func (r *smokingRoom) Say(msg *model.SmokeMessage) {
	for _, sm := range r.smokers {
		if err := sm.conn.WriteJSON(msg); err != nil {
			log.Println("say: ", err)
		}
	}
}

type smoker struct {
	conn     *websocket.Conn
	id       uint64
	username string
}

var mainRoom = smokingRoom{
	smokers: []*smoker{},
}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type msgRes struct {
	msg []byte
	err error
}

func (s *APIServer) handleWSSmoking() http.HandlerFunc {
	const op = "handle websocket smoking"
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		defer ws.Close()
		id := getCurrentID(r)
		u, err := s.storage.User().FindByID(r.Context(), id)
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		sm := &smoker{
			conn:     ws,
			id:       id,
			username: u.Username,
		}
		if err := mainRoom.letIn(sm); err != nil {
			return
		}
		defer mainRoom.kickOut(sm)

		stop := make(chan struct{})
		defer close(stop)

		msgChan := make(chan msgRes)
		defer close(msgChan)

		t := time.NewTimer(time.Duration(s.config.WSPingRate+s.config.WSPongTimeout) * time.Second)

		go listenMessages(sm, stop, msgChan)
		initMessages, err := s.storage.Smoke().Get(r.Context(), 10)
		if err != nil {
			s.logger.Error("internal error", "error", fmt.Errorf("%s: %w", op, err))
			return
		}
		go sendMessages(ws, initMessages)
		go pingSmoker(sm, stop, s.config.WSPingRate, s.logger)
		for {
			var msg []byte
			select {
			case <-t.C:
				return
			case m := <-msgChan:
				if m.err != nil {
					return
				}
				msg = m.msg
				t.Reset(time.Duration(s.config.WSPingRate+s.config.WSPongTimeout) * time.Second)
			}

			if string(msg) == "pong" {
				continue
			}
			msgStruct := model.SmokeMessage{
				UserID:   sm.id,
				Username: sm.username,
				Body:     string(msg),
			}
			msgID, err := s.storage.Smoke().Save(r.Context(), sm.id, string(msg))
			if err != nil {
				s.logger.Error("internal error", "error", fmt.Errorf("%s: %w", op, err))
				continue
			}
			msgStruct.ID = msgID
			go mainRoom.Say(&msgStruct)
		}
	}
}

func listenMessages(smoker *smoker, stop chan struct{}, msgChan chan msgRes) {
	for {
		_, msg, err := smoker.conn.ReadMessage()
		select {
		case <-stop:
			return
		default:
			if err != nil {
				msgChan <- msgRes{err: err}
				return
			}
			if len(msg) == 0 {
				continue
			}
			msgChan <- msgRes{msg: msg}
		}
	}
}

func pingSmoker(smoker *smoker, stop chan struct{}, pingRate int, logger *slog.Logger) {
	const op = "ping smoker"
	t := time.NewTicker(time.Duration(pingRate) * time.Second)
	for {
		select {
		case <-t.C:
			if err := smoker.conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
				logger.Error("internal error", "error", fmt.Errorf("%s: %w", op, err))
				return
			}
		case <-stop:
			return
		}
	}
}

func sendMessages(conn *websocket.Conn, messages []*model.SmokeMessage) {
	for _, m := range messages {
		conn.WriteJSON(m)
	}
}

func (s *APIServer) handleSmokingRoom() http.HandlerFunc {
	const op = "handle smoking room"
	return func(w http.ResponseWriter, r *http.Request) {
		t := s.templates.Lookup("smoking.html")
		if t == nil {
			err := errors.New("template not found")
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		exeData := templates.SmokingPageData(r.Context(), "Smoking Room")
		if err := t.Execute(w, exeData); err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
	}
}
