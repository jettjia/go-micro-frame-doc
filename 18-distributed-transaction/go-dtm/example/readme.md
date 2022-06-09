

# 事务原理

## XA

XA是标准二阶段提交，一阶段Prepare，二阶段Commit/Rollback。一阶段的所有Prepare成功，则Commit，否则Rollback。从修改数据开始，到Commit/Rollback结束，锁定数据长，并发度较低。

这种事务模式适合并发不高的订单转账等各项业务。



## TCC

TCC的每个子事务有两个阶段：一阶段Try；二阶段Confirm/Cancel；Try中进行资源预留，例如冻结资金；一阶段如果全部成功，则Confirm，例如进行余额扣减，解冻资金；一阶段如果有一个子事务出现失败，则Cancel，例如解冻资金；

TCC模式不长期锁数据，并发高。常用于支付、订单类业务的拆分。



## SAGA

SAGA每个子事务有正向分支和补偿分支，它在正向分支中直接修改数据，出错则补偿所有修改过的数据。

SAGA比TCC少一个分支，一致性比TCC弱。适用于积分换礼品等不涉及资金业务；也适用于对接较多第三方的长事务。



## 事务消息模式

事务消息提供了支持事务的消息接口，允许使用方把消息发送放到本地的一个事务里，保证事务的原子性。它的工作原理如下：

本地应用

- 开启本地事务
- 进行本地数据库修改
- 调用消息事务的prepare接口，预备发送消息
- 提交本地事务
- 调用消息事务的submit接口，触发消息发送

当事务管理器，只收到prepare请求，超时未收到submit请求时，调用反查接口canSubmit，询问应用程序，是否能够发送。

事务消息与本地消息方案类似，但是将创建本地消息表和轮询本地消息表的操作换成了一个反查接口，提供更加便捷的使用。

假定一个这样的场景，用户注册成功后，需要给用户赠送优惠券和一个月会员卡。这里赠送优惠券和一个月会员一定不会失败，这种情况就非常适合可靠消息事务模式。

# tcc介绍

## tcc组成

TCC分为3个阶段

- Try 阶段：尝试执行，完成所有业务检查（一致性）, 预留必须业务资源（准隔离性）
- Confirm 阶段：如果所有分支的Try都成功了，则走到Confirm阶段。Confirm真正执行业务，不作任何业务检查，只使用 Try 阶段预留的业务资源
- Cancel 阶段：如果所有分支的Try有一个失败了，则走到Cancel阶段。Cancel释放 Try 阶段预留的业务资源。



TCC分布式事务里，有3个角色，与经典的XA分布式事务一样：

- AP/应用程序，发起全局事务，定义全局事务包含哪些事务分支
- RM/资源管理器，负责分支事务各项资源的管理
- TM/事务管理器，负责协调全局事务的正确执行，包括Confirm，Cancel的执行，并处理网络异常

# saga介绍

## saga的理论来源

saga这种事务模式最早来自这篇论文：[sagas](https://link.zhihu.com/?target=https%3A//link.segmentfault.com/%3Fenc%3DqMpSS1cEyjssWhauwWBBKA%3D%3D.KSswh0TAzfOPYFzsxodk0c5uytRijyO3IbkRKXu7OSVTJOv7PePQ1HpjLq94NhON)

在这篇论文里，作者提出了将一个长事务，分拆成多个子事务，每个子事务有正向操作Ti，反向补偿操作Ci。

假如所有的子事务Ti依次成功完成，全局事务完成

假如子事务Ti失败，那么会调用Ci, Ci-1, Ci-2 ....进行补偿

论文阐述了上述这部分基本的saga逻辑之后，提出了下面几种场景的技术处理

## 回滚与重试

对于一个SAGA事务，如果执行过程中遭遇失败，那么接下来有两种选择，一种是进行回滚，另一种是重试继续。

回滚的机制相对简单一些，只需要在进行下一步之前，把下一步的操作记录到保存点就可以了。一旦出现问题，那么从保存点处开始回滚，反向执行所有的补偿操作即可。

假如有一个持续了一天的长事务，被服务器重启这类临时失败中断后，此时如果只能进行回滚，那么业务是难以接受的。 此时最好的策略是在保存点处重试并让事务继续，直到事务完成。

往前重试的支持，需要把全局事务的所有子事务事先编排好并保存，然后在失败时，重新读取未完成的进度，并重试继续执行。

## 并发执行

对于长事务而言，并发执行的特性也是至关重要的，一个串行耗时一天的长事务，在并行的支持下，可能半天就完成了，这对业务的帮助很大。

某些场景下并发执行子事务，是业务必须的要求，例如订多张及票，而机票确认时间较长时，不应当等前一个票已经确认之后，再去定下一张票，这样会导致订票成功率大幅下降。

在子事务并发执行的场景下，支持回滚与重试，挑战会更大，涉及了较复杂的保存点。



## 解决问题实例

我们以一个真实用户案例，来讲解[dtm](https://link.zhihu.com/?target=https%3A//github.com/yedf/dtm)的saga最佳实践。

问题场景：一个用户出行旅游的应用，收到一个用户出行计划，需要预定去三亚的机票，三亚的酒店，返程的机票。

要求：

1. 两张机票和酒店要么都预定成功，要么都回滚（酒店和航空公司提供了相关的回滚接口）
2. 预订机票和酒店是并发的，避免串行的情况下，因为某一个预定最后确认时间晚，导致其他的预定错过时间
3. 预定结果的确认时间可能从1分钟到1天不等

上述这些要求，正是saga事务模式要解决的问题，我们来看看dtm怎么解决（以Go语言为例）。

首先我们根据要求1，创建一个saga事务，这个saga包含三个分支，分别是，预定去三亚机票，预定酒店，预定返程机票

```text
        saga := dtmcli.NewSaga(DtmServer, gid).
            Add(Busi+"/BookTicket", Busi+"/BookTicketRevert", bookTicketInfo1).
            Add(Busi+"/BookHotel", Busi+"/BookHotelRevert", bookHotelInfo2).
            Add(Busi+"/BookTicket", Busi+"/BookTicketRevert", bookTicketBackInfo3)
```

然后我们根据要求2，让saga并发执行（默认是顺序执行）

```text
  saga.EnableConcurrent()
```

最后我们处理3里面的“预定结果的确认时间”不是即时响应的问题。由于不是即时响应，所以我们不能够让预定操作等待第三方的结果，而是提交预定请求后，就立即返回状态-进行中。我们的分支事务未完成，dtm会重试我们的事务分支，我们把重试间隔指定为1分钟。

```text
  saga.SetOptions(&dtmcli.TransOptions{RetryInterval: 60})
  saga.Submit()
// ........
func bookTicket() string {
    order := loadOrder()
    if order == nil { // 尚未下单，进行第三方下单操作
        order = submitTicketOrder()
        order.save()
    }
    order.Query() // 查询第三方订单状态
    return order.Status // 成功-SUCCESS 失败-FAILURE 进行中-ONGOING
}
```

## 高级用法

在实际应用中，还遇见过一些业务场景，需要一些额外的技巧进行处理

### 支持重试与回滚

dtm要求业务明确返回以下几个值：

- SUCCESS表示分支成功，可以进行下一步
- FAILURE 表示分支失败，全局事务失败，需要回滚
- ONGOING表示进行中，后续按照正常的间隔进行重试
- 其他表示系统问题，后续按照指数退避算法进行重试

### 部分第三方操作无法回滚

例如一个订单中的发货，一旦给出了发货指令，那么涉及线下相关操作，那么很难直接回滚。对于涉及这类情况的saga如何处理呢？

我们把一个事务中的操作分为可回滚的操作，以及不可回滚的操作。那么把可回滚的操作放到前面，把不可回滚的操作放在后面执行，那么就可以解决这类问题

```text
        saga := dtmcli.NewSaga(DtmServer, dtmcli.MustGenGid(DtmServer)).
            Add(Busi+"/CanRollback1", Busi+"/CanRollback1Revert", req).
            Add(Busi+"/CanRollback2", Busi+"/CanRollback2Revert", req).
            Add(Busi+"/UnRollback1", Busi+"/UnRollback1NoRevert", req).
            EnableConcurrent().
            AddBranchOrder(2, []int{0, 1}) // 指定step 2，需要在0，1完成后执行
```

### 超时回滚

saga属于长事务，因此持续的时间跨度很大，可能是100ms到1天，因此saga没有默认的超时时间。

dtm支持saga事务单独指定超时时间，到了超时时间，全局事务就会回滚。

```text
    saga.SetOptions(&dtmcli.TransOptions{TimeoutToFail: 1800})
```

在saga事务中，设置超时时间一定要注意，这类事务里不能够包含无法回滚的事务分支，否则超时回滚这类的分支会有问题。

### 其他分支的结果作为输入

如果极少数的实际业务不仅需要知道某些事务分支是否执行成功，还想要获得成功的详细结果数据，那么[dtm](https://link.zhihu.com/?target=https%3A//github.com/yedf/dtm)如何处理这样的需求呢？例如B分支需要A分支的执行成功返回的详细数据。

dtm的建议做法是，在ServiceA再提供一个接口，让B可以获取到相关的数据。这种方案虽然效率稍低，但是易理解已维护，开发工作量也不会太大。

PS：有个小细节请注意，尽量在你的事务外部进行网络请求，避免事务时间跨度变长，导致并发问题。



# 安装dtm 管理器

```
git clone https://github.com/yedf/dtm
cd dtm
docker-compose up -d
```

也可以直接拷贝项目里的 docker-compose.yml

```yaml
version: '3.3'
services:
  api:
    image: 'yedf/dtm'
    environment:
      IS_DOCKER: '1'
    ports:
      - '36789:36789'
      - '36790:36790'
    volumes:
      - .:/app/work
      - /etc/localtime:/etc/localtime:ro
      - /etc/timezone:/etc/timezone:ro
    command: ['/app/dtm/main', 'dev']
    working_dir: /app/work
    extra_hosts:
      - 'host.docker.internal:host-gateway'
  db:
    image: 'mysql:5.7'
    environment:
      MYSQL_ALLOW_EMPTY_PASSWORD: 1
    volumes:
      - /etc/localtime:/etc/localtime:ro
      - /etc/timezone:/etc/timezone:ro
    command:
      [
        '--character-set-server=utf8mb4',
        '--collation-server=utf8mb4_unicode_ci',
      ]
    ports:
      - '3306:3306'

```





# http-tcc

代码参考 [tcc-http/main.go](./tcc-http/main.go)

## main.go

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/yedf/dtmcli"
)

// 事务参与者的服务地址
const tccBusiAPI = "/api/busi_start"
const tccBusiPort = 8082

var tccBusi = fmt.Sprintf("http://localhost:%d%s", tccBusiPort, tccBusiAPI)

// 启动服务
func startSvr() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	qsAddRoute(app)
	log.Printf("quick start examples listening at %d", tccBusiPort)
	go app.Run(fmt.Sprintf(":%d", tccBusiPort))
	time.Sleep(100 * time.Millisecond)
}

func tccFireRequest() string {
	log.Printf("tcc transaction begin")
	dtm := "http://localhost:36789/api/dtmsvr" // 安装dtm管理的地址
	gid := dtmcli.MustGenGid(dtm) //生产全局唯一的 gid
	// TccGlobalTransaction 开启一个TCC全局事务，第一个参数为dtm的地址，第二个参数是回调函数
	err := dtmcli.TccGlobalTransaction(dtm, gid, func(tcc *dtmcli.Tcc) (resp *resty.Response, rerr error) {
		// 调用TransOut分支，三个参数分别为post的body，tryUrl，confirmUrl，cancelUrl
		// res1 为try执行的结果
		resp, rerr = tcc.CallBranch(gin.H{"amount": 30}, tccBusi+"/TransOut", tccBusi+"/TransOutConfirm", tccBusi+"/TransOutCancel")
		if rerr != nil {
			return
		}
		// 调用TransIn分支
		resp, rerr = tcc.CallBranch(gin.H{"amount": 30}, tccBusi+"/TransIn", tccBusi+"/TransInConfirm", tccBusi+"/TransInCancel")
		if rerr != nil {
			return
		}
		// 返回后，tcc会把全局事务提交，DTM会调用个分支的Confirm
		return
	})
	if err != nil {
		log.Fatalf("Tcc transaction failed: %v", err)
	}
	log.Printf("tcc %s submitted", gid)
	return gid
}

// gin路由方法，里面有各个阶段的方法触发，模拟分布式的各个调度
func qsAddRoute(app *gin.Engine) {
	app.POST(tccBusiAPI+"/TransIn", func(c *gin.Context) {
		log.Printf("TransIn ok")
		c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
	}).POST(tccBusiAPI+"/TransInConfirm", func(c *gin.Context) {
		log.Printf("TransInConfirm ok")
		c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
	}).POST(tccBusiAPI+"/TransInCancel", func(c *gin.Context) {
		log.Printf("TransInCancel ok")
		c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
	}).POST(tccBusiAPI+"/TransOut", func(c *gin.Context) {
		log.Printf("TransOut ok")
		c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
	}).POST(tccBusiAPI+"/TransOutConfirm", func(c *gin.Context) {
		log.Printf("TransOutConfirm ok")
		c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
	}).POST(tccBusiAPI+"/TransOutCancel", func(c *gin.Context) {
		log.Printf("TransOutCancel ok")
		c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
	})
}

func main() {
	startSvr()
	tccFireRequest()
	time.Sleep(1000 * time.Second)
}

```

## 启动

```shell
go run main.go
```





# http-saga

代码参考 [saga-http/main.go](./saga-http/main.go)



## main.go

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yedf/dtmcli"
)

// 事务参与者的服务地址
const sagaBusiAPI = "/api/busi_start"
const qsBusiAPI = "/api/busi_start"
const sagaBusiPort = 8082

var sagaBusi = fmt.Sprintf("http://localhost:%d%s", sagaBusiPort, sagaBusiAPI)

type TransReq struct {
	Amount         int    `json:"amount"`
	TransInResult  string `json:"transInResult"`
	TransOutResult string `json:"transOutResult"`
}

// 启动服务
func startSvr() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	qsAddRoute(app)
	log.Printf("quick start examples listening at %d", sagaBusiPort)
	go app.Run(fmt.Sprintf(":%d", sagaBusiPort))
	time.Sleep(100 * time.Millisecond)
}

func sagaFireRequest() string {
	log.Printf("saga transaction begin")
	dtm := "http://10.4.7.88:36789/api/dtmsvr" // 安装dtm管理的地址
	gid := dtmcli.MustGenGid(dtm) //生产全局唯一的 gid
	// sagaGlobalTransaction 开启一个saga全局事务，第一个参数为dtm的地址，第二个参数是回调函数

	req := &TransReq{Amount: 30}

	saga := dtmcli.NewSaga(dtm, dtmcli.MustGenGid(dtm)).
		// 添加一个TransOut的子事务，正向操作为url: qsBusi+"/TransOut"， 逆向操作为url: qsBusi+"/TransOutCompensate"
		Add(qsBusiAPI+"/TransOut", qsBusiAPI+"/TransOutCompensate", req).
		// 添加一个TransIn的子事务，正向操作为url: qsBusi+"/TransOut"， 逆向操作为url: qsBusi+"/TransInCompensate"
		Add(qsBusiAPI+"/TransIn", qsBusiAPI+"/TransInCompensate", req)

	// 提交saga事务，dtm会完成所有的子事务/回滚所有的子事务
	err := saga.Submit()

	if err != nil {
		log.Fatalf("saga transaction failed: %v", err)
	}
	log.Printf("saga %s submitted", gid)
	return gid
}

// gin路由方法，里面有各个阶段的方法触发，模拟分布式的各个调度
func qsAddRoute(app *gin.Engine) {

	app.POST(qsBusiAPI+"/TransIn", func(c *gin.Context) {
		log.Printf("TransIn ok")
		c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
	}).POST(qsBusiAPI+"/TransInCompensate", func(c *gin.Context) {
		log.Printf("TransInCompensate ok")
		c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
	}).POST(qsBusiAPI+"/TransOut", func(c *gin.Context) {
		log.Printf("TransOut ok")
		c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
	}).POST(qsBusiAPI+"/TransOutCompensate", func(c *gin.Context) {
		log.Printf("TransOutCompensate ok")
		c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
	})
}

func main() {
	startSvr()
	sagaFireRequest()
	time.Sleep(1000 * time.Second)
}
```



## 启动

```shell
go run main.go
```



