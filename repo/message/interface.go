package message

import (
	"context"
	"time"

	"github.com/img21326/fb_chat/structure/message"
)

type MessageRepoInterface interface {
	Save(context.Context, *message.Message)
	GetByID(ctx context.Context, ID uint) (*message.Message, error)
	LastsByTime(ctx context.Context, t time.Time, c int) ([]*message.Message, error)
}
