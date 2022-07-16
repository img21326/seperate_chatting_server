package message

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/google/uuid"
	ModelMessage "github.com/img21326/fb_chat/structure/message"
	"github.com/stretchr/testify/assert"
)

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("db open error")
	}
	db.AutoMigrate(&ModelMessage.Message{})
	return db
}

func TestSave(t *testing.T) {
	db := initDB()
	messageRepo := &MessageRepo{
		DB: db,
	}
	ctx := context.Background()
	msg := ModelMessage.Message{
		RoomId:  uuid.New(),
		UserId:  1,
		Message: "test",
		Time:    time.Now(),
	}
	messageRepo.Save(ctx, &msg)

	var getMsg ModelMessage.Message
	messageRepo.DB.Where(&ModelMessage.Message{Message: "test", UserId: 1}).First(&getMsg)
	assert.Equal(t, getMsg.Message, msg.Message)
	assert.Equal(t, getMsg.UserId, uint(1))
}

func TestGetByID(t *testing.T) {
	db := initDB()
	messageRepo := &MessageRepo{
		DB: db,
	}
	ctx := context.Background()
	msg := ModelMessage.Message{
		ID:      1,
		RoomId:  uuid.New(),
		UserId:  1,
		Message: "test",
		Time:    time.Now(),
	}
	messageRepo.Save(ctx, &msg)

	getMsg, err := messageRepo.GetByID(ctx, 1)
	assert.Equal(t, getMsg.Message, msg.Message)
	assert.Equal(t, getMsg.UserId, uint(1))
	assert.Equal(t, getMsg.ID, uint(1))
	assert.Nil(t, err)
}

func TestLastsByRoomID(t *testing.T) {
	db := initDB()
	messageRepo := &MessageRepo{
		DB: db,
	}
	ctx := context.Background()
	roomID := uuid.New()
	roomID2 := uuid.New()
	msgs := []*ModelMessage.Message{
		&ModelMessage.Message{
			RoomId:  roomID,
			UserId:  1,
			Message: "test",
			Time:    time.Now(),
		},
		&ModelMessage.Message{
			RoomId:  roomID,
			UserId:  1,
			Message: "test",
			Time:    time.Now(),
		},
		&ModelMessage.Message{
			RoomId:  roomID,
			UserId:  1,
			Message: "test",
			Time:    time.Now(),
		},
		&ModelMessage.Message{
			RoomId:  roomID,
			UserId:  1,
			Message: "test",
			Time:    time.Now(),
		},
		&ModelMessage.Message{
			RoomId:  roomID,
			UserId:  1,
			Message: "test",
			Time:    time.Now(),
		},
		&ModelMessage.Message{
			RoomId:  roomID2,
			UserId:  1,
			Message: "test",
			Time:    time.Now(),
		},
	}
	for _, m := range msgs {
		messageRepo.Save(ctx, m)
	}

	getMsg, err := messageRepo.LastsByRoomID(ctx, roomID, 10)
	assert.Equal(t, len(getMsg), 5)
	assert.Nil(t, err)
	getMsg, err = messageRepo.LastsByRoomID(ctx, roomID, 4)
	assert.Equal(t, len(getMsg), 4)
	assert.Nil(t, err)
	getMsg, err = messageRepo.LastsByRoomID(ctx, roomID2, 10)
	assert.Equal(t, len(getMsg), 1)
	assert.Nil(t, err)
}

func TestLastsByTime(t *testing.T) {
	db := initDB()
	messageRepo := &MessageRepo{
		DB: db,
	}
	ctx := context.Background()
	roomID := uuid.New()
	var ts []time.Time
	for index, _ := range []int{1, 2, 3, 4, 5} {
		t, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("2022-01-01 0%v:0%v:0%v", index, index, index))
		ts = append(ts, t)
	}
	msgs := []*ModelMessage.Message{
		&ModelMessage.Message{
			RoomId:  roomID,
			UserId:  1,
			Message: "test",
			Time:    ts[0],
		},
		&ModelMessage.Message{
			RoomId:  roomID,
			UserId:  1,
			Message: "test",
			Time:    ts[1],
		},
		&ModelMessage.Message{
			RoomId:  roomID,
			UserId:  1,
			Message: "test",
			Time:    ts[2],
		},
		&ModelMessage.Message{
			RoomId:  roomID,
			UserId:  1,
			Message: "test",
			Time:    ts[3],
		},
		&ModelMessage.Message{
			RoomId:  roomID,
			UserId:  1,
			Message: "test",
			Time:    ts[4],
		},
	}
	for _, m := range msgs {
		messageRepo.Save(ctx, m)
	}

	getMsg, err := messageRepo.LastsByTime(ctx, roomID, ts[3], 5)
	assert.Equal(t, len(getMsg), 3)
	assert.Nil(t, err)
}
