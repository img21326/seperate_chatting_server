package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/img21326/fb_chat/mock"
	"github.com/img21326/fb_chat/structure/user"
	"github.com/img21326/fb_chat/usecase/auth"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	c := gomock.NewController(t)
	userRepo := mock.NewMockUserRepoInterFace(c)

	user := user.User{
		FbID:  "abcd",
		Email: "abc@gmail.com",
		Name:  "Liao",
	}

	userRepo.EXPECT().FindByFbID(gomock.Any(), "abcd").Return(&user, nil)

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

	user := user.User{
		FbID:  "abcd",
		Email: "abc@gmail.com",
		Name:  "Liao",
	}

	userRepo.EXPECT().FindByFbID(gomock.Any(), "abcd").Return(&user, nil).AnyTimes()

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

	assert.Equal(t, user.FbID, getUser.FbID)
	assert.Equal(t, err, nil)
}
