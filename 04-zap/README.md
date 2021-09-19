# 官方

https://pkg.go.dev/go.uber.org/zap

https://github.com/uber-go/zap

# 案例

zap/initialize/logger.go

```go
package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}


```



首页初始化并调用

```go
package main

import (
	"go-micro-frame-doc/04-zap/initialize"
	"go.uber.org/zap"
)

func main()  {
	// 初始化 logger
	initialize.InitLogger()

	zap.S().Debugf("entry main.go", "wwwwwwwwww")
}

```

