package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"majiang/model"
)

// debian-sys-maint:ZF0kfsp5uMD2lVo7

const (
	USER     = "root"
	PASSWORD = "224488"
)

var (
	DB *gorm.DB
)

// InitDB gorm连接
func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/majiang?charset=utf8mb4&parseTime=True&loc=Local", USER, PASSWORD)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("failed to connect database")
	}
	// 设置数据库连接池参数
	sqlDB, _ := db.DB()
	// 设置数据库连接池最大连接数
	sqlDB.SetMaxOpenConns(100)
	// 连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭
	sqlDB.SetMaxIdleConns(20)
	DB = db
	AutoMigrate()
}

func AutoMigrate() {
	DB.AutoMigrate(&model.User{})
}
