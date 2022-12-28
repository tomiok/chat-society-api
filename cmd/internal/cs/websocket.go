package cs

import (
	"chat-society-api/cmd/platform/trace"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"nhooyr.io/websocket"
	"time"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 1 * time.Second

	pingPeriod = (pongWait * 9) / 10
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
		w.close(fmt.Sprintf("%s - closing reader", trace.Trace()))
	}()

	for {
		msgType, msg, err := w.Read(context.Background())

		if err != nil {
			log.Error().Msgf("%s -  %s", trace.Trace(), err.Error())
			break
		}

		log.Info().Msgf("sending message. Type: %d", msgType)

		//message handler
		messageHandler(string(msg))
	}
}

// close the handler connection
func (w *WebsocketConn) close(msg string) {
	close(w.Send) //this close cause panic FIXME
	err := w.Close(websocket.StatusNormalClosure, msg)

	if err != nil {
		log.Warn().Msgf("%s cannot close connection %s", trace.Trace(), err.Error())
	}
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) (*WebsocketConn, error) {
	conn, err := websocket.Accept(w, r, nil)
	return &WebsocketConn{
		Conn: conn,
		Send: make(chan string),
	}, err
}

// WritePump pumps messages from the hub to the websocket connection.
// A goroutine running WritePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (w *WebsocketConn) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		w.close("err")
	}()
	for {
		select {
		case message, ok := <-w.Send:
			if !ok {
				w.close("no messages")
			}

			err := w.Write(context.Background(), websocket.MessageText, []byte(message))
			if err != nil {
				w.close("cannot write " + err.Error())
			}
		case <-ticker.C:
			err := w.Ping(context.Background())
			if err != nil {
				w.close(fmt.Sprintf("no ping available -  %v", err))
				break
			}
			log.Info().Msg("pinging...")
		}
	}
}
