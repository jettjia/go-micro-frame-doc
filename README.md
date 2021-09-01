# go-micro-module

这里将使用下面的工具，分模块的方式，实现微服务。
每一个模块就像打怪一样，逐步升级，最终实现 golang 的微服务整体解决方案。
后面会增加一个完备的模板案例，来方便快速开发

```
gorm		【orm】
gin		【web服务】
grpc、proto	【rpc微服务】
zap 		【日志】
viper		【配置读取】
consul 		【服务注册和发现】
nacos		【配置中心】
grpc自带	【负载均衡】
es		【搜索】
分布式锁	        【redis实现】
幂等性		【grpc重试，分布式锁处理，token处理等】
jaeger		【链路追踪】
sentinel	【限流、熔断、降级】
kong		【网关】
amqp            【amqp，消息队列，比如：rabbitmq】
cron            【分布式定时任务】
分布式事务	【rocketmq，事务消息方式；方式2：seata-golang】
```

