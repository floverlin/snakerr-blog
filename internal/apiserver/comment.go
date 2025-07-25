package apiserver

import (
	"blog/internal/model"
	"blog/internal/storage"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *APIServer) handleAPICreateComments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := model.Comment{}
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			s.errorJSON(w, r, http.StatusBadRequest, err)
			return
		}
		id, err := s.storage.Comment().CreateComment(r.Context(), &c)
		if err != nil {
			if err == storage.ErrNotNull {
				s.errorJSON(w, r, http.StatusBadRequest, err)
			} else {
				s.errorJSON(w, r, http.StatusInternalServerError, err)
			}
			return
		}
		s.respondJSON(w, r, http.StatusOK, map[string]any{
			"id": id,
		})
	}
}

func (s *APIServer) handleAPIGetComments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDString, ok := mux.Vars(r)["postID"]
		if !ok {
			s.errorJSON(w, r, http.StatusBadRequest, nil)
			return
		}
		postID, err := strconv.ParseUint(postIDString, 10, 64)
		if err != nil {
			s.errorJSON(w, r, http.StatusBadRequest, nil)
			return
		}
		comments, err := s.storage.Comment().GetComments(r.Context(), postID)
		if err != nil && err != storage.ErrNoRows {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respondJSON(w, r, http.StatusOK, comments)
	}
}
