# 官方

https://nacos.io/zh-cn/



# 安装

```shell
docker run --name nacos-standalone -e MODE=standalone -e JVM_XMS=512m -e JVM_XMX=512m -e JVM_XMN=256m -p 8848:8848 -d nacos/nacos-server:latest
```

访问： http://10.4.7.71:8848/nacos/index.html

密码：nacos/nacos



# nacos概念简述

* 命名空间：隔离配置；比如用户，商品

* 配置集：一个完整的配置文件

* 组： 一组配置文件；比如本地、开发、测试、生产

 <img src="images/image-20210826184746965.png" alt="image-20210826184746965" style="zoom: 67%;" />



# nacos api 使用

## 后台配置-案例

1） 添加命名空间

 ![image-20210826184908558](images/image-20210826184908558.png)



2）配置管理->新建配置

 <img src="images/image-20210826185001110.png" alt="image-20210826185001110" style="zoom:50%;" />

```json
{
  "name": "user-srv",
  "host": "10.4.7.71",
  "port": "50051",
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
  "consul": {
    "host": "10.4.7.71",
    "port": 8500
  }
}
```



## golang 读取 ncaos配置

文档：https://github.com/nacos-group/nacos-sdk-go

### config/config.go

```go
package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Host string   `mapstructure:"host"`
	Port uint64   `mapstructure:"port"`
	Tags []string `mapstructure:"tags" json:"tags"`
	Name string   `mapstructure:"name" json:"name"`

	MysqlInfo  MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`
}

```



### main.go

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"go-micro-module/08-nacos/01-golang-nacos-api/config"
)

func main() {
	sc := []constant.ServerConfig{
		{
			IpAddr: "10.4.7.71",
			Port:   8848,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         "e16c905e-bd10-4f62-98cb-acc78d967b59", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "user-dev-local.json",
		Group:  "dev"})

	if err != nil {
		panic(err)
	}
	//fmt.Println(content) //字符串 - yaml
	serverConfig := config.ServerConfig{}
	//想要将一个json字符串转换成struct，需要去设置这个struct的tag
	json.Unmarshal([]byte(content), &serverConfig)
	fmt.Println(serverConfig)

	//err = configClient.ListenConfig(vo.ConfigParam{
	//	DataId: "user-web.json",
	//	Group:  "dev",
	//	OnChange: func(namespace, group, dataId, data string) {
	//		fmt.Println("配置文件变化")
	//		fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
	//	},
	//})
	//time.Sleep(3000 * time.Second)

}

```



# grpc 集成 nacos



