package repository

import "chat-society-api/cmd/internal/cs"

var Participants = make(map[string]*cs.Participant, 1000)
var ParticipantsByRoom = make(map[string][]*cs.Participant)

var Rooms = make(map[string]*cs.Room)
