package hub

import (
	"context"

	MessageHub "github.com/img21326/fb_chat/hub/message"
	PairHub "github.com/img21326/fb_chat/hub/pair"
	PubSubHub "github.com/img21326/fb_chat/hub/pubsub"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/usecase/message"
	"github.com/img21326/fb_chat/usecase/pair"
	"github.com/img21326/fb_chat/usecase/pubsub"
	"github.com/img21326/fb_chat/usecase/ws"
	"github.com/img21326/fb_chat/ws/client"
)

func StartHub(
	pubsubUsecase pubsub.SubMessageUsecaseInterface,
	pairUsecase pair.PairUsecaseInterface,
	messageUsecase message.MessageUsecaseInterface,
	wsUsecase ws.WebsocketUsecaseInterface,
) (
	pubMessageChan chan *pubmessage.PublishMessage,
	clientQueueChan chan *client.Client,
) {
	pubMessageChan = make(chan *pubmessage.PublishMessage, 4096)
	clientQueueChan = make(chan *client.Client, 4096)
	subMessageChan := make(chan *pubmessage.PublishMessage, 4096)

	ctx := context.Background()

	subHub := PubSubHub.NewSubHub(pubsubUsecase)
	pubHub := PubSubHub.NewPubHub(pubsubUsecase)
	messageHub := MessageHub.NewMessageHub(messageUsecase, wsUsecase, subMessageChan)
	pairHub := PairHub.NewPairHub(pairUsecase, pubMessageChan, clientQueueChan)

	go subHub.Run(ctx, "message", subMessageChan)
	go pubHub.Run(ctx, "message", pubMessageChan)
	go messageHub.Run(ctx)
	go pairHub.Run(ctx)

	return pubMessageChan, clientQueueChan
}
