package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Mysql    Mysql
	Redis    Redis
	RabbitMQ RabbitMQ
}

type Mysql struct {
	IP       string `yaml:"ip"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Redis struct {
	RedisAddr string `yaml:"redisAddr"`
	Password  string `yaml:"password"`
	RedisDB   string `yaml:"redisDB"`
}

type RabbitMQ struct {
	RabbitMQ         string `yaml:"rabbitMQ"`
	RabbitMQUser     string `yaml:"rabbitMQUser"`
	RabbitMQPassWord string `yaml:"rabbitMQPassWord"`
	RabbitMQHost     string `yaml:"rabbitMQHost"`
	RabbitMQPort     string `yaml:"rabbitMQPort"`
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
