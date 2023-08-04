package handler

import (
	"chat-society-api/internal/cs"
	"chat-society-api/platform/trace"
	"chat-society-api/platform/web"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Handler struct {
	ChatService *cs.ChatService
}

func NewHandler(repo cs.ChatRepository) *Handler {
	return &Handler{
		ChatService: cs.NewChatService(repo),
	}
}

func (h *Handler) AddParticipant() func(w http.ResponseWriter, r *http.Request) {
	type b struct {
		Nick string `json:"nick"`
	}

	type j struct {
		ID   string `json:"id"`
		Nick string `json:"nick"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		body := r.Body
		defer func() {
			_ = body.Close()
		}()
		var req b
		err := json.NewDecoder(body).Decode(&req)

		if err != nil {
			web.ResponseBadRequest(w, "cannot create participant")
			return
		}

		participant := h.ChatService.AddParticipant(req.Nick)

		web.ResponseCreated(w, "participant created", &j{
			ID:   participant.ID,
			Nick: participant.Nick,
		})
	}
}

func (h *Handler) AddRoom() func(w http.ResponseWriter, r *http.Request) {
	type j struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`

		Owner string `json:"owner"`

		// configuration
		IsModerated bool    `json:"isModerated"`
		Moderator   *string `json:"moderator,omitempty"`

		IsOnlyAudio bool   `json:"isOnlyAudio"`
		IsOnlyText  bool   `json:"isOnlyText"`
		IsBoth      bool   `json:"isBoth"`
		Max         int    `json:"max"`
		URL         string `json:"url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req cs.RoomReq
		body := r.Body
		defer func() {
			_ = body.Close()
		}()
		err := json.NewDecoder(body).Decode(&req)

		if err != nil {
			web.ResponseBadRequest(w, "cannot create participant")
			return
		}

		room := h.ChatService.AddRoom(req)

		web.ResponseCreated(w, "", &j{
			ID:          room.ID,
			Title:       room.Title,
			Description: room.Description,
			Owner:       room.Owner,
			IsModerated: room.IsModerated,
			Moderator:   &room.Moderator,
			IsOnlyAudio: room.IsOnlyAudio,
			IsOnlyText:  room.IsOnlyText,
			IsBoth:      false,
			Max:         room.Max,
			URL:         room.URL,
		})
	}
}

func (h *Handler) RegisterWebsocket() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		participantID := r.URL.Query().Get("participantID")
		roomID := r.URL.Query().Get("roomID")

		if participantID == "" {
			log.Error().Msg("participant ID is empty")
			return
		}

		p, err := h.ChatService.FindParticipant(participantID)

		if err != nil {
			log.Error().Msg(err.Error())
			return
		}

		client, err := cs.RegistrationHandler(w, r)
		if err != nil {
			log.Error().Msgf("%s - %s", trace.Trace(), err.Error())
			return
		}

		p.Conn = client

		room, err := h.ChatService.FindRoomByID(roomID)

		if err != nil {
			log.Error().Msg(err.Error())
			return
		}

		// join participant to the chat service
		err = h.ChatService.Join(roomID, p)

		if err != nil {
			log.Error().Msg(err.Error())
			return
		}

		go client.ReadPump(func(msg string) {
			room.Broadcast <- msg
		})
		go client.WritePump()
	}
}
