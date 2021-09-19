package global

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"

	"go-micro-frame-doc/17-es/02-grpc/config"
)

var (
	DB *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig config.NacosConfig

	EsClient *elastic.Client

	RedsyncLock *redsync.Redsync
)