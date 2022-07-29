package room

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/img21326/fb_chat/structure/room"
	"github.com/stretchr/testify/assert"
)

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("db open error")
	}
	db.AutoMigrate(&room.Room{})
	return db
}

func TestCreate(t *testing.T) {
	db := initDB()
	roomRepo := &RoomRepo{
		DB: db,
	}
	ctx := context.Background()
	r := room.Room{
		UserId1: 1,
		UserId2: 2,
	}
	roomRepo.Create(ctx, &r)

	var getRoom room.Room
	roomRepo.DB.Where(&room.Room{ID: r.ID}).First(&getRoom)
	assert.Equal(t, getRoom.UserId1, uint(1))
	assert.Equal(t, getRoom.UserId2, uint(2))
	assert.Equal(t, getRoom.Close, false)
}

func TestClose(t *testing.T) {
	db := initDB()
	roomRepo := &RoomRepo{
		DB: db,
	}
	ctx := context.Background()
	r := room.Room{
		UserId1: 1,
		UserId2: 2,
	}
	roomRepo.Create(ctx, &r)
	roomRepo.Close(ctx, r.ID)
	var getRoom room.Room
	roomRepo.DB.Where(&room.Room{ID: r.ID}).First(&getRoom)
	assert.Equal(t, getRoom.Close, true)
}

func TestFindByUserId(t *testing.T) {
	db := initDB()
	roomRepo := &RoomRepo{
		DB: db,
	}
	ctx := context.Background()
	r := room.Room{
		UserId1: 1,
		UserId2: 2,
	}
	roomRepo.Create(ctx, &r)
	getRoom, err := roomRepo.FindByUserId(ctx, 1)
	assert.Equal(t, getRoom.ID, r.ID)
	assert.Nil(t, err)
}

// func TestFindByUserIdWithClose(t *testing.T) {
// 	db := initDB()
// 	roomRepo := &RoomRepo{
// 		DB: db,
// 	}
// 	ctx := context.Background()
// 	r := room.Room{
// 		UserId1: 1,
// 		UserId2: 2,
// 	}
// 	roomRepo.Create(ctx, &r)
// 	roomRepo.Close(ctx, r.ID)
// 	getRoom, err := roomRepo.FindByUserId(ctx, 1)
// 	assert.Equal(t, err, error.RoomIsClose)
// 	assert.Nil(t, getRoom)
// }
