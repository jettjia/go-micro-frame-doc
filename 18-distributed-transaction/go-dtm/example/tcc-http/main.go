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
	dtm := "http://10.4.7.88:36789/api/dtmsvr" // 安装dtm管理的地址
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