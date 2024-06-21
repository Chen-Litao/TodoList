package cache

import (
	"ToDoList_self/config"
	"github.com/redis/go-redis/v9"
)

//var Ctx = context.Background()

var RdbTest *redis.Client

// UserFollowings 根据用户id找到他关注的人
var UserFollowings *redis.Client

// UserFollowers 根据用户id找到他的粉丝
var UserFollowers *redis.Client

func InitRedis() {
	RdbTest = redis.NewClient(&redis.Options{
		Addr:     config.ConfigVal.Redis.RedisAddr,
		Password: config.ConfigVal.Redis.Password, // no password set
		DB:       0,                               // use default DB
	})
	UserFollowings = redis.NewClient(&redis.Options{
		Addr:     config.ConfigVal.Redis.RedisAddr,
		Password: config.ConfigVal.Redis.Password,
		DB:       1,
	})
	UserFollowers = redis.NewClient(&redis.Options{
		Addr:     config.ConfigVal.Redis.RedisAddr,
		Password: config.ConfigVal.Redis.Password,
		DB:       2,
	})
}
