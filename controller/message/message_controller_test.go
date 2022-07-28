package message

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/img21326/fb_chat/mock"

	ModelMessage "github.com/img21326/fb_chat/structure/message"
	"github.com/img21326/fb_chat/structure/user"
)

func TestGetHistoryByUserID(t *testing.T) {
	c := gomock.NewController(t)

	roomID := uuid.New()
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
			RoomId:  roomID,
			UserId:  1,
			Message: "test",
			Time:    time.Now(),
		},
	}
	messageUsecase := mock.NewMockMessageUsecaseInterface(c)
	messageUsecase.EXPECT().LastByUserID(gomock.Any(), uint(1), int(20)).AnyTimes().Return(msgs, nil)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	u := &user.User{}
	u.ID = 1

	r.Use(func(ctx *gin.Context) {
		ctx.Set("user", u)
		ctx.Next()
	})

	NewMessageController(r, messageUsecase)

	req, _ := http.NewRequest("GET", "/history", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
