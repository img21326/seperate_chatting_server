package main

import (
	"github.com/gin-gonic/gin"
	"github.com/img21326/fb_chat/controller"
	"github.com/img21326/fb_chat/helper"
	"github.com/img21326/fb_chat/usecase/oauth"
)

func main() {
	server := gin.Default()

	FacebookOauth := helper.NewFacebookOauth()
	FacebookUsecase := oauth.NewFacebookOauthUsecase(FacebookOauth)

	controller.NewLoginController(server, FacebookUsecase)

	server.Run(":8081")
}
