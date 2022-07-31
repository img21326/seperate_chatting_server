package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/img21326/fb_chat/helper"
	"github.com/img21326/fb_chat/server"
	"github.com/img21326/fb_chat/structure/message"
	"github.com/img21326/fb_chat/structure/room"
	"github.com/img21326/fb_chat/structure/user"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAPI(t *testing.T) {
	Port := strconv.Itoa(randintRange(9500, 9400))

	go server.StartUpRedisServer(DB, Redis, Port)

	res, err := http.Get(URL + fmt.Sprintf(":%v", Port) + "/register?gender=male")
	assert.Nil(t, err)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	type Res struct {
		Token string `json:"token"`
		UUID  string `json:"UUID"`
	}
	var r Res
	err = json.Unmarshal(body, &r)
	assert.Nil(t, err)

	uid, err := uuid.Parse(r.UUID)
	assert.Nil(t, err)
	assert.NotNil(t, r.Token)

	var getU user.User
	err = DB.Where(user.User{UUID: uid}).First(&getU).Error
	assert.Nil(t, err)
	assert.Equal(t, uid, getU.UUID)
}

func TestMessageHistoryAPI(t *testing.T) {
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
			Close:   false,
		},
	}
	DB.Create(&fakeRoom)
	var fakeMessages []*message.Message
	for i := 1; i <= 50; i++ {
		randInt := rand.Intn(2)

		var room *room.Room
		if randInt == 1 {
			room = fakeRoom[0]
		} else {
			room = fakeRoom[1]
		}

		var usrID uint
		randInt = rand.Intn(2)
		if randInt == 1 {
			usrID = room.UserId1
		} else {
			usrID = room.UserId2
		}
		fakeMessages = append(fakeMessages, &message.Message{
			RoomId:  room.UUID,
			UserId:  usrID,
			Message: helper.RandString(15),
			Time:    time.Now().Add(time.Hour * time.Duration(i)),
		})
	}
	DB.Create(&fakeMessages)

	Port := strconv.Itoa(randintRange(9400, 9300))
	go server.StartUpRedisServer(DB, Redis, Port)

	res, err := http.Get(URL + fmt.Sprintf(":%v", Port) + fmt.Sprintf("/refresh?uuid=%v", fakeUser[0].UUID))
	assert.Nil(t, err)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	type ResToken struct {
		Token string `json:"token"`
	}
	var a ResToken
	err = json.Unmarshal(body, &a)
	assert.Nil(t, err)
	assert.NotNil(t, a.Token)

	/// refresh token end

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL+fmt.Sprintf(":%v", Port)+"/auth/history", nil)
	assert.Nil(t, err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", a.Token))
	res, err = client.Do(req)
	assert.Nil(t, err)
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	type Res struct {
		Messages []message.Message `json:"messages"`
	}
	var r Res
	err = json.Unmarshal(body, &r)

	var resultMessageID []uint
	for i := len(fakeMessages) - 1; i >= 0; i-- {
		mes := fakeMessages[i]
		if mes.RoomId == fakeRoom[1].UUID {
			resultMessageID = append(resultMessageID, mes.ID)
		}
		if len(resultMessageID) >= 20 {
			break
		}
	}
	var getMessageID []uint
	for _, mes := range r.Messages {
		getMessageID = append(getMessageID, mes.ID)
	}

	assert.Nil(t, err)
	assert.Equal(t, getMessageID, resultMessageID)
	// basic end
	randi := randintRange(len(fakeMessages)-1, 0)
	req, err = http.NewRequest("GET", URL+fmt.Sprintf(":%v", Port)+fmt.Sprintf("/auth/history?last_message_id=%v", fakeMessages[randi].ID), nil)
	assert.Nil(t, err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", a.Token))
	res, err = client.Do(req)
	assert.Nil(t, err)
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &r)

	var resultMessageID2 []uint
	for i := len(fakeMessages) - 1; i >= 0; i-- {
		mes := fakeMessages[i]
		if mes.Time.After(fakeMessages[randi].Time) {
			continue
		}
		if mes.RoomId == fakeRoom[1].UUID {
			resultMessageID2 = append(resultMessageID2, mes.ID)
		}
		if len(resultMessageID2) >= 20 {
			break
		}
	}
	var getMessageID2 []uint
	for _, mes := range r.Messages {
		getMessageID2 = append(getMessageID2, mes.ID)
	}

	assert.Nil(t, err)
	assert.Equal(t, resultMessageID2, getMessageID2)
	// with params out
}
