package ws

import (
	"context"

	"github.com/img21326/fb_chat/structure/room"
	"github.com/img21326/fb_chat/ws/client"
)

type WebsocketUsecaseInterface interface {
	Register(ctx context.Context, client *client.Client)
	UnRegister(ctx context.Context, client *client.Client)
	FindRoomByUserId(ctx context.Context, userID uint) (*room.Room, error)
}
