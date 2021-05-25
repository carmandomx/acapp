package chat

import (
	"fmt"

	"github.com/carmandomx/acapp/trace"
	"github.com/google/uuid"
)

type Room struct {
	Name string `json:"name"`

	forward chan *Message

	join chan *Client

	leave chan *Client

	clients map[*Client]bool

	Tracer trace.Tracer

	Id uuid.UUID `json:"id"`

	Private bool `json:"private"`
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.join:
			r.registerClientInRoom(client)
			r.Tracer.Trace("New client joined")

		case client := <-r.leave:
			r.unregisterClientInRoom(client)
			r.Tracer.Trace("Client left")

		case msg := <-r.forward:
			r.broadcastToClientsInRoom(msg)
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

func (room *Room) broadcastToClientsInRoom(message *Message) {
	for client := range room.clients {
		client.send <- message
	}
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
		Tracer:  trace.Off(),
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

	room.broadcastToClientsInRoom(message)
}

func (room *Room) GetId() string {
	return room.Id.String()
}
