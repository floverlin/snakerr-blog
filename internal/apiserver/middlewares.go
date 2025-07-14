package apiserver

import (
	"blog/internal/pkg"
	"blog/internal/templates"
	"bufio"
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ResponceWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *ResponceWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}

func (w *ResponceWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (s *APIServer) logMWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/static") ||
			strings.HasPrefix(r.RequestURI, "/uploads") ||
			strings.HasPrefix(r.RequestURI, "/metrics") ||
			strings.HasPrefix(r.RequestURI, "/ws") {
			next.ServeHTTP(w, r)
			return
		}
		s.metrics.requestCounter.Inc()
		start := time.Now()
		logW := &ResponceWriter{ResponseWriter: w}
		next.ServeHTTP(logW, r)
		s.logger.Info("request log",
			"method", r.Method,
			"uri", r.RequestURI,
			"code", logW.StatusCode,
			"status", http.StatusText(logW.StatusCode),
			"addr", r.RemoteAddr,
			"duration", time.Since(start))
		s.metrics.requestDuration.Observe(time.Since(start).Seconds())
	})
}

func (s *APIServer) contextMWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookie, err := r.Cookie("flash")
		if !errors.Is(err, http.ErrNoCookie) {
			ctx = context.WithValue(ctx, templates.CtxFlashKey, cookie.Value)
			deleteFlash(w)
		}
		ctx = context.WithValue(ctx, templates.CtxLocaleKey, s.locales[getLang(r)])
		ctx = context.WithValue(ctx, templates.CtxIsLoggedKey, isLogged(r))

		chats, err := s.storage.Chat().GetAllChats(r.Context(), getCurrentID(r))
		if err != nil {
			s.logger.Error("unreaded", "error", err)
		}
		cnt := 0
		for _, chat := range chats {
			if !chat.Readed {
				cnt += 1
			}
		}

		ctx = context.WithValue(ctx, templates.CtxUnreadedKey, cnt)

		ctx = context.WithValue(ctx, templates.CtxVersionKey, s.version)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (s *APIServer) refreshMWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pkg.UpdateJSImports(s.config.StaticPath, s.version)

		newTemplates, err := templates.Functions("refreshed", "templates")
		if err != nil {
			s.logger.Error("refresh", "error", err)
		}
		s.templates = newTemplates

		s.handleStatic(s.config.StaticPath)

		s.version = uuid.New().String()

		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) adminMWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if getCurrentID(r) != 1 {
			s.handleNotFound().ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
