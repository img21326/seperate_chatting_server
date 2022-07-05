package user

import (
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewFacebookOauthUsecase(db *gorm.DB) UserRepoInterFace {
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
	if err := repo.DB.Where(u.FbID, FbId).Find(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}
