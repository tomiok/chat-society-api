package main

import (
	"chat-society-api/internal/cs"
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
	chatPeople := &cs.ChatPeople{
		Participants:       make(map[string]*cs.Participant),
		Rooms:              make(map[string]*cs.Room),
		ParticipantsByRoom: make(map[string][]*cs.Participant),
	}

	mySQL, err := db.New(mysqlURI)
	if err != nil {
		panic(err)
	}
	chatRepo := repository.NewStorage(mySQL, chatPeople)
	chatHandler := handler.NewHandler(chatRepo, chatPeople)
	return &Deps{
		handler: chatHandler,
	}
}
