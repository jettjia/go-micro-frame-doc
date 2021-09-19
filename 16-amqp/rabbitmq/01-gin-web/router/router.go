package router

import (
	"github.com/gin-gonic/gin"
	"go-micro-frame-doc/16-amqp/rabbitmq/01-gin-web/api/user"
	"go-micro-frame-doc/16-amqp/rabbitmq/01-gin-web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup){
	UserRouter := Router.Group("user").Use(middlewares.Trace())
	{
		UserRouter.GET("list", user.GetUserList)
	}
}