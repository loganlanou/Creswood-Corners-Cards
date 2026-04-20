package main

import (
	"log"
	"net/http"

	"creswoodcornerscards/app"
	"creswoodcornerscards/config"
	"creswoodcornerscards/data"
)

func main() {
	cfg := config.Load()
	store, err := data.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	server, err := app.New(cfg, store)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting %s on %s", cfg.AppName, cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
