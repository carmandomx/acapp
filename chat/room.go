package chat

import (
	"fmt"

	"github.com/carmandomx/acapp/trace"
)

type Room struct {
	name string

	forward chan *Message

	join chan *client

	leave chan *client

	clients map[*client]bool

	Tracer trace.Tracer
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.join:
			r.registerClientInRoom(client)
			r.Tracer.Trace("New client joined")

		case client := <-r.leave:
			r.unregisterClientInRoom(client)
			// close(client.send)
			r.Tracer.Trace("Client left")

		case msg := <-r.forward:
			r.broadcastToClientsInRoom(msg)
		}

	}
}

func (room *Room) registerClientInRoom(client *client) {
	room.notifyClientJoined(client)
	room.clients[client] = true
}

func (room *Room) unregisterClientInRoom(client *client) {

	delete(room.clients, client)

}

func (room *Room) broadcastToClientsInRoom(message *Message) {
	for client := range room.clients {
		client.send <- message
	}
}

func (room *Room) GetName() string {
	return room.name
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

func NewRoom(name string) *Room {
	return &Room{
		name:    name,
		forward: make(chan *Message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		Tracer:  trace.Off(),
	}
}

const welcomeMessage = "%s joined the room"

func (room *Room) notifyClientJoined(client *client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  room.name,
		Message: fmt.Sprintf(welcomeMessage, client.GetName()),
	}

	room.broadcastToClientsInRoom(message)
}
