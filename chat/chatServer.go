package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/carmandomx/acapp/config"
	"github.com/carmandomx/acapp/models"
)

const PubSubGeneralChannel = "general"

type WsServer struct {
	clients        map[*Client]bool
	register       chan *Client
	unregister     chan *Client
	broadcast      chan Message
	rooms          map[*Room]bool
	users          []models.User
	roomRepository models.RoomRepository
	userRepository models.UserRepository
}

func NewWSServer(roomRepository models.RoomRepository, userRepository models.UserRepository) *WsServer {
	ws := &WsServer{
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan Message),
		rooms:          make(map[*Room]bool),
		users:          make([]models.User, 1),
		roomRepository: roomRepository,
		userRepository: userRepository,
	}

	// ws.users, _ = userRepository.GetAllUsers()

	return ws
}

// Run our websocket server, accepting various requests
func (server *WsServer) Run() {
	go server.listenPubSubChannel()
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

	user, _ := server.userRepository.FindById(int(client.Id))
	server.users = append(server.users, *user)
	server.clients[client] = true
	server.publishClientJoined(client)
	// server.listOnlineClients(client)
}

func (server *WsServer) publishClientJoined(client *Client) {
	fmt.Println(client.Id)
	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}

	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, message).Err(); err != nil {
		log.Println(err)
	}
}

func (server *WsServer) publishClientLeft(client *Client) {

	message := &Message{
		Action: UserLeftAction,
		Sender: client,
	}

	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, message).Err(); err != nil {
		log.Println(err)
	}
}

func (server *WsServer) listenPubSubChannel() {

	pubsub := config.Redis.Subscribe(ctx, PubSubGeneralChannel)
	ch := pubsub.Channel()
	for msg := range ch {

		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Printf("Error on unmarshal JSON message %s", err)
			return
		}

		switch message.Action {
		case UserJoinedAction:
			server.handleUserJoined(message)
		case UserLeftAction:
			server.handleUserLeft(message)
		case JoinRoomPrivateAction:
			server.handleUserJoinPrivate(message)
		}
	}
}

func (server *WsServer) handleUserJoinPrivate(message Message) {
	// Find client for given user, if found add the user to the room.
	targetClients := server.findClientsByID(message.Message)
	for _, targetClient := range targetClients {
		targetClient.joinRoom(message.Target.GetName(), message.Sender)
	}
}

func (server *WsServer) handleUserJoined(message Message) {
	// Add the user to the slice
	id, _ := strconv.ParseInt(message.Sender.GetId(), 0, 32)

	user, _ := server.userRepository.FindById(int(id))
	server.users = append(server.users, *user)
	server.broadcastToClients(message)
}

func (server *WsServer) handleUserLeft(message Message) {
	// Remove the user from the slice
	for i, user := range server.users {

		if user.GetId() == message.Sender.GetId() {
			server.users[i] = server.users[len(server.users)-1]
			server.users = server.users[:len(server.users)-1]
			break
		}
	}
	server.broadcastToClients(message)
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {

		if _, ok := server.clients[client]; ok {
			delete(server.clients, client)

			server.publishClientLeft(client)
		}
	}
}

func (server *WsServer) broadcastToClients(message Message) {
	for client := range server.clients {
		fmt.Println(message.Sender)
		client.send <- message.encode()
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

	fmt.Println(foundRoom)

	if foundRoom == nil {
		foundRoom = server.runRoomFromRepo(name)
	}

	return foundRoom
}

func (server *WsServer) runRoomFromRepo(name string) *Room {
	var room *Room
	dbRoom, err := server.roomRepository.FindRoomByName(name)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	room = NewRoom(dbRoom.GetName(), dbRoom.GetPrivate())

	go room.Run()
	server.rooms[room] = true

	return room
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

func (server *WsServer) createRoom(name string, private bool) *Room {
	room := NewRoom(name, private)

	server.roomRepository.AddRoom(&models.Room{
		Name:    room.Name,
		Private: room.Private,
		Id:      room.GetId(),
	})
	go room.Run()
	server.rooms[room] = true

	return room
}

func (server *WsServer) findUserById(Id string) models.User {
	var foundUser models.User

	for _, user := range server.users {
		if user.GetId() == Id {
			foundUser = user
			break
		}
	}

	return foundUser
}

func (server *WsServer) findClientsByID(ID string) []*Client {
	// Find all clients for given user ID.
	var foundClients []*Client
	for client := range server.clients {
		if client.GetId() == ID {
			foundClients = append(foundClients, client)
		}
	}

	return foundClients
}

func (server *WsServer) listOnlineClients(client *Client) {
	// Find unique users instead of returning all users.
	var uniqueUsers = make(map[string]bool)
	for _, user := range server.users {
		if user.ID == 0 {
			continue
		}
		if ok := uniqueUsers[user.GetId()]; !ok {
			message := &Message{
				Action: UserJoinedAction,
				Sender: user,
			}
			uniqueUsers[user.GetId()] = true
			client.send <- message.encode()
		}
	}
}
