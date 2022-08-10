package pubsub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/img21326/fb_chat/mock"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
)

func TestSub(t *testing.T) {
	c := gomock.NewController(t)
	ctx, _ := context.WithCancel(context.Background())
	pubUsecase := mock.NewMockSubMessageUsecaseInterface(c)

	msg := &pubmessage.PublishMessage{
		Type:    "test",
		Payload: "asd",
	}
	j, _ := json.Marshal(msg)
	pubUsecase.EXPECT().Subscribe(gomock.Any(), "test").Times(1).DoAndReturn(
		func(ctx context.Context, topic string) func() ([]byte, error) {
			return func() ([]byte, error) {
				return j, nil
			}
		})

	subHub := NewSubHub(pubUsecase)
	mc := make(chan *pubmessage.PublishMessage, 1)
	go subHub.Run(ctx, "test", mc)

	assert.Equal(t, <-mc, msg)
}
