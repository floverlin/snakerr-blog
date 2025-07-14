package apiserver

import (
	"blog/internal/model"
	"blog/internal/storage"
	"blog/internal/templates"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *APIServer) handleAllChats() http.HandlerFunc {
	const op = "handle all chats"
	return func(w http.ResponseWriter, r *http.Request) {
		chats, err := s.storage.Chat().GetAllChats(r.Context(), getCurrentID(r))
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		exeData := templates.AllChatsPageData(r.Context(), "My Chats", chats)
		s.renderPage(w, r, http.StatusOK, "my_chats", exeData)
	}
}

func (s *APIServer) handleChat() http.HandlerFunc {
	const op = "handle chat"
	return func(w http.ResponseWriter, r *http.Request) {
		dialogistID, err := getIDParam(r)
		if err != nil {
			s.handleNotFound().ServeHTTP(w, r)
			return
		}
		dialogist, err := s.storage.User().FindByID(r.Context(), dialogistID)
		if err != nil {
			if errors.Is(err, storage.ErrNoRows) {
				s.handleNotFound().ServeHTTP(w, r)
				return
			} else {
				s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
				return
			}
		}
		messages, err := s.storage.Chat().GetMessages(r.Context(), getCurrentID(r), dialogist.ID, 0, 10) // TODO
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		exeData := templates.ChatPageData(r.Context(), "Chat", messages, dialogist)
		s.renderPage(w, r, http.StatusOK, "chat", exeData)
	}
}

func getIDParam(r *http.Request) (uint64, error) {
	paramID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		return 0, errors.New("no id in params")
	}
	return paramID, nil
}

func (s *APIServer) getMessages() http.HandlerFunc {
	const op = "get messages"
	type messageDTO struct {
		Offset      int    `json:"offset"`
		Limit       int    `json:"limit"`
		DialogistID uint64 `json:"dialogist_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var dto messageDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			s.errorJSON(w, r, http.StatusBadRequest, err)
			return
		}
		messages, err := s.storage.Chat().GetMessages(
			r.Context(),
			getCurrentID(r),
			dto.DialogistID,
			dto.Offset,
			dto.Limit,
		)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
			return
		}
		if len(messages) != 0 {
			if err := s.storage.Chat().UpdateChat(r.Context(), messages[len(messages)-1], true); err != nil {
				s.errorJSON(w, r, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
				return
			}
		}
		for _, message := range messages {
			message.Type = model.NewMessage
		}
		s.respondJSON(w, r, http.StatusOK, messages)
	}
}

func (s *APIServer) markReaded() http.HandlerFunc {
	const op = "mark chat readed"
	type dto struct {
		Me        uint64 `json:"me"`
		Dialogist uint64 `json:"dialogist"`
	}
	type resp struct {
		AllReaded bool `json:"all_readed"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var dto dto
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			s.errorJSON(w, r, http.StatusBadRequest, err)
			return
		}
		message := model.Message{
			From: model.User{ID: dto.Dialogist},
			To:   model.User{ID: dto.Me},
		}
		if err := s.storage.Chat().UpdateChat(r.Context(), &message, true); err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
			return
		}
		chats, err := s.storage.Chat().GetAllChats(r.Context(), getCurrentID(r))
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
			return
		}
		resp := resp{
			AllReaded: !slices.ContainsFunc(chats, func(c *model.Chat) bool {
				return !c.Readed
			}),
		}
		s.respondJSON(w, r, http.StatusOK, resp)
	}
}

func (s *APIServer) handleWSMain() http.HandlerFunc {
	const op = "handle websocket main"
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := s.storage.User().FindByID(r.Context(), getCurrentID(r))
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
			return
		}
		if err := s.chat.HandleWS(w, r, user); err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
			return
		}
	}
}
