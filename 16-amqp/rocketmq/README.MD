# roccketmq

## 官方

http://rocketmq.apache.org/



## 安装

```shell
sudo curl -L "https://github.com/docker/compose/releases/download/1.26.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

sudo chmod +x /usr/local/bin/docker-compose
```

install/conf/broker.conf

```shell
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.


# 所属集群名字
brokerClusterName=DefaultCluster

# broker 名字，注意此处不同的配置文件填写的不一样，如果在 broker-a.properties 使用: broker-a,
# 在 broker-b.properties 使用: broker-b
brokerName=broker-a

# 0 表示 Master，> 0 表示 Slave
brokerId=0

# nameServer地址，分号分割
# namesrvAddr=rocketmq-nameserver1:9876;rocketmq-nameserver2:9876

# 启动IP,如果 docker 报 com.alibaba.rocketmq.remoting.exception.RemotingConnectException: connect to <192.168.0.120:10909> failed
# 解决方式1 加上一句 producer.setVipChannelEnabled(false);，解决方式2 brokerIP1 设置宿主机IP，不要使用docker 内部IP
brokerIP1=127.0.0.1

# 在发送消息时，自动创建服务器不存在的topic，默认创建的队列数
defaultTopicQueueNums=4

# 是否允许 Broker 自动创建 Topic，建议线下开启，线上关闭 ！！！这里仔细看是 false，false，false
autoCreateTopicEnable=true

# 是否允许 Broker 自动创建订阅组，建议线下开启，线上关闭
autoCreateSubscriptionGroup=true

# Broker 对外服务的监听端口
listenPort=10911

# 删除文件时间点，默认凌晨4点
deleteWhen=04

# 文件保留时间，默认48小时
fileReservedTime=120

# commitLog 每个文件的大小默认1G
mapedFileSizeCommitLog=1073741824

# ConsumeQueue 每个文件默认存 30W 条，根据业务情况调整
mapedFileSizeConsumeQueue=300000

# destroyMapedFileIntervalForcibly=120000
# redeleteHangedFileInterval=120000
# 检测物理文件磁盘空间
diskMaxUsedSpaceRatio=88
# 存储路径
# storePathRootDir=/home/ztztdata/rocketmq-all-4.1.0-incubating/store
# commitLog 存储路径
# storePathCommitLog=/home/ztztdata/rocketmq-all-4.1.0-incubating/store/commitlog
# 消费队列存储
# storePathConsumeQueue=/home/ztztdata/rocketmq-all-4.1.0-incubating/store/consumequeue
# 消息索引存储路径
# storePathIndex=/home/ztztdata/rocketmq-all-4.1.0-incubating/store/index
# checkpoint 文件存储路径
# storeCheckpoint=/home/ztztdata/rocketmq-all-4.1.0-incubating/store/checkpoint
# abort 文件存储路径
# abortFile=/home/ztztdata/rocketmq-all-4.1.0-incubating/store/abort
# 限制的消息大小
maxMessageSize=65536

# flushCommitLogLeastPages=4
# flushConsumeQueueLeastPages=2
# flushCommitLogThoroughInterval=10000
# flushConsumeQueueThoroughInterval=60000

# Broker 的角色
# - ASYNC_MASTER 异步复制Master
# - SYNC_MASTER 同步双写Master
# - SLAVE
brokerRole=ASYNC_MASTER

# 刷盘方式
# - ASYNC_FLUSH 异步刷盘
# - SYNC_FLUSH 同步刷盘
flushDiskType=ASYNC_FLUSH

# 发消息线程池数量
# sendMessageThreadPoolNums=128
# 拉消息线程池数量
# pullMessageThreadPoolNums=128

```

install/docker-compose.yml

```yaml
version: '3.5'
services:
  rmqnamesrv:
    image: foxiswho/rocketmq:server
    container_name: rmqnamesrv
    ports:
      - 9876:9876
    volumes:
      - ./logs:/opt/logs
      - ./store:/opt/store
    networks:
        rmq:
          aliases:
            - rmqnamesrv

  rmqbroker:
    image: foxiswho/rocketmq:broker
    container_name: rmqbroker
    ports:
      - 10909:10909
      - 10911:10911
    volumes:
      - ./logs:/opt/logs
      - ./store:/opt/store
      - ./conf/broker.conf:/etc/rocketmq/broker.conf
    environment:
        NAMESRV_ADDR: "rmqnamesrv:9876"
        JAVA_OPTS: " -Duser.home=/opt"
        JAVA_OPT_EXT: "-server -Xms256m -Xmx256m -Xmn256m"
    command: mqbroker -c /etc/rocketmq/broker.conf
    depends_on:
      - rmqnamesrv
    networks:
      rmq:
        aliases:
          - rmqbroker

  rmqconsole:
    image: styletang/rocketmq-console-ng
    container_name: rmqconsole
    ports:
      - 8080:8080
    environment:
        JAVA_OPTS: "-Drocketmq.namesrv.addr=rmqnamesrv:9876 -Dcom.rocketmq.sendMessageWithVIPChannel=false"
    depends_on:
      - rmqnamesrv
    networks:
      rmq:
        aliases:
          - rmqconsole

networks:
  rmq:
    name: rmq
    driver: bridge
```

安装

```
docker-compose up -d
```

管理后台: http://10.4.7.71:8080



## rocketmq基本概念

producer: 生产者

consumer: 消费者

broker: 暂存和传输消息，比如：邮局

nameserver: 管理broker, 比如各个邮局的管理机构

topic: 区分消息的种类，一个发送者可以发送消息给一个或多个topic，一个消息的接受者可以订阅一个或多个topic消息

message queue：相当于是topic的分区，用于并行发送和接受消息



## rocketmq 消息类型

### 安装发送的特点分

同步发送、异步发送、单向发送



### 安装使用功能分

普通消息（订阅）、顺序消息、延时消息（场景：库存超时归还）、事务消息（场景：分布式事务）



延迟消息级别：

messageDelayLevel=1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h

RocketMQ 支持定时消息，但是不支持任意时间精度，仅支持特定的 level，例如定时 5s， 10s， 1m 等。其中，level=0 级表示不延时，level=1 表示 1 级延时，level=2 表示 2 级延时，以此类推。



## rocket-go-test

https://github.com/apache/rocketmq-client-go



### go发送普通消息

rocketmq_test/simple/main.go

```go
package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"10.4.7.71:9876"}))
	if err != nil {
		panic("生成producer失败")
	}

	if err = p.Start(); err != nil {panic("启动producer失败")}

	res, err := p.SendSync(context.Background(), primitive.NewMessage("imooc1", []byte("this is imooc1")))
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
	}else{
		fmt.Printf("发送成功: %s\n", res.String())
	}

	if err = p.Shutdown(); err != nil {panic("关闭producer失败")}

}
```

### consumer消费消息

rocketmq_test/consumer/main.go

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"10.4.7.71:9876"}),
		consumer.WithGroupName("mxshop"),
	)

	if err := c.Subscribe("imooc1", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("获取到值： %v \n", msgs[i])
		}
		return consumer.ConsumeSuccess, nil
	}); err != nil {
		fmt.Println("读取消息失败")
	}
	_ = c.Start()
	//不能让主goroutine退出
	time.Sleep(time.Hour)
	_ = c.Shutdown()
}

```



### 延迟消息

rocketmq_test/delay/main.go

```go
package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"10.4.7.71:9876"}))
	if err != nil {
		panic("生成producer失败")
	}

	if err = p.Start(); err != nil {panic("启动producer失败")}

	msg := primitive.NewMessage("imooc1", []byte("this is delay message"))
	msg.WithDelayTimeLevel(3)
	res, err := p.SendSync(context.Background(), msg)
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
	}else{
		fmt.Printf("发送成功: %s\n", res.String())
	}

	if err = p.Shutdown(); err != nil {panic("关闭producer失败")}

	//支付的时候， 淘宝， 12306， 购票， 超时归还 - 定时执行逻辑
	//我可以去写一个轮询， 轮询的问题： 1. 多久执行一次轮询 30分钟
	//在12:00执行过一次， 下一次执行就是在 12:30的时候 但是12:01的时候下了单， 12:31就应该超时 13:00时候才能超时
	//那我1分钟执行一次啊， 比如我的订单量没有这么大，1分钟执行一次， 其中29次查询都是无用， 而且你还还会轮询mysql
	//rocketmq的延迟消息， 1. 时间一到就执行， 2. 消息中包含了订单编号，你只查询这种订单编号


}
```



### 事务消息

rocketmq_test/transaction/main.go

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type OrderListener struct{}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	fmt.Println("开始执行本地逻辑")
	time.Sleep(time.Second*3)
	fmt.Println("执行本地逻辑失败")
	//本地执行逻辑无缘无故失败 代码异常 宕机
	return primitive.UnknowState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("rocketmq的消息回查")
	time.Sleep(time.Second*15)
	return primitive.CommitMessageState
}

func main() {
	p, err := rocketmq.NewTransactionProducer(
		&OrderListener{},
		producer.WithNameServer([]string{"10.4.7.71:9876"}),
	)
	if err != nil {
		panic("生成producer失败")
	}

	if err = p.Start(); err != nil {panic("启动producer失败")}

	res, err := p.SendMessageInTransaction(context.Background(), primitive.NewMessage("TransTopic", []byte("this is transaction message2")))
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
	}else{
		fmt.Printf("发送成功: %s\n", res.String())
	}

	time.Sleep(time.Hour)
	if err = p.Shutdown(); err != nil {panic("关闭producer失败")}
}

```



# 实践

