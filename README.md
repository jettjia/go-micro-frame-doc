# go-micro-frame-doc

框架：https://github.com/jettjia/go-micro-frame

go-micro-frame 框架的文档。下面的这些功能模块，大部分都已经实现。文档是循序渐进的进行迭代，可以很好的学习如何搭建微服务，也可以完全自主的替换对应的模块。这里暂缺的是分布式事务部分，后续将完善分布式事务的两种实现方式（分布式事务有多种实现，这里只演示介绍中的两种）。

```
gorm		【orm】
gin		【web服务】
grpc、proto	【rpc微服务】
zap 		【日志】
viper		【配置读取】
consul 		【服务注册和发现】
nacos		【配置中心，服务注册和发现】
grpc-lb 	【负载均衡】
es		【搜索】
分布式锁	        【redis实现】
幂等性		【grpc重试，分布式锁处理，token处理等】
jaeger		【链路追踪】
sentinel	【限流、熔断、降级】
kong		【网关】
amqp            【amqp，消息队列，比如：rabbitmq】
cron            【分布式定时任务;go:go-cron,java:xxl-job】
分布式事务	【方式1：rocketmq，事务消息方式；方式2：seata-golang】
分布式mysql	【go: gaea分库分表; java: shardingsphere-proxy】
```

