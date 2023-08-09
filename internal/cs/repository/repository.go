package repository

import (
	"chat-society-api/internal/cs"
	"chat-society-api/platform/db"
	"errors"
	"github.com/rs/zerolog/log"
)

type Storage struct {
	db.StorageService
	*InMemoryStorage
}

func NewStorage(mySql *db.MySql) *Storage {
	return &Storage{
		StorageService:  mySql,
		InMemoryStorage: NewInMemoryStorage(),
	}
}

func (s *Storage) AddParticipant(p *cs.Participant) error {
	return s.Add(p)
}

func (s *Storage) FindParticipant(id string) (*cs.Participant, error) {
	return s.Find(id)
}

func (s *Storage) CreateRoom(r *cs.Room) error {
	err := s.One("insert into room (id, title, description, owner, is_moderated, moderator, is_only_audio, is_only_text, is_both, max, url, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		r.ID, r.Title, r.Description, r.Owner, r.IsModerated, r.Moderator, r.IsOnlyAudio, r.IsOnlyText, r.IsBoth, 250, r.URL, r.CreatedAt).Err()

	if err != nil {
		return err
	}
	// add the room to memory
	return s.AddRoom(r)
}

func (s *Storage) FindRoom(id string) (*cs.Room, error) {
	query := "select id, title, description, owner, is_moderated, moderator, is_only_audio, is_only_text, is_both, max, url from room where id=?"

	row := s.GetByID(query, id)
	var res cs.Room

	err := row.Scan(
		&res.ID,
		&res.Title,
		&res.Description,
		&res.Owner,
		&res.IsModerated,
		&res.Moderator,
		&res.IsOnlyAudio,
		&res.IsOnlyText,
		&res.IsBoth,
		&res.Max,
		&res.URL,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *Storage) JoinParticipant(roomID string, p *cs.Participant) error {
	_, err := s.FindRoom(roomID)

	if err != nil {
		return err
	}

	return s.AddToRoom(roomID, p)
}

func (s *Storage) GetParticipantsByRoom(roomID string) ([]*cs.Participant, error) {
	participants, ok := cs.ParticipantsByRoom[roomID]

	if !ok {
		log.Warn().Msgf("no participants in room %s", roomID)
		return nil, errors.New("")
	}

	return participants, nil
}
