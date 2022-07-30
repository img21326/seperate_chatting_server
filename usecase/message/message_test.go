package message

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/img21326/fb_chat/mock"
	"github.com/img21326/fb_chat/structure/message"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/structure/room"
	"github.com/img21326/fb_chat/ws/client"
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
		UUID:    uuid.New(),
		Close:   false,
	}

	m := []*message.Message{
		&message.Message{
			RoomId:  r.UUID,
			UserId:  1,
			Message: "test",
		},
		&message.Message{
			RoomId:  r.UUID,
			UserId:  1,
			Message: "test",
		},
		&message.Message{
			RoomId:  r.UUID,
			UserId:  1,
			Message: "test",
		},
	}

	roomRepo.EXPECT().FindByUserId(gomock.Any(), uint(1)).Times(1).Return(&r, nil)
	messageRepo.EXPECT().LastsByRoomID(gomock.Any(), r.UUID, 30).Times(1).Return(m, nil)
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
		UUID:    uuid.New(),
		Close:   false,
	}

	var ts []time.Time
	for index, _ := range []int{1, 2, 3, 4, 5} {
		t, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("2022-01-01 0%v:0%v:0%v", index, index, index))
		ts = append(ts, t)
	}
	m := []*message.Message{
		&message.Message{
			RoomId:  r.UUID,
			UserId:  1,
			Message: "test",
			Time:    ts[0],
		},
		&message.Message{
			RoomId:  r.UUID,
			UserId:  1,
			Message: "test",
			Time:    ts[1],
		},
		&message.Message{
			RoomId:  r.UUID,
			UserId:  1,
			Message: "test",
			Time:    ts[2],
		},
	}

	roomRepo.EXPECT().FindByUserId(gomock.Any(), uint(1)).Times(1).Return(&r, nil)
	messageRepo.EXPECT().GetByID(gomock.Any(), uint(3)).Times(1).Return(m[2], nil)
	messageRepo.EXPECT().LastsByTime(gomock.Any(), r.UUID, m[2].Time, 3).Times(1).Return(m, nil)
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

func TestGetOnlineClients(t *testing.T) {
	c := gomock.NewController(t)
	localOnlineRepo := mock.NewMockLocalOnlineRepoInterface(c)

	messageUsecase := MessageUsecase{
		LocalOnlineRepo: localOnlineRepo,
	}

	c1 := &client.Client{}
	c1.User.ID = 1
	c2 := &client.Client{}
	c2.User.ID = 2

	localOnlineRepo.EXPECT().FindUserByID(gomock.Any()).Times(2).
		DoAndReturn(func(clientID uint) (*client.Client, error) {
			if clientID == uint(1) {
				return c1, nil
			} else {
				return c2, nil
			}
		})
	sender, receiver := messageUsecase.GetOnlineClients(uint(1), uint(2))
	assert.Equal(t, c1, sender)
	assert.Equal(t, c2, receiver)
}

func TestHandlePairSuccessMessage(t *testing.T) {
	MessageUsecase := &MessageUsecase{}

	receiver := client.Client{
		Send: make(chan []byte, 1),
	}
	roomID := uuid.New()
	mes := &pubmessage.PublishMessage{
		SendFrom: 2,
		Payload:  roomID.String(),
	}

	jsonMessage, _ := json.Marshal(mes)

	err := MessageUsecase.HandlePairSuccessMessage(&receiver, mes)
	assert.Nil(t, err)
	assert.Equal(t, receiver.RoomId, roomID)
	assert.Equal(t, receiver.PairId, uint(2))
	assert.Equal(t, jsonMessage, <-receiver.Send)
}

func TestHandleClientOnMessage(t *testing.T) {
	MessageUsecase := &MessageUsecase{}

	saveChan := make(chan *message.Message, 1)

	sender := client.Client{
		Send: make(chan []byte, 1),
	}
	sender.User.ID = 1

	receiver := client.Client{
		Send: make(chan []byte, 1),
	}
	roomID := uuid.New()

	message := &message.Message{
		ID:      1,
		RoomId:  roomID,
		UserId:  1,
		Message: "test",
		Time:    time.Now(),
	}
	pubMes := &pubmessage.PublishMessage{
		Type:     "message",
		SendFrom: 2,
		Payload:  message,
	}

	jsonMessage, _ := json.Marshal(pubMes)

	err := MessageUsecase.HandleClientOnMessage(&sender, &receiver, pubMes, saveChan)
	getMes := <-saveChan
	assert.Nil(t, err)
	assert.Equal(t, jsonMessage, <-receiver.Send)
	assert.Equal(t, jsonMessage, <-sender.Send)
	assert.Equal(t, message.ID, getMes.ID)
}

func TestHandleLeaveMessage(t *testing.T) {
	c := gomock.NewController(t)
	roomRepo := mock.NewMockRoomRepoInterface(c)
	MessageUsecase := &MessageUsecase{
		RoomRepo: roomRepo,
	}

	roomID := uuid.New()

	roomRepo.EXPECT().Close(gomock.Any(), roomID).Times(1)

	sender := &client.Client{
		Send:   make(chan []byte, 1),
		RoomId: roomID,
		PairId: 1,
	}
	sender.User.ID = 1

	receiver := &client.Client{
		Send:   make(chan []byte, 1),
		RoomId: roomID,
		PairId: 1,
	}
	receiver.User.ID = 2

	callCount := 0
	unRegisterFunc := func(ctx context.Context, client *client.Client) {
		callCount += 1
	}

	err := MessageUsecase.HandleLeaveMessage(sender, receiver, unRegisterFunc)

	assert.Nil(t, err)
	assert.Equal(t, receiver.RoomId, uuid.Nil)
	assert.Equal(t, sender.RoomId, uuid.Nil)
	assert.Equal(t, receiver.PairId, uint(0))
	assert.Equal(t, sender.PairId, uint(0))
	assert.Equal(t, callCount, 2)
}
