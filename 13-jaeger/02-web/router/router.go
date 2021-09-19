package router

import (
	"github.com/gin-gonic/gin"
	"go-micro-frame-doc/13-jaeger/02-web/api/user"
	"go-micro-frame-doc/13-jaeger/02-web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup){
	UserRouter := Router.Group("user").Use(middlewares.Trace())
	{
		UserRouter.GET("list", user.GetUserList)
	}
}