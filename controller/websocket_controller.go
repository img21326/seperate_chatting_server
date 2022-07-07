package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/img21326/fb_chat/entity/ws"
	"github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/repo/room"
	"github.com/img21326/fb_chat/usecase/pair"
)

type WebsocketController struct {
	OnlineHub    *OnlineHub
	PairHub      *PairHub
	MessageQueue *MessageQueue
	PairUsecase  pair.PairUsecaseInterface
	WSUpgrader   websocket.Upgrader
}

func NewWebsocketController(e *gin.Engine, pairUsecase pair.PairUsecaseInterface) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	var onlineHub = OnlineHub{
		Register:   make(chan *ws.Client),
		Unregister: make(chan *ws.Client),
	}

	var pairHub = PairHub{
		Add:    make(chan *ws.Client),
		Delete: make(chan *ws.Client),
	}

	var messageQueue = MessageQueue{
		SaveMessage: make(chan *message.MessageModel),
		Close:       make(chan uuid.UUID),
	}

	controller := &WebsocketController{
		OnlineHub:    &onlineHub,
		PairHub:      &pairHub,
		MessageQueue: &messageQueue,
		WSUpgrader:   upgrader,
		PairUsecase:  pairUsecase,
	}
	go controller.OnlineHub.run()
	go controller.PairHub.run()
	go controller.MessageQueue.run()
	e.GET("/ws", controller.WS)
}

func (c *WebsocketController) WS(ctx *gin.Context) {
	conn, _ := c.WSUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	fb_id := ctx.MustGet("fb_id").(string)
	user, err := c.PairUsecase.FindUserByFbID(fb_id)
	if err != nil {
		log.Printf("connect ws can't find user: %v", err)
		return
	}

	client := ws.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
		User: *user,
	}
	room, _ := c.PairUsecase.FindRoomByUserId(user.ID)
	if (room != nil) && (!room.Close) {
		client.RoomId = room.ID
		var pairId uint
		if room.UserId1 == user.ID {
			pairId = room.UserId2
		} else {
			pairId = room.UserId1
		}
		pairClient, err := c.PairUsecase.FindOnlineUserByFbID(pairId)
		if err == nil {
			client.PairClient = pairClient
			pairClient.PairClient = &client
		}
	} else {
		c.PairHub.Add <- &client
		client.Send <- []byte("pairing")
	}

	go client.ReadPump(c.MessageQueue.SaveMessage, c.MessageQueue.Close, c.PairHub.Delete)
	go client.WritePump()

}

type OnlineHub struct {
	Register    chan *ws.Client
	Unregister  chan *ws.Client
	PairUsecase pair.PairUsecaseInterface
}

func (h *OnlineHub) run() {
	for {
		select {
		case client := <-h.Register:
			h.PairUsecase.RegisterOnline(client)
		case client := <-h.Unregister:
			h.PairUsecase.UnRegisterOnline(client)
		}
	}
}

type MessageQueue struct {
	SaveMessage chan *message.MessageModel
	Close       chan uuid.UUID
	PairUsecase pair.PairUsecaseInterface
}

func (q *MessageQueue) run() {
	for {
		select {
		case roomId := <-q.Close:
			q.PairUsecase.CloseRoom(roomId)
		case message := <-q.SaveMessage:
			q.PairUsecase.SaveMessage(message)
		}
	}
}

type PairHub struct {
	Add         chan *ws.Client
	Delete      chan *ws.Client
	PairUsecase pair.PairUsecaseInterface
}

func (h *PairHub) run() {
	for {
		select {
		case client := <-h.Add:
			pairClient, err := h.PairUsecase.GetFirstQueueUser(client.User.Gender)
			if err != nil {
				h.PairUsecase.AddUserToQueue(client)
				return
			}

			room := &room.Room{
				UserId1: client.User.ID,
				UserId2: client.PairClient.User.ID,
				Close:   false,
			}
			err = h.PairUsecase.CreateRoom(room)
			if err != nil {
				log.Printf("create chat room err: %v", err)
				client.Send <- []byte("pairError")
				pairClient.Send <- []byte("pairError")
			}
			client.PairClient = pairClient
			client.RoomId = room.ID
			pairClient.PairClient = client
			pairClient.RoomId = room.ID

			client.Send <- []byte("connect")
			pairClient.Send <- []byte("connect")
		case client := <-h.Delete:
			h.PairUsecase.DeleteuserFromQueue(client)
		}
	}
}
