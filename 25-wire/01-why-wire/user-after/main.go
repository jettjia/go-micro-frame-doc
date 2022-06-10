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
