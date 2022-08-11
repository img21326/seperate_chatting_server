package pair

import (
	"context"
	"fmt"
	"log"

	"github.com/img21326/fb_chat/repo/online"
	RoomRepo "github.com/img21326/fb_chat/repo/room"
	"github.com/img21326/fb_chat/repo/wait"
	errorStruct "github.com/img21326/fb_chat/structure/error"
	pubmessage "github.com/img21326/fb_chat/structure/pub_message"
	"github.com/img21326/fb_chat/structure/room"
	"github.com/img21326/fb_chat/ws/client"
)

type RedisPairUsecase struct {
	WaitRepo   wait.WaitRepoInterface
	OnlineRepo online.OnlineRepoInterface
	RoomRepo   RoomRepo.RoomRepoInterface
}

func NewRedisSubUsecase(
	waitRepo wait.WaitRepoInterface, onlineRepo online.OnlineRepoInterface, roomRepo RoomRepo.RoomRepoInterface,
) PairUsecaseInterface {
	return &RedisPairUsecase{
		WaitRepo:   waitRepo,
		OnlineRepo: onlineRepo,
		RoomRepo:   roomRepo,
	}
}

func (u *RedisPairUsecase) getInsertQueueName(client *client.Client) string {
	return fmt.Sprintf("%v_%v", client.User.Gender, client.WantToFind)
}

func (u *RedisPairUsecase) getPairQueueName(client *client.Client) string {
	return fmt.Sprintf("%v_%v", client.WantToFind, client.User.Gender)
}

func (u *RedisPairUsecase) TryToPair(ctx context.Context, client *client.Client) (newRoom *room.Room, err error) {
	pairStat := false
	for {
		// concurrence error
		if u.WaitRepo.Len(ctx, u.getPairQueueName(client)) < 1 {
			return nil, errorStruct.QueueSmallerThan1
		}
		pairClientID, err := u.WaitRepo.Pop(ctx, u.getPairQueueName(client))
		if err != nil {
			return nil, err
		}
		// concurrence error

		// 如果使用者已經下線,則在找下一個
		if !u.OnlineRepo.CheckUserOnline(ctx, pairClientID) {
			continue
		}
		newRoom = &room.Room{
			UserId1: client.User.ID,
			UserId2: pairClientID,
			Close:   false,
		}
		pairStat = true
		log.Printf("[PairUsecase] pair success: %v & %v\n", newRoom.UserId1, newRoom.UserId2)
		break
	}
	if pairStat {
		return newRoom, nil
	}
	return nil, errorStruct.PairNotSuccess
}

func (u *RedisPairUsecase) AddToQueue(ctx context.Context, client *client.Client) {
	u.WaitRepo.Add(ctx, u.getInsertQueueName(client), client.User.ID)
	log.Printf("[PairUsecase] add queue user: %v\n", client.User.ID)
}

func (u *RedisPairUsecase) PairSuccess(ctx context.Context, room *room.Room) (m1 *pubmessage.PublishMessage,
	m2 *pubmessage.PublishMessage, err error) {
	err = u.RoomRepo.Create(ctx, room)
	if err != nil {
		log.Printf("create chat room err: %v", err)
		m1 = &pubmessage.PublishMessage{
			Type:     "pairError",
			SendFrom: room.UserId1,
			SendTo:   room.UserId2,
			Payload:  "create room error",
		}
		m2 = &pubmessage.PublishMessage{
			Type:     "pairError",
			SendFrom: room.UserId1,
			SendTo:   room.UserId2,
			Payload:  "create room error",
		}
	} else {
		m1 = &pubmessage.PublishMessage{
			Type:     "pairSuccess",
			SendFrom: room.UserId1,
			SendTo:   room.UserId2,
			Payload:  room.UUID.String(),
		}
		m2 = &pubmessage.PublishMessage{
			Type:     "pairSuccess",
			SendFrom: room.UserId2,
			SendTo:   room.UserId1,
			Payload:  room.UUID.String(),
		}
	}
	return
}
