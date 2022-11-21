package main

import (
	"chat-society-api/cmd/internal/cs/handler"
	"chat-society-api/cmd/internal/cs/repository"
)

type Deps struct {
	handler *handler.Handler
}

func buildDeps() *Deps {
	chatRepo := repository.NewChatRepository()
	chatHandler := handler.NewHandler(chatRepo)
	return &Deps{
		handler: chatHandler,
	}
}
