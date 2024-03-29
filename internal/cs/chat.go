package cs

import (
	"chat-society-api/platform/trace"
	"fmt"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
)

const seed = 9999

type Message struct {
	Sender  string
	Message string
}

type ChatService struct {
	ChatRepository
	*ChatPeople
}

type ChatRepository interface {
	AddParticipant(p *Participant) error
	FindParticipant(id string) (*Participant, error)
	Login(nick, pass string) (string, error)

	CreateRoom(r *Room) error
	FindRoom(id string) (*Room, error)

	JoinParticipant(roomID string, p *Participant) error
	GetParticipantsByRoom(roomID string) ([]*Participant, error)
	GetAllRooms() ([]Room, error)
}

func NewChatService(r ChatRepository, chatPeople *ChatPeople) *ChatService {
	return &ChatService{
		r,
		chatPeople,
	}
}

// Participant is the actor of the room. Could join only one room at once.
type Participant struct {
	ID       string  // uid
	Nick     string  // the name given + some random numer
	Password string  // the pass :D
	Gender   *string // optional gender

	CurrentRoom string

	Conn *WebsocketConn

	JoinedAt time.Time // TODO define this or move up to current room
}

func (c *ChatService) AddParticipant(nick string) *Participant {
	participant := &Participant{
		ID:       generateRandomRune(),
		Nick:     NickGenerator(nick),
		JoinedAt: time.Now(),
	}
	_ = c.ChatRepository.AddParticipant(participant)
	return participant
}

func (c *ChatService) FindParticipant(id string) (*Participant, error) {
	return c.ChatRepository.FindParticipant(id)
}

func (c *ChatService) Join(roomID string, p *Participant) error {
	return c.ChatRepository.JoinParticipant(roomID, p)
}

func NickGenerator(nick string) string {
	src := rand.NewSource(time.Now().UnixMilli())
	i := rand.New(src).Intn(seed)
	return fmt.Sprintf("%s#%d", nick, i)
}

// Room

type Room struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Owner string `json:"owner"`

	// configuration
	IsModerated bool   `json:"is_moderated"`
	Moderator   string `json:"moderator"`

	IsOnlyAudio bool `json:"is_only_audio"`
	IsOnlyText  bool `json:"is_only_text"`
	IsBoth      bool `json:"is_both"`

	Max          int                    `json:"max"` // 250 hardcoded max number js
	Participants map[string]Participant `json:"-"`

	Broadcast chan string `json:"-"`
	URL       string      `json:"url"`
	CreatedAt time.Time   `json:"created_at"`
}

type RoomReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (c *ChatService) AddRoom(roomReq RoomReq) *Room {
	roomID := generateRandomRune()
	ws := fmt.Sprintf("ws://localhost:9001/join?roomID=%s&participantID={}", roomID)
	room := &Room{
		ID:           roomID,
		Title:        roomReq.Title,
		Description:  roomReq.Description,
		Owner:        "system",
		IsModerated:  false,
		Moderator:    "",
		IsOnlyAudio:  false,
		IsOnlyText:   true,
		IsBoth:       false,
		Max:          250,
		Participants: make(map[string]Participant, 250),
		Broadcast:    make(chan string),
		URL:          ws,
		CreatedAt:    time.Now(),
	}
	_ = c.ChatRepository.CreateRoom(room)

	go c.BroadcastMessage(roomID, room.Broadcast)
	return room
}

func (c *ChatService) BroadcastMessage(roomID string, ch chan string) {
	for {
		select {
		case message := <-ch:
			c.applyBroadcast(roomID, message)
		}
	}

}

func (c *ChatService) applyBroadcast(roomID, message string) {
	participants, err := c.GetParticipantsByRoom(roomID)

	if err != nil {
		log.Warn().Msgf("%s - cannot get participants %s", trace.Trace(), err.Error())
		return
	}
	for _, participant := range participants {
		part := participant
		go func() {
			part.Conn.Send <- message
		}()
	}
}

func (c *ChatService) FindRoomByID(id string) (*Room, error) {
	_, err := c.FindRoom(id)
	if err != nil {
		return nil, err
	}

	return c.Rooms[id], nil
}

func (c *ChatService) GetAllRooms() ([]Room, error) {
	rooms, err := c.ChatRepository.GetAllRooms()

	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (c *ChatService) Login(nick, pass string) (string, error) {
	return c.ChatRepository.Login(nick, pass)
}

var runes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func generateRandomRune() string {
	randRune := make([]rune, 11)

	for i := range randRune {
		randRune[i] = runes[rand.Intn(len(runes))]
	}
	return string(randRune)
}

type ChatPeople struct {
	// List of connected participants
	Participants map[string]*Participant
	// List of rooms
	Rooms map[string]*Room
	// roomID and every participant in that room
	ParticipantsByRoom map[string][]*Participant
}
