package global

import (
	ut "github.com/go-playground/universal-translator"

	"go-micro-module/20-temp/web/config"
	"go-micro-module/20-temp/web/proto"
)

var (
	Trans ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient proto.UserClient

	NacosConfig *config.NacosConfig = &config.NacosConfig{}
)
