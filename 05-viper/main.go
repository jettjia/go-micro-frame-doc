package main

import (
	"fmt"
	"go-micro-frame-doc/05-viper/global"
	"go-micro-frame-doc/05-viper/initialize"
)

func main() {
	// 初始化 logger
	initialize.InitLogger()

	//初始化配置文件
	initialize.InitConfig()

	// 输出配置文件内容
	fmt.Println(global.ServerConfig.Name)
	fmt.Println(global.ServerConfig.Port)
	fmt.Println(global.ServerConfig.UserSrvInfo)
}
