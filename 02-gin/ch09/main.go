package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	router := gin.Default()
	//优雅退出
	//LoadHTMLFiles会将指定的目录下的文件加载好， 相对目录
	//为什么我们通过goland运行main.go的时候并没有生成main.exe文件
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Println(dir)
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/**/*")
	//router.LoadHTMLFiles("templates/index.tmpl", "templates/goods.html")

	//如果没有在模板中使用define定义 那么我们就可以使用默认的文件名来找
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "myindex.tmpl", gin.H{
			"title": "慕课网",
		})
	})

	router.GET("/goods/list", func(c *gin.Context) {
		c.HTML(http.StatusOK, "goods/list.html", gin.H{
			"title": "慕课网",
		})
	})

	router.GET("/users/list", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/list.html", gin.H{
			"title": "慕课网",
		})
	})

	router.GET("/goods", func(c *gin.Context) {
		c.HTML(http.StatusOK, "goods.html", gin.H{
			"name": "微服务开发",
		})
	})

	router.Run(":8083")
}
