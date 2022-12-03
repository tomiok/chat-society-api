package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"time"
)

func main() {
	run()
}

func run() {
	deps := buildDeps()

	r := chi.NewRouter()

	r.Post("/participants", deps.handler.AddParticipant())
	r.Post("/rooms", deps.handler.AddRoom())
	r.Get("/ws", deps.handler.RegisterWebsocket())

	port := os.Getenv("PORT")

	if port == "" {
		port = "9001"
	}

	srv := &http.Server{
		Addr: ":" + port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	server := server{srv}
	server.Start()
}
