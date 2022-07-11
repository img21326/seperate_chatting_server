package user

import "github.com/img21326/fb_chat/structure/user"

type UserRepoInterFace interface {
	Create(u *user.User) error
	FindByFbID(FbId string) (*user.User, error)
}
