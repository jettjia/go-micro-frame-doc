package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"

	"go-micro-frame-doc/25-wire/05-wire-interface/internal/config"
)

// var Provider = wire.NewSet(NewDb)

// 这里我们加入了Dao, 并且绑定了 IDao 和Dao
// 将 New 方法声明为 Provider,表示New方法可以创建一个被别人依赖的对象
var Provider = wire.NewSet(NewDb, NewDao, wire.Bind(new(IDao), new(*Dao))) // 这里将接口和实现进行绑定

func NewDb(cfg *config.Config) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", cfg.Database.Dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
