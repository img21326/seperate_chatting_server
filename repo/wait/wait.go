package wait

import (
	"errors"
	"sync"

	"github.com/img21326/fb_chat/entity/ws"
)

type WaitRepo struct {
	ClientMap map[string][]*ws.Client
	lock      *sync.Mutex
}

func NewWaitRepo() WaitRepoInterface {
	return &WaitRepo{
		ClientMap: make(map[string][]*ws.Client),
		lock:      &sync.Mutex{},
	}
}

func (r *WaitRepo) Add(client *ws.Client) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()
	r.ClientMap[client.User.Gender] = append(r.ClientMap[client.User.Gender], client)
}

func (r *WaitRepo) Remove(client *ws.Client) {
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

func (r *WaitRepo) GetFirst(gender string) (client *ws.Client, err error) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()
	if _, isExist := r.ClientMap[gender]; !isExist {
		err = errors.New("QueueIsEmpty")
		return
	}
	if len(r.ClientMap[gender]) < 1 {
		err = errors.New("QueueIsEmpty")
		return
	}
	client = r.ClientMap[gender][0]
	r.ClientMap[gender] = r.ClientMap[gender][1:]
	return
}
