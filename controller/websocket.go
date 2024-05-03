package controller

import (
	"chat_system/websocket"
	"net/http"

	"github.com/google/uuid"
	websocketPkg "github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

type WebSocketController struct {
	hub websocket.HubInterface
}

func StartWebSocketController(router *gin.Engine, hub websocket.HubInterface) {
	controller := &WebSocketController{
		hub: hub,
	}
	router.GET("/ws", controller.handleWebSocket)
}

func (c *WebSocketController) handleWebSocket(ctx *gin.Context) {
	upGrader := websocketPkg.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		panic(err)
	}

	websocket.NewClient(conn, func(client *websocket.Client) {
		client.WithID(uuid.New().String())
		client.WithContext(ctx)
		client.WithHub(c.hub)
	})
}
