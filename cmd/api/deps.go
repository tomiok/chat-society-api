package main

import (
	"chat-society-api/internal/cs/handler"
	"chat-society-api/internal/cs/repository"
	"chat-society-api/platform/db"
	_ "github.com/go-sql-driver/mysql"
)

var mysqlURI = "tomi:tomi@tcp(localhost:3306)/chats_dev?parseTime=true"

type Deps struct {
	handler *handler.Handler
}

func buildDeps() *Deps {
	mySQL, err := db.New(mysqlURI)
	if err != nil {
		panic(err)
	}
	chatRepo := repository.NewStorage(mySQL)
	chatHandler := handler.NewHandler(chatRepo)
	return &Deps{
		handler: chatHandler,
	}
}
