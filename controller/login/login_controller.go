package login

import (
	"fmt"
	"log"

	"github.com/img21326/fb_chat/structure/user"

	"github.com/gin-gonic/gin"
	"github.com/img21326/fb_chat/usecase/auth"
)

type LoginController struct {
	AuthUsecase auth.AuthUsecaseInterFace
}

func NewLoginController(e gin.IRoutes, authUsecase auth.AuthUsecaseInterFace) {
	controller := &LoginController{
		AuthUsecase: authUsecase,
	}

	// e.GET("/login", controller.Login)
	// e.GET("/oauth", controller.Redirect)
	e.GET("/register", controller.Register)
}

func (c *LoginController) Register(ctx *gin.Context) {
	gender, ok := ctx.GetQuery("gender")
	if !ok {
		log.Print("[LoginController] Register without gender params")
		ctx.JSON(410, gin.H{
			"error": "should add params with gender",
		})
		return
	}
	if gender != "male" && gender != "female" {
		log.Print("[LoginController] Register without gender params")
		ctx.JSON(410, gin.H{
			"error": "gender should be male or female",
		})
		return
	}
	newUser := &user.User{
		Gender: gender,
	}

	token, err := c.AuthUsecase.GenerateToken(ctx, newUser)
	if err != nil {
		log.Printf("[LoginController] Register Generate token err: %v", err)
		ctx.JSON(500, gin.H{
			"error": fmt.Sprint("server generate token error"),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"token": token,
		"uuid":  newUser.UUID,
	})
	return
}

// func (c *LoginController) Login(ctx *gin.Context) {
// 	loginUrl := c.OauthUsecase.GetLoginURL()
// 	ctx.Redirect(http.StatusFound, loginUrl)
// }

// func (c *LoginController) Redirect(ctx *gin.Context) {
// 	key := ctx.Request.FormValue("state")
// 	token := ctx.Request.FormValue("code")
// 	oauthToken, err := c.OauthUsecase.GetRedirectToken(ctx, key, token)
// 	if err != nil {
// 		ctx.JSON(500, gin.H{
// 			"error": fmt.Sprintf("%v", err),
// 		})
// 		return
// 	}
// 	u, err := c.OauthUsecase.GetUser(oauthToken.Token)
// 	if err != nil {
// 		ctx.JSON(500, gin.H{
// 			"error": fmt.Sprintf("%v", err),
// 		})
// 		return
// 	}
// 	userModel := &user.User{

// 	}
// 	jwtToken, err := c.AuthUsecase.GenerateToken(ctx, userModel)
// 	if err != nil {
// 		ctx.JSON(500, gin.H{
// 			"error": fmt.Sprintf("%v", err),
// 		})
// 		return
// 	}
// 	ctx.JSON(200, gin.H{
// 		"token": jwtToken,
// 	})
// }
