package ws

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/img21326/fb_chat/mock"
	"github.com/posener/wstest"
	"github.com/stretchr/testify/assert"
)

func TestGetHistoryByUserID(t *testing.T) {
	c := gomock.NewController(t)

	wsUsecase := mock.NewMockWebsocketUsecaseInterface(c)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	NewWebsocketController(r, wsUsecase, nil, nil, nil)
	d := wstest.NewDialer(r)

	ws, _, err := d.Dial("ws://localhost/ws", nil)

	if err != nil {
		t.Errorf("connect ws error: %v", err)
	}

	_, p, err := ws.ReadMessage()

	assert.Equal(t, string(p[:]), `{'error': 'NotSetToken'}`)
	assert.Nil(t, err)
}
