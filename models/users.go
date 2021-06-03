package models

import (
	"strconv"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Password string `json:"password"`
	Username string `json:"username"`
	Name     string `json:"name"`
	PhotoUrl string `json:"photo_url"`
}

type IUser interface {
	GetId() string
	GetName() string
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRepository interface {
	FindById(Id int) (*User, error)
	FindByEmail(e string) (*User, error)
	Create(u *User) error
	GetAllUsers() ([]User, error)
	Delete(Id string) error
}

func (user User) GetId() string {
	return strconv.FormatUint(uint64(user.ID), 10)
}

func (user User) GetName() string {
	return user.Name
}
