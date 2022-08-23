package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	// 初始化
	s := gocron.NewScheduler(time.Local) // 使用系统的本地时区
	//s:=gocron.NewScheduler(time.UTC) // 使用UTC时区
	log.Println("start")

	// 示例一：每2秒钟执行一次
	s.Every(2).Seconds().Do(func() {
		log.Println("执行了。。。。2s")
	})

	// 示例二：WaitForSchedule() 的使用
	// 默认情况下会立即执行一次，使用 WaitForSchedule() 可禁止这种情况，等到下次才执行
	s.Every(5).Seconds().Do(func() {
		log.Println("啦啦啦111")
	})
	s.Every(5).Seconds().WaitForSchedule().Do(func() {
		log.Println("啦啦啦222")
	})

	// 示例三：通过 crontab表达式来执行
	// 标准的crontab格式，最小单位是分
	s.Cron("*/1 * * * *").Do(task)
	// 最小单位是秒的crontab表达式
	s.CronWithSeconds("*/1 * * * * *").Do(task)

	// 示例四：指定时间运行
	s.Every(1).Sunday().At("00:30").Do(task)
	s.Every(1).Day().At("10:00").Do(task)

	// 示例五：SingletonMode() 单例模式
	// 如果之前的任务尚未完成，单例模式将阻止新任务启动
	s.Every("2").Seconds().SingletonMode().Do(task)

	// 示例六：带有参数的任务
	s.Every(1).Seconds().Do(taskWithParams, 2, "test")

	// 异步启动
	s.StartAsync()

	// 同步启动，阻塞进程
	s.StartBlocking()

}

func task() {
	log.Println("hello")
}

func taskWithParams(a int, b string) {
	log.Println(a, b)
}
