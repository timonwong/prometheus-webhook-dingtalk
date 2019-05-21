package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var GormDB = new(gorm.DB)
var (
	SAVE_TO_MYSQL  bool
	MYSQL_USER     string
	MYSQL_PASSWD   string
	MYSQL_HOST     string
	MYSQL_PORT     int
	MYSQL_DATABASE string
)

func DBinit() {
	if !SAVE_TO_MYSQL {
		fmt.Println("告警记录不做持久化")
		return
	}
	GormDB = newConn()
	GormDB.DB().SetMaxIdleConns(100)
	GormDB.DB().SetMaxOpenConns(200)
	GormDB.LogMode(true)
}

func newConn() *gorm.DB {
	gormdb, err := gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		MYSQL_USER,
		MYSQL_PASSWD,
		MYSQL_HOST,
		MYSQL_PORT,
		MYSQL_DATABASE,
	))

	if err != nil {
		panic(err)
	}
	return gormdb
}
