package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/img21326/fb_chat/usecase/oauth"
)

type LoginController struct {
	OauthUsecase oauth.UsecaseOauthInterFace
}

func NewLoginController(e *gin.Engine, oauthUsecase oauth.UsecaseOauthInterFace) {
	controller := &LoginController{
		OauthUsecase: oauthUsecase,
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
			"error": err,
		})
	}
	user, err := c.OauthUsecase.GetUser(oauthToken.Token)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err,
		})
	}
	ctx.JSON(200, gin.H{
		"user": user,
	})
}
