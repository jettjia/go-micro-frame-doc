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
