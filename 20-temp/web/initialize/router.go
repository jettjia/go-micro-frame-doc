package initialize

import (
	"github.com/gin-gonic/gin"
	"go-micro-module/20-temp/web/middlewares"
	"go-micro-module/20-temp/web/router"
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
	router.InitTestRouter(ApiGroup)

	return Router
}