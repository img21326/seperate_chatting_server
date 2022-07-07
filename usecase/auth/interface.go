package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/img21326/fb_chat/repo/user"
)

type JwtConfig struct {
	Key            []byte
	ExpireDuration time.Duration
}

type AuthClaims struct {
	jwt.StandardClaims
	User user.UserModel
}

type AuthUsecaseInterFace interface {
	VerifyToken(token string) (user *user.UserModel, err error)
	GenerateToken(user *user.UserModel) (string, error)
}
