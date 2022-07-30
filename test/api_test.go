package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/img21326/fb_chat/server"
	"github.com/img21326/fb_chat/structure/user"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAPI(t *testing.T) {
	go server.StartUpRedisServer(DB, Redis, Port)

	res, err := http.Get(URL + "/register?gender=male")
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
