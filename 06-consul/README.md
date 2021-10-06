# 官方

https://github.com/hashicorp/consul



# 安装和配置

## 安装

```shell
docker run --name consul -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp consul consul agent -dev -client=0.0.0.0


# 访问http
127.0.0.1:8500

# 访问dns [测试]
# consul提供dns功能，可以让我们通过dig命令来测试，dns的端口是8600，命令如下
yum install bind-utils
dig @10.4.7.71 -p 8600 consul.service.consul SRV
```



```
win10 安装
https://www.consul.io/downloads
解压到 E:\dev\dev-win\consul
cmd: consul
```



## consul 的 api 接口

https://www.consul.io/api-docs/agent/service

可以通过上面的 api 去调用 consul的接口



# 案例



## 目录说明

```shell
├── config				# 解析配置
│   └── config.go
├── config-dev.yaml		# 开发配置
├── config-prod.yaml	# 生产配置
├── global				# 全局变量，比如db, logger 等
│   └── global.go
├── handler				# 具体业务
│   ├── base.go
│   └── user.go
├── initialize			# 初始化，比如初始化 db, logger, config等
│   ├── config.go
│   ├── db.go
│   └── logger.go
├── main.go				# 程序入口，初始化 initialize里的包，或者 grpc启动时候注册到 consul
├── model				# model，定义struct
│   └── user.go
├── proto				# proto
│   ├── user.pb.go
│   └── user.proto
├── tests				# 调用 grpc 测试
│   └── user.go	
├── README.md
└── utils				# 工具类
 ├── addr.go			# 生成动态的端口号
 └── register		# 注册中心封装
     └── consul
         └── register.go

```



## 步骤

```shell
1) 创建基础的目录，入下面所示
├── config
├── global	
├── handler
├── initialize
├── model
├── proto
├── README.md
├── tests
└── utils
├── main.go
├── config-dev.yaml
├── config-prod.yaml

2) 逐步实现目录说明中的研发
2.1）定义 config-dev.yaml, config-prod.yaml
2.2) config/config.go 结构体解析 yaml配置
2.3) 初始化 initialize/config, db, logger
2.4) global 定义全局可以访问的 db, config, logger
2.5) 具体的业务，定义model, handler实现业务
2.6) utils工具类，获取服务端口；consul注册中心方法封装
2.7) main.go 中，调用 initialize中的初始化方法，并且启动 grpc的服务，把服务注册到 consul中

3) 在linux控制台启动程序，然后在 consul 控制台中去查看是否注册成功
```



## grpc 接入 consul

下面是整个实现的方法

### proto

```protobuf
syntax = "proto3";
import "google/protobuf/empty.proto";
option go_package = ".;proto";

service User{
    rpc GetUserList(PageInfo) returns (UserListResponse); // 用户列表
    rpc GetUserByMobile(MobileRequest) returns (UserInfoResponse); //通过mobile查询用户
    rpc GetUserById(IdRequest) returns (UserInfoResponse); //通过id查询用户
    rpc CreateUser(CreateUserInfo) returns (UserInfoResponse); // 添加用户
    rpc UpdateUser(UpdateUserInfo) returns (google.protobuf.Empty); // 更新用户
    rpc CheckPassWord(PasswordCheckInfo) returns (CheckResponse); //检查密码
}

message PasswordCheckInfo {
    string password = 1;
    string encryptedPassword = 2;
}

message CheckResponse{
    bool success = 1;
}

message PageInfo {
    uint32 pn = 1;
    uint32 pSize = 2;
}

message MobileRequest{
    string mobile = 1;
}

message IdRequest {
    int32 id = 1;
}

message CreateUserInfo {
    string nickName = 1;
    string passWord = 2;
    string mobile = 3;
}

message UpdateUserInfo {
    int32 id = 1;
    string nickName = 2;
    string gender = 3;
    uint64 birthDay = 4;
}

message UserInfoResponse {
    int32 id = 1;
    string passWord = 2;
    string mobile = 3;
    string nickName = 4;
    uint64 birthDay = 5;
    string gender = 6;
    int32 role = 7;
}

message UserListResponse {
    int32 total = 1;
    repeated UserInfoResponse data = 2;
}
```

```shell
protoc --go_out=plugins=grpc:./ ./user.proto
```



### config-dev.yaml

```yaml
name: 'user-srv'
host: '10.4.7.71'
port: '50051'
tags:
  - 'imooc'
  - 'user'
  - 'srv'

mysql:
  host: '10.4.7.71'
  port: 3307
  user: 'root'
  password: 'root'
  db: 'mxshop_user_srv'

consul:
  host: '10.4.7.71'
  port: 8500
```

config-prod.yaml

```yaml
name: 'user-srv'
host: '10.4.7.71'
tags:
  - 'imooc'
  - 'user'
  - 'srv'

mysql:
  host: '10.4.7.71'
  port: 3307
  user: 'root'
  password: 'root'
  db: 'mxshop_user_srv'

consul:
  host: '10.4.7.71'
  port: 8500
```



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



### initialize/config.go

```go
package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"go-micro-frame-doc/06-consul/global"
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
	configFileName := fmt.Sprintf("06-consul/%s-prod.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("06-consul/%s-dev.yaml", configFilePrefix)
	}

	// 读取文件配置内容
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	// 把内容设置到全局变量的 ServerConfig中
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
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



### initialize/db.go

```go
package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go-micro-frame-doc/06-consul/global"
)

func InitDB() {
	c := global.ServerConfig.MysqlInfo
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		zap.S().Error("db connect err:", err.Error())
		//panic(err.Error())
	}
}

```



### global/global.go

```go
package global

import (
	"gorm.io/gorm"

	"go-micro-frame-doc/06-consul/config"
)

var (
	DB *gorm.DB
	ServerConfig config.ServerConfig
)
```

### model/user.go

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID int32 `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

/*
1. 密文 2. 密文不可反解
	1. 对称加密
	2. 非对称加密
	3. md5 信息摘要算法
	密码如果不可以反解，用户找回密码
*/
type User struct {
	BaseModel
	Mobile string `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string `gorm:"type:varchar(100);not null"`
	NickName string `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender string `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女, male表示男'"`
	Role int `gorm:"column:role;default:1;type:int comment '1表示普通用户, 2表示管理员'"`
}
```



### handler/user.go

```go
package handler

import (
	"context"
	"fmt"

	"go-micro-frame-doc/06-consul/global"
	"go-micro-frame-doc/06-consul/model"
	"go-micro-frame-doc/06-consul/proto"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func ModelToRsponse(user model.User) proto.UserInfoResponse {
	//在grpc的message中字段有默认值，你不能随便赋值nil进去，容易出错
	//这里要搞清， 哪些字段是有默认值
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		PassWord: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
		Mobile:   user.Mobile,
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

// 获取用户列表
func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	fmt.Println("我被调用了")
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	for _, user := range users {
		userInfoRsp := ModelToRsponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}

	return rsp, nil
}

```

handler/base.go

```go
package handler

import "gorm.io/gorm"

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

```



### utils/addr.go

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



### utils/consul

utils/register/consul/register.go

```go
package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

type Registry struct {
	Host string
	Port int
}

type RegistryClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceId string) error
}

func NewRegistryClient(host string, port int) RegistryClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}

func (r *Registry) Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应grpc的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", address, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
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

func (r *Registry) DeRegister(serviceId string) error {
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



### main.go

```go
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"go-micro-frame-doc/06-consul/global"
	"go-micro-frame-doc/06-consul/handler"
	"go-micro-frame-doc/06-consul/initialize"
	"go-micro-frame-doc/06-consul/proto"
	"go-micro-frame-doc/06-consul/utils"
	"go-micro-frame-doc/06-consul/utils/register/consul"
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

	zap.S().Info(global.ServerConfig)

	/////////////////////////////////
	// 启动grpc，并注册到 consul
	server := grpc.NewServer()
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
		zap.S().Info("注销成功:")
	}

}

```



### test

在linux中运行如下命令，然后在 consul后台查看是否已经注册  http://10.4.7.71:8500/

```shell
 go run 06-consul/main.go
```

#### 



