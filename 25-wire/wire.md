# wire

官方：github.com/google/wire



```text
go get github.com/google/wire/cmd/wire
go get github.com/google/wire
```



# why wire

## server.go

模拟的是业务代码

```go
package main

type Config struct {
	DbSource string
}

func NewConfig() *Config {
	return &Config{
		DbSource: "root:root@tcp(127.0.0.1:3306)/test",
	}
}

type DB struct {
	table string
}

func NewDB(cfg *Config) *DB {
	return &DB{
		table: "test_table",
	}
}

func (db *DB) Find() string {
	return "db info string"
}
```



## wire.go

这里是真正的wire代码。处理依赖对象的生成

```go
//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
)

//go:generate wire
func InitApp() (*App, error) {
	wire.Build(NewConfig, NewDB, NewApp) //调用 wire.Build 方法，传入所有的依赖对象 以及构建最终对象的函数 得到目标对象
	return &App{}, nil                   // 这里的返回没有实际的意义，只需要符合函数的签名即可，生成的 wire_gen会帮你实现
}
```

执行 wire，会生成 wire_gen的代码



## main.go

```go
package main

import "fmt"

type App struct {
	db *DB
}

func NewApp(db *DB) *App {
	return &App{db: db}
}

func main() {
	app, err := InitApp()
	if err != nil {
		panic(err)
	}

	result := app.db.Find()
	fmt.Println(result)
}
```



# 概念

两个概念： Provider 和 Injector

* Provider: 负责创建对象的方法，比如上文 控制反转示例的 NewDB(提供DB对象)和 NewConfig(提供Config对象)方法。
* Injector: 负责根据对象的依赖，依次构造依赖对象，最终构造目的对象的方法，比如上文中 控制反转示例的InitApp方法。

在上文中，NewConfig和NewDB都是 provider, wire_gen.go中的 InitApp 函数是 injector，可以看到 Injector通过依赖顺序调用 provider来生成我们需要的对象App.



# 快速入门

现在来做一个比较完整的demo,来进一步学习wire中的知识。

demo文件结构如下：

```shell
|--cmd
	|-- main.go
	|-- wire.go
|--config
	|-- app.json
|--internal
	|-- config
		|-- config.go
	|-- db
		|-- db.go
```

安装mysql

```shell
docker run -itd -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root --name test_mysql mysql:5.7
```

```shell
CREATE TABLE `order` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `price` decimal(10,2) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```



具体代码参考：02-wire-demo



# wire绑定struct

核心语法 wire.Struct(new(FooMessage), "*)

```go
func InitApp() (*FooMessage, error) {
	// 通过 wire.Struct 来指定那些字段要被注入到结构体中
	// 这里的 Msg、Ber 代表要导入的字段
	// 如果你要全部导入，可以这样写： wire.Struct(new(FooMessage), "*)
	wire.Build(ProvideMessage, ProvideBeer, wire.Struct(new(FooMessage), "Msg", "Ber"))
	return &FooMessage{}, nil
}
```

具体代码参考：03-wire-struct

# wire 绑定值

核心语法

```
wire.Build(config.Provider, wire.Value("demo string1"), wire.Value(config.String2("demo string 2")), NewApp)
```

```go
func InitApp() (*App, error) {
	// 绑定值用 wire.Value 进行绑定
	// 这里首先绑定了一个 string 类型的值
	// 然后绑定了 String2 类型的值，因为本例子需要绑定两个 string 类型的值。
	// 如果都用了 string 那么注入的时候，wire 无法区分具体的 string， 所以另外一个 string 使用自定义string类型
	wire.Build(config.Provider, wire.Value("demo string1"), wire.Value(config.String2("demo string 2")), NewApp)
	return &App{}, nil
}
```

具体代码参考：04-wire-value

# wire 绑定接口



# wire 清理函数 cleanup



# wire 工程化实践

