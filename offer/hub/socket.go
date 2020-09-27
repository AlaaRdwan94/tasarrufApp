package hub

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// pongwait is the time the server awaits for a pong message
var pongWait = 60 * time.Second

// pingPeriod is the interval for sending a ping message
var pingPeriod = (pongWait * 9) / 10

// wsupgrader upgrades HTTP/HTTPS connection to WS connection
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// AuthMessage the first message expected from the client upon connection
type AuthMessage struct {
	Token string `json:"token"`
}

func (h *usersHub) ConnectWebSocket(w http.ResponseWriter, r *http.Request) {
	wsupgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var authMessage AuthMessage
	err = conn.ReadJSON(&authMessage)
	if err != nil {
		conn.WriteMessage(websocket.CloseMessage, []byte("error:invalid token"))
		conn.Close()
		return
	}
	if authMessage.Token == "" {
		conn.WriteMessage(websocket.CloseMessage, []byte("error: empty token"))
		conn.Close()
		return
	}
}
