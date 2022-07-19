package wait

import (
	"context"
	"errors"
	"sync"
)

type WaitRepo struct {
	ClientMap map[string][]uint
	lock      *sync.Mutex
}

func NewLocalWaitRepo() WaitRepoInterface {
	return &WaitRepo{
		ClientMap: make(map[string][]uint),
		lock:      &sync.Mutex{},
	}
}

func (r *WaitRepo) Add(ctx context.Context, queueName string, clientID uint) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()
	r.ClientMap[queueName] = append(r.ClientMap[queueName], clientID)
}

func (r *WaitRepo) Remove(ctx context.Context, queueName string, clientID uint) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()
	var index int
	stat := false
	for i := range r.ClientMap[queueName] {
		if r.ClientMap[queueName][i] == clientID {
			index = i
			stat = true
			break
		}
	}
	if !stat {
		return
	}
	r.ClientMap[queueName] = append(r.ClientMap[queueName][:index], r.ClientMap[queueName][index+1:]...)
}

func (r *WaitRepo) Len(ctx context.Context, queueName string) int {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()
	return len(r.ClientMap[queueName])
}

func (r *WaitRepo) Pop(ctx context.Context, queueName string) (clientID uint, err error) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()
	if _, isExist := r.ClientMap[queueName]; !isExist {
		err = errors.New("QueueIsEmpty")
		return
	}
	if len(r.ClientMap[queueName]) < 1 {
		err = errors.New("QueueIsEmpty")
		return
	}
	clientID = r.ClientMap[queueName][0]
	r.ClientMap[queueName] = r.ClientMap[queueName][1:]
	return
}
