package apiserver

import (
	"blog/internal/templates"
	"fmt"
	"net/http"
	"strconv"
)

func (s *APIServer) handleGlobalNews() http.HandlerFunc {
	const op = "handle global news"
	return func(w http.ResponseWriter, r *http.Request) {
		pageString := r.URL.Query().Get("page")
		if pageString == "" {
			pageString = "1"
		}
		page, err := strconv.Atoi(pageString)
		if err != nil || page < 0 {
			page = 1
		}
		id := getCurrentID(r)
		posts, n, err := s.storage.Post().GetAllPaginate(r.Context(), id, page, s.config.Visual.PostPerPage)
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		totalPages := calculateTotalPages(n, s.config.Visual.PostPerPage)
		exeData := templates.NewsPageData(r.Context(), "Global News", "global_news", page, totalPages, posts)
		s.renderPage(w, r, http.StatusOK, "news", exeData)
	}
}

func (s *APIServer) handleMyNews() http.HandlerFunc {
	const op = "handle my news"
	return func(w http.ResponseWriter, r *http.Request) {
		pageString := r.URL.Query().Get("page")
		if pageString == "" {
			pageString = "1"
		}
		page, err := strconv.Atoi(pageString)
		if err != nil || page < 0 {
			page = 1
		}
		id := getCurrentID(r)
		posts, n, err := s.storage.Post().GetMyPaginate(r.Context(), id, page, s.config.Visual.PostPerPage)
		if err != nil {
			s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
			return
		}
		totalPages := calculateTotalPages(n, s.config.Visual.PostPerPage)
		exeData := templates.NewsPageData(r.Context(), "My News", "my_news", page, totalPages, posts)
		s.renderPage(w, r, http.StatusOK, "news", exeData)
	}
}

func calculateTotalPages(n int, pagesPerPage int) int {
	totalPages := (n / pagesPerPage)
	if totalPages == 0 {
		totalPages = 1
	} else if n%pagesPerPage != 0 {
		totalPages += 1
	}
	return totalPages
}
