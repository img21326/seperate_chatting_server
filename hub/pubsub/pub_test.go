package pubsub

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/img21326/fb_chat/mock"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
)

func TestPub(t *testing.T) {
	c := gomock.NewController(t)
	ctx, cancel := context.WithCancel(context.Background())
	pubUsecase := mock.NewMockSubMessageUsecaseInterface(c)

	msg := &pubmessage.PublishMessage{}
	pubUsecase.EXPECT().Publish(gomock.Any(), "test", msg).Times(1)

	pubHub := NewPubHub(pubUsecase)
	mc := make(chan *pubmessage.PublishMessage, 1)

	go pubHub.Run(ctx, "test", mc)
	mc <- msg
	cancel()
	time.Sleep(time.Duration(3) * time.Second)
}
