package dao

import (
	"ToDoList_self/config"
	"ToDoList_self/repository/db/model"
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

var DB *gorm.DB

func InitMysql() {
	path := strings.Join([]string{config.ConfigVal.Mysql.User, ":", config.ConfigVal.Mysql.Password, "@tcp(", config.ConfigVal.Mysql.IP, ":",
		config.ConfigVal.Mysql.Port, ")/", config.ConfigVal.Mysql.Database, "?charset=utf8mb4&parseTime=true"}, "")
	fmt.Println(path)
	db, err := gorm.Open(mysql.Open(path), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败：", err)
		return
	}
	sqlDB, err := db.DB()
	//设置空闲连接池
	sqlDB.SetMaxIdleConns(10)
	//设置最大连接数
	sqlDB.SetMaxOpenConns(100)
	//设置连接时间
	sqlDB.SetConnMaxLifetime(time.Hour)
	err = db.AutoMigrate(&model.User{},
		&model.Task{},
		&model.Follow{})
	if err != nil {
		fmt.Println("数据库迁移失败：", err)
		return
	}
	DB = db
}
func NewDBClient(ctx context.Context) *gorm.DB {
	db := DB
	//加上上下文关联
	return db.WithContext(ctx)
}
