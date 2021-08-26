package global

import (
	ut "github.com/go-playground/universal-translator"

	"go-micro-module/06-consul-discovery/config"
	"go-micro-module/06-consul-discovery/proto"
)

var (
	Trans ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient proto.UserClient
)
