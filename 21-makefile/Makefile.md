# 教程

http://c.biancheng.net/view/7097.html



# Makefile

我们可以把`Makefile`简单理解为它定义了一个项目文件的编译规则。借助`Makefile`我们在编译过程中不再需要每次手动输入编译的命令和编译的参数，可以极大简化项目编译过程。同时使用`Makefile`也可以在项目中确定具体的编译规则和流程，很多开源项目中都会定义`Makefile`文件。

本文不会详细介绍`Makefile`的各种规则，只会给出Go项目中常用的`Makefile`示例。关于`Makefile`的详细内容推荐阅读[Makefile教程](http://c.biancheng.net/view/7097.html)。

### 规则概述

`Makefile`由多条规则组成，每条规则主要由两个部分组成，分别是依赖的关系和执行的命令。

其结构如下所示：

```makefile
[target] ... : [prerequisites] ...
<tab>[command]
    ...
    ...
```

其中：

- targets：规则的目标
- prerequisites：可选的要生成 targets 需要的文件或者是目标。
- command：make 需要执行的命令（任意的 shell 命令）。可以有多条命令，每一条命令占一行。

举个例子：

```makefile
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o xx
```

## 示例

```makefile
.PHONY: all build run gotool clean help

BINARY="bluebell"

all: gotool build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}

run:
	@go run ./

gotool:
	go fmt ./
	go vet ./

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

help:
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"
```

其中：

- `BINARY="bluebell"`是定义变量。
- `.PHONY`用来定义伪目标。不创建目标文件，而是去执行这个目标下面的命令。

------



# go-makefile

```makefile
PKG := "github.com/jettjia/example"
PKG_LIST := $(shell go list ${PKG}/...)
APP=example
DOCKER_IMG=xxxxx
VERSION=1.0.0

.PHONY: tidy
tidy:
	$(eval files=$(shell find . -name go.mod))
	@set -e; \
	for file in ${files}; do \
		goModPath=$$(dirname $$file); \
		cd $$goModPath; \
		go mod tidy; \
		cd -; \
	done

.PHONY: fmt
fmt:
	@go fmt ${PKG_LIST}

init: # install golint
	@go install golang.org/x/lint/golint@latest

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

.PHONY: vet
vet: ## Vet the files
	@go vet ${PKG_LIST}

.PHONY: test
test:
	@go test -cover ./...

race: ## Run tests with data race detector
	@go test -race ${PKG_LIST}

.PHONY: test-coverage
test-coverage:
	@go test ./... -v -coverprofile=report/cover 2>&1 | go-junit-report > report/ut_report.xml
	@gocov convert report/cover | gocov-xml > report/coverage.xml
	@gocov convert report/cover.out | gocov-html > report/coverage.html

.PHONY: docker-image
docker-image:
	@docker build -t ${DOCKER_IMG}:v1.0.0 .

```

