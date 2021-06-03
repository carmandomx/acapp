package chat

import (
	"encoding/json"
	"log"

	"github.com/carmandomx/acapp/models"
)

type Message struct {
	Action  string       `json:"action"`
	Message string       `json:"message"`
	Target  *Room        `json:"target"`
	Sender  models.IUser `json:"sender"`
}

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const UserJoinedAction = "user-join"
const UserLeftAction = "user-left"
const JoinRoomPrivateAction = "join-room-private"
const RoomJoinedAction = "room-joined"
const ListUsersOnRoom = "users-on-room"

func (i Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}

func (message *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	msg := &struct {
		Sender Client `json:"sender"`
		*Alias
	}{
		Alias: (*Alias)(message),
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	message.Sender = &msg.Sender
	return nil
}
