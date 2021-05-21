package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Password string `json:"password"`
	Username string `json:"username"`
	Name     string `json:"name"`
	PhotoUrl string `json:"photo_url"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRepository interface {
	FindById(Id int) (*User, error)
	FindByEmail(e string) (*User, error)
	Create(u *User) error
	// Update(u *User) error
	Delete(Id string) error
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.Password, err = HashPassword(u.Password)

	if err != nil {
		return err
	}

	return nil

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
