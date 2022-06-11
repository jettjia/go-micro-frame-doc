package main

import (
	"fmt"
	"go-micro-frame-doc/25-wire/04-wire-value/internal/config"
	"log"
)

type App struct {
	Config *config.Config
}

func NewApp(config *config.Config) *App {
	return &App{Config: config}
}

func main() {
	app, err := InitApp()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("输出数据库配置", app.Config.Database.Dsn)
}
