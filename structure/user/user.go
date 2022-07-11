package user

import (
	"time"

	"gorm.io/gorm"
)

// Model
type User struct {
	gorm.Model `gorm:"index:idx_name,unique"`
	FbID       string
	Name       string
	Email      string
	Gender     string
	FbLink     string
	Birth      time.Time
}
