package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/carmandomx/acapp/config"
	"github.com/carmandomx/acapp/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	socket   *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	rooms    map[*Room]bool
	Name     string `json:"name"`
	Id       uint   `json:"id"`
}

var (
	newline = []byte{'\n'}
)

const (
	writeWait = 10 * time.Second
	// Max time till next pong from peer
	pongWait   = 60 * time.Second * 5
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

type contextKey string

const UserContextKey = contextKey("user")

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize, Subprotocols: []string{"token"}}

func newClient(conn *websocket.Conn, wsServer *WsServer, name string, id uint) *Client {
	return &Client{
		socket:   conn,
		wsServer: wsServer,
		send:     make(chan []byte),
		rooms:    make(map[*Room]bool),
		Name:     name,
		Id:       id,
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

		_, jsonMsg, err := c.socket.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		c.handleNewMessage(jsonMsg)

		// var msg *Message
		// err := c.socket.ReadJSON(&msg)
		// fmt.Println(msg)
		// if err != nil {
		// 	fmt.Printf("error on readPump: %s", err.Error())
		// 	return
		// }
		// c.handleNewMessage(*msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
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

func ServeWS(wsServer *WsServer, w http.ResponseWriter, r *http.Request, id uint) {

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

	client := newClient(sock, wsServer, name[0], id)

	go client.writePump()
	go client.readPump()

	wsServer.register <- client

}

func (client *Client) handleNewMessage(jsonMsg []byte) {

	var message Message

	if err := json.Unmarshal(jsonMsg, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

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

	// target := client.wsServer.findClientByID(message.Message)
	target := client.wsServer.findUserById(message.Message)

	// create unique room name combined to the two IDs
	roomName := message.Message + strconv.FormatUint(uint64(client.Id), 10)

	joinedRoom := client.joinRoom(roomName, target)

	if joinedRoom != nil {
		client.inviteTargetUser(target, joinedRoom)
	}

	// client.joinRoom(roomName, target)
	// target.joinRoom(roomName, client)

}

func (client *Client) notifyRoomJoined(room *Room, sender models.IUser) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
		Sender: sender,
	}

	client.send <- message.encode()
}

func (client *Client) joinRoom(roomName string, sender models.IUser) *Room {

	room := client.wsServer.findRoomByName(roomName)
	if room == nil {
		room = client.wsServer.createRoom(roomName, sender != nil)
	}

	// Don't allow to join private rooms through public room message
	if sender == nil && room.Private {
		return nil
	}

	if !client.isInRoom(room) {
		client.rooms[room] = true
		room.join <- client
		client.notifyRoomJoined(room, sender)
	}
	return room

}

func (client *Client) inviteTargetUser(target models.User, room *Room) {
	inviteMessage := &Message{
		Action:  JoinRoomPrivateAction,
		Message: target.GetId(),
		Target:  room,
		Sender:  client,
	}

	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, inviteMessage.encode()).Err(); err != nil {
		log.Println(err)
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

func (client *Client) GetId() string {
	return strconv.FormatUint(uint64(client.Id), 10)
}
