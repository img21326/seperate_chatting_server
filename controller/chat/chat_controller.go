package chat

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	errStruct "github.com/img21326/fb_chat/structure/error"
	"github.com/img21326/fb_chat/structure/user"
	"github.com/img21326/fb_chat/usecase/message"
	"github.com/img21326/fb_chat/usecase/ws"

	"gorm.io/gorm"
)

type ChatController struct {
	MessageUsecase message.MessageUsecaseInterface
	WsUsecase      ws.WebsocketUsecaseInterface
}

func NewChatController(e gin.IRoutes, messageUsecase message.MessageUsecaseInterface, wsUsecase ws.WebsocketUsecaseInterface) {
	controller := &ChatController{
		MessageUsecase: messageUsecase,
		WsUsecase:      wsUsecase,
	}

	e.GET("/inroom", controller.InRoom)
	e.GET("/history", controller.GetHistory)
}

func (c *ChatController) InRoom(ctx *gin.Context) {
	user := ctx.MustGet("user").(*user.User)
	_, err := c.WsUsecase.FindRoomByUserId(ctx, user.ID)
	if err != nil && err != gorm.ErrRecordNotFound && err != errStruct.RoomIsClose {
		log.Printf("find room error: %v", err)
		ctx.JSON(500, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}
	if err == gorm.ErrRecordNotFound || err == errStruct.RoomIsClose {
		ctx.JSON(200, gin.H{
			"status": false,
		})
		return
	}
	ctx.JSON(200, gin.H{
		"status": true,
	})
}

func (c *ChatController) GetHistory(ctx *gin.Context) {
	user := ctx.MustGet("user").(*user.User)
	lastMessageIdstr, ok := ctx.GetQuery("last_message_id")
	if !ok {
		messages, err := c.MessageUsecase.LastByUserID(ctx, user.ID, 20)
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": fmt.Sprintf("%v", err),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"messages": messages,
		})
	} else {
		lastMessageId, _ := strconv.Atoi(lastMessageIdstr)
		messages, err := c.MessageUsecase.LastByMessageID(ctx, user.ID, uint(lastMessageId), 20)
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": fmt.Sprintf("%v", err),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"messages": messages,
		})
	}
}
