package main

import (
	"log"
	"net/http"

	"gitlab.com/horo-go/horo"
	"gitlab.com/lattetalk/lattetalk/config"
	"gitlab.com/lattetalk/lattetalk/router"
	"gitlab.com/lattetalk/lattetalk/ws"
)

func main() {
	log.SetPrefix("LATTETALK  ")

	// Run WebSocket server
	go ws.Listen(":5050")

	// Run HTTP server
	r := newEngine()
	router.RegisterRoutes(r)
	log.Fatal(http.ListenAndServe(":5000", r))
}

func newEngine() *horo.Engine {
	cfg := horo.DefaultEngineConfig()
	cfg.BasePathPrefix = config.Config.Prefix
	return horo.NewEngine(cfg)
}
