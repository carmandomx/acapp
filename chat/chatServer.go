package chat

import "fmt"

type WsServer struct {
	clients    map[*client]bool
	register   chan *client
	unregister chan *client
	broadcast  chan *Message
	rooms      map[*Room]bool
}

func NewWSServer() *WsServer {
	return &WsServer{
		clients:    make(map[*client]bool),
		register:   make(chan *client),
		unregister: make(chan *client),
		broadcast:  make(chan *Message),
		rooms:      make(map[*Room]bool),
	}
}

// Run our websocket server, accepting various requests
func (server *WsServer) Run() {
	for {
		select {

		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)

		case message := <-server.broadcast:
			server.broadcastToClients(message)
		}

	}
}

func (server *WsServer) registerClient(client *client) {
	server.clients[client] = true
	fmt.Println(server.clients)
}

func (server *WsServer) unregisterClient(client *client) {
	delete(server.clients, client)
}

func (server *WsServer) broadcastToClients(message *Message) {
	for client := range server.clients {
		client.send <- message
	}
}

func (server *WsServer) findRoomByName(name string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetName() == name {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) createRoom(name string) *Room {
	room := NewRoom(name)
	go room.Run()
	server.rooms[room] = true

	return room
}
