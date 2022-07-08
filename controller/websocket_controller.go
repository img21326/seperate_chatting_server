package controller

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/img21326/fb_chat/helper"
	"github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/repo/user"
	"github.com/img21326/fb_chat/usecase/auth"
	hubUsecase "github.com/img21326/fb_chat/usecase/hub"
	"github.com/img21326/fb_chat/ws/client"
	"github.com/img21326/fb_chat/ws/hub"
	"gorm.io/gorm"
)

type WebsocketController struct {
	OnlineHub    *hub.OnlineHub
	PairHub      *hub.PairHub
	MessageQueue *hub.MessageQueue
	HubUsecase   hubUsecase.HubUsecaseInterface
	AuthUsecase  auth.AuthUsecaseInterFace
	WSUpgrader   websocket.Upgrader
}

func NewWebsocketController(e *gin.Engine, hubUsecase hubUsecase.HubUsecaseInterface, authUsecase auth.AuthUsecaseInterFace, redis *redis.Client) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	var messageQueue = hub.MessageQueue{
		SendMessage: make(chan *message.MessageModel, 4096),
		Close:       make(chan uuid.UUID, 1024),
		HubUsecase:  hubUsecase,
	}

	var onlineHub = hub.OnlineHub{
		Register:     make(chan *client.Client, 1024),
		Unregister:   make(chan *client.Client, 1024),
		ReceiveChan:  make(chan message.PublishMessage, 1024),
		PublishChan:  make(chan message.PublishMessage, 1024),
		MessageQueue: &messageQueue,
		HubUsecase:   hubUsecase,
	}

	var pairHub = hub.PairHub{
		Add:        make(chan *client.Client, 1024),
		Delete:     make(chan *client.Client, 1024),
		OnlineHub:  &onlineHub,
		HubUsecase: hubUsecase,
	}

	var subHub = hub.SubHub{
		OnlineHub: &onlineHub,
		Redis:     redis,
	}

	controller := &WebsocketController{
		OnlineHub:    &onlineHub,
		PairHub:      &pairHub,
		MessageQueue: &messageQueue,
		WSUpgrader:   upgrader,
		HubUsecase:   hubUsecase,
		AuthUsecase:  authUsecase,
	}
	ctx := context.Background()
	go subHub.MessageController(ctx)
	go controller.OnlineHub.Run()
	go controller.PairHub.Run()
	go controller.MessageQueue.Run()
	e.GET("/ws", controller.WS)
}

func (c *WebsocketController) WS(ctx *gin.Context) {
	token := ctx.Query("token")
	id, _ := strconv.Atoi(token)
	// user, err := c.AuthUsecase.VerifyToken(token)
	// if err != nil {
	// 	log.Printf("token error: %v", err)
	// 	return
	// }

	m := gorm.Model{
		ID: uint(id),
	}
	user := &user.UserModel{
		Model:  m,
		FbID:   helper.RandString(16),
		Name:   helper.RandString(5),
		Gender: "male",
	}
	log.Printf("new ws connection: %v", user.Name)
	room, err := c.HubUsecase.FindRoomByUserId(user.ID)
	if err != nil && err != gorm.ErrRecordNotFound && err.Error() != "RoomIsClosed" {
		log.Printf("find room error: %v", err)
		return
	}
	conn, err := c.WSUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("ws error: %v", err)
		return
	}
	client := client.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
		User: *user,
	}
	c.OnlineHub.Register <- &client
	if room != nil {
		log.Printf("new ws connection: %v in room %v", user.Name, room.ID)
		client.RoomId = room.ID
		if room.UserId1 == client.User.ID {
			client.PairId = room.UserId2
		} else {
			client.PairId = room.UserId1
		}
		client.Send <- []byte("{'type': 'inRoom'}")
	} else {
		log.Printf("new ws connection: %v new pairing", user.Name)
		want, ok := ctx.GetQuery("want")
		if !ok {
			log.Printf("ws not set want param")
			client.Conn.Close()
			return
		}
		client.WantToFind = want
		c.PairHub.Add <- &client
		client.Send <- []byte("{'type': 'paring'}")
	}

	go client.ReadPump(c.OnlineHub.PublishChan, c.OnlineHub.Unregister, c.PairHub.Delete)
	go client.WritePump()

}
