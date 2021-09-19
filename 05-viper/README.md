# 官方

github.com/spf13/viper

# 案例

## 步骤说明

> 1、实际开发中，会把配置文件都放到 config-X.yaml 中进行管理
>
> 2、config 把配置文件中的字段，转成对应的 struct 结构体
>
> 3、global 定义全局变量，方便去调用 config 中的内容
>
> 4、initialize/config.go 定义初始化获取配置的方法
>
> 5、main.go 中 初始化调用配置，并且输出测试的数据

## config-dev.yaml

```yaml
name: 'user-web'
port: 8021
user_srv:
  host: '127.0.0.1'
  port: 50051
```

## config/config.go

```go
package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name        string        `mapstructure:"name" json:"name"`
	Port        int           `mapstructure:"port" json:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
}
```

## global/global.go

```go
package global

import (
	"go-micro-frame-doc/05-viper/config"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
```

## initialize/config.go

```go
package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"go-micro-frame-doc/05-viper/global"
	"go.uber.org/zap"

	"github.com/spf13/viper"
)

// 读取环境变量的配置
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("05-viper/%s-prod.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("05-viper/%s-dev.yaml", configFilePrefix)
	}

	// 读取文件配置内容
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	// 把内容设置到全局变量的 ServerConfig中
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	//zap.S().Infof("配置信息: &v", global.ServerConfig)

	//viper的功能 - 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置信息发生了变化: &v", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&global.ServerConfig)
	})
}
```



## main.go

```go
package main

import (
	"fmt"
	"go-micro-frame-doc/05-viper/global"
	"go-micro-frame-doc/05-viper/initialize"
)

func main() {
	// 初始化 logger
	initialize.InitLogger()

	//初始化配置文件
	initialize.InitConfig()

	// 输出配置文件内容
	fmt.Println(global.ServerConfig.Name)
	fmt.Println(global.ServerConfig.Port)
	fmt.Println(global.ServerConfig.UserSrvInfo)
}

```

