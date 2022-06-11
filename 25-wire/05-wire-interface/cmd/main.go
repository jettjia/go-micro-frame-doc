package main

import (
	"fmt"
	"go-micro-frame-doc/25-wire/05-wire-interface/internal/db"
)

type App struct {
	dao db.IDao
}

func NewApp(dao db.IDao) *App {
	return &App{dao: dao}
}

func main() {
	app, _ := InitApp() // 使用 wire 生成的 injector 方法获取 app对象
	version, _ := app.dao.Version()
	fmt.Println(version)
}
