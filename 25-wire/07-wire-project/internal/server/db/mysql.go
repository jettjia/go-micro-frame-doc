package db

import (
	"context"
	"database/sql"
	"github.com/google/wire"

	_ "github.com/go-sql-driver/mysql"

	"go-micro-frame-doc/25-wire/07-wire-project/internal/server/config"
)

var Provider = wire.NewSet(NewDb)

func NewDb(ctx context.Context, cfg *config.Config) (db *sql.DB, cleanup func(), err error) {
	db, err = sql.Open("mysql", cfg.Database.Dsn)
	if err != nil {
		return
	}

	if err = db.Ping(); err != nil {
		return
	}
	return db, func() {
		db.Close()
	}, nil
}
