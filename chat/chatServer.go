package chat

import "fmt"

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	rooms      map[*Room]bool
}

func NewWSServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
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

func (server *WsServer) registerClient(client *Client) {
	server.notifyClientJoined(client)
	server.listOnlineUsers(client)
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {

		delete(server.clients, client)
		server.notifyClientLeft(client)
	}
}

func (server *WsServer) broadcastToClients(message *Message) {
	for client := range server.clients {
		fmt.Println(message.Action)
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

func (server *WsServer) findRoomById(id string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetId() == id {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) findClientByID(ID string) *Client {
	var foundClient *Client
	for client := range server.clients {
		if client.Id.String() == ID {
			foundClient = client
			break
		}
	}

	return foundClient
}

func (server *WsServer) createRoom(name string, private bool) *Room {
	room := NewRoom(name, private)
	go room.Run()
	server.rooms[room] = true

	return room
}

func (server *WsServer) notifyClientJoined(client *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}

	server.broadcastToClients(message)
}

func (server *WsServer) notifyClientLeft(client *Client) {
	message := &Message{
		Action: UserLeftAction,
		Sender: client,
	}

	server.broadcastToClients(message)
}

func (server *WsServer) listOnlineUsers(client *Client) {
	for existingClient := range server.clients {
		message := &Message{
			Action: UserJoinedAction,
			Sender: existingClient,
		}

		client.send <- message
	}
}
