package router

import (
	"github.com/gin-gonic/gin"
	"go-micro-frame-doc/14-sentinel/02-gin-sentinel/api/user"
	"go-micro-frame-doc/14-sentinel/02-gin-sentinel/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup){
	UserRouter := Router.Group("user").Use(middlewares.Trace())
	{
		UserRouter.GET("list", user.GetUserList)
	}
}