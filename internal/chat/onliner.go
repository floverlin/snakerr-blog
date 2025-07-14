package chat

import (
	"blog/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type onliner struct {
	ID      int
	WS      *websocket.Conn
	User    *model.User
	Chat    *Chat
	Timeout chan struct{}
	Done    chan struct{}
	Mutex   *sync.Mutex
}

func (o *onliner) handleWS() {
	defer o.goOffline()
	o.sendInfo()

	go o.ping()
	go o.listen()

	<-o.Done
}

func (o *onliner) goOffline() {
	o.Chat.removeOnliner(o)
}

func (o *onliner) ping() {
	t := time.NewTimer(o.Chat.pingRate)
	defer t.Stop()
	defer close(o.Timeout)
	if err := o.WS.SetReadDeadline(time.Now().Add(o.Chat.pingRate + o.Chat.pongTimeout)); err != nil {
		o.Chat.logger.Error("ping", "set read deadline", err)
	}
	for {
		select {
		case <-t.C:
			o.Mutex.Lock()
			if err := o.WS.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
				o.Chat.logger.Error("ping", "write ping", err)
				return
			}
			o.Mutex.Unlock()
		case <-o.Timeout:
			t.Reset(o.Chat.pingRate)
			if err := o.WS.SetReadDeadline(time.Now().Add(o.Chat.pingRate + o.Chat.pongTimeout)); err != nil {
				o.Chat.logger.Error("ping", "reset read deadline", err)
				return
			}
		case <-o.Done:
			return
		}
	}
}

func (o *onliner) listen() {
	defer close(o.Done)
	for {
		_, data, err := o.WS.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); !ok {
				o.Chat.logger.Error("read message", "error", err)
			}
			return
		}
		o.Timeout <- struct{}{}
		if string(data) == "pong" {
			continue
		}
		message := model.Message{}
		if err := json.Unmarshal(data, &message); err != nil {
			o.Chat.logger.Error("listen", "unmarshal json", err)
		}
		go o.handleMessage(&message)
	}
}

func (o *onliner) sendInfo() {
	info := model.Message{
		Type: model.Info,
		From: model.User{},
		To:   *o.User,
		Body: "",
	}
	o.Mutex.Lock()
	defer o.Mutex.Unlock()
	_ = o.WS.NetConn().SetWriteDeadline(time.Now().Add(5 * time.Second))
	_ = o.WS.WriteJSON(info)
}

func (o *onliner) handleMessage(message *model.Message) {
	if o.User.ID != message.From.ID {
		o.Chat.logger.Warn("message from wrong id",
			"real id", o.User.ID,
			"in message", message.From.ID)
		return
	}
	// Save to DB
	if err := o.Chat.storage.AddMessage(context.TODO(), message); err != nil {
		fmt.Println(err)
	}
	// Update chat info
	if err := o.Chat.storage.UpdateChat(context.TODO(), message, false); err != nil {
		log.Println(err)
	}

	message.From = *o.User

	go o.Chat.sendMessage(message.To.ID, message)
	if o.User.ID != message.To.ID { // not the self chat
		go o.Chat.sendMessage(message.From.ID, message)
	}
}
