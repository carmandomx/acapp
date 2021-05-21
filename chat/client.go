package chat

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn
	wsServer *WsServer
	send     chan *Message
	rooms    map[*Room]bool
	Name     string `json:"name"`
}

const (

	// Max time till next pong from peer
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize, Subprotocols: []string{"token"}}

func newClient(conn *websocket.Conn, wsServer *WsServer, name string) *client {
	return &client{
		socket:   conn,
		wsServer: wsServer,
		send:     make(chan *Message),
		rooms:    make(map[*Room]bool),
		Name:     name,
	}

}

func (c *client) readPump() {
	defer func() {
		c.disconnect()
	}()
	c.socket.SetReadLimit(maxMessageSize)
	c.socket.SetReadDeadline(time.Now().Add(pongWait))
	c.socket.SetPongHandler(func(string) error { c.socket.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg *Message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		c.handleNewMessage(*msg)
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.socket.Close()
	}()
	for msg := range c.send {
		fmt.Println(msg)
		err := c.socket.WriteJSON(msg)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (c *client) disconnect() {
	c.wsServer.unregister <- c
	for room := range c.rooms {
		room.leave <- c
	}
	close(c.send)
	c.socket.Close()
}

func ServeWS(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {

	name, ok := r.URL.Query()["name"]

	if !ok || len(name[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	sock, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("ServeHTTP:", err)
		return
	}

	client := newClient(sock, wsServer, name[0])

	go client.writePump()
	client.readPump()

	wsServer.register <- client
	fmt.Println("me ejecute hasta aqui")
}

func (client *client) handleNewMessage(message Message) {

	// Attach the client object as the sender of the messsage.
	message.Sender = client

	switch message.Action {
	case SendMessageAction:
		// The send-message action, this will send messages to a specific room now.
		// Which room wil depend on the message Target
		roomName := message.Target
		// Use the ChatServer method to find the room, and if found, broadcast!
		if room := client.wsServer.findRoomByName(roomName); room != nil {
			room.forward <- &message
		}
	// We delegate the join and leave actions.
	case JoinRoomAction:
		client.handleJoinRoomMessage(message)

	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)
	}
}

func (client *client) handleJoinRoomMessage(message Message) {
	roomName := message.Message

	room := client.wsServer.findRoomByName(roomName)
	if room == nil {
		room = client.wsServer.createRoom(roomName)
	}

	client.rooms[room] = true

	room.join <- client
}

func (client *client) handleLeaveRoomMessage(message Message) {
	room := client.wsServer.findRoomByName(message.Message)

	delete(client.rooms, room)

	room.leave <- client
}

func (client *client) GetName() string {
	return client.Name
}
