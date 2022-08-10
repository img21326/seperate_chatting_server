package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	res, err := http.Get(URL + fmt.Sprintf(":%v", Port) + fmt.Sprintf("/refresh?uuid=%v", fakeUser[0].UUID))
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	type ResToken struct {
		Token string `json:"token"`
	}
	var a ResToken
	_ = json.Unmarshal(body, &a)

	c, _, err := websocket.DefaultDialer.Dial(WSURL+fmt.Sprintf(":%v", Port)+fmt.Sprintf("/ws?token=%v", a.Token), nil)
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
	res, err := http.Get(URL + fmt.Sprintf(":%v", Port) + fmt.Sprintf("/refresh?uuid=%v", fakeUser[0].UUID))
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	type ResToken struct {
		Token string `json:"token"`
	}
	var a ResToken
	_ = json.Unmarshal(body, &a)

	c, _, err := websocket.DefaultDialer.Dial(WSURL+fmt.Sprintf(":%v", Port)+fmt.Sprintf("/ws?token=%v", a.Token), nil)
	assert.Nil(t, err)
	_, message, err := c.ReadMessage()
	c.Close()
	assert.Nil(t, err)
	assert.Equal(t, string(message[:]), `{'error': 'NotSetWantParams'}`)
}

func TestUserParing(t *testing.T) {
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
	res, err := http.Get(URL + fmt.Sprintf(":%v", Port) + fmt.Sprintf("/refresh?uuid=%v", fakeUser[0].UUID))
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	type ResToken struct {
		Token string `json:"token"`
	}
	var a ResToken
	_ = json.Unmarshal(body, &a)

	c, _, err := websocket.DefaultDialer.Dial(WSURL+fmt.Sprintf(":%v", Port)+fmt.Sprintf("/ws?token=%v&want=female", a.Token), nil)
	assert.Nil(t, err)
	_, message, err := c.ReadMessage()
	c.Close()
	assert.Nil(t, err)
	assert.Equal(t, string(message[:]), `{'type': 'Paring'}`)
}
