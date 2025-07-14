package apiserver

import (
	"blog/internal/chat"
	"blog/internal/config"
	"blog/internal/locales"
	"blog/internal/logger"
	"blog/internal/pkg"
	"blog/internal/storage"
	"blog/internal/templates"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	GET  = http.MethodGet
	POST = http.MethodPost
)

type APIServer struct {
	storage   storage.Storage
	router    *mux.Router
	server    *http.Server
	logger    *slog.Logger
	templates *template.Template
	locales   locales.Locales
	config    config.ServerCf
	metrics   serverMetrics
	chat      *chat.Chat
	version   string
}

func New(
	router *mux.Router,
	server *http.Server,
	storage storage.Storage,
	templates *template.Template,
	locales locales.Locales,
	config config.ServerCf,
) *APIServer {
	logger := slog.New(logger.NewColorHandler(slog.LevelDebug))
	return &APIServer{
		storage:   storage,
		router:    server.Handler.(*mux.Router),
		server:    server,
		logger:    logger,
		templates: templates,
		locales:   locales,
		config:    config,
		chat: chat.New(
			storage.Chat(),
			logger,
			time.Duration(config.WSPingRate)*time.Second,
			time.Duration(config.WSPongTimeout)*time.Second,
		),
		version: uuid.New().String(),
	}
}

func (s *APIServer) Run() error {
	s.configureRouter()
	s.registerMetrics()
	pkg.UpdateJSImports(s.config.StaticPath, s.version)
	s.logger.Debug("server started", "on local", "http://localhost"+s.server.Addr)
	go s.chat.LogOnliners()
	return s.server.ListenAndServe()
}

func (s *APIServer) configureRouter() {
	if s.config.Enviroment == "dev" {
		s.router.Use(s.refreshMWare)
	}
	s.router.Use(s.logMWare, s.authMWare, s.contextMWare)

	s.handleStatic(s.config.StaticPath)
	s.router.PathPrefix("/uploads/").
		Handler(
			http.StripPrefix(
				"/uploads/",
				http.FileServer(
					http.Dir(s.config.UploadsPath))))

	s.router.NotFoundHandler = s.logMWare(s.authMWare(s.contextMWare(s.handleNotFound())))

	s.router.Handle("/metrics", promhttp.Handler()).Methods(GET)

	s.router.Handle("/login", s.handleLogin()).Methods(GET, POST)
	s.router.Handle("/register", s.handleRegister()).Methods(GET, POST)

	s.router.Handle("/logout", s.handleLogout()).Methods(GET)
	s.router.Handle("/set_lang", s.handleSetLang()).Methods(GET)

	s.router.Handle("/", s.handleHello()).Methods(GET)

	auth := s.router.NewRoute().Subrouter()
	auth.Use(s.authSecureMWare)

	auth.Handle("/admin", s.adminMWare(s.handleAdmin())).Methods(GET)

	auth.Handle("/my_page", s.handleMyPage()).Methods(GET, POST)
	auth.Handle("/error", s.handleError()).Methods(GET)
	auth.Handle("/edit", s.handleEditUserPage()).Methods(GET, POST)
	auth.Handle("/upload_avatar", s.handleUploadAvatar()).Methods(POST)
	auth.Handle("/user/{id}", s.handleUserPage()).Methods(GET, POST)
	auth.Handle("/global_news", s.handleGlobalNews()).Methods(GET)
	auth.Handle("/my_news", s.handleMyNews()).Methods(GET)
	auth.Handle("/smoking", s.handleSmokingRoom()).Methods(GET)
	auth.Handle("/my_chats", s.handleAllChats()).Methods(GET)
	auth.Handle("/chat/{id}", s.handleChat()).Methods(GET)
	auth.Handle("/snake", s.handleSnake()).Methods(GET)
	auth.Handle("/snake/leaderboard", s.handleSnakeLeaderboard()).Methods(GET)
	auth.Handle("/my_followers", s.handleMyFollowers()).Methods(GET)
	auth.Handle("/my_follows", s.handleMyFollows()).Methods(GET)

	api := auth.PathPrefix("/api").Subrouter()

	api.Handle("/chat/messages", s.getMessages()).Methods(POST)
	api.Handle("/chat/mark_readed", s.markReaded()).Methods(POST)
	api.Handle("/follow", s.handleAPIFollow()).Methods(POST)
	api.Handle("/time", s.handleAPITime()).Methods(GET)
	api.Handle("/like", s.handleAPILike()).Methods(POST)
	api.Handle("/snake/record", s.handleAPISnakeSaveRecord()).Methods(POST)
	api.Handle("/snake/personal", s.handleAPISnakeGetPersonalRecord()).Methods(GET)
	api.Handle("/snake/global", s.handleAPISnakeGetGlobalRecord()).Methods(GET)

	api.Handle("/comments", s.handleAPICreateComments()).Methods(POST)
	api.Handle("/comments/{postID}", s.handleAPIGetComments()).Methods(GET)

	ws := auth.PathPrefix("/ws").Subrouter()

	ws.Handle("/smoking", s.handleWSSmoking()).Methods(GET)
	ws.Handle("/main", s.handleWSMain()).Methods(GET)

}

func (s *APIServer) respondJSON(w http.ResponseWriter, _ *http.Request, code int, data any) {
	resp, err := json.Marshal(data)
	if err != nil {
		s.logger.Error("respond marshal", "error", err)
		resp = make([]byte, 0)
	}
	w.WriteHeader(code)
	w.Write(resp)
}

func (s *APIServer) errorJSON(w http.ResponseWriter, r *http.Request, code int, err error) {
	data := map[string]string{
		"error":  err.Error(),
		"status": http.StatusText(code),
	}
	if code/100 == 5 {
		s.logger.Error("internal error", "error", err)
		delete(data, "error")
	}
	s.respondJSON(w, r, code, data)
}

func (s *APIServer) errorPageReferer(w http.ResponseWriter, r *http.Request, err error) {
	s.setFlash(w, err.Error())
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func (s *APIServer) setFlash(w http.ResponseWriter, text string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "flash",
		Value:    text,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Duration(s.config.FlashTimeout) * time.Second),
		MaxAge:   s.config.FlashTimeout,
	})
}

func deleteFlash(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "flash",
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
}

func (s *APIServer) errorPage(w http.ResponseWriter, r *http.Request, err error) {
	s.metrics.requestErrors.Inc()
	s.logger.Error("internal error", "error", err)
	exeData := templates.ErrorPageData(
		r.Context(),
		http.StatusInternalServerError,
		http.StatusText(http.StatusInternalServerError),
	)
	t := s.templates.Lookup("error.html")
	if t == nil {
		s.logger.Error("template not found")
		return
	}
	if err := t.Execute(w, exeData); err != nil {
		s.logger.Error("execute error template", "error", err)
		return
	}
}

func (s *APIServer) handleStatic(readPath string) {
	files, err := os.ReadDir(readPath)
	if err != nil {
		s.logger.Error("failed to read static files:", "error", err)
		return
	}
	for _, file := range files {
		fileName := file.Name()
		if file.IsDir() {
			s.handleStatic(path.Join(readPath, fileName))
		}
		s.router.HandleFunc("/"+path.Join(readPath, fileName), func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, path.Join(readPath, fileName))
		}).Methods(GET)
	}
}

func (s *APIServer) renderPage(w http.ResponseWriter, r *http.Request, code int, templateName string, data any) {
	const op = "render page"
	t := s.templates.Lookup(templateName + ".html")
	if t == nil {
		err := fmt.Errorf("template \"%s\" not found", templateName)
		s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}
	w.WriteHeader(code)
	if err := t.Execute(w, data); err != nil {
		s.errorPage(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}
}

func (s *APIServer) handleAdmin() http.HandlerFunc {
	// const op = "handle admin page"
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./database/data.db")
	}
}
