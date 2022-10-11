package pubsub

import (
	"context"
	"sync"

	ChannelClose "github.com/img21326/fb_chat/helper/chan_close"
	errorStruct "github.com/img21326/fb_chat/structure/error"
	"github.com/img21326/fb_chat/structure/pub"
)

type LocalPubSubRepo struct {
	SubscribeMap map[string][]*ChannelClose.ChannelClose[*pub.ReceiveMessage]
	Lock         *sync.Mutex
}

func NewLocalPubSubRepo() PubSubRepoInterface {
	return &LocalPubSubRepo{
		SubscribeMap: make(map[string][]*ChannelClose.ChannelClose[*pub.ReceiveMessage]),
		Lock:         new(sync.Mutex),
	}
}

func (repo *LocalPubSubRepo) Sub(ctx context.Context, topic string) func() ([]byte, error) {
	defer repo.Lock.Unlock()
	repo.Lock.Lock()
	c := make(chan *pub.ReceiveMessage, 1024)
	chanClose := ChannelClose.NewChanClose(c)
	repo.SubscribeMap[topic] = append(repo.SubscribeMap[topic], chanClose)

	returnChan := make(chan *pub.ReceiveMessage, 1024)
	go func(ctx context.Context, ReturnChan chan *pub.ReceiveMessage) {
		defer func() {
			close(ReturnChan)
			chanClose.Close()
		}()
		for {
			select {
			case <-ctx.Done():
				rm := &pub.ReceiveMessage{Error: errorStruct.ChannelClosed}
				ReturnChan <- rm
				return
			default:
				msg := chanClose.Pop()
				ReturnChan <- msg
			}
		}
	}(ctx, returnChan)
	return func() ([]byte, error) {
		rm := <-returnChan
		return rm.Payload, rm.Error
	}
}

func (repo *LocalPubSubRepo) Pub(ctx context.Context, topic string, message []byte) error {
	rm := &pub.ReceiveMessage{Payload: message, Error: nil}
	for _, ch := range repo.SubscribeMap[topic] {
		if !ch.Closed {
			_ = ch.Push(rm)
		}
	}
	return nil
}
