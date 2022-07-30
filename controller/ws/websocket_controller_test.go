package ws

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/img21326/fb_chat/mock"
	"github.com/img21326/fb_chat/structure/room"
	userModel "github.com/img21326/fb_chat/structure/user"
	"github.com/img21326/fb_chat/ws/client"
	"github.com/posener/wstest"
	"github.com/stretchr/testify/assert"
)

func TestWithNotSetToken(t *testing.T) {

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	NewWebsocketController(r, nil, nil, nil, nil)
	d := wstest.NewDialer(r)

	ws, _, err := d.Dial("ws://localhost/ws", nil)

	if err != nil {
		t.Errorf("connect ws error: %v", err)
	}

	_, p, err := ws.ReadMessage()

	assert.Equal(t, string(p[:]), `{'error': 'NotSetToken'}`)
	assert.Nil(t, err)
}

func TestWithVerifyError(t *testing.T) {
	c := gomock.NewController(t)
	authUsecase := mock.NewMockAuthUsecaseInterFace(c)
	authUsecase.EXPECT().VerifyToken("1").Return(nil, errors.New("VerifyTokenError"))

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	NewWebsocketController(r, nil, authUsecase, nil, nil)
	d := wstest.NewDialer(r)

	ws, _, err := d.Dial("ws://localhost/ws?token=1", nil)

	if err != nil {
		t.Errorf("connect ws error: %v", err)
	}

	_, p, err := ws.ReadMessage()

	assert.Equal(t, string(p[:]), `{'error': 'VerifyTokenError'}`)
	assert.Nil(t, err)
}

func TestWithInRoom(t *testing.T) {
	c := gomock.NewController(t)
	authUsecase := mock.NewMockAuthUsecaseInterFace(c)
	wsUsecase := mock.NewMockWebsocketUsecaseInterface(c)

	usr := userModel.User{}
	usr.UUID = uuid.New()
	usr.ID = 1
	authUsecase.EXPECT().VerifyToken("1").Return(&usr, nil)
	ro := &room.Room{
		UUID:    uuid.New(),
		UserId1: 1,
		UserId2: 2,
		Close:   false,
	}
	wsUsecase.EXPECT().FindRoomByUserId(gomock.Any(), uint(1)).Return(ro, nil)
	wsUsecase.EXPECT().Register(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, cli *client.Client) {
		assert.Equal(t, cli.User.ID, uint(1))
	})

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	NewWebsocketController(r, wsUsecase, authUsecase, nil, nil)
	d := wstest.NewDialer(r)

	ws, _, err := d.Dial("ws://localhost/ws?token=1", nil)

	if err != nil {
		t.Errorf("connect ws error: %v", err)
	}

	_, p, err := ws.ReadMessage()
	assert.Equal(t, string(p[:]), `{'type': 'InRoom'}`)
	assert.Nil(t, err)
}

func TestWithoutRoomWithoutParams(t *testing.T) {
	c := gomock.NewController(t)
	authUsecase := mock.NewMockAuthUsecaseInterFace(c)
	wsUsecase := mock.NewMockWebsocketUsecaseInterface(c)

	usr := userModel.User{}
	usr.UUID = uuid.New()
	usr.ID = 1
	authUsecase.EXPECT().VerifyToken("1").Return(&usr, nil)

	wsUsecase.EXPECT().FindRoomByUserId(gomock.Any(), uint(1)).Return(nil, nil)
	wsUsecase.EXPECT().Register(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, cli *client.Client) {
		assert.Equal(t, cli.User.ID, uint(1))
	})

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	NewWebsocketController(r, wsUsecase, authUsecase, nil, nil)
	d := wstest.NewDialer(r)

	ws, _, err := d.Dial("ws://localhost/ws?token=1", nil)

	if err != nil {
		t.Errorf("connect ws error: %v", err)
	}

	_, p, err := ws.ReadMessage()
	assert.Equal(t, string(p[:]), `{'error': 'NotSetWantParams'}`)
	assert.Nil(t, err)
}

func TestWithoutRoom(t *testing.T) {
	c := gomock.NewController(t)
	authUsecase := mock.NewMockAuthUsecaseInterFace(c)
	wsUsecase := mock.NewMockWebsocketUsecaseInterface(c)

	usr := userModel.User{}
	usr.UUID = uuid.New()
	usr.ID = 1
	authUsecase.EXPECT().VerifyToken("1").Return(&usr, nil)

	wsUsecase.EXPECT().FindRoomByUserId(gomock.Any(), uint(1)).Return(nil, nil)
	wsUsecase.EXPECT().Register(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, cli *client.Client) {
		assert.Equal(t, cli.User.ID, uint(1))
	})

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	NewWebsocketController(r, wsUsecase, authUsecase, nil, make(chan *client.Client, 1))
	d := wstest.NewDialer(r)

	ws, _, err := d.Dial("ws://localhost/ws?token=1&want=male", nil)

	if err != nil {
		t.Errorf("connect ws error: %v", err)
	}

	_, p, err := ws.ReadMessage()
	assert.Equal(t, string(p[:]), `{'type': 'Paring'}`)
	assert.Nil(t, err)
}
