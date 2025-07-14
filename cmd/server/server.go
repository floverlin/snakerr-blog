package server

import (
	"blog/internal/apiserver"
	"blog/internal/config"
	"blog/internal/locales"
	"blog/internal/storage/sqlite"
	templ "blog/internal/templates"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func Execute() {
	var port int
	var env string
	var cfg string
	flag.IntVar(&port, "p", 8080, "port")
	flag.StringVar(&env, "e", "dev", "enviroment")
	flag.StringVar(&cfg, "c", "config.json", "config path")
	flag.Parse()

	var config config.Config
	f, err := os.ReadFile(cfg)
	if err != nil {
		log.Fatal("osen config: ", err)
	}
	if err := json.Unmarshal(f, &config); err != nil {
		log.Fatal("unmarshal config: ", err)
	}
	if env != "dev" && env != "prod" {
		log.Fatal("wrong enviroment")
	}
	config.Server.Enviroment = env
	config.Server.Secret = os.Getenv("SECRET")

	router := mux.NewRouter()
	server := &http.Server{
		Handler:           router,
		Addr:              fmt.Sprint(":", port),
		IdleTimeout:       time.Duration(config.Server.Timeouts) * time.Second,
		ReadTimeout:       time.Duration(config.Server.Timeouts) * time.Second,
		WriteTimeout:      time.Duration(config.Server.Timeouts) * time.Second,
		ReadHeaderTimeout: time.Duration(config.Server.Timeouts) * time.Second,
	}
	db, err := sql.Open("sqlite3", path.Join(config.Database.Path, config.Database.Name))
	if err != nil {
		log.Fatal("open db: ", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("ping db: ", err)
	}
	sqlite.MustMigration(db, config.Database.Migrations)
	storage := sqlite.New(db)
	templates, err := templ.Functions("init", config.Templates.Path)
	if err != nil {
		log.Fatal("parse templates: ", err)
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	locales := locales.MustLocales(config.Locales.Path)
	go func() {
		err := apiserver.New(router, server, storage, templates, locales, config.Server).Run()
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("server closed")
			return
		}
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatal("run server: ", err)
		}
	}()
	go Prom()
	<-interrupt
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("shutdown: ", err)
	}
}
