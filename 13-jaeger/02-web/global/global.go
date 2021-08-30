package global

import (
	ut "github.com/go-playground/universal-translator"

	"go-micro-module/13-jaeger/02-web/config"
	"go-micro-module/13-jaeger/02-web/proto"
)

var (
	Trans ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient proto.UserClient

	NacosConfig *config.NacosConfig = &config.NacosConfig{}
)
