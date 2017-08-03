package project

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
	"time"
)

var (
	// Time allowed to write the file to the client.
    writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	UpdateChannel = make(chan Project)
)

type message struct {
	Action string `json:"action"`
	Project Project `json:"project"`
}

func RegisterWebsocket(r *mux.Router) {
	log.Println("ws on /ws")
	r.HandleFunc("/ws", wsHandler)
}

func wsHandler(writer http.ResponseWriter, request *http.Request) {

	log.Println("Handle ws request")

	ws, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		} else {
			log.Println("Error service ws: ", err)
		}

		return
	}

	go listenForChangedProjects(ws, UpdateChannel)
	go answerPong(ws)
}

func listenForChangedProjects(conn *websocket.Conn, projectUpdate <-chan Project) {

	pingTicker := time.NewTicker(pingPeriod)

	defer func() {
		pingTicker.Stop()
		conn.Close()
	}()

	for {
		select {
		case p := <- projectUpdate:
			log.Println("Send project update to ws")

			conn.WriteJSON(message{Action:"update", Project:p})
			break
		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("Error ping ws: ", err)
				return
			}
		}
	}
}

func answerPong(conn *websocket.Conn) {
	defer func() {
		conn.Close()
	}()

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Println(message)
	}
}
