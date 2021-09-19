# gin调用 grpc微服务

​		这里 grpc 里的服务都是注册在 consul里的，我们通过下面完备的案例，可以利用 gin 从 consul 注册中心里获取服务。

​		这是利用 gin 实现的 微服务的web层，没有把 上一章的微服务直接暴露成 api的方式；当然 grpc 利用 gateway也是可以暴露成接口的。这里这样设计是为了程序后面更好的解耦，也可以更好的扩展。

​		对于 gin 的使用，可以参考上面的关于 gin的章节。



## 目录说明

```shell
├── api			# api 逻辑
├── config		# 配置读取的 struct定义
├── forms		# 表单参数验证
├── global		# 全局访问对象定义
├── initialize	        # 初始化访问定义
├── middlewares		# 中间件
├── models		# db model定义
├── proto		# proto
├── README.md
├── router		# 路由定义
└── utils		# 工具类
├── main.go		# main启动
├── config-debug.yaml	# 配置
├── config-pro.yaml
```



## 开发步骤

```
1) 定义好上诉的文件

2） 根据上面的文件，进行逻辑实现
```



## 具体实现

### proto

这里的proto，拷贝 server里的

### config-debug.yaml

```yaml
name: 'user-web'
port: 8021

user_srv:
  host: '10.4.7.71'
  port: 50051
  name: 'user-srv'

redis:
  host: '127.0.0.1'
  port: 6379

consul:
  host: '10.4.7.71'
  port: 8500
```

config-pro.yaml

```yaml
name: 'user-web'
port: 8021

user_srv:
  host: '10.4.7.71'
  port: 50051
  name: 'user-srv'

redis:
  host: '127.0.0.1'
  port: 6379

consul:
  host: '10.4.7.71'
  port: 8500
```



### 初始化 logger (zap)

```go
package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

```



### 初始化 viper

从 yaml 配置中，解析配置到 struct中

#### config/config.go

```go
package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`
	Port   int    `mapstructure:"port" json:"port"`
	Expire int    `mapstructure:"expire" json:"expire"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name        string        `mapstructure:"name" json:"name"`
	Host        string        `mapstructure:"host" json:"host"`
	Tags        []string      `mapstructure:"tags" json:"tags"`
	Port        int           `mapstructure:"port" json:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	RedisInfo   RedisConfig   `mapstructure:"redis" json:"redis"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul"`
}

```

#### global/global.go

```go
package global

import "go-micro-frame-doc/07-consul-discovery/config"

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
)

```



#### initialize/config.go

```go
package initialize

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"go-micro-frame-doc/07-consul-discovery/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig(){
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user-web/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("user-web/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息: &v", global.ServerConfig)
	fmt.Printf("%V", v.Get("name"))


}
```



### 初始化 gin

#### api/user/user.go

```go
package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserList(ctx *gin.Context) {
	obj := []int{1, 3, 5}
	ctx.JSON(http.StatusOK, obj)
}
```



#### router/user.go

```go
package router

import (
	"github.com/gin-gonic/gin"
	api "go-micro-frame-doc/07-consul-discovery/api/user"
)

func InitUserRouter(Router *gin.RouterGroup){
	UserRouter := Router.Group("user")
	{
		UserRouter.GET("list", api.GetUserList)
	}
}
```



#### initialize/router.go

```go
package initialize

import (
	"github.com/gin-gonic/gin"
	"go-micro-frame-doc/07-consul-discovery/middlewares"
	"go-micro-frame-doc/07-consul-discovery/router"
	"net/http"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"code":http.StatusOK,
			"success":true,
		})
	})

	//配置跨域
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup)

	return Router
}
```



#### main.go

```go
package main

import (
	"fmt"

	"go.uber.org/zap"

	"go-micro-frame-doc/07-consul-discovery/global"
	"go-micro-frame-doc/07-consul-discovery/initialize"
)

func main()  {
	//1. 初始化logger
	initialize.InitLogger()

	//2. 初始化配置文件
	initialize.InitConfig()

	//3. 初始化routers
	Router := initialize.Routers()
	//4. 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}
	zap.S().Debugf("启动服务器, 端口： %d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil{
		zap.S().Panic("启动失败:", err.Error())
	}
}

```

启动项目，浏览器访问：http://127.0.0.1:8021/u/v1/user/list



### 初始化 grpc调用

步骤：

封装一个调用 grpc的方法;

api 接口去调用远程的 grpc server;

main.go 中初始化 sever的连接；



#### global/global.go

```go
package global

import (
	ut "github.com/go-playground/universal-translator"

	"go-micro-frame-doc/07-consul-discovery/config"
	"go-micro-frame-doc/07-consul-discovery/proto"
)

var (
	Trans ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient proto.UserClient
)

```



#### initialize/srv_conn.go

```go
package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"go-micro-frame-doc/07-consul-discovery/global"
	"go-micro-frame-doc/07-consul-discovery/proto"
)

func InitSrvConn(){
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}

```



#### api/user.go

```go
package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go-micro-frame-doc/07-consul-discovery/global"
	reponse "go-micro-frame-doc/07-consul-discovery/global/response"
	"go-micro-frame-doc/07-consul-discovery/proto"
)

func removeTopStruct(fileds map[string]string) map[string]string{
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}


func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg:": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})
			}
			return
		}
	}
}

func HandleValidatorError(c *gin.Context, err error){
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg":err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}


func GetUserList(ctx *gin.Context) {
	// 获取请求参数
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	// grpc远程调用
	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 组装返回结果
	reMap := gin.H{
		"total": rsp.Total,
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		user := reponse.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Birthday: reponse.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		result = append(result, user)
	}

	reMap["data"] = result
	ctx.JSON(http.StatusOK, reMap)
}
```

#### main.go

```go
package main

import (
	"fmt"

	"go.uber.org/zap"

	"go-micro-frame-doc/07-consul-discovery/global"
	"go-micro-frame-doc/07-consul-discovery/initialize"
)

func main()  {
	// 初始化logger
	initialize.InitLogger()

	// 初始化配置文件
	initialize.InitConfig()

	// 初始化routers
	Router := initialize.Routers()

	// 初始化srv的连接
	initialize.InitSrvConn()

	zap.S().Debugf("启动web服务的端口： %d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil{
		zap.S().Panic("启动失败:", err.Error())
	}
}

```



### test

启动 main.go, 访问端口测试:http://127.0.0.1:8021/u/v1/user/list



# gin-web 也注册到 consul

## 思路

> 1， 能够生成端口号
>
> 2，能够 注册到 consul 中的方法
>
> 3，main.go 启动的时候，也同时注册服务到 consul中

## config-debug.yaml

```yaml
name: 'user-web'
host: '10.4.7.71'
port: 8021
tags:
  - 'imooc'
  - 'user'
  - 'web'

env:
  'dev'

user_srv:
  host: '10.4.7.71'
  port: 50051
  name: 'user-srv'

redis:
  host: '127.0.0.1'
  port: 6379

consul:
  host: '10.4.7.71'
  port: 8500
```

config-pro.yaml

```go
name: 'user-web'
host: '10.4.7.71'
port: 8021
tags:
  - 'imooc'
  - 'user'
  - 'web'
env:
  'pro'

user_srv:
  host: '10.4.7.71'
  port: 50051
  name: 'user-srv'

redis:
  host: '127.0.0.1'
  port: 6379

consul:
  host: '10.4.7.71'
  port: 8500
```



### config/config.go

```go
package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`
	Port   int    `mapstructure:"port" json:"port"`
	Expire int    `mapstructure:"expire" json:"expire"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name string   `mapstructure:"name" json:"name"`
	Host string   `mapstructure:"host" json:"host"`
	Tags []string `mapstructure:"tags" json:"tags"`
	Port int      `mapstructure:"port" json:"port"`
	Env  string   `mapstructure:"env" json:"port"`

	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	RedisInfo   RedisConfig   `mapstructure:"redis" json:"redis"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul"`
}

```



## utils

在utils里创建两个方法

1） addr.go里的，生成 端口号的方法

```go
package utils

import (
	"net"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port,  nil
}

```



2) 注册 grpc 服务到 consul的方法

07-consul-discovery\utils\register\consul\register.go

```go
package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

type Registry struct{
	Host string
	Port int
}

type RegistryClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceId string) error
}

func NewRegistryClient(host string, port int) RegistryClient{
	return &Registry{
		Host: host,
		Port: port,
	}
}

func (r *Registry)Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		HTTP: fmt.Sprintf("http://%s:%d/health", address, port),
		Timeout: "5s",
		Interval: "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return nil
}

func (r *Registry)DeRegister(serviceId string) error{
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	err = client.Agent().ServiceDeregister(serviceId)
	return err
}
```



## main.go

```go
package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"go-micro-frame-doc/07-consul-discovery/global"
	"go-micro-frame-doc/07-consul-discovery/initialize"
	"go-micro-frame-doc/07-consul-discovery/utils"
	"go-micro-frame-doc/07-consul-discovery/utils/register/consul"
)

func main() {
	// 初始化logger
	initialize.InitLogger()

	// 初始化配置文件
	initialize.InitConfig()

	// 初始化routers
	Router := initialize.Routers()

	// 初始化srv的连接
	initialize.InitSrvConn()

	/////////////////////////////////////////////
	// 随机生成 port, 如果是本地开发环境端口号固定，线上环境启动获取端口号
	if global.ServerConfig.Env != "dev" {
		global.ServerConfig.Port, _ = utils.GetFreePort()
	}

	//注册服务健康检查
	server := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 注册服务到 consul中
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := register_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}
	/////////////////////////////////////////////

	// 启动 web服务
	zap.S().Debugf("启动web服务的端口： %d", global.ServerConfig.Port)
	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	} else {
		zap.S().Info("注销成功")
	}
}

```



## test

启动 web服务，然后在 consul中确认是否注册成功

```shell
 go run 07-consul-discovery/main.go
```

http://10.4.7.71:8500/ui/dc1/services



访问接口  http://10.4.7.71:8021/u/v1/user/list

