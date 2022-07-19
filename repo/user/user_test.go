package user

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/img21326/fb_chat/structure/user"
)

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("db open error")
	}
	db.AutoMigrate(&user.User{})
	return db
}

func TestCreate(t *testing.T) {
	db := initDB()
	userRepo := &UserRepo{DB: db}
	u := user.User{
		FbID:  "fb",
		Name:  "name",
		Email: "email@email.com",
	}
	ctx := context.Background()
	userRepo.Create(ctx, &u)

	var getU user.User
	err := userRepo.DB.Where(user.User{FbID: "fb"}).First(&getU).Error
	assert.Nil(t, err)
	assert.Equal(t, u.ID, getU.ID)
}

func TestFindByFbID(t *testing.T) {
	db := initDB()
	userRepo := &UserRepo{DB: db}
	u := user.User{
		FbID:  "fb",
		Name:  "name",
		Email: "email@email.com",
	}
	ctx := context.Background()
	userRepo.Create(ctx, &u)

	getU, err := userRepo.FindByFbID(ctx, "fb")
	assert.Nil(t, err)
	assert.Equal(t, u.ID, getU.ID)
}
