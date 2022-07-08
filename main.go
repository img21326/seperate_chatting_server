package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/go-redis/redis/v8"
	"github.com/img21326/fb_chat/controller"
	"github.com/img21326/fb_chat/helper"
	"github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/repo/online"
	"github.com/img21326/fb_chat/repo/room"
	"github.com/img21326/fb_chat/repo/user"
	"github.com/img21326/fb_chat/repo/wait"
	"github.com/img21326/fb_chat/usecase/auth"
	"github.com/img21326/fb_chat/usecase/hub"
	"github.com/img21326/fb_chat/usecase/oauth"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Panicf("Open db error: %v", err)
	}
	db.AutoMigrate(&user.UserModel{}, &message.MessageModel{}, &room.Room{})
	return db
}

func initRedis() *redis.Client { // 實體化redis.Client 並返回實體的位址
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}

func main() {

	db := initDB()
	redis := initRedis()

	messageRepo := message.NewMessageRepo(db, redis)
	onlineRepo := online.NewOnlineRepo()
	roomRepo := room.NewRoomRepo(db)
	userRepo := user.NewUserRepo(db)
	waitRepo := wait.NewWaitRepo()

	server := gin.Default()

	FacebookOauth := helper.NewFacebookOauth()
	FacebookUsecase := oauth.NewFacebookOauthUsecase(FacebookOauth)

	Hubsecase := hub.NewHubUsecase(userRepo, messageRepo, onlineRepo, roomRepo, waitRepo)

	jwtConfig := auth.JwtConfig{
		Key:            []byte("secret168"),
		ExpireDuration: time.Hour * 24,
	}
	AuthUsecase := auth.NewAuthUsecase(jwtConfig, userRepo)
	// jwtMiddleware := jwt.NewJWTValidMiddleware(AuthUsecase)

	// jwtRoute := server.Group("/auth")
	// jwtRoute.Use(jwtMiddleware.ValidHeaderToken)

	controller.NewLoginController(server, FacebookUsecase, AuthUsecase)
	controller.NewWebsocketController(server, Hubsecase, AuthUsecase, redis)

	server.Run(":8081")
}
