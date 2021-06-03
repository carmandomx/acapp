package repositories

import (
	"github.com/carmandomx/acapp/models"
	"gorm.io/gorm"
)

type RoomRepo struct {
	db *gorm.DB
}

func NewRoomRepo(db *gorm.DB) *RoomRepo {
	return &RoomRepo{
		db: db,
	}
}

func (r *RoomRepo) AddRoom(room *models.Room) error {
	q := r.db.Create(&room)

	if q.Error != nil {
		return q.Error
	}

	return nil
}

func (r *RoomRepo) FindRoomByName(name string) (*models.Room, error) {
	var room *models.Room
	q := r.db.Where(&models.Room{
		Name: name,
	}).First(&room)

	if q.Error != nil {
		return nil, q.Error
	}

	return room, nil
}
