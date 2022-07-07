package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/img21326/fb_chat/entity/ws"
	"github.com/img21326/fb_chat/helper"
	"github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/repo/room"
	"github.com/img21326/fb_chat/repo/user"
	"github.com/img21326/fb_chat/usecase/auth"
	"github.com/img21326/fb_chat/usecase/pair"
)

type WebsocketController struct {
	OnlineHub    *OnlineHub
	PairHub      *PairHub
	MessageQueue *MessageQueue
	PairUsecase  pair.PairUsecaseInterface
	AuthUsecase  auth.AuthUsecaseInterFace
	WSUpgrader   websocket.Upgrader
}

func NewWebsocketController(e *gin.Engine, pairUsecase pair.PairUsecaseInterface, authUsecase auth.AuthUsecaseInterFace) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	var onlineHub = OnlineHub{
		Register:    make(chan *ws.Client),
		Unregister:  make(chan *ws.Client),
		PairUsecase: pairUsecase,
	}

	var pairHub = PairHub{
		Add:         make(chan *ws.Client),
		Delete:      make(chan *ws.Client),
		PairUsecase: pairUsecase,
	}

	var messageQueue = MessageQueue{
		SaveMessage: make(chan *message.MessageModel),
		Close:       make(chan uuid.UUID),
		PairUsecase: pairUsecase,
	}

	controller := &WebsocketController{
		OnlineHub:    &onlineHub,
		PairHub:      &pairHub,
		MessageQueue: &messageQueue,
		WSUpgrader:   upgrader,
		PairUsecase:  pairUsecase,
		AuthUsecase:  authUsecase,
	}

	go controller.OnlineHub.run()
	go controller.PairHub.run()
	go controller.MessageQueue.run()
	e.GET("/ws", controller.WS)
}

func (c *WebsocketController) WS(ctx *gin.Context) {
	// token := ctx.Query("token")
	// user, err := c.AuthUsecase.VerifyToken(token)
	// if err != nil {
	// 	log.Printf("token error: %v", err)
	// 	return
	// }
	user := &user.UserModel{
		FbID:   helper.RandString(16),
		Name:   helper.RandString(5),
		Gender: "male",
	}
	room, err := c.PairUsecase.FindRoomByUserId(user.ID)
	if err != nil {
		log.Printf("find room error: %v", err)
		return
	}
	conn, err := c.WSUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("ws error: %v", err)
		return
	}
	client := ws.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
		User: *user,
	}

	log.Print("A")
	if (room.ID != uuid.Nil) && (!room.Close) {
		log.Print("B")
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
		log.Print("C")
		want, ok := ctx.GetQuery("want")
		if !ok {
			log.Printf("ws not set want param")
			client.Conn.Close()
			return
		}
		client.WantToFind = want
		c.PairHub.Add <- &client
		client.Send <- []byte("pairing")
	}
	log.Print("D")

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
			// 先試著配對看看
			pairClient, err := h.PairUsecase.GetFirstQueueUser(client)
			if err != nil {
				//如果配對失敗 就加入等待中
				h.PairUsecase.AddUserToQueue(client)
			} else {
				// 以下為配對成功所做的事
				room := &room.Room{
					UserId1: client.User.ID,
					UserId2: pairClient.User.ID,
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
			}
		case client := <-h.Delete:
			h.PairUsecase.DeleteuserFromQueue(client)
		}
	}
}
