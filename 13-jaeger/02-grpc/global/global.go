package global

import (
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"

	"go-micro-frame-doc/13-jaeger/02-grpc/config"
)

var (
	DB *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig config.NacosConfig

	RedsyncLock *redsync.Redsync
)