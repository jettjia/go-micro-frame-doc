package global

import (
	"gorm.io/gorm"

	"go-micro-module/06-consul/config"
)

var (
	DB *gorm.DB
	ServerConfig config.ServerConfig
)