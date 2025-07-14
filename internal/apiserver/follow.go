package apiserver

import (
	"blog/internal/templates"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (s *APIServer) handleAPIFollow() http.HandlerFunc {
	type dto struct {
		Type string `json:"type"`
		To   uint64 `json:"to"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.errorJSON(w, r, http.StatusBadRequest, err)
			return
		}
		if req.Type != "follow" && req.Type != "unfollow" {
			s.errorJSON(w, r, http.StatusBadRequest, errors.New("wrong type"))
			return
		}
		toU, err := s.storage.User().FindByID(r.Context(), req.To)
		if err != nil {
			s.errorJSON(w, r, http.StatusBadRequest, err)
			return
		}
		fromID := getCurrentID(r)
		fromU, err := s.storage.User().FindByID(r.Context(), fromID)
		if err != nil {
			s.errorJSON(w, r, http.StatusBadRequest, err)
			return
		}
		switch req.Type {
		case "follow":
			err = s.storage.Follow().Follow(r.Context(), fromU, toU)
		case "unfollow":
			err = s.storage.Follow().Unfollow(r.Context(), fromU, toU)
		}
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (s *APIServer) handleMyFollowers() http.HandlerFunc {
	const op = "handle my followers"
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.storage.Follow().GetFollowers(r.Context(), getCurrentID(r))
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
		}
		exeData := templates.UsersListPageData(r.Context(), "My Followers", users)
		s.renderPage(w, r, http.StatusOK, "users_list", exeData)
	}
}

func (s *APIServer) handleMyFollows() http.HandlerFunc {
	const op = "handle my follows"
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.storage.Follow().GetFollows(r.Context(), getCurrentID(r))
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
		}
		exeData := templates.UsersListPageData(r.Context(), "My Follows", users)
		s.renderPage(w, r, http.StatusOK, "users_list", exeData)
	}
}
