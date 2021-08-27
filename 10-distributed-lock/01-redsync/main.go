package main

import (
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"time"
)

func main() {
	//这里的变量哪些可以放到global中， redis的配置是否应该在nacos中
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "10.4.7.71:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	// Obtain a new mutex by using the same name for all instances wanting the
	// same lock.
	mutexname := "goods_421" //421号商品

	mutex := rs.NewMutex(mutexname)

	fmt.Println("开始获取锁")
	if err := mutex.Lock(); err != nil {
		panic(err)
	}

	fmt.Println("获取锁成功")

	fmt.Println("处理实际的业务---start")
	time.Sleep(time.Second * 7) // 模拟业务处理，需要消耗7S;
	fmt.Println("处理实际的业务---end")

	fmt.Println("开始释放锁")
	if ok, err := mutex.Unlock(); !ok || err != nil {
		panic("unlock failed")
	}
	fmt.Println("释放锁成功")

}
