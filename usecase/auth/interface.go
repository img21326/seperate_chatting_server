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
	UserFbID string
}

type UsecaseAuthInterFace interface {
	GetUserByToken(token string) (*user.UserModel, error)
	GenerateToken(user *user.UserModel) (string, error)
}
