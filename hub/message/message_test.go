package message

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/img21326/fb_chat/mock"
	"github.com/img21326/fb_chat/structure/message"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/ws/client"
)

func TestSaveMessage(t *testing.T) {
	c := gomock.NewController(t)

	mes := &message.Message{ID: 1}
	ctx, cancel := context.WithCancel(context.Background())

	messageUsecase := mock.NewMockMessageUsecaseInterface(c)
	messageUsecase.EXPECT().Save(gomock.Any(), mes).Times(1).Do(
		func(ctx context.Context, message *message.Message) {
			cancel()
		})

	mesHub := MessageHub{
		SaveMessageChan: make(chan *message.Message, 1),
		MessageUsecase:  messageUsecase,
	}
	mesHub.SaveMessageChan <- mes

	mesHub.Run(ctx)
}

func TestReceivePairSuccess(t *testing.T) {
	c := gomock.NewController(t)

	mes1 := &pubmessage.PublishMessage{Type: "pairSuccess", SendTo: 1}
	mes2 := &pubmessage.PublishMessage{Type: "pairSuccess", SendTo: 2}
	ctx, cancel := context.WithCancel(context.Background())

	messageUsecase := mock.NewMockMessageUsecaseInterface(c)

	sender := client.Client{}
	sender.User.ID = 1
	receiver := client.Client{}
	receiver.User.ID = 2
	messageUsecase.EXPECT().GetOnlineClients(gomock.Any(), gomock.Any()).Times(2).DoAndReturn(
		func(senderID uint, receiverID uint) (*client.Client, *client.Client) {
			if receiverID == uint(1) {
				return nil, &sender
			} else {
				cancel()
				return nil, &receiver
			}
		})
	messageUsecase.EXPECT().HandlePairSuccessMessage(gomock.Any(), gomock.Any()).Times(2).DoAndReturn(
		func(c *client.Client, m *pubmessage.PublishMessage) error {
			if c.User.ID == 1 || c.User.ID == 2 {
				return nil
			}
			return errors.New("OutOfUserID")
		})

	mesHub := MessageHub{
		ReceiveMessageChan: make(chan *pubmessage.PublishMessage, 2),
		MessageUsecase:     messageUsecase,
	}
	mesHub.ReceiveMessageChan <- mes1
	mesHub.ReceiveMessageChan <- mes2

	mesHub.Run(ctx)
}

func TestClientOnMessage(t *testing.T) {
	c := gomock.NewController(t)

	mes1 := &pubmessage.PublishMessage{Type: "message", SendFrom: 1, SendTo: 2}

	ctx, cancel := context.WithCancel(context.Background())

	messageUsecase := mock.NewMockMessageUsecaseInterface(c)

	sender := client.Client{}
	sender.User.ID = 1
	receiver := client.Client{}
	receiver.User.ID = 2
	messageUsecase.EXPECT().GetOnlineClients(gomock.Any(), gomock.Any()).Times(1).DoAndReturn(
		func(senderID uint, receiverID uint) (*client.Client, *client.Client) {
			return &sender, &receiver
		})
	messageUsecase.EXPECT().HandleClientOnMessage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).DoAndReturn(
		func(getSender *client.Client, getReceiver *client.Client, receiveMessage *pubmessage.PublishMessage, saveMessageChan chan *message.Message) error {
			assert.Equal(t, sender, getSender)
			assert.Equal(t, receiver, getReceiver)
			assert.Equal(t, receiveMessage, mes1)
			cancel()
			return nil
		})

	mesHub := MessageHub{
		ReceiveMessageChan: make(chan *pubmessage.PublishMessage, 1),
		SaveMessageChan:    make(chan *message.Message, 1),
		MessageUsecase:     messageUsecase,
	}
	mesHub.ReceiveMessageChan <- mes1

	mesHub.Run(ctx)
}
