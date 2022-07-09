package wait

import (
	"errors"
	"sync"

	"github.com/img21326/fb_chat/ws/client"
)

type WaitRepo struct {
	ClientMap map[string][]*client.Client
	lock      *sync.Mutex
}

func NewLocalWaitRepo() WaitRepoInterface {
	return &WaitRepo{
		ClientMap: make(map[string][]*client.Client),
		lock:      &sync.Mutex{},
	}
}

func (r *WaitRepo) Add(client *client.Client) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()
	r.ClientMap[client.User.Gender] = append(r.ClientMap[client.User.Gender], client)
}

func (r *WaitRepo) Remove(client *client.Client) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()
	var index int
	stat := false
	for i := range r.ClientMap[client.User.Gender] {
		if r.ClientMap[client.User.Gender][i].User.FbID == client.User.FbID {
			index = i
			stat = true
			break
		}
	}
	if !stat {
		return
	}
	r.ClientMap[client.User.Gender] = append(r.ClientMap[client.User.Gender][:index], r.ClientMap[client.User.Gender][index+1:]...)
}

func (r *WaitRepo) GetFirst(client *client.Client) (rclient *client.Client, err error) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()
	if _, isExist := r.ClientMap[client.WantToFind]; !isExist {
		err = errors.New("QueueIsEmpty")
		return
	}
	if len(r.ClientMap[client.WantToFind]) < 1 {
		err = errors.New("QueueIsEmpty")
		return
	}
	stat := false
	var index int
	for i := range r.ClientMap[client.WantToFind] {
		if r.ClientMap[client.WantToFind][i].WantToFind == client.User.Gender {
			index = i
			stat = true
			break
		}
	}
	if !stat {
		err = errors.New("NotFoundPairUser")
		return
	}
	rclient = r.ClientMap[client.WantToFind][index]
	r.ClientMap[client.WantToFind] = append(r.ClientMap[client.WantToFind][:index], r.ClientMap[client.WantToFind][index+1:]...)
	return
}
