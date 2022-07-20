package ws

import (
	"context"

	localonline "github.com/img21326/fb_chat/repo/local_online"
	"github.com/img21326/fb_chat/repo/online"
	RepoRoom "github.com/img21326/fb_chat/repo/room"
	"github.com/img21326/fb_chat/structure/room"
	"github.com/img21326/fb_chat/ws/client"
)

type RedisWebsocketUsecase struct {
	LocalOnlineRepo localonline.LocalOnlineRepoInterface
	OnlineRepo      online.OnlineRepoInterface
	RoomRepo        RepoRoom.RoomRepoInterface
}

func NewRedisWebsocketUsecase(
	localOnlineRepo localonline.LocalOnlineRepoInterface, onlineRepo online.OnlineRepoInterface,
	roomRepo RepoRoom.RoomRepoInterface,
) WebsocketUsecaseInterface {
	return &RedisWebsocketUsecase{
		LocalOnlineRepo: localOnlineRepo,
		OnlineRepo:      onlineRepo,
		RoomRepo:        roomRepo,
	}
}

func (u *RedisWebsocketUsecase) FindRoomByUserId(ctx context.Context, userID uint) (*room.Room, error) {
	return u.RoomRepo.FindByUserId(ctx, userID)
}

func (u *RedisWebsocketUsecase) UnRegister(ctx context.Context, client *client.Client) {
	u.LocalOnlineRepo.UnRegister(client)
	u.OnlineRepo.UnRegister(ctx, client.User.ID)
	client.CtxCancel()
}

func (u *RedisWebsocketUsecase) Register(ctx context.Context, client *client.Client) {
	u.LocalOnlineRepo.Register(client)
	u.OnlineRepo.Register(ctx, client.User.ID)
}
