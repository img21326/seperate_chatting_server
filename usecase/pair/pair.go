package pair

import (
	"github.com/google/uuid"
	"github.com/img21326/fb_chat/entity/ws"
	"github.com/img21326/fb_chat/repo/message"
	"github.com/img21326/fb_chat/repo/online"
	"github.com/img21326/fb_chat/repo/room"
	"github.com/img21326/fb_chat/repo/user"
	"github.com/img21326/fb_chat/repo/wait"
)

type PairUsecase struct {
	UserRepo    user.UserRepoInterFace
	MessageRepo message.MessageRepoInterface
	OnlineRepo  online.OnlineRepoInterface
	RoomRepo    room.RoomRepoInterface
	WaitRepo    wait.WaitRepoInterface
}

func NewPairUsecase(userRepo user.UserRepoInterFace, messageRepo message.MessageRepoInterface,
	onlineRepo online.OnlineRepoInterface, roomRepo room.RoomRepoInterface, waitRepo wait.WaitRepoInterface,
) PairUsecaseInterface {
	return &PairUsecase{
		UserRepo:    userRepo,
		MessageRepo: messageRepo,
		OnlineRepo:  onlineRepo,
		RoomRepo:    roomRepo,
		WaitRepo:    waitRepo,
	}
}

func (u *PairUsecase) RegisterOnline(client *ws.Client) {
	u.OnlineRepo.Register(client)
}

func (u *PairUsecase) UnRegisterOnline(client *ws.Client) {
	u.OnlineRepo.UnRegister(client)
}

func (u *PairUsecase) FindOnlineUserByFbID(userId uint) (*ws.Client, error) {
	return u.OnlineRepo.FindUserByFbID(userId)
}

func (u *PairUsecase) FindUserByFbID(FbId string) (*user.UserModel, error) {
	return u.UserRepo.FindByFbID(FbId)
}

func (u *PairUsecase) GetFirstQueueUser(client *ws.Client) (*ws.Client, error) {
	return u.WaitRepo.GetFirst(client)
}

func (u *PairUsecase) AddUserToQueue(client *ws.Client) {
	u.WaitRepo.Add(client)
}

func (u *PairUsecase) DeleteuserFromQueue(client *ws.Client) {
	u.WaitRepo.Remove(client)
}

func (u *PairUsecase) CreateRoom(room *room.Room) error {
	return u.RoomRepo.Create(room)
}

func (u *PairUsecase) CloseRoom(uuid uuid.UUID) error {
	return u.RoomRepo.Close(uuid)
}

func (u *PairUsecase) FindRoomByUserId(userId uint) (*room.Room, error) {
	return u.RoomRepo.FindByUserId(userId)
}

func (u *PairUsecase) SaveMessage(message *message.MessageModel) {
	u.MessageRepo.Save(message)
}
