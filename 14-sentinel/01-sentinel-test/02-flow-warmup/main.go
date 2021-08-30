package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
)

func main() {
	//先初始化sentinel
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("初始化sentinel 异常: %v", err)
	}

	var globalTotal int
	var passTotal int
	var blockTotal int
	ch := make(chan struct{})

	//配置限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			TokenCalculateStrategy: flow.WarmUp, //冷启动策略
			ControlBehavior:        flow.Reject, //直接拒绝
			Threshold:              1000,
			WarmUpPeriodSec:        30, //预热的时间长度，该字段仅仅对 WarmUp 的TokenCalculateStrategy生效
		},
	})

	if err != nil {
		log.Fatalf("加载规则失败: %v", err)
	}

	//我会在每一秒统计一次，这一秒只能 你通过了多少，总共有多少， block了多少, 每一秒会产生很多的block
	for i := 0; i < 100; i++ {
		go func() {
			for {
				globalTotal++
				e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					//fmt.Println("限流了")
					blockTotal++
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					passTotal++
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					e.Exit()
				}
			}
		}()
	}

	go func() {
		var oldTotal int //过去1s总共有多少个
		var oldPass int  //过去1s总共pass多少个
		var oldBlock int //过去1s总共block多少个
		for {
			oneSecondTotal := globalTotal - oldTotal
			oldTotal = globalTotal

			oneSecondPass := passTotal - oldPass
			oldPass = passTotal

			oneSecondBlock := blockTotal - oldBlock
			oldBlock = blockTotal

			time.Sleep(time.Second)
			fmt.Printf("total:%d, pass:%d, block:%d\n", oneSecondTotal, oneSecondPass, oneSecondBlock)
		}
	}()

	<-ch
}
