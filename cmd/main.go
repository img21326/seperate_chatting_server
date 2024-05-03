package main

import (
	"chat_system/consumer/handler"
	"chat_system/controller"
	"chat_system/websocket"
	"chat_system/websocket/pubsub"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	chatsystem "chat_system"
	consumerPkg "chat_system/consumer"

	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	redis := redis.NewClient(&redis.Options{
		Addr:     chatsystem.REDIS_ADDR,
		Password: "",
		DB:       3,
	})

	err := redis.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	hub := websocket.NewHub(ctx, func(h *websocket.Hub) {
		h.SetPubSub(pubsub.NewRedisPubSub(ctx, redis))
	})

	router := gin.Default()

	controller.StartWebSocketController(router, hub)

	baseHandler := handler.NewBaseHandler(hub, redis)
	joinRoomHandler := handler.NewJoinRoomHandler(baseHandler, redis)
	leaveRoomHandler := handler.NewLeaveRoomHandler(baseHandler, redis)
	onMessageHandler := handler.NewOnMessageHandler(baseHandler, hub)
	consumer := consumerPkg.NewPairingConsumer(joinRoomHandler, leaveRoomHandler, onMessageHandler)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg := hub.ConsumeMsg()
				err := consumer.ConsumeMsg(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

	srv := &http.Server{
		Addr:    os.Getenv("PORT"),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	<-ctx.Done()
	log.Println("Server exiting")
}
