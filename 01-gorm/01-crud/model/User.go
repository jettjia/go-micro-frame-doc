package model

type User struct {
	Id     int    `gorm:"primary_key" json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender int    `json:"gender"` //1:男、2:女
}
