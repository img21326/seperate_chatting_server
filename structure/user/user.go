package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model
type User struct {
	gorm.Model
	UUID   uuid.UUID `gorm:"index:idx_name,unique"`
	Gender string
	// FbID       string
	// Name       string
	// Email      string
	// FbLink     string
	// Birth      time.Time
}
