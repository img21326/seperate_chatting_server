package message

import (
	"time"

	"github.com/google/uuid"
)

// Modal
type Message struct {
	ID      uint      `gorm:"primarykey"`
	RoomId  uuid.UUID `gorm:"Index"`
	UserId  uint
	Message string
	Time    time.Time
}
