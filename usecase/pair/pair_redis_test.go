package pair

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/img21326/fb_chat/mock"
	errorStruct "github.com/img21326/fb_chat/structure/error"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/structure/room"
	"github.com/img21326/fb_chat/ws/client"
	"github.com/stretchr/testify/assert"
)

func TestTryToPairWithLenSmallerThan1(t *testing.T) {
	c := gomock.NewController(t)
	waitRepo := mock.NewMockWaitRepoInterface(c)

	waitRepo.EXPECT().Len(gomock.Any(), "female_male").Times(1).Return(0)

	pairUsecase := RedisPairUsecase{
		WaitRepo: waitRepo,
	}
	ctx := context.Background()
	client := client.Client{}
	client.User.ID = 1
	client.User.Gender = "male"
	client.WantToFind = "female"
	room, err := pairUsecase.TryToPair(ctx, &client)
	assert.Nil(t, room)
	assert.Equal(t, err, errorStruct.QueueSmallerThan1)
}

func TestTryToPairWithSuccess(t *testing.T) {
	c := gomock.NewController(t)
	waitRepo := mock.NewMockWaitRepoInterface(c)
	onlineRepo := mock.NewMockOnlineRepoInterface(c)

	waitRepo.EXPECT().Len(gomock.Any(), "female_male").Times(1).Return(1)
	waitRepo.EXPECT().Pop(gomock.Any(), "female_male").Times(1).Return(uint(2), nil)
	onlineRepo.EXPECT().CheckUserOnline(gomock.Any(), uint(2)).Times(1).Return(true)

	pairUsecase := RedisPairUsecase{
		WaitRepo:   waitRepo,
		OnlineRepo: onlineRepo,
	}

	ctx := context.Background()
	client := client.Client{}
	client.User.ID = 1
	client.User.Gender = "male"
	client.WantToFind = "female"
	room, err := pairUsecase.TryToPair(ctx, &client)
	assert.Nil(t, err)
	assert.NotNil(t, room)
}

func TestTryToPairWithUserNotOnline(t *testing.T) {
	c := gomock.NewController(t)
	waitRepo := mock.NewMockWaitRepoInterface(c)
	onlineRepo := mock.NewMockOnlineRepoInterface(c)

	waitRepo.EXPECT().Len(gomock.Any(), "female_male").AnyTimes().Return(1)
	pairID := uint(0)
	waitRepo.EXPECT().Pop(gomock.Any(), "female_male").AnyTimes().
		DoAndReturn(func(ctx context.Context, queuename string) (uint, error) {
			pairID += 1
			return pairID, nil
		})
	onlineRepo.EXPECT().CheckUserOnline(gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(
			func(ctx context.Context, pairID uint) bool {
				if pairID == 3 {
					return false
				} else {
					return true
				}
			})

	pairUsecase := RedisPairUsecase{
		WaitRepo:   waitRepo,
		OnlineRepo: onlineRepo,
	}

	ctx := context.Background()
	client := client.Client{}
	client.User.ID = 1
	client.User.Gender = "male"
	client.WantToFind = "female"
	room, err := pairUsecase.TryToPair(ctx, &client)
	assert.Nil(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, room.UserId2, pairID)
}

func TestAddToQueue(t *testing.T) {
	c := gomock.NewController(t)
	waitRepo := mock.NewMockWaitRepoInterface(c)

	waitRepo.EXPECT().Add(gomock.Any(), "male_female", uint(1)).Times(1)
	pairUsecase := RedisPairUsecase{
		WaitRepo: waitRepo,
	}

	ctx := context.Background()
	client := client.Client{}
	client.User.ID = 1
	client.User.Gender = "male"
	client.WantToFind = "female"
	pairUsecase.AddToQueue(ctx, &client)
}

func TestPairSuccess(t *testing.T) {
	c := gomock.NewController(t)
	roomRepo := mock.NewMockRoomRepoInterface(c)

	room := room.Room{
		ID:      uuid.New(),
		UserId1: 1,
		UserId2: 2,
	}
	roomRepo.EXPECT().Create(gomock.Any(), &room).Times(1).Return(nil)

	pairUsecase := RedisPairUsecase{
		RoomRepo: roomRepo,
	}

	ctx := context.Background()

	m1, m2, err := pairUsecase.PairSuccess(ctx, &room)
	assert.Nil(t, err)
	assert.Equal(t, m1, &pubmessage.PublishMessage{
		Type:     "pairSuccess",
		SendFrom: room.UserId1,
		SendTo:   room.UserId2,
		Payload:  room.ID,
	})
	assert.Equal(t, m2, &pubmessage.PublishMessage{
		Type:     "pairSuccess",
		SendFrom: room.UserId2,
		SendTo:   room.UserId1,
		Payload:  room.ID,
	})
}
