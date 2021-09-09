package router

import (
	"github.com/gin-gonic/gin"
	"go-micro-module/20-temp/web/api/ceshi"
	"go-micro-module/20-temp/web/middlewares"
)

func InitTestRouter(Router *gin.RouterGroup){
	UserRouter := Router.Group("test").Use(middlewares.Trace())
	{
		UserRouter.GET("send-mq", ceshi.SendMq)
	}
}