package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketController struct {
	WSUpgrader websocket.Upgrader
}

func NewWebsocketController(e *gin.Engine) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	controller := WebsocketController{
		WSUpgrader: upgrader,
	}

	fmt.Printf("%+v", controller)

	// e.GET("/ws", controller.WS)
}

// func (c *WebsocketController) WS(ctx *gin.Context) {
// 	token := ctx.Query("token")
// 	id, _ := strconv.Atoi(token)
// 	// user, err := c.AuthUsecase.VerifyToken(token)
// 	// if err != nil {
// 	// 	log.Printf("token error: %v", err)
// 	// 	return
// 	// }

// 	m := gorm.Model{
// 		ID: uint(id),
// 	}
// 	user := &user.UserModel{
// 		Model:  m,
// 		FbID:   helper.RandString(16),
// 		Name:   helper.RandString(5),
// 		Gender: "male",
// 	}
// 	log.Printf("new ws connection: %v", user.Name)
// 	room, err := c.HubUsecase.FindRoomByUserId(user.ID)
// 	if err != nil && err != gorm.ErrRecordNotFound && err.Error() != "RoomIsClosed" {
// 		log.Printf("find room error: %v", err)
// 		return
// 	}
// 	conn, err := c.WSUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
// 	if err != nil {
// 		log.Printf("ws error: %v", err)
// 		return
// 	}
// 	client := client.Client{
// 		Conn: conn,
// 		Send: make(chan []byte, 256),
// 		User: *user,
// 	}
// 	c.OnlineHub.Register <- &client
// 	if room != nil {
// 		log.Printf("new ws connection: %v in room %v", user.Name, room.ID)
// 		client.RoomId = room.ID
// 		if room.UserId1 == client.User.ID {
// 			client.PairId = room.UserId2
// 		} else {
// 			client.PairId = room.UserId1
// 		}
// 		client.Send <- []byte("{'type': 'inRoom'}")
// 	} else {
// 		log.Printf("new ws connection: %v with new pairing", user.Name)
// 		want, ok := ctx.GetQuery("want")
// 		if !ok {
// 			log.Printf("ws not set want param")
// 			client.Conn.Close()
// 			return
// 		}
// 		client.WantToFind = want
// 		c.PairHub.Add <- &client
// 		client.Send <- []byte("{'type': 'paring'}")
// 	}

// 	go client.ReadPump(c.OnlineHub.PublishChan, c.OnlineHub.Unregister, c.PairHub.Delete)
// 	go client.WritePump()

// }
