package global

import (
	ut "github.com/go-playground/universal-translator"

	"go-micro-frame-doc/07-consul-discovery/config"
	"go-micro-frame-doc/07-consul-discovery/proto"
)

var (
	Trans ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient proto.UserClient
)
