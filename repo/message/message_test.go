package message

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/img21326/fb_chat/helper"
	"github.com/img21326/fb_chat/structure/message"
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

	room1 := uuid.New()
	room2 := uuid.New()

	var fakeMessages []*message.Message
	for i := 1; i <= 150; i++ {

		usrID := uint(rand.Intn(2))
		var roomUID uuid.UUID
		if rand.Intn(2) == 1 {
			roomUID = room1
		} else {
			roomUID = room2
		}
		fakeMessages = append(fakeMessages, &message.Message{
			RoomId:  roomUID,
			UserId:  usrID,
			Message: helper.RandString(15),
			Time:    time.Now().Add(time.Hour * time.Duration(i)),
		})
	}
	db.Create(&fakeMessages)

	var resultMessageID2 []uint
	for i := len(fakeMessages) - 1; i >= 0; i-- {
		mes := fakeMessages[i]
		if mes.Time.After(fakeMessages[20].Time) {
			continue
		}
		if mes == fakeMessages[20] {
			continue
		}
		if mes.RoomId == room1 {
			resultMessageID2 = append(resultMessageID2, mes.ID)
		}
		if len(resultMessageID2) >= 20 {
			break
		}
	}

	getMsg, err := messageRepo.LastsByTime(ctx, room1, fakeMessages[20].Time, 20)
	var getMessageID []uint
	for _, mes := range getMsg {
		getMessageID = append(getMessageID, mes.ID)
	}
	assert.Equal(t, resultMessageID2, getMessageID)
	assert.Nil(t, err)
}
