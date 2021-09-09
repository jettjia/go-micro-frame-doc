package global

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"

	"go-micro-module/20-temp/grpc/config"
)

var (
	DB *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig config.NacosConfig
	EsClient *elastic.Client
	RedsyncLock *redsync.Redsync
)