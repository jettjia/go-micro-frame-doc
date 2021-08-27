# 安装 redis

生产中，一般redis 搭建的是集群，可以参考 云原生中的介绍，有搭建 redis 集群。

搭建的集群中，节点数最好是奇数个。这里直接用一个节点来测试了

```shell
docker run -itd --name redis -p 6379:6379 redis
```



# 实现分布式锁-redsync

https://github.com/go-redsync/redsync/blob/master/examples

main.go

```go
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

```



# grpc-集成 redsync

## 项目准备

这里用08 章节中的 grpc 微服务，可以直接拷贝即可



## 实践

### nacos 增加 redis 的配置

```go
{
  "name": "user-srv",
  "host": "10.4.7.71",
  "port": 50051,
  "tags": [
    "imooc",
    "user",
    "srv"
  ],
  "mysql": {
    "host": "10.4.7.71",
    "port": 3307,
    "user": "root",
    "password": "root",
    "db": "mxshop_user_srv"
  },
  "redis": {
    "host": "10.4.7.71",
    "port": 6379,
    "user": "",
    "password": ""
  },
  "consul": {
    "host": "10.4.7.71",
    "port": 8500
  }
}
```



### config/config.go 定义 redis 的 struct

```go
type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ServerConfig struct {
	Name string   `mapstructure:"name" json:"name"`
	Host string   `mapstructure:"host" json:"host"`
	Port uint64   `mapstructure:"port" json:"port"`
	Tags []string `mapstructure:"tags" json:"tags"`

	MysqlInfo  MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`
	RedisConfig ConsulConfig `mapstructure:"redis" json:"redis"`
}
```

### initialize/redsyncLock.go

```go

```



### global/global.go

```

```



###  main.go

### test

tests/user.go

```go

```



```
1) 运行项目
2） tests/user.go中测试
```



