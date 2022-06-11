package db

import "database/sql"

// IDao 为接口， Dao为实现， Dao 依赖于 IDao

// 接口声明
type IDao interface {
	Version() (string, error)
}

// 默认实现
type Dao struct {
	db *sql.DB
}

// 生成dao对象的方法
func NewDao(db *sql.DB) *Dao {
	return &Dao{db: db}
}

// 在Dao中实现 IDao的接口
func (d *Dao) Version() (string, error) {
	var version string
	row := d.db.QueryRow("SELECT  VERSION()")
	if err := row.Scan(&version); err != nil {
		return "", err
	}

	return version, nil
}
