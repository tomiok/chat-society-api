package repository

import (
	"chat-society-api/internal/cs"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
)

type InMemoryStorage struct {
	*cs.ChatPeople
}

func NewInMemoryStorage(chatPeople *cs.ChatPeople) *InMemoryStorage {
	return &InMemoryStorage{
		chatPeople,
	}
}

func (c *InMemoryStorage) Add(p *cs.Participant) error {
	c.Participants[p.ID] = p
	return nil
}

func (c *InMemoryStorage) Find(id string) (*cs.Participant, error) {
	p, ok := c.Participants[id]

	if !ok {
		log.Error().Msgf("cannot find participant with id %s", id)
		return nil, errors.New("participant not found in memory")
	}

	return p, nil
}

func (c *InMemoryStorage) FindRoom(id string) (*cs.Room, error) {
	r, ok := c.Rooms[id]

	if !ok {
		log.Error().Msgf("cannot find participant with id %s", id)
		return nil, errors.New("participant not found in memory")
	}

	return r, nil
}

func (c *InMemoryStorage) AddToRoom(roomID string, p *cs.Participant) error {
	participants, ok := c.ParticipantsByRoom[roomID]

	if !ok {
		return fmt.Errorf("room %s is not in our server", roomID)
	}

	participants = append(participants, p)
	c.ParticipantsByRoom[roomID] = participants
	return nil
}

func (c *InMemoryStorage) GetByRoom(roomID string) ([]*cs.Participant, error) {
	participants, ok := c.ParticipantsByRoom[roomID]

	if !ok {
		log.Warn().Msgf("no participants in room %s", roomID)
		return nil, errors.New("")
	}

	return participants, nil
}

// AddRoom create a room in memory.
func (c *InMemoryStorage) AddRoom(r *cs.Room) error {
	c.Rooms[r.ID] = r
	c.ParticipantsByRoom[r.ID] = make([]*cs.Participant, 0)
	return nil
}

// GetParticipantsByRoom given a room, return all the participants. Maybe debug only.
func (c *InMemoryStorage) GetParticipantsByRoom(roomID string) ([]*cs.Participant, error) {
	participants, ok := c.ParticipantsByRoom[roomID]

	if !ok {
		log.Warn().Msgf("no participants in room %s", roomID)
		return nil, errors.New("")
	}

	return participants, nil
}
