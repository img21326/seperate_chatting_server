package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/img21326/fb_chat/structure/user"
	"github.com/img21326/fb_chat/usecase/auth"
	"github.com/img21326/fb_chat/usecase/oauth"
)

type LoginController struct {
	OauthUsecase oauth.OauthUsecaseInterFace
	AuthUsecase  auth.AuthUsecaseInterFace
}

func NewLoginController(e *gin.Engine, oauthUsecase oauth.OauthUsecaseInterFace, authUsecase auth.AuthUsecaseInterFace) {
	controller := &LoginController{
		OauthUsecase: oauthUsecase,
		AuthUsecase:  authUsecase,
	}

	e.GET("/login", controller.Login)
	e.GET("/oauth", controller.Redirect)
}

func (c *LoginController) Login(ctx *gin.Context) {
	loginUrl := c.OauthUsecase.GetLoginURL()
	ctx.Redirect(http.StatusFound, loginUrl)
}

func (c *LoginController) Redirect(ctx *gin.Context) {
	key := ctx.Request.FormValue("state")
	token := ctx.Request.FormValue("code")
	oauthToken, err := c.OauthUsecase.GetRedirectToken(ctx, key, token)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}
	u, err := c.OauthUsecase.GetUser(oauthToken.Token)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}
	userModel := &user.User{
		FbID:   u.ID,
		Name:   u.Name,
		Email:  u.Email,
		Gender: u.Gender,
		FbLink: u.Link,
		Birth:  u.Birth,
	}
	jwtToken, err := c.AuthUsecase.GenerateToken(userModel)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"token": jwtToken,
	})
}
