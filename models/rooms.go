package models

type Room struct {
	Id      string `gorm:"primaryKey"`
	Name    string
	Private bool
}

type IRoom interface {
	GetId() uint
	GetName() string
	GetPrivate() bool
}

type RoomRepository interface {
	AddRoom(room *Room) error
	FindRoomByName(name string) (*Room, error)
}

func (room *Room) GetId() string {
	return room.Id
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetPrivate() bool {
	return room.Private
}
