package message

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageModel struct {
	gorm.Model
	RoomId  uuid.UUID
	UserId  uint
	Message string
}

type MessageRepoInterface interface {
	Save(*MessageModel)
}
