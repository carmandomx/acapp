package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/carmandomx/acapp/auth"
	"github.com/carmandomx/acapp/models"
	"github.com/carmandomx/acapp/repositories"
	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
	userRepo models.UserRepository
}

func NewBaseHandler(userRepo models.UserRepository) *BaseHandler {
	return &BaseHandler{
		userRepo: userRepo,
	}
}

func (h *BaseHandler) CreateUser(c *gin.Context) {
	var u models.User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "some error",
		})
		return
	}
	_, err = h.userRepo.FindByEmail(u.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already exists",
		})
		return
	}

	err = h.userRepo.Create(&u)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "some error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": u.ID,
	})
}

func (h *BaseHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := h.userRepo.Delete(id)

	if err == nil {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "could not delete id:" + id,
	})

}

func (h *BaseHandler) Login(c *gin.Context) {
	var loginData models.Login

	err := c.ShouldBindJSON(&loginData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	u, err := h.userRepo.FindByEmail(loginData.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
	}

	ok := repositories.CheckPasswordHash(loginData.Password, u.Password)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	td, err := auth.NewTokenService().CreateToken(strconv.FormatUint(uint64(u.ID), 10), u.Username)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
	}
	tokens := map[string]string{
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
		"userId":        strconv.FormatUint(uint64(u.ID), 10),
	}

	c.JSON(http.StatusOK, tokens)
}
