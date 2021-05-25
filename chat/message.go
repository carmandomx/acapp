package chat

import (
	"time"
)

type Message struct {
	Action  string    `json:"action"`
	Message string    `json:"message"`
	Target  *Room     `json:"target"`
	Sender  *Client   `json:"sender"`
	When    time.Time `json:"when"`
}

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const UserJoinedAction = "user-join"
const UserLeftAction = "user-left"
const JoinRoomPrivateAction = "join-room-private"
const RoomJoinedAction = "room-joined"
