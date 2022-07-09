package message

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MessageModel struct {
	ID      uint `gorm:"primarykey"`
	RoomId  uuid.UUID
	UserId  uint
	Message string
	Time    time.Time
}

type MessageRepoInterface interface {
	Save(context.Context, *MessageModel)
	Publish(context.Context, []byte) error
}
