package message

import (
	"context"

	"github.com/img21326/fb_chat/structure/message"
)

type MessageUsecaseInterface interface {
	SetMessageChan(m chan *message.Message)
	Last(ctx context.Context, lastMessageID uint, c int) ([]*message.Message, error)
	Save(*message.Message)
	Run(ctx context.Context)
}
