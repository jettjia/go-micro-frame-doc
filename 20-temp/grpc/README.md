# 目录说明

```
├── amqp		// amqp
│   ├── consumer
│   └── producer
├── config              // 配置读取struct 
│   └── config.go
├── global              // 全局变量定义，比如db等
│   └── global.go
├── handler             // 逻辑处理 
│   ├── base.go
│   └── user.go
├── initialize          // 初始化
│   ├── config.go
│   ├── db.go
│   ├── es.go
│   ├── jaeger.go
│   ├── logger.go
│   └── redsyncLock.go
├── model               // 数据库模型定义
│   ├── es_user.go
│   └── user.go
├── proto               // proto
│   ├── user.pb.go
│   └── user.proto
├── tests               // test
│   └── user.go
└── utils               // 工具
    ├── addr.go
    ├── amqpRabbit
    ├── otgrpc
    └── register
├── config-dev.yaml     // nacos配置
├── config-prod.yaml
├── main.go             // 启动
├── nacos-grpc.json.example // nacos配置模板
├── README.md
```



# 功能说明

## 功能模块

​		这里集成了非常常用的功能模块，比如mysql, redis, rabbitmq, elasticsearch；如果需要额外的扩展，参考如何扩展功能模块章节的步骤。当然开发中，微服务非常常见的如：配置中心、注册/发现中心、链路追踪、网关等也集成在案例中。



## 如何扩展功能模块

​		扩展功能模块在开发中是非常常见的，我们只需要了解了目录说明章节中，对于开发中的目录规划。然后针对性的改造我们的项目代码就可以了。这里我以扩展 redis为例的具体步骤提供参考思路。

1）nacos 配置 

2）config 定义 redis的 struct

3) initialize 初始化 （如果需要定义初始化的话）

4) global 定义全局访问变量（如果需要）

5）main.go 完成初始化

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



2)  运行 main.go

```
go run 20-temp/grpc/main.go
```



3) 开发步骤

1）proto

2) handler

3) test
