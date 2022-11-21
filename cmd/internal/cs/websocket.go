package cs

import (
	"chat-society-api/cmd/internal/platform/trace"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 5 * time.Minute

	// Time allowed to read the next pong message from the peer.
	pongWait = 1 * time.Hour

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 2

	pingPeriod = (pongWait * 9) / 10
)

var (
	u = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type WebsocketConn struct {
	*websocket.Conn

	Send chan string
}

// ReadPump pumps messages from the websocket connection to the hub.
// The application runs ReadPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (w *WebsocketConn) ReadPump(messageHandler func(msg string)) {
	defer func() {
		w.close()
	}()
	w.Conn.SetReadLimit(maxMessageSize)
	_ = w.Conn.SetReadDeadline(time.Now().Add(pongWait))
	w.Conn.SetPongHandler(
		func(string) error {
			return w.Conn.SetReadDeadline(time.Now().Add(pongWait))
		},
	)
	for {
		_, s, err := w.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Warn().Msgf("%s unexpected close error", trace.Trace())
			}
			log.Warn().Msgf("%s cannot read message %s", trace.Trace(), err.Error())
			break
		}

		//message handler
		messageHandler(string(s))
	}
}

// close the handler connection
func (w *WebsocketConn) close() {
	close(w.Send)
	err := w.Close()

	if err != nil {
		log.Warn().Msgf("%s cannot close connection %s", trace.Trace(), err.Error())
	}
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) (*WebsocketConn, error) {
	u.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	responseHeader := make(http.Header)
	conn, err := u.Upgrade(w, r, responseHeader)

	return &WebsocketConn{
		Conn: conn,
		Send: make(chan string),
	}, err
}

// WritePump pumps messages from the hub to the websocket connection.
//
// A goroutine running WritePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (w *WebsocketConn) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		w.close()
	}()
	for {
		select {
		case message, ok := <-w.Send:
			_ = w.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = w.WriteMessage(websocket.CloseMessage, []byte{})
				log.Warn().Msg("websocket is closed")
				return
			}

			err := w.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				return
			}
		case <-ticker.C:
			_ = w.SetWriteDeadline(time.Now().Add(writeWait))
			if err := w.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
