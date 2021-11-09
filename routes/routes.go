package routes

import (
	"gin_bluebell/controllers"
	"gin_bluebell/logger"
	"gin_bluebell/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter( model string) *gin.Engine {
	//读取配置是否启用发布模式
	if model == gin.ReleaseMode{
		gin.SetMode(gin.ReleaseMode)
	}
	//注册翻译
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	//注册业务路由
	v1 := r.Group("/api/v1")
	v1.POST("/signup",controllers.SignUpHandler) //登录
	v1.POST("/login",controllers.LoginHandler) //注册
	v1.GET("/refresh_token", controllers.RefreshTokenHandler) //刷新token

	v1.POST("/ping",middlewares.JWTAuthMiddleware(),func(c *gin.Context) {
		c.JSON(http.StatusOK,"TOKEN有效")
	})
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.GET("/community",controllers.CommunityHandler)
		v1.GET("/community/:id",controllers.CommunityDetailHandler)
	}
	return r
}


