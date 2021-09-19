package initialize

import (
	"github.com/gin-gonic/gin"
	"go-micro-frame-doc/13-jaeger/02-web/middlewares"
	"go-micro-frame-doc/13-jaeger/02-web/router"
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