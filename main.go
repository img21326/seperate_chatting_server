package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/go-redis/redis/v8"
	"github.com/img21326/fb_chat/controller"
	"github.com/img21326/fb_chat/helper"
	RepoLocal "github.com/img21326/fb_chat/repo/local_online"
	RepoMessage "github.com/img21326/fb_chat/repo/message"
	RepoOnline "github.com/img21326/fb_chat/repo/online"
	RepoPubSub "github.com/img21326/fb_chat/repo/pubsub"
	RepoRoom "github.com/img21326/fb_chat/repo/room"
	RepoUser "github.com/img21326/fb_chat/repo/user"
	RepoWait "github.com/img21326/fb_chat/repo/wait"
	ModelMessage "github.com/img21326/fb_chat/structure/message"
	ModelRoom "github.com/img21326/fb_chat/structure/room"
	ModelUser "github.com/img21326/fb_chat/structure/user"
	"github.com/img21326/fb_chat/usecase/auth"
	"github.com/img21326/fb_chat/usecase/message"
	"github.com/img21326/fb_chat/usecase/oauth"
	"github.com/img21326/fb_chat/usecase/pair"
	"github.com/img21326/fb_chat/usecase/sub"
	"github.com/img21326/fb_chat/usecase/ws"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Panicf("Open db error: %v", err)
	}
	db.AutoMigrate(&ModelUser.User{}, &ModelMessage.Message{}, &ModelRoom.Room{})
	return db
}

func initRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "139.162.125.28:6379",
		Password: "",
		DB:       5,
	})
	return client
}

func initRedisRepo(db *gorm.DB, redis *redis.Client) (messageRepo RepoMessage.MessageRepoInterface,
	localOnlineRepo RepoLocal.OnlineRepoInterface, onlineRepo RepoOnline.OnlineRepoInterface, roomRepo RepoRoom.RoomRepoInterface,
	userRepo RepoUser.UserRepoInterFace, waitRepo RepoWait.WaitRepoInterface, pubSubRepo RepoPubSub.PubSubRepoInterface) {
	messageRepo = RepoMessage.NewMessageRepo(db)
	localOnlineRepo = RepoLocal.NewOnlineRepo()
	onlineRepo = RepoOnline.NewOnlineRedisRepo(redis)
	roomRepo = RepoRoom.NewRoomRepo(db)
	userRepo = RepoUser.NewUserRepo(db)
	waitRepo = RepoWait.NewRedisWaitRepo(redis)
	pubSubRepo = RepoPubSub.NewPubSubRepo(redis)
	return
}

func main() {

	// if helper.GetEnv("HOST_NAME", "") == "" {
	// 	panic("Please set env HOST_NAME")
	// }

	db := initDB()
	redis := initRedis()

	messageRepo, localOnlineRepo, onlineRepo, roomRepo, userRepo, waitRepo, pubSubRepo := initRedisRepo(db, redis)

	//For AuthUsecase
	FacebookOauth := helper.NewFacebookOauth()
	FacebookUsecase := oauth.NewFacebookOauthUsecase(FacebookOauth)

	jwtConfig := auth.JwtConfig{
		Key:            []byte("secret168"),
		ExpireDuration: time.Hour * 24,
	}
	AuthUsecase := auth.NewAuthUsecase(jwtConfig, userRepo)

	// For Websocket
	wsUsecase := ws.NewRedisWebsocketUsecase(localOnlineRepo, onlineRepo, roomRepo)
	subUsecase := sub.NewRedisSubUsecase(pubSubRepo)
	pairUsecase := pair.NewRedisSubUsecase(waitRepo, onlineRepo, roomRepo)
	messageUsecase := message.NewMessageUsecase(messageRepo)

	// jwtMiddleware := jwt.NewJWTValidMiddleware(AuthUsecase)
	// jwtRoute := server.Group("/auth")
	// jwtRoute.Use(jwtMiddleware.ValidHeaderToken)

	server := gin.Default()
	controller.NewLoginController(server, FacebookUsecase, AuthUsecase)
	controller.NewWebsocketController(server, wsUsecase, subUsecase, pairUsecase, messageUsecase)

	port := os.Args[1]

	server.Run(fmt.Sprintf(":%v", port))
}
