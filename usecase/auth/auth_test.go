package auth_test

import (
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

	userRepo.EXPECT().FindByFbID("abcd").Return(&user, nil)

	AuthUsecase := auth.NewAuthUsecase(
		auth.JwtConfig{
			Key:            []byte("TEST"),
			ExpireDuration: time.Hour * 24,
		},
		userRepo,
	)

	token, err := AuthUsecase.GenerateToken(&user)

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

	userRepo.EXPECT().FindByFbID("abcd").Return(&user, nil).AnyTimes()

	AuthUsecase := auth.NewAuthUsecase(
		auth.JwtConfig{
			Key:            []byte("TEST"),
			ExpireDuration: time.Hour * 24,
		},
		userRepo,
	)

	token, _ := AuthUsecase.GenerateToken(&user)

	userFbId, err := AuthUsecase.VerifyToken(token)

	assert.Equal(t, user.FbID, userFbId)
	assert.Equal(t, err, nil)
}
