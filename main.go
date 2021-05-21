package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/carmandomx/acapp/chat"
	"github.com/carmandomx/acapp/config"
	"github.com/carmandomx/acapp/controllers"
	"github.com/carmandomx/acapp/middleware"
	"github.com/carmandomx/acapp/models"
	"github.com/carmandomx/acapp/repositories"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	wsServer := chat.NewWSServer()
	go wsServer.Run()
	// room := chat.NewRoom()
	// room.Tracer = trace.New(os.Stdout)
	// go room.Run()
	db := config.ConnectDB()
	userRepo := repositories.NewUserRepo(db)
	userH := controllers.NewBaseHandler(userRepo)
	db.AutoMigrate(&models.User{})
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/login", userH.Login)
	r.POST("/users", userH.CreateUser)
	authorized := r.Group("/")
	authorized.Use(middleware.TokenAuthMiddleware())
	authorized.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authorized.DELETE("/users/:id", userH.DeleteUser)
	r.GET("/ws", func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Sec-WebSocket-Protocol")
		strArr := strings.Split(tokenString, ", ")
		decoded, err := url.QueryUnescape(strArr[1])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		token, err := jwt.Parse(decoded, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("ACCESS_SECRET")), nil
		})

		if err != nil {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		chat.ServeWS(wsServer, c.Writer, c.Request)
	})
	r.Run()
}
