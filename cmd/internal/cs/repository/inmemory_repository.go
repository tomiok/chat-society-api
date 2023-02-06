package repository

import (
	"chat-society-api/cmd/internal/cs"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
)

type inMemoryStorage struct{}

func newInMemoryStorage() *inMemoryStorage {
	return &inMemoryStorage{}
}

func (c *inMemoryStorage) Add(p *cs.Participant) error {
	Participants[p.ID] = p
	return nil
}

func (c *inMemoryStorage) Find(id string) (*cs.Participant, error) {
	p, ok := Participants[id]

	if !ok {
		log.Error().Msgf("cannot find participant with id %s", id)
		return nil, errors.New("participant not found")
	}

	return p, nil
}

func (c *inMemoryStorage) AddToRoom(roomID string, p *cs.Participant) error {
	participants, ok := ParticipantsByRoom[roomID]

	if !ok {
		return fmt.Errorf("room %s is not in our server", roomID)
	}

	participants = append(participants, p)
	ParticipantsByRoom[roomID] = participants
	return nil
}

func (c *inMemoryStorage) GetByRoom(roomID string) ([]*cs.Participant, error) {
	participants, ok := ParticipantsByRoom[roomID]

	if !ok {
		log.Warn().Msgf("no participants in room %s", roomID)
		return nil, errors.New("")
	}

	return participants, nil
}

// AddRoom create a room in memory.
func (c *inMemoryStorage) AddRoom(r *cs.Room) error {
	Rooms[r.ID] = r
	ParticipantsByRoom[r.ID] = make([]*cs.Participant, 1000)
	return nil
}

// GetParticipantsByRoom given a room, return all the participants. Maybe debug only.
func (c *inMemoryStorage) GetParticipantsByRoom(roomID string) ([]*cs.Participant, error) {
	participants, ok := ParticipantsByRoom[roomID]

	if !ok {
		log.Warn().Msgf("no participants in room %s", roomID)
		return nil, errors.New("")
	}

	return participants, nil
}
