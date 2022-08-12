package test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/img21326/fb_chat/server"
	"github.com/img21326/fb_chat/structure/room"
	"github.com/img21326/fb_chat/structure/user"
	"github.com/stretchr/testify/assert"
)

func TestConnectWithoutSetToken(t *testing.T) {
	Port := strconv.Itoa(randintRange(9550, 9500))
	go server.StartUpRedisServer(DB, Redis, Port)

	c, _, err := websocket.DefaultDialer.Dial(WSURL+fmt.Sprintf(":%v", Port)+"/ws", nil)
	assert.Nil(t, err)
	_, message, err := c.ReadMessage()
	c.Close()
	assert.Nil(t, err)
	assert.Equal(t, string(message[:]), `{'error': 'NotSetToken'}`)
}

func TestConnectWithTokenUnVerify(t *testing.T) {
	Port := strconv.Itoa(randintRange(9650, 9600))
	go server.StartUpRedisServer(DB, Redis, Port)

	c, _, err := websocket.DefaultDialer.Dial(WSURL+fmt.Sprintf(":%v", Port)+"/ws?token=abc", nil)
	assert.Nil(t, err)
	_, message, err := c.ReadMessage()
	c.Close()
	assert.Nil(t, err)
	assert.Equal(t, string(message[:]), `{'error': 'token contains an invalid number of segments'}`)
}

func TestUserInRoom(t *testing.T) {
	fakeUser := []*user.User{
		&user.User{
			UUID:   uuid.New(),
			Gender: "male",
		},
		&user.User{
			UUID:   uuid.New(),
			Gender: "female",
		},
		&user.User{
			UUID:   uuid.New(),
			Gender: "female",
		},
	}
	DB.Create(&fakeUser)
	fakeRoom := []*room.Room{
		&room.Room{
			UserId1: fakeUser[0].ID,
			UserId2: fakeUser[1].ID,
			UUID:    uuid.New(),
			Close:   false,
		},
		&room.Room{
			UserId1: fakeUser[0].ID,
			UserId2: fakeUser[2].ID,
			UUID:    uuid.New(),
			Close:   true,
		},
	}
	DB.Create(&fakeRoom)

	Port := strconv.Itoa(randintRange(9700, 9650))
	go server.StartUpRedisServer(DB, Redis, Port)
	token := getUserToken(Port, fakeUser[0].UUID)

	c, _, err := websocket.DefaultDialer.Dial(WSURL+fmt.Sprintf(":%v", Port)+fmt.Sprintf("/ws?token=%v", token), nil)
	assert.Nil(t, err)
	_, message, err := c.ReadMessage()
	c.Close()
	assert.Nil(t, err)
	assert.Equal(t, string(message[:]), `{'type': 'InRoom'}`)
}

func TestUserParingErrorWithNotSetParams(t *testing.T) {
	fakeUser := []*user.User{
		&user.User{
			UUID:   uuid.New(),
			Gender: "male",
		},
		&user.User{
			UUID:   uuid.New(),
			Gender: "female",
		},
		&user.User{
			UUID:   uuid.New(),
			Gender: "female",
		},
	}
	DB.Create(&fakeUser)
	fakeRoom := []*room.Room{
		&room.Room{
			UserId1: fakeUser[0].ID,
			UserId2: fakeUser[1].ID,
			UUID:    uuid.New(),
			Close:   true,
		},
		&room.Room{
			UserId1: fakeUser[0].ID,
			UserId2: fakeUser[2].ID,
			UUID:    uuid.New(),
			Close:   true,
		},
	}
	DB.Create(&fakeRoom)

	Port := strconv.Itoa(randintRange(9700, 9650))
	go server.StartUpRedisServer(DB, Redis, Port)
	token := getUserToken(Port, fakeUser[0].UUID)

	c, _, err := websocket.DefaultDialer.Dial(WSURL+fmt.Sprintf(":%v", Port)+fmt.Sprintf("/ws?token=%v", token), nil)
	assert.Nil(t, err)
	_, message, err := c.ReadMessage()
	c.Close()
	assert.Nil(t, err)
	assert.Equal(t, string(message[:]), `{'error': 'NotSetWantParams'}`)
}

func TestUserInParing(t *testing.T) {
	fakeUser := []*user.User{
		&user.User{
			UUID:   uuid.New(),
			Gender: "male",
		},
		&user.User{
			UUID:   uuid.New(),
			Gender: "female",
		},
		&user.User{
			UUID:   uuid.New(),
			Gender: "female",
		},
	}
	DB.Create(&fakeUser)
	fakeRoom := []*room.Room{
		&room.Room{
			UserId1: fakeUser[0].ID,
			UserId2: fakeUser[1].ID,
			UUID:    uuid.New(),
			Close:   true,
		},
		&room.Room{
			UserId1: fakeUser[0].ID,
			UserId2: fakeUser[2].ID,
			UUID:    uuid.New(),
			Close:   true,
		},
	}
	DB.Create(&fakeRoom)

	Port := strconv.Itoa(randintRange(9700, 9650))
	go server.StartUpRedisServer(DB, Redis, Port)
	token := getUserToken(Port, fakeUser[0].UUID)

	c, _, err := websocket.DefaultDialer.Dial(WSURL+fmt.Sprintf(":%v", Port)+fmt.Sprintf("/ws?token=%v&want=female", token), nil)
	assert.Nil(t, err)
	_, message, err := c.ReadMessage()
	c.Close()
	assert.Nil(t, err)
	assert.Equal(t, string(message[:]), `{'type': 'Paring'}`)
}

func TestUserParingSuccess(t *testing.T) {
	fakeUser := []*user.User{
		&user.User{
			UUID:   uuid.New(),
			Gender: "male",
		},
		&user.User{
			UUID:   uuid.New(),
			Gender: "female",
		},
	}
	DB.Create(&fakeUser)

	Port := strconv.Itoa(randintRange(9750, 9700))
	go server.StartUpRedisServer(DB, Redis, Port)
	user1Token := getUserToken(Port, fakeUser[0].UUID)
	user2Token := getUserToken(Port, fakeUser[1].UUID)
	c1, _, err := websocket.DefaultDialer.Dial(WSURL+fmt.Sprintf(":%v", Port)+fmt.Sprintf("/ws?token=%v&want=female", user1Token), nil)
	assert.Nil(t, err)
	go func() {
		c2, _, err := websocket.DefaultDialer.Dial(WSURL+fmt.Sprintf(":%v", Port)+fmt.Sprintf("/ws?token=%v&want=male", user2Token), nil)
		assert.Nil(t, err)
		for {
			c2.ReadMessage()
		}
	}()
	index := 0
	for {
		_, msg, err := c1.ReadMessage()
		t.Log(string(msg[:]))
		t.Log(index)
		assert.Nil(t, err)
		if index == 0 {
			assert.Equal(t, string(msg[:]), `{'type': 'Paring'}`)
			index += 1
			continue
		}
		if index == 1 {
			type Res struct {
				Type     string    `json:"type"`
				SendFrom int       `json:"sendFrom"`
				SendTo   int       `json:"sendTo"`
				Payload  uuid.UUID `json:"payload"`
			}
			var res Res
			err = json.Unmarshal(msg, &res)
			assert.Nil(t, err)
			assert.Equal(t, res.Type, "pairSuccess")
			assert.NotNil(t, res.SendFrom)
			assert.NotNil(t, res.SendTo)
			assert.NotNil(t, res.Payload)
			break
		}
	}

}
