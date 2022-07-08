package message

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PublishMessage struct {
	Type     string      `json:"type"`
	SendFrom uint        `json:"sendFrom"`
	SendTo   uint        `json:"sendTo"`
	Payload  interface{} `json:"payload"`
}

type SendToUserMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type MessageModel struct {
	gorm.Model
	RoomId  uuid.UUID
	UserId  uint
	Message string
	Time    time.Time
}

type MessageRepoInterface interface {
	Save(context.Context, *MessageModel)
	Publish(context.Context, []byte) error
}
