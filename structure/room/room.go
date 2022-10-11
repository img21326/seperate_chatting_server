package room

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Modal
type Room struct {
	gorm.Model
	UUID    uuid.UUID `gorm:"index:idx_id",json:"uuid"`
	UserId1 uint      `gorm:"index:idx_user1_id",json:"user_id1"`
	UserId2 uint      `gorm:"index:idx_user2_id",json:"user_id2"`
	Close   bool      `json:"close"`
}
