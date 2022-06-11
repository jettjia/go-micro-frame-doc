//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"go-micro-frame-doc/25-wire/07-wire-project/internal/data"
	"go-micro-frame-doc/25-wire/07-wire-project/internal/server/config"
	"go-micro-frame-doc/25-wire/07-wire-project/internal/server/db"
)

//go:generate wire
func InitApp(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(config.Provider, db.Provider, data.OrderSet, NewApp))
}
