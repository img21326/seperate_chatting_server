package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/img21326/fb_chat/mock"
	"github.com/img21326/fb_chat/structure/user"
	"github.com/img21326/fb_chat/usecase/auth"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGenerateToken(t *testing.T) {
	c := gomock.NewController(t)
	userRepo := mock.NewMockUserRepoInterFace(c)

	uid := uuid.New()
	user := user.User{
		UUID: uid,
	}

	userRepo.EXPECT().FindByID(gomock.Any(), uid.String()).Return(&user, nil)

	AuthUsecase := auth.NewAuthUsecase(
		auth.JwtConfig{
			Key:            []byte("TEST"),
			ExpireDuration: time.Hour * 24,
		},
		userRepo,
	)

	ctx := context.Background()
	token, err := AuthUsecase.GenerateToken(ctx, &user)

	assert.NotEqual(t, token, "")
	assert.Equal(t, err, nil)
}

func TestGetUserByToken(t *testing.T) {
	c := gomock.NewController(t)
	userRepo := mock.NewMockUserRepoInterFace(c)

	uid := uuid.New()
	user := user.User{
		UUID: uid,
	}

	userRepo.EXPECT().FindByID(gomock.Any(), uid.String()).Return(&user, nil).AnyTimes()

	AuthUsecase := auth.NewAuthUsecase(
		auth.JwtConfig{
			Key:            []byte("TEST"),
			ExpireDuration: time.Hour * 24,
		},
		userRepo,
	)

	ctx := context.Background()
	token, _ := AuthUsecase.GenerateToken(ctx, &user)

	getUser, err := AuthUsecase.VerifyToken(token)

	assert.Equal(t, user.UUID, getUser.UUID)
	assert.Equal(t, err, nil)
}

func TestGetUserWithNotFoundErr(t *testing.T) {
	c := gomock.NewController(t)
	userRepo := mock.NewMockUserRepoInterFace(c)

	uid := uuid.New()
	u := &user.User{
		UUID: uid,
	}

	userRepo.EXPECT().FindByID(gomock.Any(), uid.String()).Times(1).Return(nil, gorm.ErrRecordNotFound)
	userRepo.EXPECT().Create(gomock.Any(), u).Times(1).Do(func(ctx context.Context, u *user.User) {
		u.ID = 1
	})

	AuthUsecase := auth.NewAuthUsecase(
		auth.JwtConfig{
			Key:            []byte("TEST"),
			ExpireDuration: time.Hour * 24,
		},
		userRepo,
	)

	ctx := context.Background()
	token, _ := AuthUsecase.GenerateToken(ctx, u)

	getUser, err := AuthUsecase.VerifyToken(token)

	assert.Equal(t, u.UUID, getUser.UUID)
	assert.Equal(t, err, nil)
}
