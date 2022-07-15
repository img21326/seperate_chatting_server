package message

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/img21326/fb_chat/structure/message"
)

type MessageRepoInterface interface {
	Save(context.Context, *message.Message)
	GetByID(ctx context.Context, ID uint) (*message.Message, error)
	LastsByRoomID(ctx context.Context, roomID uuid.UUID, c int) ([]*message.Message, error)
	LastsByTime(ctx context.Context, roomID uuid.UUID, t time.Time, c int) ([]*message.Message, error)
}
