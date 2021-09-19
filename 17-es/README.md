# es



## 官方

github.com/olivere/elastic/v7



## 安装

```
参考：https://github.com/deviantony/docker-elk
```

http://10.4.7.71:5601/

elastic/hZksYkpkcweABXu68qh0



# 实践

这里我们拷贝 13-jaeger/02-grpc 来改造



## 思路

> 1）nacos 增加 es的配置
>
> 2） config/config.go

## nacos 配置

user-srv.json

```json
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
  },
  "jaeger": {
    "host": "10.4.7.71",
    "port": 6831,
    "name" : "mxshopall"
  },
  "es" : {
    "host": "10.4.7.71",
    "port": 9200,
    "user" : "elastic",
    "password" : "hZksYkpkcweABXu68qh0"
  }
}
```



## config/config.go 

```go
type EsConfig struct {
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

	MysqlInfo   MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	ConsulInfo  ConsulConfig `mapstructure:"consul" json:"consul"`
	RedisConfig ConsulConfig `mapstructure:"redis" json:"redis"`
	JaegerInfo  JaegerConfig `mapstructure:"consul" json:"jaeger"`
	EsInfo      EsConfig     `mapstructure:"es" json:"es"`
}
```



## global/global.go

```go
EsClient *elastic.Client
```



## initialize/es.go

```go
package initialize

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic/v7"

	"go-micro-frame-doc/17-es/02-grpc/global"
	"go-micro-frame-doc/17-es/02-grpc/model"
)

func InitEs() {
	//初始化连接
	host := fmt.Sprintf("http://%s:%d", global.ServerConfig.EsInfo.Host, global.ServerConfig.EsInfo.Port)
	logger := log.New(os.Stdout, "mxshop", log.LstdFlags)
	var err error
	global.EsClient, err = elastic.NewClient(
		elastic.SetURL(host), elastic.SetSniff(false),
		elastic.SetBasicAuth(global.ServerConfig.EsInfo.User, global.ServerConfig.EsInfo.Password),
		elastic.SetTraceLog(logger),
	)
	if err != nil {
		panic(err)
	}

	//新建mapping和index
	exists, err := global.EsClient.IndexExists(model.EsUser{}.GetIndexName()).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exists {
		_, err = global.EsClient.CreateIndex(model.EsUser{}.GetIndexName()).BodyString(model.EsUser{}.GetMapping()).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
}

```



## model/es_user.go

```go
package model

type EsUser struct {
	ID       int32  `json:"id"`
	Mobile   string `json:"mobile"`
	NickName string `json:"nickname"`
	Gender   int32  `json:"gender"`
}

func (EsUser) GetIndexName() string {
	return "user"
}

func (EsUser) GetMapping() string {
	userMapping := `
	{
		"mappings" : {
			"properties" : {
				"Mobile" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"Mobile" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"NickName" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"gender" : {
					"type" : "integer"
				}
			}
		}
	}`
	return userMapping
}

```



## main.go

```go
	// 初始化es
	initialize.InitEs()
```



## test





