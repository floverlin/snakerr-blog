package apiserver

import (
	"blog/internal/model"
	"blog/internal/pkg"
	"blog/internal/storage"
	"blog/internal/templates"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *APIServer) handleMyPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := getCurrentID(r)
		newReq := mux.SetURLVars(r, map[string]string{"id": strconv.FormatUint(id, 10)})
		s.handleUserPage().ServeHTTP(w, newReq)
	}
}

func (s *APIServer) handleUserPage() http.HandlerFunc {
	const op = "handle user page"
	return func(w http.ResponseWriter, r *http.Request) {
		paramID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			s.handleNotFound().ServeHTTP(w, r)
			return
		}
		pageString := r.URL.Query().Get("page")
		if pageString == "" {
			pageString = "1"
		}
		page, err := strconv.Atoi(pageString)
		if err != nil || page < 0 {
			page = 1
		}
		id := getCurrentID(r)
		u, err := s.storage.User().FindByID(r.Context(), paramID)
		if err != nil {
			if errors.Is(err, storage.ErrNoRows) {
				s.handleNotFound().ServeHTTP(w, r)
				return

			} else {
				s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
				return
			}
		}
		var isFollowed bool
		var title string
		if id == paramID {
			title = "My Page"
			if r.Method == http.MethodPost {
				s.handleMyPageCreatePost().ServeHTTP(w, r)
				return
			}
		} else {
			title = "User Page"
			isFollowed, err = s.storage.Follow().IsFollower(r.Context(), id, paramID)
			if err != nil {
				s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
				return
			}
			if r.Method == http.MethodPost {
				s.handleNotFound().ServeHTTP(w, r)
				return
			}
		}
		posts, n, err := s.storage.Post().GetByIDPaginate(r.Context(), id, paramID, page, s.config.Visual.PostPerPage)
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}

		subscribers, subscribes, err := s.storage.Follow().Count(r.Context(), paramID)
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}

		totalPages := calculateTotalPages(n, s.config.Visual.PostPerPage)

		csrfToken := pkg.SetAndGetCSRFCookie(w, r)
		exeData := templates.UserPageData(
			r.Context(),
			title,
			fmt.Sprintf("user/%d", u.ID),
			page,
			totalPages,
			posts,
			u,
			subscribers,
			subscribes,
			isFollowed,
			csrfToken,
		)
		s.renderPage(w, r, http.StatusOK, "user", exeData)
	}
}

func (s *APIServer) handleMyPageCreatePost() http.HandlerFunc {
	const op = "handle my page create post"
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		vals := r.Form
		if !(vals.Has("title") && vals.Has("body")) {
			s.errorPageReferer(w, r, errors.New("some values are empty"))
			return
		}
		if !pkg.CheckAndDeleteCSRF(w, r, vals) {
			s.errorPageReferer(w, r, errors.New("wrong csrf token"))
			return
		}
		p := model.Post{
			UserID: getCurrentID(r),
			Title:  vals.Get("title"),
			Body:   vals.Get("body"),
		}
		if err := p.ValidateBeforeCreate(); err != nil {
			s.errorPageReferer(w, r, err)
			return
		}
		_, err = s.storage.Post().Create(r.Context(), &p)
		if err != nil {
			if errors.Is(err, storage.ErrNotNull) {
				s.errorPageReferer(w, r, errors.New("some values are empty"))
				return
			} else {
				s.errorPage(w, r, err)
				return
			}
		}
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
	}
}

func (s *APIServer) handleEditUserPage() http.HandlerFunc {
	const op = "handle edit user page"
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
				return
			}
			vals := r.Form
			if !(vals.Has("username") && vals.Has("description")) {
				s.errorPageReferer(w, r, errors.New("some values are empty"))
				return
			}
			if !pkg.CheckAndDeleteCSRF(w, r, vals) {
				s.errorPageReferer(w, r, errors.New("wrong scrf token"))
				return
			}
			id := getCurrentID(r)
			u, err := s.storage.User().FindByID(r.Context(), id)
			if err != nil {
				s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
				return
			}
			u.Username = vals.Get("username")
			u.Description = vals.Get("description")
			if err := u.ValidateBeforeUpdate(); err != nil {
				s.errorPageReferer(w, r, err)
				return
			}
			err = s.storage.User().Update(r.Context(), u)
			if err != nil {
				if errors.Is(err, storage.ErrNotNull) {
					s.errorPageReferer(w, r, errors.New("some values are empty"))
				} else {
					s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
					return
				}
			}
			http.Redirect(w, r, "/my_page", http.StatusSeeOther)
			return
		}
		id := getCurrentID(r)
		u, err := s.storage.User().FindByID(r.Context(), id)
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		csrfToken := pkg.SetAndGetCSRFCookie(w, r)
		exeData := templates.EditUserPageData(r.Context(), "Edit Profile", u, csrfToken)
		s.renderPage(w, r, http.StatusOK, "edit_user", exeData)
	}
}
