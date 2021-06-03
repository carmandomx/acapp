package chat

import (
	"context"
	"fmt"
	"log"

	"github.com/carmandomx/acapp/config"
	"github.com/google/uuid"
)

var ctx = context.Background()

type Room struct {
	Name string `json:"name"`

	forward chan *Message

	join chan *Client

	leave chan *Client

	clients map[*Client]bool

	Id uuid.UUID `json:"id"`

	Private bool `json:"private"`
}

func (r *Room) Run() {
	go r.subscribeToRoomMessages()
	for {
		select {
		case client := <-r.join:
			r.registerClientInRoom(client)

		case client := <-r.leave:
			r.unregisterClientInRoom(client)

		case msg := <-r.forward:
			r.publishRoomMessage(msg)
		}

	}
}

func (room *Room) registerClientInRoom(client *Client) {
	if !room.Private {
		room.notifyClientJoined(client)

	}
	room.clients[client] = true

}

func (room *Room) unregisterClientInRoom(client *Client) {

	delete(room.clients, client)

}

func (room *Room) GetName() string {
	return room.Name
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

func NewRoom(name string, private bool) *Room {
	return &Room{
		Name:    name,
		forward: make(chan *Message),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		clients: make(map[*Client]bool),
		Id:      uuid.New(),
		Private: private,
	}
}

const welcomeMessage = "%s joined the room"

func (room *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  room,
		Message: fmt.Sprintf(welcomeMessage, client.GetName()),
	}

	room.publishRoomMessage(message)
}

func (room *Room) GetId() string {
	return room.Id.String()
}

func (room *Room) publishRoomMessage(message *Message) {
	err := config.Redis.Publish(ctx, room.GetName(), message).Err()

	if err != nil {
		log.Println("Error al publicar mensaje")
		log.Println(err)
	}
}

func (room *Room) subscribeToRoomMessages() {
	pubsub := config.Redis.Subscribe(ctx, room.GetName())

	ch := pubsub.Channel()

	for msg := range ch {
		fmt.Println(msg)
		room.broadcastToClientsInRoom([]byte(msg.Payload))
	}
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}
