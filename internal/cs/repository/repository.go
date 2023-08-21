package repository

import (
	"chat-society-api/internal/cs"
	"chat-society-api/platform/db"
	"errors"
	"github.com/rs/zerolog/log"

	"golang.org/x/crypto/bcrypt"
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
	// add in memory.
	err := s.Add(p)
	if err != nil {
		return err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(p.Password), 14)

	if err != nil {
		return err
	}

	// add in DB.
	err = s.Save("insert into participants (nick, password) values (?,?)", p.Nick, string(bytes))

	if err != nil {
		return err
	}

	return nil
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
	participants, ok := s.ParticipantsByRoom[roomID]

	if !ok {
		log.Warn().Msgf("no participants in room %s", roomID)
		return nil, errors.New("")
	}

	return participants, nil
}

func (s *Storage) GetAllRooms() ([]cs.Room, error) {
	rows, err := s.Many("select id, title, description, owner, is_moderated, is_only_text, is_only_audio, is_both, max, url from room")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []cs.Room
	for rows.Next() {
		var room cs.Room
		err = rows.Scan(
			&room.ID,
			&room.Title,
			&room.Description,
			&room.Owner,
			&room.IsModerated,
			&room.IsOnlyText,
			&room.IsOnlyAudio,
			&room.IsBoth,
			&room.Max,
			&room.URL,
		)

		if err != nil {
			log.Warn().Err(err).Msg("cannot get room")
			continue
		}

		result = append(result, room)
	}

	// TODO add binding with in-memory rooms here.
	return result, nil
}

func (s *Storage) Login(nick, passwordHash string) (string, error) {
	row := s.One("select password from participants where nick=?", nick)

	var password string
	err := row.Scan(&password)

	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))

	if err != nil {
		return "", err
	}

	return "tokenOK", nil
}
