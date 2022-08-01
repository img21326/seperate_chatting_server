package test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/img21326/fb_chat/server"
	"github.com/stretchr/testify/assert"
)

func TestConnectWithoutSetToken(t *testing.T) {
	Port := strconv.Itoa(randintRange(9600, 9500))
	go server.StartUpRedisServer(DB, Redis, Port)

	c, _, err := websocket.DefaultDialer.Dial(WSURL+fmt.Sprintf(":%v", Port)+"/ws", nil)
	assert.Nil(t, err)
	_, message, err := c.ReadMessage()
	c.Close()
	assert.Nil(t, err)
	assert.Equal(t, string(message[:]), `{'error': 'NotSetToken'}`)
}
