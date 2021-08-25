package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
func main() {
	//实例化一个gin的server对象
	r := gin.Default()
	r.GET("/ping", pong)
	r.Run(":8083") // listen and serve on 0.0.0.0:8080
}
