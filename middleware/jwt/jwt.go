package jwt

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/img21326/fb_chat/usecase/auth"
)

type JWTValidMiddleware struct {
	authUsecase auth.AuthUsecaseInterFace
}

type JWTValidMiddlewareInterface interface {
	ValidHeaderToken(c *gin.Context)
}

func NewJWTValidMiddleware(u auth.AuthUsecaseInterFace) JWTValidMiddlewareInterface {
	return &JWTValidMiddleware{
		authUsecase: u,
	}
}

func (j *JWTValidMiddleware) ValidHeaderToken(c *gin.Context) {
	auth_string := c.GetHeader("Authorization")
	token := strings.Split(auth_string, "Bearer ")[1]
	user, err := j.authUsecase.VerifyToken(token)
	if err != nil {
		if err.Error() == "tokenExpired" {
			c.JSON(401, gin.H{
				"status": false,
				"msg":    fmt.Sprintf("jwt token expired"),
			})
			c.AbortWithStatus(401)
			return
		} else {
			c.JSON(500, gin.H{
				"status": false,
				"msg":    fmt.Sprintf("jwt valid error: %v", err),
			})
			c.AbortWithStatus(500)
			return
		}
	}
	c.Set("user", user)
	c.Next()
}
