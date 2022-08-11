package pair

import (
	"context"
	"log"

	errorStruct "github.com/img21326/fb_chat/structure/error"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/structure/room"
	PairUsecase "github.com/img21326/fb_chat/usecase/pair"
	"github.com/img21326/fb_chat/ws/client"
)

type PairHub struct {
	PubMessageChan   chan *pubmessage.PublishMessage
	InsertClientChan chan *client.Client
	PairSuccessChan  chan *room.Room
	PairUsecase      PairUsecase.PairUsecaseInterface
}

func NewPairHub(pairUsecase PairUsecase.PairUsecaseInterface,
	pubMessageChan chan *pubmessage.PublishMessage,
	insertClientChan chan *client.Client) *PairHub {
	return &PairHub{
		InsertClientChan: insertClientChan,
		PubMessageChan:   pubMessageChan,
		PairSuccessChan:  make(chan *room.Room, 1024),
		PairUsecase:      pairUsecase,
	}
}

func (h *PairHub) HandleSuccess(room *room.Room) {
	c := context.Background()
	m1, m2, err := h.PairUsecase.PairSuccess(c, room)
	if err != nil {
		log.Printf("[PairHub] error: %+v, room: %+v", err, room)
	}
	h.PubMessageChan <- m1
	h.PubMessageChan <- m2
	log.Printf("[PairHub] pair user: %v %v in room: %v\n", room.UserId1, room.UserId2, room.ID)
}

func (h *PairHub) Run(ctx context.Context) {
	log.Printf("[PairHub] start")
	if h.PubMessageChan == nil {
		log.Printf("[PairHub] not set message chan")
	}
	for {
		select {
		case <-ctx.Done():
			// 關閉後就不再讓任何人加入
			close(h.InsertClientChan)
			close(h.PairSuccessChan)
			log.Printf("[PairHub] close channel\n")
			n := len(h.PairSuccessChan)
			for i := 0; i < n; i++ {
				h.HandleSuccess(<-h.PairSuccessChan)
			}
			log.Printf("[PairHub] finished room queue\n")
			return
		case client := <-h.InsertClientChan:
			// 先試著配對看看
			c := context.Background()
			room, err := h.PairUsecase.TryToPair(c, client)
			if err != nil {
				if err == errorStruct.PairNotSuccess || err == errorStruct.QueueSmallerThan1 {
					h.PairUsecase.AddToQueue(c, client)
				} else {
					log.Printf("[PairHub] error: %+v", err)
				}
				continue
			}
			h.PairSuccessChan <- room
		case room := <-h.PairSuccessChan:
			h.HandleSuccess(room)
		}
	}
}
