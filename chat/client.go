package chat

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	socket   *websocket.Conn
	wsServer *WsServer
	send     chan *Message
	rooms    map[*Room]bool
	Name     string    `json:"name"`
	Id       uuid.UUID `json:"id"`
}

const (

	// Max time till next pong from peer
	pongWait   = 60 * time.Second * 5
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize, Subprotocols: []string{"token"}}

func newClient(conn *websocket.Conn, wsServer *WsServer, name string) *Client {
	return &Client{
		socket:   conn,
		wsServer: wsServer,
		send:     make(chan *Message),
		rooms:    make(map[*Room]bool),
		Name:     name,
		Id:       uuid.New(),
	}

}

func (c *Client) readPump() {
	defer func() {
		c.disconnect()
	}()
	c.socket.SetReadLimit(maxMessageSize)
	c.socket.SetReadDeadline(time.Now().Add(pongWait))
	c.socket.SetPongHandler(func(string) error { c.socket.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg *Message
		err := c.socket.ReadJSON(&msg)
		fmt.Println(msg)
		if err != nil {
			fmt.Printf("error on readPump: %s", err.Error())
			return
		}
		c.handleNewMessage(*msg)
	}
}

func (c *Client) writePump() {
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

func (c *Client) disconnect() {
	c.wsServer.unregister <- c
	for room := range c.rooms {
		room.leave <- c
	}
	fmt.Println("cerre el chan send")
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
	go client.readPump()

	wsServer.register <- client

}

func (client *Client) handleNewMessage(message Message) {

	// Attach the client object as the sender of the messsage.
	message.Sender = client

	switch message.Action {
	case SendMessageAction:
		// The send-message action, this will send messages to a specific room now.
		// Which room wil depend on the message Target
		roomName := message.Target.GetId()
		// Use the ChatServer method to find the room, and if found, broadcast!
		if room := client.wsServer.findRoomById(roomName); room != nil {
			room.forward <- &message
		}
	// We delegate the join and leave actions.
	case JoinRoomAction:
		client.handleJoinRoomMessage(message)

	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)

	case JoinRoomPrivateAction:
		client.handleJoinRoomPrivateMessage(message)
	}

}

func (client *Client) handleJoinRoomMessage(message Message) {
	roomName := message.Message
	client.joinRoom(roomName, nil)

}

func (client *Client) handleLeaveRoomMessage(message Message) {
	room := client.wsServer.findRoomById(message.Message)
	if room == nil {
		return
	}

	delete(client.rooms, room)

	room.leave <- client
}

func (client *Client) handleJoinRoomPrivateMessage(message Message) {

	target := client.wsServer.findClientByID(message.Message)
	if target == nil {
		return
	}

	// create unique room name combined to the two IDs
	roomName := message.Message + client.Id.String()

	client.joinRoom(roomName, target)
	target.joinRoom(roomName, client)

}

func (client *Client) notifyRoomJoined(room *Room, sender *Client) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
		Sender: sender,
	}

	client.send <- &message
}

func (client *Client) joinRoom(roomName string, sender *Client) {

	room := client.wsServer.findRoomByName(roomName)
	if room == nil {
		room = client.wsServer.createRoom(roomName, sender != nil)
	}

	// Don't allow to join private rooms through public room message
	if sender == nil && room.Private {
		return
	}

	if !client.isInRoom(room) {

		client.rooms[room] = true
		room.join <- client

		client.notifyRoomJoined(room, sender)
	}

}

func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}

	return false
}

func (client *Client) GetName() string {
	return client.Name
}
