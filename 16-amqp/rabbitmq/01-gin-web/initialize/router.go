package initialize

import (
	"github.com/gin-gonic/gin"
	"go-micro-module/16-amqp/rabbitmq/01-gin-web/middlewares"
	"go-micro-module/16-amqp/rabbitmq/01-gin-web/router"
	"net/http"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"code":http.StatusOK,
			"success":true,
		})
	})

	//配置跨域
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup)

	return Router
}