package cache

import (
	"ToDoList_self/repository/db/model"
	"fmt"
	"strconv"
)

const (
	RankKey = "rank"
)

// TaskViewKey 点击数的key
func TaskViewKey(id uint) string {
	return fmt.Sprintf("view:task:%s", strconv.Itoa(int(id)))
}

func View(Task *model.Task) uint64 {
	countStr, _ := RedisClient.Get(TaskViewKey(Task.ID)).Result()
	count, _ := strconv.ParseUint(countStr, 10, 64)
	return count
}

// AddView
func AddView(Task *model.Task) {
	//增加点击数
	RedisClient.Incr(TaskViewKey(Task.ID))                      // 增加视频点击数
	RedisClient.ZIncrBy(RankKey, 1, strconv.Itoa(int(Task.ID))) // 增加排行点击数
}
