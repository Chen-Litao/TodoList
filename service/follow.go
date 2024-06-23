package service

import (
	"ToDoList_self/repository/cache"
	"ToDoList_self/repository/db/dao"
	"ToDoList_self/repository/mq"
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var FollowSrvIns *FollowSrv
var FollowSrvOnce sync.Once

type FollowSrv struct {
}

func GetFollowSrv() *FollowSrv {
	TaskSrvOnce.Do(func() {
		//FollowSrvIns = &FollowSrv{}
	})
	return FollowSrvIns
}
func CacheTimeGenerator() time.Duration {
	// 先设置随机数 - 这里比较重要
	rand.Seed(time.Now().Unix())
	// 再设置缓存时间
	// 10 + [0~20) 分钟的随机时间
	return time.Duration((10 + rand.Int63n(20)) * int64(time.Minute))
}

// FollowAction 关注操作的业务
func (followService *FollowSrv) FollowAction(ctx context.Context, userId int64, targetId int64) (bool, error) {

	followDao := dao.NewFollowDao(ctx)
	follow, err := followDao.FindEverFollowing(userId, targetId)
	// 寻找SQL 出错。
	if nil != err {
		return false, err
	}
	// 获取关注的消息队列
	followAddMQ := mq.SimpleFollowAddMQ
	// 曾经关注过，只需要update一下followed即可。
	if nil != follow {
		//发送消息队列
		err := followAddMQ.PublishSimpleFollow(fmt.Sprintf("%d-%d-%s", userId, targetId, "update"))
		if err != nil {
			return false, err
		}
		//更新Redis
		followService.AddToRDBWhenFollow(ctx, userId, targetId)
		return true, nil

	}
	//发送消息队列
	err = followAddMQ.PublishSimpleFollow(fmt.Sprintf("%d-%d-%s", userId, targetId, "insert"))
	if err != nil {
		return false, err
	}
	followService.AddToRDBWhenFollow(ctx, userId, targetId)
	return true, nil
}

func (followService *FollowSrv) AddToRDBWhenFollow(ctx context.Context, userId int64, targetId int64) {
	followDao := dao.NewFollowDao(ctx)
	// 尝试给following数据库追加user关注target的记录
	stringID := strconv.FormatInt(userId, 10)
	keyCnt1 := cache.UserFollowings.Exists(context.Background(), stringID)

	if keyCnt1.Err() != nil {
		log.Println(keyCnt1.Err().Error())
	}

	// 只判定键是否不存在，若不存在即从数据库导入
	if keyCnt1.Val() <= 0 {
		userFollowingsId, _, err := followDao.GetFollowingsInfo(userId)
		if err != nil {
			log.Println(err.Error())
			return
		}
		ImportToRDBFollowing(ctx, userId, userFollowingsId)
	}
	// 数据库导入到redis结束后追加记录
	cache.UserFollowings.SAdd(ctx, strconv.FormatInt(userId, 10), targetId)

	// 尝试给follower数据库追加target的粉丝有user的记录
	keyCnt2 := cache.UserFollowings.Exists(context.Background(), strconv.FormatInt(targetId, 10))

	if keyCnt2.Err() != nil {
		log.Println(keyCnt2.Err().Error())
	}
	//
	if keyCnt2.Val() <= 0 {
		//获取target的粉丝，直接刷新，关注时刷新target的粉丝
		userFollowersId, _, err := followDao.GetFollowersInfo(targetId)
		if err != nil {
			log.Println(err.Error())
			return
		}
		ImportToRDBFollower(ctx, targetId, userFollowersId)
	}

	cache.UserFollowers.SAdd(ctx, strconv.FormatInt(targetId, 10), userId)
}

// ImportToRDBFollowing 将登陆用户的关注id列表导入到following数据库中
func ImportToRDBFollowing(ctx context.Context, userId int64, ids []int64) {
	// 将传入的userId及其关注用户id更新至redis中
	for _, id := range ids {
		cache.UserFollowings.SAdd(ctx, strconv.FormatInt(userId, 10), int(id))
	}

	cache.UserFollowings.Expire(ctx, strconv.FormatInt(userId, 10), CacheTimeGenerator())
}

// ImportToRDBFollower 将登陆用户的关注id列表导入到follower数据库中
func ImportToRDBFollower(ctx context.Context, userId int64, ids []int64) {
	// 将传入的userId及其粉丝id更新至redis中
	for _, id := range ids {
		cache.UserFollowers.SAdd(ctx, strconv.FormatInt(userId, 10), int(id))
	}

	cache.UserFollowers.Expire(ctx, strconv.FormatInt(userId, 10), CacheTimeGenerator())
}

// CancelFollowAction 取关操作的业务
func (followService *FollowSrv) CancelFollowAction(ctx context.Context, userId int64, targetId int64) (bool, error) {

	// 获取取关的消息队列
	followDelMQ := mq.SimpleFollowDelMQ
	followDao := dao.NewFollowDao(ctx)
	follow, err := followDao.FindEverFollowing(userId, targetId)
	// 寻找 SQL 出错。
	if nil != err {
		return false, err
	}
	// 曾经关注过，只需要update一下cancel即可。
	if nil != follow {
		err := followDelMQ.PublishSimpleFollow(fmt.Sprintf("%d-%d-%s", userId, targetId, "update"))
		if err != nil {
			return false, err
		}
		DelToRDBWhenCancelFollow(userId, targetId)
		return true, nil

	}
	// 没有关注关系
	return false, nil
}
func DelToRDBWhenCancelFollow(userId int64, targetId int64) {
	// 当a取关b时，redis的三个关注数据库会有以下操作
	cache.UserFollowings.SRem(cache.Ctx, strconv.FormatInt(userId, 10), targetId)

	cache.UserFollowers.SRem(cache.Ctx, strconv.FormatInt(targetId, 10), userId)

}

func (followService *FollowSrv) GetFollowings(userId int64) ([]User, error) {
	// 调用集成redis的关注用户获取接口获取关注用户id和关注用户数量
	userFollowingsId, userFollowingsCnt, err := GetFollowingsByRedis(userId)
	if nil != err {
		log.Println(err.Error())
	}
	// 根据关注用户数量创建空用户结构体数组
	userFollowings := make([]User, userFollowingsCnt)

	// 传入buildtype调用用户构建函数构建关注用户数组
	err1 := followService.BuildUser(userId, userFollowings, userFollowingsId, 0)

	if nil != err1 {
		log.Println(err1.Error())
	}

	return userFollowings, nil

}

func (followService *FollowSrv) GetFollower(ctx context.Context, userId int64) ([]User, error) {
	userFollowersId, userFollowingsCnt, err := GetFollowerByRedis(ctx, userId)
	if nil != err {
		log.Println(err.Error())
	}
	userFollowers := make([]User, userFollowingsCnt)
	err1 := followService.BuildUser(userId, userFollowers, userFollowersId, 0)
	if nil != err1 {
		log.Println(err1.Error())
	}

	return userFollowers, nil
}

// GetFollowingsByRedis 从redis获取登陆用户关注列表
func GetFollowerByRedis(ctx context.Context, userId int64) ([]int64, int64, error) {
	followDao := dao.NewFollowDao(cache.Ctx)
	// 判定键是否存在
	keyCnt, err := cache.UserFollowers.Exists(ctx, strconv.FormatInt(userId, 10)).Result()

	if err != nil {
		log.Println(err.Error())
	}

	// 若键存在，获取缓存数据后返回
	if keyCnt > 0 {
		ids := cache.UserFollowers.SMembers(ctx, strconv.FormatInt(userId, 10)).Val()
		idsInt64, _ := convertToInt64Array(ids)

		return idsInt64, int64(len(idsInt64)), nil
	} else {
		// 键不存在，获取数据库数据，更新缓存并返回
		userFollowersId, userFollowersCnt, err1 := followDao.GetFollowersInfo(userId)
		if err1 != nil {
			log.Println(err1.Error())
		}
		ImportToRDBFollowing(ctx, userId, userFollowersId)
		return userFollowersId, userFollowersCnt, nil
	}

}

// BuildUser 根据传入的id列表和空user数组，构建业务所需user数组并返回
func (followService *FollowSrv) BuildUser(userId int64, users []User, ids []int64, buildtype int) error {
	folowDao := dao.NewFollowDao(cache.Ctx)

	// 遍历传入的用户id，组成user结构体
	for i := 0; i < len(ids); i++ {

		// 用户id赋值
		users[i].Id = ids[i]

		// 用户name赋值
		var err1 error
		users[i].Name, err1 = folowDao.GetUserName(ids[i])
		if nil != err1 {
			log.Println(err1)
			return err1
		}

		// 用户关注数赋值
		var err2 error
		users[i].FollowCount, err2 = folowDao.GetFollowingCnt(ids[i])
		if nil != err2 {
			log.Println(err2.Error())
			return err2
		}

		// 用户粉丝数赋值
		var err3 error
		users[i].FollowerCount, err3 = folowDao.GetFollowerCnt(ids[i])
		if nil != err3 {
			log.Println(err3.Error())
			return err3
		}

	}
	return nil
}

// GetFollowingsByRedis 从redis获取登陆用户关注列表
func GetFollowingsByRedis(userId int64) ([]int64, int64, error) {
	followDao := dao.NewFollowDao(cache.Ctx)
	// 判定键是否存在
	keyCnt, err := cache.UserFollowings.Exists(cache.Ctx, strconv.FormatInt(userId, 10)).Result()

	if err != nil {
		log.Println(err.Error())
	}

	// 若键存在，获取缓存数据后返回
	if keyCnt > 0 {
		ids := cache.UserFollowings.SMembers(cache.Ctx, strconv.FormatInt(userId, 10)).Val()
		idsInt64, _ := convertToInt64Array(ids)

		return idsInt64, int64(len(idsInt64)), nil
	} else {
		// 键不存在，获取数据库数据，更新缓存并返回
		userFollowingsId, userFollowingsCnt, err1 := followDao.GetFollowingsInfo(userId)
		if err1 != nil {
			log.Println(err1.Error())
		}
		ImportToRDBFollowing(cache.Ctx, userId, userFollowingsId)
		return userFollowingsId, userFollowingsCnt, nil
	}

}

func convertToInt64Array(strArr []string) ([]int64, error) {
	int64Arr := make([]int64, len(strArr))
	for i, str := range strArr {
		int64Val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		int64Arr[i] = int64Val
	}
	return int64Arr, nil
}
