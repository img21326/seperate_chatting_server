package room

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Modal
type Room struct {
	gorm.Model
	UUID    uuid.UUID `json:"uuid"`
	UserId1 uint      `json:"user_id1"`
	UserId2 uint      `json:"user_id2"`
	Close   bool      `json:"close"`
}
