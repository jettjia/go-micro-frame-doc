# golang-all

这里将使用下面的工具，分模块的方式，实现微服务。

```
gorm		【orm】
gin			【web服务】
grpc、proto	【rpc微服务】
zap 		【日志】
viper		【配置读取】
consul 		【服务注册和发现】
nacos		【配置中心】
grpc自带		【负载均衡】
es			【搜索】
分布式锁	【redis实现】
分布式事务	【rocketmq，事务消息方式；方式2：seata-golang】
幂等性		【grpc重试，分布式锁处理，token处理等】
jaeger		【链路追踪】
sentinel	【限流、熔断、降级】
kong		【网关】
```

