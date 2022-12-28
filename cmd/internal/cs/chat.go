package cs

import (
	"chat-society-api/cmd/internal/platform/trace"
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
}

type ChatRepository interface {
	AddParticipant(p *Participant) error
	FindParticipant(id string) (*Participant, error)

	AddRoom(r *Room) error
	FindRoom(id string) (*Room, error)

	JoinParticipant(roomID string, p *Participant) error
	GetParticipantsByRoom(roomID string) ([]*Participant, error)
}

func NewChatService(r ChatRepository) *ChatService {
	return &ChatService{
		r,
	}
}

// Participant is the actor of the room. Could join only one room at once.
type Participant struct {
	ID     string  // uid
	Nick   string  // the name given + some random numer
	Gender *string // optional gender

	CurrentRoom string

	Conn *WebsocketConn

	JoinedAt time.Time // TODO define this or move up to current room
}

var i int

func (c *ChatService) AddParticipant(nick string) *Participant {

	participant := &Participant{
		ID:       generateRandomRune(),
		Nick:     NickGenerator(nick),
		JoinedAt: time.Now(),
	}
	_ = c.ChatRepository.AddParticipant(participant)
	i++
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
	ID          string
	Title       string
	Description string

	Owner string

	// configuration
	IsModerated bool
	Moderator   *string

	IsOnlyAudio bool
	IsOnlyText  bool
	IsBoth      bool

	Max int // 250 hardcoded max number

	Participants map[string]Participant

	Broadcast chan string

	URL string

	CreatedAt time.Time
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
		Moderator:    nil,
		IsOnlyAudio:  false,
		IsOnlyText:   true,
		IsBoth:       false,
		Max:          250,
		Participants: make(map[string]Participant, 250),
		Broadcast:    make(chan string),
		URL:          ws,
		CreatedAt:    time.Now(),
	}
	_ = c.ChatRepository.AddRoom(room)

	go c.BroadcastMessage(roomID, room.Broadcast)
	return room
}

func (c *ChatService) BroadcastMessage(roomID string, ch chan string) {
	go func() {
		for {
			select {
			case message := <-ch:
				c.applyBroadcast(roomID, message)
			}
		}
	}()
}

func (c *ChatService) applyBroadcast(roomID, message string) {
	participants, err := c.GetParticipantsByRoom(roomID)

	if err != nil {
		log.Warn().Msgf("%s - cannot get participants %s", trace.Trace(), err.Error())
		return
	}
	for _, participant := range participants {
		participant := participant
		go func() {
			participant.Conn.Send <- message
		}()
	}
}

func (c *ChatService) FindRoomByID(id string) (*Room, error) {
	return c.FindRoom(id)
}

var runes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func generateRandomRune() string {
	randRune := make([]rune, 11)

	for i := range randRune {
		rand.Seed(time.Now().UnixNano())

		randRune[i] = runes[rand.Intn(len(runes))]
	}
	return string(randRune)
}
