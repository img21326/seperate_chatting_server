package hub

import (
	"fmt"
	"log"

	"github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/repo/room"
	"github.com/img21326/fb_chat/usecase/hub"
	"github.com/img21326/fb_chat/ws/client"
)

type PairHub struct {
	Add                chan *client.Client
	Delete             chan *client.Client
	PublishMessageChan chan message.PublishMessage
	HubUsecase         hub.HubUsecaseInterface
}

func (h *PairHub) Run() {
	log.Printf("[pairHub] start")
	for {
		select {
		case client := <-h.Add:
			// 先試著配對看看
			pairClient, err := h.HubUsecase.GetFirstQueueUser(client)
			if err != nil {
				//如果配對失敗 就加入等待中
				h.HubUsecase.AddUserToQueue(client)
				log.Printf("[pairHub] add user to queue: %v\n", client.User.ID)
			} else {
				// 以下為配對成功所做的事
				room := &room.Room{
					UserId1: client.User.ID,
					UserId2: pairClient.User.ID,
					Close:   false,
				}
				err = h.HubUsecase.CreateRoom(room)
				if err != nil {
					log.Printf("create chat room err: %v", err)
					m1 := message.PublishMessage{
						Type:     "pairError",
						SendFrom: 0,
						SendTo:   pairClient.User.ID,
						Payload:  fmt.Sprintf("%v", err),
					}
					m2 := message.PublishMessage{
						Type:     "pairError",
						SendFrom: 0,
						SendTo:   client.User.ID,
						Payload:  fmt.Sprintf("%v", err),
					}
					h.PublishMessageChan <- m1
					h.PublishMessageChan <- m2
				}
				m1 := message.PublishMessage{
					Type:     "pairSuccess",
					SendFrom: client.User.ID,
					SendTo:   pairClient.User.ID,
					Payload:  room.ID,
				}
				m2 := message.PublishMessage{
					Type:     "pairSuccess",
					SendFrom: pairClient.User.ID,
					SendTo:   client.User.ID,
					Payload:  room.ID,
				}
				h.PublishMessageChan <- m1
				h.PublishMessageChan <- m2
				log.Printf("[pairHub] pair user: %v %v in room: %v\n", client.User.ID, pairClient.User.ID, room.ID)
			}
		case client := <-h.Delete:
			h.HubUsecase.DeleteuserFromQueue(client)
			log.Printf("[pairHub] delete queue user: %v\n", client.User.ID)
		}
	}
}