package message

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/img21326/fb_chat/mock"
	"github.com/img21326/fb_chat/structure/message"
	"github.com/img21326/fb_chat/structure/room"
	"github.com/stretchr/testify/assert"
)

func TestLastByUserID(t *testing.T) {
	c := gomock.NewController(t)
	roomRepo := mock.NewMockRoomRepoInterface(c)
	messageRepo := mock.NewMockMessageRepoInterface(c)

	messageUsecase := MessageUsecase{
		RoomRepo:    roomRepo,
		MessageRepo: messageRepo,
	}

	r := room.Room{
		UserId1: 1,
		UserId2: 2,
		ID:      uuid.New(),
		Close:   false,
	}

	m := []*message.Message{
		&message.Message{
			RoomId:  r.ID,
			UserId:  1,
			Message: "test",
		},
		&message.Message{
			RoomId:  r.ID,
			UserId:  1,
			Message: "test",
		},
		&message.Message{
			RoomId:  r.ID,
			UserId:  1,
			Message: "test",
		},
	}

	roomRepo.EXPECT().FindByUserId(gomock.Any(), uint(1)).Times(1).Return(&r, nil)
	messageRepo.EXPECT().LastsByRoomID(gomock.Any(), r.ID, 30).Times(1).Return(m, nil)
	ctx := context.Background()
	message, err := messageUsecase.LastByUserID(ctx, 1, 30)

	assert.Nil(t, err)
	assert.Equal(t, message, m)
}

func TestLastByMessageID(t *testing.T) {
	c := gomock.NewController(t)
	roomRepo := mock.NewMockRoomRepoInterface(c)
	messageRepo := mock.NewMockMessageRepoInterface(c)

	messageUsecase := MessageUsecase{
		RoomRepo:    roomRepo,
		MessageRepo: messageRepo,
	}

	r := room.Room{
		UserId1: 1,
		UserId2: 2,
		ID:      uuid.New(),
		Close:   false,
	}

	var ts []time.Time
	for index, _ := range []int{1, 2, 3, 4, 5} {
		t, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("2022-01-01 0%v:0%v:0%v", index, index, index))
		ts = append(ts, t)
	}
	m := []*message.Message{
		&message.Message{
			RoomId:  r.ID,
			UserId:  1,
			Message: "test",
			Time:    ts[0],
		},
		&message.Message{
			RoomId:  r.ID,
			UserId:  1,
			Message: "test",
			Time:    ts[1],
		},
		&message.Message{
			RoomId:  r.ID,
			UserId:  1,
			Message: "test",
			Time:    ts[2],
		},
	}

	roomRepo.EXPECT().FindByUserId(gomock.Any(), uint(1)).Times(1).Return(&r, nil)
	messageRepo.EXPECT().GetByID(gomock.Any(), uint(3)).Times(1).Return(m[2], nil)
	messageRepo.EXPECT().LastsByTime(gomock.Any(), r.ID, m[2].Time, 3).Times(1).Return(m, nil)
	ctx := context.Background()
	message, err := messageUsecase.LastByMessageID(ctx, 1, 3, 3)

	assert.Nil(t, err)
	assert.Equal(t, message, m)
}

func TestSave(t *testing.T) {
	c := gomock.NewController(t)
	messageRepo := mock.NewMockMessageRepoInterface(c)

	messageUsecase := MessageUsecase{
		MessageRepo: messageRepo,
	}

	m := message.Message{
		RoomId:  uuid.New(),
		UserId:  1,
		Message: "test",
	}

	messageRepo.EXPECT().Save(gomock.Any(), &m).Times(1)
	ctx := context.Background()
	messageUsecase.Save(ctx, &m)
}
