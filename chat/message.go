package chat

import (
	"time"
)

type Message struct {
	Action  string    `json:"action"`
	Message string    `json:"message"`
	Target  string    `json:"target"`
	Sender  *client   `json:"sender"`
	When    time.Time `json:"when"`
}

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
