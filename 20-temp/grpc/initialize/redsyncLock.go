package initialize

import (
	"fmt"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

	"go-micro-module/20-temp/grpc/global"
)

func InitRedsyncLock() {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisConfig.Host, global.ServerConfig.RedisConfig.Port),
		PoolSize: 5,
		MinIdleConns: 10,
	})

	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	global.RedsyncLock = redsync.New(pool)
}

