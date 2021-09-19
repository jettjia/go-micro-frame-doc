package router

import (
	"github.com/gin-gonic/gin"
	"go-micro-frame-doc/16-amqp/rabbitmq/01-gin-web/api/ceshi"
	"go-micro-frame-doc/16-amqp/rabbitmq/01-gin-web/middlewares"
)

func InitTestRouter(Router *gin.RouterGroup){
	UserRouter := Router.Group("test").Use(middlewares.Trace())
	{
		UserRouter.GET("send-mq", ceshi.SendMq)
	}
}