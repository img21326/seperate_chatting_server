package chat

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/img21326/fb_chat/mock"
	"gorm.io/gorm"

	errStruct "github.com/img21326/fb_chat/structure/error"
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
	messageUsecase.EXPECT().LastByUserID(gomock.Any(), uint(1), gomock.Any()).Return(msgs, nil)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	u := &user.User{}
	u.ID = 1

	r.Use(func(ctx *gin.Context) {
		ctx.Set("user", u)
		ctx.Next()
	})

	NewChatController(r, messageUsecase, nil)

	type Res struct {
		Message []*ModelMessage.Message `json:"messages"`
	}
	jsonMsg, err := json.Marshal(&Res{Message: msgs})
	if err != nil {
		t.Errorf("marsh json err: %v", err)
	}

	req, _ := http.NewRequest("GET", "/history", nil)
	r.ServeHTTP(w, req)

	body := w.Body.Bytes()

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonMsg, body)
}

func TestGetHistoryBylastMessageId(t *testing.T) {
	c := gomock.NewController(t)

	roomID := uuid.New()
	msgs := []*ModelMessage.Message{
		&ModelMessage.Message{
			RoomId:  roomID,
			UserId:  1,
			Message: "test",
			Time:    time.Now(),
		},
	}
	messageUsecase := mock.NewMockMessageUsecaseInterface(c)
	messageUsecase.EXPECT().LastByMessageID(gomock.Any(), uint(1), uint(1), gomock.Any()).Return(msgs, nil)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	u := &user.User{}
	u.ID = 1

	r.Use(func(ctx *gin.Context) {
		ctx.Set("user", u)
		ctx.Next()
	})

	NewChatController(r, messageUsecase, nil)

	req, _ := http.NewRequest("GET", "/history?last_message_id=1", nil)
	r.ServeHTTP(w, req)

	type Res struct {
		Message []*ModelMessage.Message `json:"messages"`
	}
	jsonMsg, err := json.Marshal(&Res{Message: msgs})

	body := w.Body.Bytes()
	if err != nil {
		t.Errorf("read body err: %v", err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonMsg, body)
}

func TestInRoomWithClose(t *testing.T) {
	c := gomock.NewController(t)

	wsUsecase := mock.NewMockWebsocketUsecaseInterface(c)
	wsUsecase.EXPECT().FindRoomByUserId(gomock.Any(), uint(1)).Return(nil, errStruct.RoomIsClose)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	u := &user.User{}
	u.ID = 1

	r.Use(func(ctx *gin.Context) {
		ctx.Set("user", u)
		ctx.Next()
	})

	NewChatController(r, nil, wsUsecase)

	req, _ := http.NewRequest("GET", "/inroom", nil)
	r.ServeHTTP(w, req)

	type Res struct {
		Status bool `json:"status"`
	}
	jsonMsg, err := json.Marshal(&Res{Status: false})

	body := w.Body.Bytes()
	if err != nil {
		t.Errorf("read body err: %v", err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonMsg, body)
}

func TestInRoomWithNotFound(t *testing.T) {
	c := gomock.NewController(t)

	wsUsecase := mock.NewMockWebsocketUsecaseInterface(c)
	wsUsecase.EXPECT().FindRoomByUserId(gomock.Any(), uint(1)).Return(nil, gorm.ErrRecordNotFound)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	u := &user.User{}
	u.ID = 1

	r.Use(func(ctx *gin.Context) {
		ctx.Set("user", u)
		ctx.Next()
	})

	NewChatController(r, nil, wsUsecase)

	req, _ := http.NewRequest("GET", "/inroom", nil)
	r.ServeHTTP(w, req)

	type Res struct {
		Status bool `json:"status"`
	}
	jsonMsg, err := json.Marshal(&Res{Status: false})

	body := w.Body.Bytes()
	if err != nil {
		t.Errorf("read body err: %v", err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonMsg, body)
}

func TestInRoomWithNotClose(t *testing.T) {
	c := gomock.NewController(t)

	wsUsecase := mock.NewMockWebsocketUsecaseInterface(c)
	wsUsecase.EXPECT().FindRoomByUserId(gomock.Any(), uint(1)).Return(nil, nil)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	u := &user.User{}
	u.ID = 1

	r.Use(func(ctx *gin.Context) {
		ctx.Set("user", u)
		ctx.Next()
	})

	NewChatController(r, nil, wsUsecase)

	req, _ := http.NewRequest("GET", "/inroom", nil)
	r.ServeHTTP(w, req)

	type Res struct {
		Status bool `json:"status"`
	}
	jsonMsg, err := json.Marshal(&Res{Status: true})

	body := w.Body.Bytes()
	if err != nil {
		t.Errorf("read body err: %v", err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonMsg, body)
}
