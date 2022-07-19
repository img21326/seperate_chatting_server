package ws

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/img21326/fb_chat/mock"
	"github.com/img21326/fb_chat/structure/room"
	"github.com/img21326/fb_chat/ws/client"
	"github.com/stretchr/testify/assert"
)

func TestFindRoomByUserId(t *testing.T) {
	c := gomock.NewController(t)
	roomRepo := mock.NewMockRoomRepoInterface(c)
	r := room.Room{
		UserId1: 1,
		UserId2: 2,
		ID:      uuid.New(),
		Close:   false,
	}
	roomRepo.EXPECT().FindByUserId(gomock.Any(), uint(1)).Times(1).Return(&r, nil)

	wsUsecase := RedisWebsocketUsecase{
		RoomRepo: roomRepo,
	}
	ctx := context.Background()
	room, err := wsUsecase.FindRoomByUserId(ctx, uint(1))
	assert.Nil(t, err)
	assert.Equal(t, room, &r)
}

func TestRegister(t *testing.T) {
	c := gomock.NewController(t)

	localOnlineRepo := mock.NewMockLocalOnlineRepoInterface(c)
	onlineRepo := mock.NewMockOnlineRepoInterface(c)

	callTime := 0
	mockClient := &client.Client{}
	mockClient.User.ID = 1
	localOnlineRepo.EXPECT().Register(mockClient).DoAndReturn(
		func(c *client.Client) {
			callTime += 1
		})
	onlineRepo.EXPECT().Register(gomock.Any(), uint(1)).DoAndReturn(
		func(ctx context.Context, clientID uint) {
			callTime += 1
		})

	wsUsecase := RedisWebsocketUsecase{
		OnlineRepo:      onlineRepo,
		LocalOnlineRepo: localOnlineRepo,
	}

	ctx := context.Background()
	wsUsecase.Register(ctx, mockClient)
	assert.Equal(t, callTime, 2)
}

func TestUnRegister(t *testing.T) {
	c := gomock.NewController(t)

	localOnlineRepo := mock.NewMockLocalOnlineRepoInterface(c)
	onlineRepo := mock.NewMockOnlineRepoInterface(c)

	callTime := 0
	mockClient := &client.Client{}
	mockClient.User.ID = 1
	mockClient.CtxCancel = func() {
		callTime += 1
	}
	localOnlineRepo.EXPECT().UnRegister(mockClient).DoAndReturn(
		func(c *client.Client) {
			callTime += 1
		})
	onlineRepo.EXPECT().UnRegister(gomock.Any(), uint(1)).DoAndReturn(
		func(ctx context.Context, clientID uint) {
			callTime += 1
		})

	wsUsecase := RedisWebsocketUsecase{
		OnlineRepo:      onlineRepo,
		LocalOnlineRepo: localOnlineRepo,
	}

	ctx := context.Background()
	wsUsecase.UnRegister(ctx, mockClient)
	assert.Equal(t, callTime, 3)
}
