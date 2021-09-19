package global

import (
	"gorm.io/gorm"

	"go-micro-frame-doc/06-consul/config"
)

var (
	DB *gorm.DB
	ServerConfig config.ServerConfig
)