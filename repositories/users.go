package repositories

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/carmandomx/acapp/models"
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
	r := u.db.Delete(&User{}, Id)
	fmt.Println(r.RowsAffected)
	fmt.Println(r.Error)
	if r.Error != nil || r.RowsAffected == 0 {
		return errors.New("Did not delete id:" + Id)
	}

	return nil
}

func (user User) GetId() string {
	return strconv.FormatUint(uint64(user.ID), 10)
}

func (user User) GetName() string {
	return user.Name
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.Password, err = HashPassword(u.Password)

	if err != nil {
		return err
	}

	return nil

}

func (u *UserRepo) GetAllUsers() ([]models.User, error) {
	var users []models.User

	res := u.db.Find(&users)

	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
