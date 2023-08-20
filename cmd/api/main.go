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
	r.Get("/rooms", deps.handler.GetRooms())
	r.Get("/ws", deps.handler.RegisterWebsocket())

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

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
