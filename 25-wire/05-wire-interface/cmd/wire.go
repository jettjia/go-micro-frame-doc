//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	
	"go-micro-frame-doc/25-wire/05-wire-interface/internal/config"
	"go-micro-frame-doc/25-wire/05-wire-interface/internal/db"
)

//go:generate wire
func InitApp() (*App, error) {
	wire.Build(config.Provider, db.Provider, NewApp)
	return &App{}, nil
}
