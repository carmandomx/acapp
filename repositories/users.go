package repositories

import (
	"errors"
	"fmt"

	"github.com/carmandomx/acapp/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) FindById(Id int) (*models.User, error) {
	var user *models.User
	r := u.db.First(&user, Id)

	if r.Error != nil {
		return nil, r.Error
	}

	return user, nil
}

func (u *UserRepo) FindByEmail(e string) (*models.User, error) {
	var user *models.User
	if r := u.db.Where("username = ?", e).First(&user); r.Error != nil {
		return nil, r.Error
	}

	return user, nil
}

func (u *UserRepo) Create(user *models.User) error {

	r := u.db.Omit("PhotoUrl").Create(&user)

	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (u *UserRepo) Delete(Id string) error {
	r := u.db.Delete(&models.User{}, Id)
	fmt.Println(r.RowsAffected)
	fmt.Println(r.Error)
	if r.Error != nil || r.RowsAffected == 0 {
		return errors.New("Did not delete id:" + Id)
	}

	return nil
}
