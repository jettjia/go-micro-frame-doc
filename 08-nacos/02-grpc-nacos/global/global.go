package global

import (
	"gorm.io/gorm"

	"go-micro-module/08-nacos/02-grpc-nacos/config"
)

var (
	DB *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig config.NacosConfig
)