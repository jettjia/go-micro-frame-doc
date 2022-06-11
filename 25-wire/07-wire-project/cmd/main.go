package main

import (
	"context"
	"database/sql"
	"fmt"
	"go-micro-frame-doc/25-wire/07-wire-project/internal/biz"
	"go-micro-frame-doc/25-wire/07-wire-project/internal/server/config"
	"log"
)

type App struct {
	conf      *config.Config
	db        *sql.DB
	OrderRepo biz.OrderRepo
}

func NewApp(conf *config.Config, db *sql.DB, orderRepo biz.OrderRepo) *App {
	return &App{
		conf:      conf,
		db:        db,
		OrderRepo: orderRepo,
	}
}

func main() {
	ctx := context.Background()
	app, cleanup, err := InitApp(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer cleanup()

	// 测插入
	var o = biz.Order{
		Name:  "ins1",
		Price: 22,
	}
	ID, err := app.OrderRepo.Create(ctx, &o)
	if err != nil {
		fmt.Println("插入失败:", err)
	}
	fmt.Printf("查询成功, %+v", ID)

	// 测试查询
	order, err := app.OrderRepo.Find(ctx, 1)
	if err != nil {
		fmt.Println("order find error:", err)
		return
	}
	fmt.Printf("查询成功, %+v", order)
	fmt.Println()
}
