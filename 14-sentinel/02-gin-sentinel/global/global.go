package global

import (
	ut "github.com/go-playground/universal-translator"

	"go-micro-frame-doc/14-sentinel/02-gin-sentinel/config"
	"go-micro-frame-doc/14-sentinel/02-gin-sentinel/proto"
)

var (
	Trans ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient proto.UserClient

	NacosConfig *config.NacosConfig = &config.NacosConfig{}
)
