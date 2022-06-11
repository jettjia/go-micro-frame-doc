//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"go-micro-frame-doc/25-wire/06-wire-cleanup/internal/config"
)

//go:generate wire
func InitApp() (*App, func(), error) {
	wire.Build(config.Provider, NewApp)
	return &App{}, func() {}, nil
}
