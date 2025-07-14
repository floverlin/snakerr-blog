package apiserver

import (
	"encoding/json"
	"net/http"
)

func (s *APIServer) handleAPILike() http.HandlerFunc {
	type dto struct {
		PostID uint64 `json:"post_id"`
		Type   string `json:"type"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := getCurrentID(r)
		var req dto
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.errorJSON(w, r, http.StatusBadRequest, err)
			return
		}
		liked, err := s.storage.Like().IsLiked(r.Context(), id, req.PostID)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}
		disliked, err := s.storage.Like().IsDisliked(r.Context(), id, req.PostID)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}
		if req.Type == "like" {
			if liked {
				if err := s.storage.Like().Unlike(r.Context(), id, req.PostID); err != nil {
					s.errorJSON(w, r, http.StatusInternalServerError, err)
					return
				}
				liked = false
			} else {
				if disliked {
					if err := s.storage.Like().Undislike(r.Context(), id, req.PostID); err != nil {
						s.errorJSON(w, r, http.StatusInternalServerError, err)
						return
					}
					disliked = false
				}
				if err := s.storage.Like().Like(r.Context(), id, req.PostID); err != nil {
					s.errorJSON(w, r, http.StatusInternalServerError, err)
					return
				}
				liked = true
			}
		} else if req.Type == "dislike" {
			if disliked {
				if err := s.storage.Like().Undislike(r.Context(), id, req.PostID); err != nil {
					s.errorJSON(w, r, http.StatusInternalServerError, err)
					return
				}
				disliked = false
			} else {
				if liked {
					if err := s.storage.Like().Unlike(r.Context(), id, req.PostID); err != nil {
						s.errorJSON(w, r, http.StatusInternalServerError, err)
						return
					}
					liked = false
				}
				if err := s.storage.Like().Dislike(r.Context(), id, req.PostID); err != nil {
					s.errorJSON(w, r, http.StatusInternalServerError, err)
					return
				}
				disliked = true
			}
		}
		likes, dislikes, err := s.storage.Like().Count(r.Context(), req.PostID)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}
		data := map[string]any{
			"liked":    liked,
			"disliked": disliked,
			"likes":    likes,
			"dislikes": dislikes,
		}
		s.respondJSON(w, r, http.StatusOK, data)
	}
}
