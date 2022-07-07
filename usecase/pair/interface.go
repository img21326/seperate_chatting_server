package pair

import (
	"github.com/google/uuid"
	"github.com/img21326/fb_chat/entity/ws"
	"github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/repo/room"
	"github.com/img21326/fb_chat/repo/user"
)

type PairUsecaseInterface interface {
	RegisterOnline(client *ws.Client)
	UnRegisterOnline(client *ws.Client)
	FindOnlineUserByFbID(userId uint) (*ws.Client, error)

	FindUserByFbID(FbId string) (*user.UserModel, error)

	GetFirstQueueUser(gender string) (*ws.Client, error)
	AddUserToQueue(client *ws.Client)
	DeleteuserFromQueue(client *ws.Client)

	CreateRoom(*room.Room) (uuid.UUID, error)
	CloseRoom(uuid uuid.UUID) error
	FindRoomByUserId(userId uint) (*room.Room, error)

	SaveMessage(*message.MessageModel)
}
