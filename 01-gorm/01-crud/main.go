package main

import (
	"fmt"
	"go-micro-frame-doc/01-gorm/01-crud/model"
)


func main() {
	// 自动生成表结构
	//autoMigrate()

	// create
	//add()

	// simple select
	//findOne()

	// update
	updateOne()
}


func updateOne() {
	var user model.User

	db := GetDb()
	db.First(&user, 1)

	// 更改
	user.Name = "lisi"
	db.Save(&user)

}

func findOne() {
	var user model.User

	db := GetDb()
	db.First(&user, 1)
	fmt.Println(user)
}

func deleteOne() {
	db := GetDb()
	db.Delete(&model.User{}, 1)
}

func add() {
	user := model.User{
		Name:   "zhangsan",
		Age:    18,
		Gender: 1,
	}

	// create
	db := GetDb()
	result := db.Create(&user)
	fmt.Println(result.RowsAffected) // 返回插入记录的条数
}


func autoMigrate() {
	GetDb().AutoMigrate(&model.User{})
}