package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/img21326/fb_chat/repo/user"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	JwtConfig JwtConfig
	UserRepo  user.UserRepoInterFace
}

func NewAuthUsecase(jwtConfig JwtConfig, userRepo user.UserRepoInterFace) AuthUsecaseInterFace {
	return &AuthUsecase{
		JwtConfig: jwtConfig,
		UserRepo:  userRepo,
	}
}

func (u *AuthUsecase) VerifyToken(token string) (user *user.UserModel, err error) {
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

func (u *AuthUsecase) GenerateToken(user *user.UserModel) (string, error) {
	findUser, err := u.UserRepo.FindByFbID(user.FbID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}
	if err == gorm.ErrRecordNotFound {
		err = u.UserRepo.Create(user)
		if err != nil {
			return "", err
		}
		findUser, err = u.UserRepo.FindByFbID(user.FbID)
		if err != nil {
			return "", err
		}
	}
	jwtExpireAt := time.Now().Add(u.JwtConfig.ExpireDuration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   findUser.FbID,
			ExpiresAt: jwtExpireAt,
		},
		User: *findUser,
	})
	tokenString, err := token.SignedString(u.JwtConfig.Key)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
