package message

import (
	"context"

	"github.com/img21326/fb_chat/structure/message"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/ws/client"
)

type MessageUsecaseInterface interface {
	// For API
	LastByUserID(ctx context.Context, userID uint, c int) (messages []*message.Message, err error) // 從使用者找房間 最後聊天紀錄
	LastByMessageID(ctx context.Context, userID uint, lastMessageID uint, c int) ([]*message.Message, error)
	// For API End

	// For Hub
	Save(ctx context.Context, message *message.Message)
	GetOnlineClients(senderID uint, receiverID uint) (sender *client.Client, receiver *client.Client)
	HandlePairSuccessMessage(receiver *client.Client, receiveMessage *pubmessage.PublishMessage) error
	HandleClientOnMessage(sender *client.Client, receiver *client.Client, receiveMessage *pubmessage.PublishMessage, saveMessageChan chan *message.Message) error
	HandleLeaveMessage(sender *client.Client, receiver *client.Client, unRegisterFunc func(ctx context.Context, client *client.Client)) error
	// For Hub End
}
