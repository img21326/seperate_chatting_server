package ws

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/img21326/fb_chat/mock"
	"github.com/img21326/fb_chat/structure/room"
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

func TestUnRegister(t *testing.T) {
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
