package hub_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/img21326/fb_chat/mock"
	"github.com/img21326/fb_chat/ws/client"
	"github.com/img21326/fb_chat/ws/hub"
)

func TestRegister(t *testing.T) {
	c := gomock.NewController(t)
	hubUsecase := mock.NewMockHubUsecaseInterface(c)

	hubUsecase.EXPECT().
		RegisterOnline(&client.Client{}).
		Return()

	onlineHub := hub.OnlineHub{
		Register:   make(chan *client.Client, 1),
		HubUsecase: hubUsecase,
	}

	client := &client.Client{}

	go onlineHub.Run()
	onlineHub.Register <- client

}
