package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	repo "github.com/img21326/fb_chat/repo/user"
	model "github.com/img21326/fb_chat/structure/user"
)

type AuthUsecase struct {
	JwtConfig JwtConfig
	UserRepo  repo.UserRepoInterFace
}

func NewAuthUsecase(jwtConfig JwtConfig, userRepo repo.UserRepoInterFace) AuthUsecaseInterFace {
	return &AuthUsecase{
		JwtConfig: jwtConfig,
		UserRepo:  userRepo,
	}
}

func (u *AuthUsecase) VerifyToken(token string) (user *model.User, err error) {
	var claims AuthClaims
	t, err := jwt.ParseWithClaims(token, &claims, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", jwtToken.Header["alg"])
		}
		return u.JwtConfig.Key, nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	user = &claims.User
	return
}

func (u *AuthUsecase) GenerateToken(ctx context.Context, user *model.User) (token string, err error) {
	var usr *model.User
	if user.UUID == uuid.Nil {
		user.UUID = uuid.New()
		err = u.UserRepo.Create(ctx, user)
		if err != nil {
			return "", err
		}
		usr = user
	} else {
		usr, err = u.UserRepo.FindByID(ctx, user.UUID.String())
		if err != nil {
			return "", err
		}
	}
	jwtExpireAt := time.Now().Add(u.JwtConfig.ExpireDuration).Unix()
	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.UUID.String(),
			ExpiresAt: jwtExpireAt,
		},
		User: *usr,
	})
	token, err = jwtClaims.SignedString(u.JwtConfig.Key)
	if err != nil {
		return "", err
	}

	return

}
