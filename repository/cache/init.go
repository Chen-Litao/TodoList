package cache

import (
	"ToDoList_self/config"
	"ToDoList_self/pkg/log"
	"github.com/go-redis/redis"
	"strconv"
)

var RedisClient *redis.Client

func InitRedis() {
	db, _ := strconv.ParseUint(config.ConfigVal.Redis.RedisDB, 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:     config.ConfigVal.Redis.RedisAddr,
		Password: config.ConfigVal.Redis.Password, // no password set
		DB:       int(db),                         // use default DB
		PoolSize: 100,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.LoggerObj.Error(err)
		return
	}
	RedisClient = client
}
