package router

import (
	"github.com/gin-gonic/gin"
	api "go-micro-module/06-consul-discovery/api/user"
)

func InitUserRouter(Router *gin.RouterGroup){
	UserRouter := Router.Group("user")
	{
		UserRouter.GET("list", api.GetUserList)
	}
}