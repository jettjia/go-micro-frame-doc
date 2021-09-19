package global

import (
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"

	"go-micro-frame-doc/10-distributed-lock/02-grpc-redsync/config"
)

var (
	DB *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig config.NacosConfig

	RedsyncLock *redsync.Redsync
)