package apiserver

import (
	"blog/internal/templates"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
)

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time := time.Now().Unix()
		exeData := templates.HelloPageData(r.Context(), "Hello", time)
		s.renderPage(w, r, http.StatusOK, "hello", exeData)
	}
}

func (s *APIServer) handleError() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.errorPage(w, r, errors.New("debug: error"))
	}
}

func (s *APIServer) handleUploadAvatar() http.HandlerFunc {
	const op = "handle upload avatar"
	return func(w http.ResponseWriter, r *http.Request) {
		id := getCurrentID(r)
		file, _, err := r.FormFile("avatar")
		if err != nil {
			s.errorPageReferer(w, r, errors.New("file upload error"))
			return
		}
		defer file.Close()
		name := uuid.New().String()
		out, err := os.Create(path.Join(s.config.UploadsPath, name+".jpg"))
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		defer out.Close()
		u, err := s.storage.User().FindByID(r.Context(), id)
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		u.Avatar = name
		if err := s.storage.User().Update(r.Context(), u); err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		if _, err := io.Copy(out, file); err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		http.Redirect(w, r, "/my_page", http.StatusSeeOther)
	}
}

func (s *APIServer) handleAPITime() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timeData := time.Now().Unix()
		data := map[string]any{
			"time": timeData,
		}
		s.respondJSON(w, r, http.StatusOK, data)
	}
}

func (s *APIServer) handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exeData := templates.NotFoundPageData(r.Context())
		s.renderPage(w, r, http.StatusNotFound, "not_found", exeData)
	}
}

func (s *APIServer) handleSetLang() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.URL.Query().Get("lang")
		if lang == "" {
			lang = "en"
		}
		cookie := &http.Cookie{
			Name:    "lang",
			Value:   lang,
			Expires: time.Now().Add(24 * 365 * time.Hour),
			MaxAge:  int(24 * 365 * time.Hour / time.Second),
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
	}
}

func getLang(r *http.Request) string {
	lang, err := r.Cookie("lang")
	if err != nil {
		return "en"
	}
	return lang.Value
}
