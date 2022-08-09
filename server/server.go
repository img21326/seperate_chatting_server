package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	ControllerChat "github.com/img21326/fb_chat/controller/chat"
	ControllerLogin "github.com/img21326/fb_chat/controller/login"
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
	ctx, cancel := context.WithCancel(context.Background())
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

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

	pubChan, queueChan := hub.StartHub(ctx, subUsecase, pairUsecase, messageUsecase, wsUsecase)

	handler := gin.Default()
	jwtMiddleware := jwt.NewJWTValidMiddleware(authUsecase)
	chatRoute := handler.Group("/chat")
	chatRoute.Use(jwtMiddleware.ValidHeaderToken)

	ControllerLogin.NewLoginController(handler, authUsecase)
	ControllerChat.NewChatController(chatRoute, messageUsecase, wsUsecase)
	ControllerWS.NewWebsocketController(handler, wsUsecase, authUsecase, pubChan, queueChan)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-shutdownChan
	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	<-ctx.Done()
	log.Println("timeout of 10 seconds.")
	log.Println("Server exiting")
}
