package config

import (
	"ToDoList_self/repository/db/model"
	"fmt"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"strings"
	"time"
)

type Config struct {
	Mysql Mysql
}

type Mysql struct {
	IP       string
	Port     string
	User     string
	Password string
	Database string
}

var DB *gorm.DB

func Loade() {
	dataBase, err := os.ReadFile("config/config.yaml")
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return
	}
	config := Config{}
	err = yaml.Unmarshal(dataBase, &config)
	if err != nil {
		fmt.Println("解析yaml失败：", err)
		return
	}
	path := strings.Join([]string{config.Mysql.User, ":", config.Mysql.Password, "@tcp(", config.Mysql.IP, ":",
		config.Mysql.Port, ")/", config.Mysql.Database, "?charset=utf8mb4"}, "")
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
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		fmt.Println("数据库迁移失败：", err)
		return
	}
	DB = db
}
