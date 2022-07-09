package hub

import (
	"context"

	"github.com/google/uuid"
	"github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/repo/room"
	"github.com/img21326/fb_chat/repo/user"
	"github.com/img21326/fb_chat/ws/client"
	"github.com/img21326/fb_chat/ws/messageType"
)

type HubUsecaseInterface interface {
	RegisterOnline(client *client.Client)
	UnRegisterOnline(client *client.Client)
	FindOnlineUserByUserID(userId uint) (*client.Client, error)

	FindUserByFbID(FbId string) (*user.UserModel, error)

	GetFirstQueueUser(client *client.Client) (*client.Client, error)
	AddUserToQueue(client *client.Client)
	DeleteuserFromQueue(client *client.Client)

	CreateRoom(*room.Room) error
	CloseRoom(uuid uuid.UUID) error
	FindRoomByUserId(userId uint) (*room.Room, error)

	SaveMesssage(context.Context, *message.MessageModel)
	SendMessage(context.Context, messageType.PublishMessage) error
}
