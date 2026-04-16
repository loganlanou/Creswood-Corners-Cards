package handler

import (
	"net/http"
	"sync"

	"creswoodcornerscards/internal/app"
	"creswoodcornerscards/internal/config"
	"creswoodcornerscards/internal/data"
)

var (
	once    sync.Once
	handler http.Handler
	initErr error
)

func initHandler() {
	cfg := config.Load()
	if cfg.DatabasePath == "./cards.db" {
		cfg.DatabasePath = "/tmp/cards.db"
	}

	store, err := data.Open(cfg)
	if err != nil {
		initErr = err
		return
	}

	server, err := app.New(cfg, store)
	if err != nil {
		initErr = err
		return
	}

	handler = server.Handler()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initHandler)
	if initErr != nil {
		http.Error(w, initErr.Error(), http.StatusInternalServerError)
		return
	}
	handler.ServeHTTP(w, r)
}
