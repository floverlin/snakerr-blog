package apiserver

import (
	"blog/internal/pkg"
	"blog/internal/templates"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (s *APIServer) handleSnake() http.HandlerFunc {
	const op = "handle snake"
	return func(w http.ResponseWriter, r *http.Request) {
		token := pkg.SetAndGetCSRFCookie(w, r)
		globalBest, err := s.storage.Snake().GetGlobalBest(r.Context())
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		personalBest, err := s.storage.Snake().GetPersonalBest(r.Context(), getCurrentID(r))
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		exeData := templates.SnakeGamePageData(
			r.Context(),
			"Snake Game",
			personalBest,
			globalBest,
			token,
		)
		s.renderPage(w, r, http.StatusOK, "snake", exeData)
	}
}

func (s *APIServer) handleAPISnakeSaveRecord() http.HandlerFunc {
	const op = "handle api snake save record"
	type dto struct {
		Record int    `json:"record"`
		Token  string `json:"csrf_token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.errorJSON(w, r, http.StatusBadRequest, errors.New(""))
			return
		}
		_ = r.ParseForm()
		vals := r.Form
		vals.Set("csrf_token", req.Token)
		if _, ok := pkg.CheckCSRF(vals, r); !ok {
			s.errorJSON(w, r, http.StatusBadRequest, errors.New("wrong csrf token"))
			return
		}
		personalBest, err := s.storage.Snake().GetPersonalBest(r.Context(), getCurrentID(r))
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
			return
		}
		if personalBest >= req.Record {
			s.errorJSON(w, r, http.StatusBadRequest, errors.New("personal best is greater than given record"))
			return
		}
		if err := s.storage.Snake().Save(r.Context(), getCurrentID(r), req.Record); err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (s *APIServer) handleAPISnakeGetPersonalRecord() http.HandlerFunc {
	const op = "handle snake get personal record"
	return func(w http.ResponseWriter, r *http.Request) {
		personalBest, err := s.storage.Snake().GetPersonalBest(r.Context(), getCurrentID(r))
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
			return
		}
		s.respondJSON(w, r, http.StatusOK, map[string]any{"record": personalBest})
	}
}

func (s *APIServer) handleAPISnakeGetGlobalRecord() http.HandlerFunc {
	const op = "handle api snake get global record"
	return func(w http.ResponseWriter, r *http.Request) {
		globalBest, err := s.storage.Snake().GetGlobalBest(r.Context())
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
			return
		}
		s.respondJSON(w, r, http.StatusOK, map[string]any{"record": globalBest})
	}
}

func (s *APIServer) handleSnakeLeaderboard() http.HandlerFunc {
	const op = "handle snake leaderboard"
	return func(w http.ResponseWriter, r *http.Request) {

		users, err := s.storage.Snake().Leaders(r.Context(), 10)
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		exeData := templates.LeaderboardPageData(r.Context(), "Leaderboard", users)
		s.renderPage(w, r, http.StatusOK, "leaderboard", exeData)
	}
}
