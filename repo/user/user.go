package user

import (
	"fmt"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepoInterFace {
	return &UserRepo{
		DB: db,
	}
}

func (repo *UserRepo) Create(u *UserModel) error {
	if err := repo.DB.Create(&u).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) FindByFbID(FbId string) (u *UserModel, err error) {
	err = repo.DB.Where("fb_id = ?", FbId).First(&u).Error
	if err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}
	return u, nil
}
