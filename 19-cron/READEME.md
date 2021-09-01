# 分布式调度工具

1) golang

go-cron

https://github.com/ouqiang/gocron



2) java

xxl-job

https://www.xuxueli.com/index.html

# xxl-job 执行器

https://github.com/xxl-job/xxl-job-executor-go

很多公司java与go开发共存，java中有xxl-job做为任务调度引擎，为此也出现了go执行器(客户端)，使用起来比较简单



## 安装

https://www.xuxueli.com/xxl-job/#其他：Docker%20镜像方式搭建调度中心：



1) 数据库拷贝

`/xxl-job/doc/db/tables_xxl_job.sql`

2) 安装

```shell
docker run -e PARAMS="--spring.datasource.url=jdbc:mysql://10.4.7.71:3307/xxl_job?useUnicode=true&characterEncoding=UTF-8&autoReconnect=true&serverTimezone=Asia/Shanghai --spring.datasource.username=root --spring.datasource.password=root --xxl.job.accessToken=1234567890olkjhhj" -p 8088:8080 -v /tmp:/data/applogs --name xxl-job-admin   -d xuxueli/xxl-job-admin:2.3.0
```

http://10.4.7.71:8088/xxl-job-admin

默认登录账号 “admin/123456”



**注意**

此为业务执行器代码访问admin的token

```
--xxl.job.accessToken=1234567890olkjhhj
```

此为执行失败发送邮箱通知

```
--spring.mail.host=smtp.qq.com
--spring.mail.port=25
--spring.mail.username=xxxx@qq.com
--spring.mail.password=xxxx
--spring.mail.properties.mail.smtp.auth=true
--spring.mail.properties.mail.smtp.starttls.enable=true
--spring.mail.properties.mail.smtp.starttls.required=true
--spring.mail.properties.mail.smtp.socketFactory.class=javax.net.ssl.SSLSocketFactory
```

参考：

```
docker run -e PARAMS="--spring.datasource.url=jdbc:mysql://127.0.0.1:3306/xxl_job?Unicode=true&characterEncoding=UTF-8 \
--spring.datasource.username=root \
--spring.datasource.password=123456 \
--spring.mail.host=smtp.qq.com \
--spring.mail.port=25 \
--spring.mail.username=xxxx@qq.com \
--spring.mail.password=xxxx \
--spring.mail.properties.mail.smtp.auth=true \
--spring.mail.properties.mail.smtp.starttls.enable=true \
--spring.mail.properties.mail.smtp.starttls.required=true \
--spring.mail.properties.mail.smtp.socketFactory.class=javax.net.ssl.SSLSocketFactory \
--xxl.job.accessToken=1234567890olkjhhj" \
-p 8080:8080 -v d:/tmp:/data/applogs \
--name xxl-job-admin --restart=always  -d xuxueli/xxl-job-admin:2.1.2
```





# xxl-job-executor的gin中间件

执行器项目地址:  https://github.com/xxl-job/xxl-job-executor-go

main.go

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-middleware/xxl-job-executor"
	"github.com/xxl-job/xxl-job-executor-go"
	"github.com/xxl-job/xxl-job-executor-go/example/task"
	"log"
)

const Port = "9999"

func main() {
	//初始化执行器
	exec := xxl.NewExecutor(
		xxl.ServerAddr("http://10.4.7.71:8088/xxl-job-admin"),
		xxl.AccessToken("1234567890olkjhhj"),            //请求令牌(默认为空)
		xxl.ExecutorIp("10.4.7.71"),    //可自动获取
		xxl.ExecutorPort(Port),         //默认9999（此处要与gin服务启动port必需一至）
		xxl.RegistryKey("golang-jobs"), //执行器名称
	)
	exec.Init()
	defer exec.Stop()
	//添加到gin路由
	r := gin.Default()
	xxl_job_executor_gin.XxlJobMux(r, exec)

	//注册gin的handler
	r.GET("ping", func(cxt *gin.Context) {
		cxt.JSON(200, "pong")
	})

	//注册任务handler
	exec.RegTask("task.test", task.Test)
	exec.RegTask("task.test2", task.Test2)
	exec.RegTask("task.panic", task.Panic)

	log.Fatal(r.Run(":" + Port))
}

```

# xxl-job-admin配置

### 添加执行器

执行器管理->新增执行器,执行器列表如下：



```undefined
AppName     名称      注册方式    OnLine      机器地址    操作
golang-jobs golang执行器   自动注册    无
```

### 添加任务

任务管理->新增(注意，使用BEAN模式，JobHandler与RegTask名称一致)



```css
1   测试panic BEAN：task.panic * 0 * * * ? admin   STOP    
2   测试耗时任务  BEAN：task.test2 * * * * * ? admin   STOP    
3   测试golang    BEAN：task.test      * * * * * ? admin   STOP
```
