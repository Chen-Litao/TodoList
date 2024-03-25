package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Mysql Mysql
	Redis Redis
}

type Mysql struct {
	IP       string
	Port     string
	User     string
	Password string
	Database string
}

type Redis struct {
	RedisAddr string
	Password  string
	RedisDB   string
}

var ConfigVal *Config

func LoadeConf() {
	dataBase, err := os.ReadFile("config/config.yaml")
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return
	}
	err = yaml.Unmarshal(dataBase, &ConfigVal)
	if err != nil {
		fmt.Println("解析yaml失败：", err)
		return
	}
}
