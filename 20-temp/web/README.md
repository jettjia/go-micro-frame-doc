# 目录说明

```
├── amqp                // amqp
│   ├── consumer
│   └── producer
├── api                 // api逻辑
│   ├── ceshi
│   └── user
├── config              // 配置读取
│   └── config.go
├── forms               // 表单验证
│   └── user.go
├── global              // 全局global定义
│   ├── global.go
│   └── response
├── initialize          // 初始化
│   ├── config.go
│   ├── router.go
│   ├── srv_conn.go
│   └── zap.go
├── middlewares         // 中间件
│   ├── cors.go
│   └── tracing.go
├── models              // 定义请求，返回
├── proto               // proto
│   ├── user.pb.go
│   └── user.proto
├── router              // 路由
│   ├── router.go
│   └── test.go
└── utils               // 工具
    ├── addr.go
    ├── amqpRabbit
    ├── otgrpc
    └── register
├── main.go             // 入口
├── config-debug.yaml   // 配置
├── config-pro.yaml
├── nacos-grpc.json.example // nacos配置模板
├── README.md
```



# 功能说明

​		这里是实现grpc的具体接口，也就是grpc的 client。有种设计，是可以直接把grpc的 gateway暴露成端口的，这里用 web端来暴露端口的方式。两种设计都是可行的，这里为了能更好的扩展程序，把web进行了单独的处理。并且把 web成也注入到了 consul 注册中心中，然后用 Kong网关来配置负载均衡，web 调用grpc 又是负载均衡，实现了微服务的调用。

# 使用说明

​		

1)	使用的时候，要先把如下工具打开：

​		mysql

​		redis

​		es

​		rabbitmq

​		nacos	http://10.4.7.71:8848/nacos/#/login

​		consul	http://10.4.7.71:8500/

​		jaeger	http://10.4.7.71:16686/

​		konga	http://10.4.7.71:1337/#!/login

2）运行 grpc服务

3)  运行 main.go

```
go run 20-temp/grpc/main.go
```



3) 开发步骤

