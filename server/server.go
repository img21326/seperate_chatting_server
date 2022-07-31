package server

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	ControllerLogin "github.com/img21326/fb_chat/controller/login"
	ControllerMessage "github.com/img21326/fb_chat/controller/message"
	ControllerWS "github.com/img21326/fb_chat/controller/ws"
	"github.com/img21326/fb_chat/hub"
	"github.com/img21326/fb_chat/middleware/jwt"
	RepoLocal "github.com/img21326/fb_chat/repo/local_online"
	RepoMessage "github.com/img21326/fb_chat/repo/message"
	RepoOnline "github.com/img21326/fb_chat/repo/online"
	RepoPubSub "github.com/img21326/fb_chat/repo/pubsub"
	RepoRoom "github.com/img21326/fb_chat/repo/room"
	RepoUser "github.com/img21326/fb_chat/repo/user"
	RepoWait "github.com/img21326/fb_chat/repo/wait"
	"github.com/img21326/fb_chat/usecase/auth"
	"github.com/img21326/fb_chat/usecase/message"
	"github.com/img21326/fb_chat/usecase/pair"
	"github.com/img21326/fb_chat/usecase/pubsub"
	"github.com/img21326/fb_chat/usecase/ws"
	"gorm.io/gorm"
)

func initRedisRepo(db *gorm.DB, redis *redis.Client) (messageRepo RepoMessage.MessageRepoInterface,
	localOnlineRepo RepoLocal.LocalOnlineRepoInterface, onlineRepo RepoOnline.OnlineRepoInterface, roomRepo RepoRoom.RoomRepoInterface,
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

func StartUpRedisServer(db *gorm.DB, redis *redis.Client, port string) {
	messageRepo, localOnlineRepo, onlineRepo, roomRepo, userRepo, waitRepo, pubSubRepo := initRedisRepo(db, redis)
	jwtConfig := auth.JwtConfig{
		Key:            []byte("secret168"),
		ExpireDuration: time.Hour * 24,
	}

	authUsecase := auth.NewAuthUsecase(jwtConfig, userRepo)
	wsUsecase := ws.NewRedisWebsocketUsecase(localOnlineRepo, onlineRepo, roomRepo)
	subUsecase := pubsub.NewRedisSubUsecase(pubSubRepo)
	pairUsecase := pair.NewRedisSubUsecase(waitRepo, onlineRepo, roomRepo)
	messageUsecase := message.NewMessageUsecase(messageRepo, roomRepo, localOnlineRepo)

	pubChan, queueChan := hub.StartHub(subUsecase, pairUsecase, messageUsecase, wsUsecase)

	server := gin.Default()
	jwtMiddleware := jwt.NewJWTValidMiddleware(authUsecase)
	jwtRoute := server.Group("/auth")
	jwtRoute.Use(jwtMiddleware.ValidHeaderToken)

	ControllerLogin.NewLoginController(server, authUsecase)
	ControllerMessage.NewMessageController(jwtRoute, messageUsecase)
	ControllerWS.NewWebsocketController(server, wsUsecase, authUsecase, pubChan, queueChan)

	server.Run(fmt.Sprintf(":%v", port))
}
