# 链路追踪

## 介绍

| 种类       | 语言   | 网站                          |
| ---------- | ------ | ----------------------------- |
| Zipkin     | java   | https://zipkin.io/            |
| Skywalking | java   |                               |
| Jaeger     | golang | https://www.jaegertracing.io/ |



## opentracing原理

https://www.yuque.com/baxiang/ms/qciuaq

# jaeger

官方：https://github.com/jaegertracing/jaeger

## jaeger 安装

```shell
docker run \
--rm \
--name jaeger \
-p 6831:6831/udp \
-p 16686:16686 \
-p 16685:16685 \
jaegertracing/all-in-one:latest

```

```
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:latest
```

浏览器访问：http://10.4.7.71:16686/search



## jaeger-client

### 发送单个span

```go
package main

import (
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func main() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "10.4.7.71:6831",
		},
		ServiceName: "mxshop",
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)


	span := opentracing.StartSpan("go-grpc-web")
	time.Sleep(time.Second)
	defer span.Finish()
}

```

运行 main.go，然后在 http://10.4.7.71:16686/ 查看是否有数据



### 发送嵌套 span

```go
package main

import (
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/uber/jaeger-client-go"

	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func main() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "10.4.7.71:6831",
		},
		ServiceName: "mxshop",
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	defer closer.Close()
	if err != nil {
		panic(err)
	}

	parentSpan := tracer.StartSpan("main")

	span := tracer.StartSpan("funcA", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Millisecond*500)
	span.Finish()

	span2 := tracer.StartSpan("funcB", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Millisecond*1000)
	span2.Finish()

	parentSpan.Finish()

}

```



### grpc-test

#### proto

```
syntax = "proto3";
option go_package = ".;proto";
service Greeter {
    rpc SayHello (HelloRequest) returns (HelloReply);
}
message HelloRequest {
    string name = 1;
}

message HelloReply {
    string message = 1;
}
```

#### otgrpc 拷贝到 项目下

https://github.com/grpc-ecosystem/grpc-opentracing

#### client.go

```go
package main

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go-micro-frame-doc/13-jaeger/01-jaeger-test/03-grpc-test/otgrpc"
	"google.golang.org/grpc"

	"go-micro-frame-doc/13-jaeger/01-jaeger-test/03-grpc-test/proto"
)

func main() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "10.4.7.71:6831",
		},
		ServiceName: "mxshop-test",
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer) // 设置成全局
	defer closer.Close()

	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &proto.HelloRequest{Name: "bobby"})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)
}

```

运行 client.go，然后 jaeger后台查看



#### server.go

```go
package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"net"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go-micro-frame-doc/13-jaeger/01-jaeger-test/03-grpc-test/otgrpc"

	"go-micro-frame-doc/13-jaeger/01-jaeger-test/03-grpc-test/proto"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	// 假设这里去调用了其他三方
	parentSpan := opentracing.SpanFromContext(ctx)
	shopCartSpan := opentracing.GlobalTracer().StartSpan("select_shopcart", opentracing.ChildOf(parentSpan.Context()))
	fmt.Println("我调用了其他三方的接口")
	shopCartSpan.Finish()

	return &proto.HelloReply{
		Message: "hello, " + req.Name,
	}, nil
}

func main() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "10.4.7.71:6831",
		},
		ServiceName: "mxshop-test",
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	g := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	proto.RegisterGreeterServer(g, &Server{})
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	err = g.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}

```



# grpc 集成 jaeger

## 说明

这里是 grpc，也就是服务端 接入 jaeger，我们这里使用的是 10-distributed-lock 这一章已经完善的代码进行改造

实现代码参考 02-grpc





## 实践

### otgrpc拷贝到 utils下



### nacos 配置

user-srv.json

```
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
    "name" : "mxshop"
  } 
}
```



### config/config.go

```go
type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
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
}
```



### initialize/jaeger.go

```go
package initialize

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"go-micro-frame-doc/13-jaeger/02-grpc/global"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func InitJaeger(){
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", global.ServerConfig.JaegerInfo.Host, global.ServerConfig.JaegerInfo.Port),
		},
		ServiceName: global.ServerConfig.JaegerInfo.Name,
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}

	opentracing.SetGlobalTracer(tracer)

	defer closer.Close()
}
```



### main.go

```go
package main

import (
	"flag"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"go-micro-frame-doc/13-jaeger/03-grpc/utils/otgrpc"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"go-micro-frame-doc/13-jaeger/03-grpc/global"
	"go-micro-frame-doc/13-jaeger/03-grpc/handler"
	"go-micro-frame-doc/13-jaeger/03-grpc/initialize"
	"go-micro-frame-doc/13-jaeger/03-grpc/proto"
	"go-micro-frame-doc/13-jaeger/03-grpc/utils"
	"go-micro-frame-doc/13-jaeger/03-grpc/utils/register/consul"
)

func main() {
	// 判断是否生成 随机的 微服务端口号
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口")
	flag.Parse()
	zap.S().Info("ip: ", *IP)
	if *Port == 0{
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info("port: ", *Port)

	// 初始化 logger
	initialize.InitLogger()

	//初始化配置文件
	initialize.InitConfig()

	// 初始化db
	initialize.InitDB()

	// 初始化 redsyncLock
	initialize.InitRedsyncLock()

	// 初始化jaeger
	tracer := initialize.InitJaeger()

	zap.S().Info(global.ServerConfig)

	/////////////////////////////////
	// 启动grpc，并注册到 consul，并且使用 jaeger
	opentracing.SetGlobalTracer(tracer)
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))

	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//启动服务
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//服务注册
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err = register_client.Register(global.ServerConfig.Host, *Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}
	zap.S().Debugf("启动服务器, 端口： %d", *Port)

	/////////////////////////////////

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	}else{
		zap.S().Info("注销成功")
	}

}

```



### handler/user.go

```go
// 获取用户列表
func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	fmt.Println("我被调用了")

	fmt.Println("输出，db连接是否变了", global.ServerConfig.MysqlInfo)

	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	// 这里模拟，有调用第三方的服务
	parentSpan := opentracing.SpanFromContext(ctx)
	getUserSpan := opentracing.GlobalTracer().StartSpan("get_user", opentracing.ChildOf(parentSpan.Context()))
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)
	getUserSpan.Finish()

	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	for _, user := range users {
		userInfoRsp := ModelToRsponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}

	return rsp, nil
}
```



# gin 集成 jaeger

## 说明

这里我们使用 07-consul-discovery 章节的代码来实现。 实现代码参考：02-web

## 实践

### gin 接入 nacos

#### 思路

> 1）在 nacos 后台 增加 配置
>
> 2）config-debug.yaml 连接 nacos 配置
>
> 3）global/global.go 定义全局变量 NacosConfig
>
> 4）config/config.go 定义 NacosConfig 结构体
>
> 5）initialize/config.go 这里之前是从本地 yaml读取的配置，现在改为从 nacos 读取配置信息

#### nacos 后台 配置

```json
{
  "name": "user-web",
  "host": "10.4.7.71",
  "port": 8021,
  "tags": [
    "imooc",
    "user",
    "web"
  ],
  "env": "dev",
  "user_srv": {
    "host": "10.4.7.71",
    "port": 50051,
    "name": "user-srv"
  },
  "redis": {
    "host": "127.0.0.1",
    "port": 6379
  },
  "consul": {
    "host": "10.4.7.71",
    "port": 8500
  }
}
```

#### config-debug.yaml

```yaml
host: '10.4.7.71'
port: 8848
namespace: 'e16c905e-bd10-4f62-98cb-acc78d967b59'
user: 'nacos'
password: 'nacos'
dataid: 'user-web.json'
group: 'dev'
```

config-pro.yaml

```yaml
host: '10.4.7.71'
port: 8848
namespace: 'e16c905e-bd10-4f62-98cb-acc78d967b59'
user: 'nacos'
password: 'nacos'
dataid: 'user-web.json'
group: 'pro'
```



#### global/global.go

```go
package global

import (
	ut "github.com/go-playground/universal-translator"

	"go-micro-frame-doc/13-jaeger/02-web/config"
	"go-micro-frame-doc/13-jaeger/02-web/proto"
)

var (
	Trans ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	UserSrvClient proto.UserClient

	NacosConfig *config.NacosConfig = &config.NacosConfig{}
)

```



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
	Name string   `mapstructure:"name" json:"name"`
	Host string   `mapstructure:"host" json:"host"`
	Tags []string `mapstructure:"tags" json:"tags"`
	Port int      `mapstructure:"port" json:"port"`
	Env  string   `mapstructure:"env" json:"port"`

	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	RedisInfo   RedisConfig   `mapstructure:"redis" json:"redis"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host" `
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
```



#### initialize/config.go

```go
package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"go-micro-frame-doc/13-jaeger/02-web/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig(){
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("13-jaeger/02-web/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("13-jaeger/02-web/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	////////////////////////////////
	// 从 nacos 中读取配置

	// 把内容设置到全局变量的 NacosConfig
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}

	//从nacos中读取配置信息
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
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
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})

	if err != nil {
		panic(err)
	}
	//想要将一个json字符串转换成struct，需要去设置这个struct的tag
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败： %s", err.Error())
	}

	// 监听 nacos 配置变化
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
		OnChange: func(namespace, group, dataId, data string) {
			// todo, 这里能获取到 nacos 变化的值，但是没能赋值到 ServerConfig中
			fmt.Println("nacos中的配置", data)
			// 这里输出的格式是：  { "name": "user-srv", "host": "10.4.7.71" }
			err = json.Unmarshal([]byte(data), &global.ServerConfig)
			if err != nil {
				zap.S().Errorf("配置中心文件改变后，解析 Json失败")
			}
			zap.S().Infof("nacos 改变后配置：", &global.ServerConfig)
		},
	})
	if err != nil {
		zap.S().Errorf("配置中心文件变化，解析失败!")
	}

	zap.S().Infof("从nacos读取到的全部配置如下：", &global.ServerConfig)
	////////////////////////////////
}
```



#### test

```shell
# 运行服务，然后查看 consul能否注册成功
go run 13-jaeger/02-web/main.go
```



### gin 接入 jaeger

#### 思路

> 这里是用拦截器的方式，注入 jaeger

#### otgrpc包拷贝到 utils下



#### nacos 增加 jaeger配置

```json
{
  "name": "user-web",
  "host": "10.4.7.71",
  "port": 8021,
  "tags": [
    "imooc",
    "user",
    "web"
  ],
  "env": "dev",
  "user_srv": {
    "host": "10.4.7.71",
    "port": 50051,
    "name": "user-srv"
  },
  "redis": {
    "host": "127.0.0.1",
    "port": 6379
  },
  "consul": {
    "host": "10.4.7.71",
    "port": 8500
  },
  "jaeger": {
    "host": "10.4.7.71",
    "port": 6831,
    "name" : "mxshop"
  } 
}
```



#### config/config.go

```go
type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
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
	JaegerInfo  JaegerConfig   `mapstructure:"consul" json:"jaeger"`
}
```

#### middlewares/tracing

```go
package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	"go-micro-frame-doc/13-jaeger/02-web/global"
)

func Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: fmt.Sprintf("%s:%d", global.ServerConfig.JaegerInfo.Host, global.ServerConfig.JaegerInfo.Port),
			},
			ServiceName: global.ServerConfig.JaegerInfo.Name,
		}

		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			panic(err)
		}
		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()

		startSpan := tracer.StartSpan(ctx.Request.URL.Path)
		defer startSpan.Finish()

		ctx.Set("tracer", tracer)
		ctx.Set("parentSpan", startSpan)
		ctx.Next()
	}
}

```



#### initialize/srv_conn.go

```go
package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/opentracing/opentracing-go"
	"go-micro-frame-doc/13-jaeger/02-web/utils/otgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"go-micro-frame-doc/13-jaeger/02-web/global"
	"go-micro-frame-doc/13-jaeger/02-web/proto"
)

func InitSrvConn(){
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())), // 注入 jaeger
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}

```



#### router/router.go

在路由上，加上 middlewares.Trace ；添加链路追踪

```go
package router

import (
	"github.com/gin-gonic/gin"
	api "go-micro-frame-doc/13-jaeger/02-web/api/user"
	"go-micro-frame-doc/13-jaeger/02-web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup){
    // 使用 jaeger
	UserRouter := Router.Group("user").Use(middlewares.Trace())
	{
		UserRouter.GET("list", api.GetUserList)
	}
}
```



#### api/user.go

```go
func GetUserList(ctx *gin.Context) {

	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	// grpc远程调用，传递ginContext
	rsp, err := global.UserSrvClient.GetUserList(context.WithValue(context.Background(), "ginContext", ctx), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})

	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}


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



#### test

```
# 我们这里是web端，可以先启动 一个服务端来测试
go run 13-jaeger/02-grpc/main.go

# 启动 web 端
go run 13-jaeger/02-web/main.go
# 访问列表页，看是否有内容写入到 jaeger后台
http://10.4.7.71:34299/u/v1/user/list

http://10.4.7.71:16686/search
```





