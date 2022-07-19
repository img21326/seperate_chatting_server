package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	model "github.com/img21326/fb_chat/structure/user"
)

type JwtConfig struct {
	Key            []byte
	ExpireDuration time.Duration
}

type AuthClaims struct {
	jwt.StandardClaims
	User model.User
}

type AuthUsecaseInterFace interface {
	VerifyToken(token string) (user *model.User, err error)
	GenerateToken(ctx context.Context, user *model.User) (string, error)
}
