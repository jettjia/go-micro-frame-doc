package router

import (
	"github.com/gin-gonic/gin"
	"go-micro-module/20-temp/web/api/user"
	"go-micro-module/20-temp/web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup){
	UserRouter := Router.Group("user").Use(middlewares.Trace())
	{
		UserRouter.GET("list", user.GetUserList)
	}
}