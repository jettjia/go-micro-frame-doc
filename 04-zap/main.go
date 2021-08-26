package main

import (
	"go-micro-module/04-zap/initialize"
	"go.uber.org/zap"
)

func main()  {
	// 初始化 logger
	initialize.InitLogger()

	zap.S().Debugf("entry main.go", "wwwwwwwwww")
}
