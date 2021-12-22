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