package global

import (
	ut "github.com/go-playground/universal-translator"

	"go-micro-module/16-amqp/rabbitmq/01-gin-web/config"
	"go-micro-module/16-amqp/rabbitmq/01-gin-web/proto"
)

var (
	Trans ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient proto.UserClient

	NacosConfig *config.NacosConfig = &config.NacosConfig{}
)
