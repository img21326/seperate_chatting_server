package localonline

import (
	"sync"
	"testing"

	errStruct "github.com/img21326/fb_chat/structure/error"
	"github.com/img21326/fb_chat/structure/user"
	"github.com/img21326/fb_chat/ws/client"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRegister(t *testing.T) {
	onlineRepo := OnlineRepo{
		ClientMap: make(map[uint]*client.Client),
		lock:      &sync.Mutex{},
	}
	assert.Equal(t, len(onlineRepo.ClientMap), 0)
	m := gorm.Model{ID: 1}
	user := user.User{
		Model: m,
	}
	c := &client.Client{User: user}
	onlineRepo.Register(c)
	assert.Equal(t, len(onlineRepo.ClientMap), 1)
	assert.Equal(t, onlineRepo.ClientMap[1].User.ID, user.ID)
}

func TestUnRegister(t *testing.T) {
	onlineRepo := OnlineRepo{
		ClientMap: make(map[uint]*client.Client),
		lock:      &sync.Mutex{},
	}
	m := gorm.Model{ID: 1}
	user := user.User{
		Model: m,
	}
	c := &client.Client{User: user}
	onlineRepo.ClientMap[1] = c
	assert.Equal(t, len(onlineRepo.ClientMap), 1)
	onlineRepo.UnRegister(c)
	assert.Equal(t, len(onlineRepo.ClientMap), 0)
}

func TestFindUserByID(t *testing.T) {
	onlineRepo := OnlineRepo{
		ClientMap: make(map[uint]*client.Client),
		lock:      &sync.Mutex{},
	}
	m := gorm.Model{ID: 1}
	user := user.User{
		Model: m,
	}
	c := &client.Client{User: user}
	onlineRepo.Register(c)
	getClient, err := onlineRepo.FindUserByID(1)
	assert.Equal(t, nil, err)
	assert.Equal(t, c, getClient)
}

func TestFindUserByIDWithError(t *testing.T) {
	onlineRepo := OnlineRepo{
		ClientMap: make(map[uint]*client.Client),
		lock:      &sync.Mutex{},
	}
	u, err := onlineRepo.FindUserByID(5)
	assert.Nil(t, u)
	assert.Equal(t, err, errStruct.RecordNotFound)
}
