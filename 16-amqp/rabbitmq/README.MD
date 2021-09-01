# 参考

https://github.com/valinurovam/garagemq



# 安装



```shell
docker run -d --name rabbitmq3.7.7 -p 5672:5672 -p 15672:15672 -v /docker-data/rabbitmq:/var/lib/rabbitmq  -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=123456 rabbitmq:3.7.7-management
```

后台：http://10.4.7.71:15672

admin, 123456

# gin 集成 amqp

## 思路

> 这里用的是13-jaeger/02-web 章节进行的改造

## nacos配置



## config/config.go



## initialize/rabbit.go



## main.go



## amqp/producer



## amqp/consumer

