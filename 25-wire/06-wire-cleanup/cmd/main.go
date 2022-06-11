package main

import (
	"fmt"
	"log"

	"go-micro-frame-doc/25-wire/06-wire-cleanup/internal/config"
)

type App struct {
	Config *config.Config
}

func NewApp(config *config.Config) *App {
	return &App{Config: config}
}

func main() {
	app, cleanup, err := InitApp()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup() // 处理需要关闭的资源
	fmt.Println("输出数据配置：", app.Config.Database.Dsn)
}
