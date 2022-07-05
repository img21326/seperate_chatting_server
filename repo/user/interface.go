package user

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model `gorm:"index:idx_name,unique"`
	FbID       string
	Name       string
	Email      string
	Gender     string
	FbLink     string
	Birth      time.Time
}

type UserRepoInterFace interface {
	Create(u *UserModel) error
	FindByFbID(FbId string) (*UserModel, error)
}
