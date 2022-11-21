package repository

import (
	"chat-society-api/cmd/internal/cs"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
)

var Participants = make(map[string]*cs.Participant, 1000)
var Rooms = make(map[string]*cs.Room, 1000)

var ParticipantsByRoom = make(map[string][]*cs.Participant)

type ChatRepository struct{}

func NewChatRepository() *ChatRepository {
	return &ChatRepository{}
}

func (c *ChatRepository) AddParticipant(p *cs.Participant) error {
	Participants[p.ID] = p
	return nil
}

func (c *ChatRepository) AddRoom(r *cs.Room) error {
	Rooms[r.ID] = r
	ParticipantsByRoom[r.ID] = []*cs.Participant{}
	return nil
}

func (c *ChatRepository) FindParticipant(id string) (*cs.Participant, error) {
	p, ok := Participants[id]

	if !ok {
		log.Error().Msgf("cannot find participant with id %s", id)
		return nil, errors.New("participant not found")
	}

	return p, nil
}

func (c *ChatRepository) FindRoom(id string) (*cs.Room, error) {
	r, ok := Rooms[id]

	if !ok {
		log.Error().Msgf("cannot find room with id %s", id)
		return nil, errors.New("cannot find room")
	}

	return r, nil
}

func (c *ChatRepository) JoinParticipant(roomID string, p *cs.Participant) error {
	participants, ok := ParticipantsByRoom[roomID]

	if !ok {
		return fmt.Errorf("room %s is not in our server", roomID)
	}

	participants = append(participants, p)
	ParticipantsByRoom[roomID] = participants
	return nil
}

func (c *ChatRepository) GetParticipantsByRoom(roomID string) ([]*cs.Participant, error) {
	participants, ok := ParticipantsByRoom[roomID]

	if !ok {
		log.Warn().Msgf("no participants in room %s", roomID)
		return nil, errors.New("")
	}

	return participants, nil
}
