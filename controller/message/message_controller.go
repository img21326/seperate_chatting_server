package message

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/img21326/fb_chat/structure/user"
	"github.com/img21326/fb_chat/usecase/message"
)

type MessageController struct {
	MessageUsecase message.MessageUsecaseInterface
}

func NewMessageController(e gin.IRoutes, messageUsecase message.MessageUsecaseInterface) {
	controller := &MessageController{
		MessageUsecase: messageUsecase,
	}

	e.GET("/history", controller.GetHistory)
}

func (c *MessageController) GetHistory(ctx *gin.Context) {
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
