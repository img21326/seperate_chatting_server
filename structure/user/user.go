package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model
type User struct {
	gorm.Model `gorm:"index:idx_name,unique"`
	UUID       uuid.UUID
	Gender     string
	// FbID       string
	// Name       string
	// Email      string
	// FbLink     string
	// Birth      time.Time
}
