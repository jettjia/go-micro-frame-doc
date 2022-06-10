//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go-micro-frame-doc/25-wire/02-wire-demo/internal/config"
	"go-micro-frame-doc/25-wire/02-wire-demo/internal/db"
)

//go:generate wire
func InitApp() (*App, error) {
	// 写法1（参考Kratos框架写法）
	//panic(wire.Build(config.Provider, db.Provider, NewApp))

	// 写法2（参考wire官方写法）
	wire.Build(config.Provider, db.Provider, NewApp)
	return &App{}, nil
}
