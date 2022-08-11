package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func randate() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func randintRange(max int, min int) int {
	return rand.Intn(max-min) + min
}

func getUserToken(port string, uid uuid.UUID) string {
	res, err := http.Get(URL + fmt.Sprintf(":%v", port) + fmt.Sprintf("/refresh?uuid=%v", uid))
	if err != nil {
		log.Fatal("get token error with get")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("get token error with read body")
	}
	type ResToken struct {
		Token string `json:"token"`
	}
	var a ResToken
	err = json.Unmarshal(body, &a)
	if err != nil {
		log.Fatal("get token error with decode json")
	}
	return a.Token
}
