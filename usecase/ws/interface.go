package ws

import (
	"context"

	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/structure/room"
	"github.com/img21326/fb_chat/ws/client"
)

type WebsocketUsecaseInterface interface {
	Run(ctx context.Context)
	Register(client *client.Client)
	UnRegister(client *client.Client)
	ReceiveMessage(message *pubmessage.PublishMessage)
	FindRoomByUserId(ctx context.Context, userID uint) (*room.Room, error)
}
